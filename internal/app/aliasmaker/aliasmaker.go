package aliasmaker

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/Schalure/urlalias/internal/app/aliaslogger/zaplogger"
	"github.com/Schalure/urlalias/internal/app/models/aliasentity"
)

const aliasKeyLen int = 9

// //go:generate mockgen -destination=../mocks/mock_loggerer.go -package=mocks github.com/Schalure/urlalias/internal/app/aliasmaker Loggerer
// type Loggerer interface {
// 	Info(args ...interface{})
// 	Infow(msg string, keysAndValues ...interface{})
// 	Errorw(msg string, keysAndValues ...interface{})
// 	Fatalw(msg string, keysAndValues ...interface{})
// 	Close()
// }

// Access interface to storage
//
//go:generate mockgen -destination=../mocks/mock_storager.go -package=mocks github.com/Schalure/urlalias/internal/app/aliasmaker Storager
type Storager interface {
	CreateUser() (uint64, error)
	Save(ctx context.Context, urlAliasNode *aliasentity.AliasURLModel) error
	SaveAll(ctx context.Context, urlAliasNodes []aliasentity.AliasURLModel) error
	FindByShortKey(ctx context.Context, shortKey string) (*aliasentity.AliasURLModel, error)
	FindByLongURL(ctx context.Context, longURL string) (*aliasentity.AliasURLModel, error)
	FindByUserID(ctx context.Context, userID uint64) ([]aliasentity.AliasURLModel, error)
	MarkDeleted(ctx context.Context, aliasesID []uint64) error
	GetLastShortKey() string
	IsConnected() bool
	Close() error
}

type Deleter struct {
	userID  uint64
	aliases []string
}

// Type of service
type AliasMakerServise struct {
	logger  *zaplogger.ZapLogger
	storage Storager

	deleterCh chan Deleter
	lastKey   string
}

// Constructor
func New(s Storager, l *zaplogger.ZapLogger) (*AliasMakerServise, error) {

	return &AliasMakerServise{
		storage: s,
		logger:  l,

		lastKey:   s.GetLastShortKey(),
		deleterCh: make(chan Deleter, 50),
	}, nil
}

// GetOriginalURL returns original url by shortKey. If original url not found or was deleted, return error
func (s *AliasMakerServise) GetOriginalURL(ctx context.Context, shortKey string) (string, error) {

	c, cancel := context.WithTimeout(ctx, time.Second*1)
	defer cancel()

	node, err := s.storage.FindByShortKey(c, shortKey)
	if err != nil {
		s.logger.Infow(
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

// AddNewURL add new URL to service and return alias entity
func (s *AliasMakerServise) GetShortKey(ctx context.Context, userID uint64, originalURL string) (string, error) {

	ctxFind, cancelFind := context.WithTimeout(ctx, time.Second*1)
	defer cancelFind()
	node, err := s.storage.FindByLongURL(ctxFind, originalURL)
	if err != nil {
		node, err := s.NewAliasEntity(userID, originalURL)
		if err != nil {
			s.logger.Errorw("error by create new short key", "error", err, "last key", s.lastKey)
			return "", ErrInternal
		}

		ctxSave, cancelSave := context.WithTimeout(ctx, time.Second*1)
		defer cancelSave()
		err = s.storage.Save(ctxSave, node)
		if err != nil {
			s.logger.Errorw("error by save new entity of alias", "error", err, "last key", s.lastKey)
			return "", ErrInternal
		}
		return node.ShortKey, nil
	}
	return node.ShortKey, ErrConflictURL
}

// GetBatchShortURL create batch of aliases and return batch of short keys
func (s *AliasMakerServise) GetBatchShortURL(ctx context.Context, userID uint64, batchOriginalURL []string) ([]string, error) {

	batchShortURL := make([]string, len(batchOriginalURL))
	var batchNodesToSave []aliasentity.AliasURLModel

	for i, originalURL := range batchOriginalURL {
		ctxFind, cancelFind := context.WithTimeout(ctx, time.Second*1)
		node, err := s.storage.FindByLongURL(ctxFind, originalURL)
		cancelFind()
		if err != nil {
			node, err = s.NewAliasEntity(userID, originalURL)
			if err != nil {
				s.logger.Errorw("error by create new short key", "error", err, "last key", s.lastKey)
				return nil, ErrInternal
			}
			batchNodesToSave = append(batchNodesToSave, *node)
		}
		batchShortURL[i] = node.ShortKey
	}

	ctxSaveAll, cancelSaveAll := context.WithTimeout(ctx, time.Second*1)
	defer cancelSaveAll()
	if err := s.storage.SaveAll(ctxSaveAll, batchNodesToSave); err != nil {
		s.logger.Errorw("can't save all URLs", "error", err)
		return nil, ErrInternal
	}

	return batchShortURL, nil
}

// Create new user
func (s *AliasMakerServise) CreateUser() (uint64, error) {

	userID, err := s.storage.CreateUser()
	if err != nil {
		return 0, err
	}
	return userID, nil
}

// GetUserAliases returns all aliases which user created
func (s *AliasMakerServise) GetUserAliases(ctx context.Context, userID uint64) ([]aliasentity.AliasURLModel, error) {

	ctxGetAliases, cancelGetAliases := context.WithTimeout(ctx, time.Second*1)
	defer cancelGetAliases()

	nodes, err := s.storage.FindByUserID(ctxGetAliases, userID)
	if err != nil {
		s.logger.Errorw("can't found aliases by user ID", "error", err, "user ID", userID)
		return nil, ErrInternal
	}
	return nodes, nil
}

// Add aliases to delete
func (s *AliasMakerServise) AddAliasesToDelete(ctx context.Context, userID uint64, aliases ...string) error {

	select {
	case <-ctx.Done():
		s.logger.Infow("AddAliasesToDelete: context Done", "userID", userID, "aliases", aliases)
		return fmt.Errorf("can't create a delete request, try again later")
	case s.deleterCh <- Deleter{userID: userID, aliases: aliases}:
		s.logger.Infow("AddAliasesToDelete: add aliases to delete", "userID", userID, "aliases", aliases)
	}
	return nil
}

func (s *AliasMakerServise) IsDatabaseActive() bool {

	return s.storage.IsConnected()
}

// Create new URL pair
func (s *AliasMakerServise) NewAliasEntity(userID uint64, longURL string) (*aliasentity.AliasURLModel, error) {

	newAliasKey, err := createAliasKey(s.lastKey)
	if err != nil {
		return nil, err
	}
	s.lastKey = newAliasKey
	return &aliasentity.AliasURLModel{
		LongURL:  longURL,
		ShortKey: newAliasKey,
		UserID:   userID,
	}, nil
}

func (s *AliasMakerServise) deleteWorker(ctx context.Context) {

	go func() {
		for {
			select {
			case <-ctx.Done():
				s.logger.Info("deleteWorker stopped by ctx.Done()")
				return
			case deleter := <-s.deleterCh:
				s.deleteAliases(ctx, deleter.userID, deleter.aliases)
			}
		}
	}()
}

func (s *AliasMakerServise) deleteAliases(ctx context.Context, userID uint64, shortKeys []string) []string {

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	inputCh := func() chan string {
		inputCh := make(chan string)
		go func() {
			defer close(inputCh)
			for i, shortKey := range shortKeys {
				select {
				case <-ctx.Done():
					s.logger.Errorw("func DeleteUserURLs: context deadline", "nums ellements added to inputCh", i)
					return
				case inputCh <- shortKey:
				}
			}
		}()
		return inputCh
	}()

	//	get nodes from DB
	resultChannels := func() []chan aliasentity.AliasURLModel {

		numWorkers := runtime.NumCPU()
		resultChannels := make([]chan aliasentity.AliasURLModel, numWorkers)

		for i := 0; i < numWorkers; i++ {
			resultChannels[i] = func() chan aliasentity.AliasURLModel {

				resultCh := make(chan aliasentity.AliasURLModel)

				go func(resultCh chan aliasentity.AliasURLModel) {

					defer close(resultCh)
					for shortKey := range inputCh {
						node, err := s.storage.FindByShortKey(ctx, shortKey)
						if err != nil {
							s.logger.Infow("func DeleteUserURLs: can't Storage.FindByShortKey", "shortKey", shortKey)
							break
						}
						s.logger.Info(node)
						select {
						case <-ctx.Done():
							s.logger.Errorw("func DeleteUserURLs: context deadline", "nums ellements added to work", i)
							return
						case resultCh <- *node:
							s.logger.Infow("func DeleteUserURLs: write to resultCh", "shortKey", shortKey)
						}
					}
				}(resultCh)
				return resultCh

			}()
		}
		return resultChannels
	}()

	//	get aliases id to mark deleted
	outCh := func() chan aliasentity.AliasURLModel {

		var wg sync.WaitGroup
		outCh := make(chan aliasentity.AliasURLModel)

		for _, result := range resultChannels {
			wg.Add(1)
			go func(result chan aliasentity.AliasURLModel) {
				defer wg.Done()
				for aliasNode := range result {
					select {
					case <-ctx.Done():
						s.logger.Errorw("func DeleteUserURLs: context deadline")
						return
					case outCh <- aliasNode:
					}
				}
			}(result)
		}

		//	wait all gorutins
		go func() {
			wg.Wait()
			close(outCh)
		}()
		return outCh
	}()

	//	mark deleted
	aliasesID := make([]uint64, 0)
	deleteAliases := make([]string, 0)
	for aliasNode := range outCh {
		if aliasNode.UserID != userID {
			s.logger.Infow(
				"Can't delete alias due to ID mismatch",
				"expected user ID", userID,
				"actual user ID", aliasNode.UserID,
				"alias ID", aliasNode.ID,
				"original URL", aliasNode.LongURL,
			)
			continue
		}
		aliasesID = append(aliasesID, aliasNode.ID)
		deleteAliases = append(deleteAliases, aliasNode.ShortKey)
		s.logger.Infow(
			"DeleteUserURLs choose to delete",
			"user ID", aliasNode.UserID,
			"alias ID", aliasNode.ID,
			"original URL", aliasNode.LongURL,
		)
	}

	ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	go func() {
		<-ctx.Done()
		if ctx.Err() == context.DeadlineExceeded {
			s.logger.Info("DeleteUserURLs context deadline while updating DB")
		}
	}()

	err := s.storage.MarkDeleted(ctx, aliasesID)
	if err != nil {
		s.logger.Info(err)
	}
	return deleteAliases
}

func (s *AliasMakerServise) Run(ctx context.Context) {
	s.deleteWorker(ctx)
}

// Stop service and full release
func (s *AliasMakerServise) Stop() {

	s.storage.Close()
	s.logger.Close()
}

// Make short alias from URL
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
