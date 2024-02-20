package aliasmaker

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Schalure/urlalias/internal/app/models/aliasentity"
)

const aliasKeyLen int = 9

type Loggerer interface {
	Info(args ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	Fatalw(msg string, keysAndValues ...interface{})
	Close()
}

// Access interface to storage
type Storager interface {
	CreateUser() (uint64, error)
	Save(ctx context.Context, urlAliasNode *aliasentity.AliasURLModel) error
	SaveAll(ctx context.Context, urlAliasNode []aliasentity.AliasURLModel) error
	FindByShortKey(ctx context.Context, shortKey string) (*aliasentity.AliasURLModel, error)
	FindByLongURL(ctx context.Context, longURL string) (*aliasentity.AliasURLModel, error)
	FindByUserID(ctx context.Context, userID uint64) ([]aliasentity.AliasURLModel, error)
	MarkDeleted(ctx context.Context, aliasesID []uint64) error
	GetLastShortKey() string
	IsConnected() bool
	Close() error
}


// Type of service
type AliasMakerServise struct {

	Logger  Loggerer
	Storage Storager

	deleter           *deleter
	aliasesToDeleteCh chan struct {
		userID  uint64
		aliases []string
	}

	lastKey string
}



//	Constructor
func New(s Storager, l Loggerer) (*AliasMakerServise, error) {

	aliasesToDeleteCh := make(chan struct {
		userID  uint64
		aliases []string
	}, 50)
	deleter := newDeleter(cancel, s, l, aliasesToDeleteCh)
	deleter.run(ctx)

	return &AliasMakerServise{
		Storage: s,
		Logger:  l,

		lastKey: s.GetLastShortKey(),

		deleter:           deleter,
		aliasesToDeleteCh: aliasesToDeleteCh,
	}, nil
}


//	GetOriginalURL returns original url by shortKey. If original url not found or was deleted, return error
func (s *AliasMakerServise) GetOriginalURL(ctx context.Context, shortKey string) (string, error) {

	c, cancel := context.WithTimeout(ctx, time.Second * 1)
	defer cancel()

	node, err := s.Storage.FindByShortKey(c, shortKey)
	if err != nil {
		s.Logger.Infow(
			"original url not found", 
			"short key", shortKey, 
			"error", err,
		)
		return "", ErrURLNotFound
	}

	if node.DeletedFlag {
		return "", ErrURLWasDeleted
	}

	return node.LongURL, nil
}


//	AddNewURL add new URL to service and return alias entity
func (s *AliasMakerServise) GetShortKey(ctx context.Context, userID uint64, originalURL string) (string, error) {

	ctxFind, cancelFind := context.WithTimeout(ctx, time.Second * 1)
	defer cancelFind()
	node, err := s.Storage.FindByLongURL(ctxFind, originalURL)
	if err != nil {
		node, err := s.NewAliasEntity(userID, originalURL)
		if err != nil {
			s.Logger.Errorw("error by create new short key", "error", err, "last key", s.lastKey)
			return "", ErrInternal
		}

		ctxSave, cancelSave := context.WithTimeout(ctx, time.Second * 1)
		defer cancelSave()
		err = s.Storage.Save(ctxSave, node)
		if err != nil {
			s.Logger.Errorw("error by save new entity of alias", "error", err, "last key", s.lastKey)
			return "", ErrInternal
		}
		return node.ShortKey, nil
	}
	return node.ShortKey, ErrConflictURL
}


//	GetBatchShortURL create batch of aliases and return batch of short keys
func (s *AliasMakerServise) GetBatchShortURL(ctx context.Context, userID uint64, batchOriginalURL []string) ([]string, error) {

	batchShortURL := make([]string, len(batchOriginalURL))
	var batchNodesToSave []aliasentity.AliasURLModel

	for i, originalURL := range batchOriginalURL {
		ctxFind, cancelFind := context.WithTimeout(ctx, time.Second * 1)
		node, err := s.Storage.FindByLongURL(ctxFind, originalURL)
		cancelFind()
		if err != nil {
			node, err = s.NewAliasEntity(userID, originalURL)
			if err != nil {
				s.Logger.Errorw("error by create new short key", "error", err, "last key", s.lastKey)
				return nil, ErrInternal
			}
			batchNodesToSave = append(batchNodesToSave, *node)
		}
		batchShortURL[i] = node.ShortKey
	}

	ctxSaveAll, cancelSaveAll := context.WithTimeout(ctx, time.Second * 1)
	defer cancelSaveAll()
	if err := s.Storage.SaveAll(ctxSaveAll, batchNodesToSave); err != nil {
		s.Logger.Errorw("can't save all URLs", "error", err)
		return nil, ErrInternal
	}

	return batchShortURL, nil
}


//	GetUserAliases returns all aliases which user created
func (s *AliasMakerServise) GetUserAliases(ctx context.Context, userID uint64) ([]aliasentity.AliasURLModel, error) {

	ctxGetAliases, cancelGetAliases := context.WithTimeout(ctx, time.Second * 1)
	defer cancelGetAliases()

	nodes, err := s.Storage.FindByUserID(ctxGetAliases, userID)
	if err != nil {
		s.Logger.Errorw("can't found aliases by user ID", "error", err, "user ID", userID)
		return nil, ErrInternal
	}
	return nodes, nil
}


//	Create new URL pair
func (s *AliasMakerServise) NewAliasEntity(userID uint64, longURL string) (*aliasentity.AliasURLModel, error) {

	newAliasKey, err := createAliasKey(s.lastKey)
	if err != nil {
		return nil, err
	}
	s.lastKey = newAliasKey
	return &aliasentity.AliasURLModel{
		LongURL:  longURL,
		ShortKey: newAliasKey,
		UserID: userID,
	}, nil
}


//	Add aliases to delete
func (s *AliasMakerServise) AddAliasesToDelete(ctx context.Context, userID uint64, aliases ...string) error {

	select {
	case <-ctx.Done():
		s.Logger.Infow("AddAliasesToDelete: context Done","userID", userID,"aliases", aliases)
		return fmt.Errorf("can't create a delete request, try again later")
	case s.aliasesToDeleteCh <- struct {
		userID  uint64
		aliases []string
	}{userID: userID, aliases: aliases}:
		s.Logger.Infow("AddAliasesToDelete: add aliases to delete", "userID", userID, "aliases", aliases)
	}
	return nil
}


//	Create new user
func (s *AliasMakerServise) CreateUser() (uint64, error) {

	userID, err := s.Storage.CreateUser()
	if err != nil {
		return 0, err
	}
	return userID, nil
}


//	Stop service and full release
func (s *AliasMakerServise) Stop() {

	s.deleter.stop()
	s.Storage.Close()
	s.Logger.Close()
}


//	Make short alias from URL
func createAliasKey(lastKey string) (string, error) {

	var charset = []string{
		"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
		"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
		"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
	}

	if lastKey == "" {
		lastKey = strings.Repeat("0", aliasKeyLen) //"000000000"
		return lastKey, nil
	}

	newKey := strings.Split(lastKey, "")
	if len(newKey) != aliasKeyLen {
		return "", fmt.Errorf("a non-valid key was received from the repository: %s", lastKey)
	}

	for i := aliasKeyLen - 1; i > 0; i-- {
		for n, char := range charset {
			if newKey[i] == char {
				if n == len(charset)-1 {
					newKey[i] = charset[0]
					break
				} else {
					newKey[i] = charset[n+1]
					lastKey = strings.Join(newKey, "")
					return lastKey, nil
				}
			}
		}
	}
	return "", fmt.Errorf("it is impossible to generate a new string because the storage is full")
}
