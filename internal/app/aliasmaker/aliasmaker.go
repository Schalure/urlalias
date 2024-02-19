package aliasmaker

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Schalure/urlalias/cmd/shortener/config"
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
	Save(urlAliasNode *aliasentity.AliasURLModel) error
	SaveAll(urlAliasNode []aliasentity.AliasURLModel) error
	FindByShortKey(ctx context.Context, shortKey string) (*aliasentity.AliasURLModel, error)
	FindByLongURL(longURL string) *aliasentity.AliasURLModel
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


// --------------------------------------------------
//
//	Constructor
func New(c *config.Configuration, s Storager, l Loggerer) (*AliasMakerServise, error) {

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
func (s *AliasMakerServise) GetShortURL(ctx context.Context, userID uint64, originalURL string) (string, error) {


}

// --------------------------------------------------
//
//	Create new URL pair
func (s *AliasMakerServise) NewPairURL(longURL string) (*aliasentity.AliasURLModel, error) {

	newAliasKey, err := s.createAliasKey()
	if err != nil {
		return nil, err
	}

	return &aliasentity.AliasURLModel{
		LongURL:  longURL,
		ShortKey: newAliasKey,
	}, nil
}

// --------------------------------------------------
//
//	Create new user
func (s *AliasMakerServise) CreateUser() (uint64, error) {

	userID, err := s.Storage.CreateUser()
	if err != nil {
		return 0, err
	}
	return userID, nil
}

// --------------------------------------------------
//
//	Create alias by originalURL
func (s *AliasMakerServise) CreateAlias(userID uint64, originalURL string) (*aliasentity.AliasURLModel, int, error) {

	var err error

	node := s.Storage.FindByLongURL(originalURL)
	if node == nil {
		if node, err = s.NewPairURL(originalURL); err != nil {
			s.Logger.Info(err.Error())
			return nil, http.StatusBadRequest, err
		}
		node.UserID = userID
		if err = s.Storage.Save(node); err != nil {
			s.Logger.Info(err.Error())
			return nil, http.StatusBadRequest, err
		}
		return node, http.StatusCreated, nil
	}
	return node, http.StatusConflict, nil
}

// --------------------------------------------------
//
//	Add aliases to delete
func (s *AliasMakerServise) AddAliasesToDelete(ctx context.Context, userID uint64, aliases ...string) error {

	select {
	case <-ctx.Done():
		s.Logger.Infow(
			"AddAliasesToDelete: context Done",
			"userID", userID,
			"aliases", aliases,
		)
		return fmt.Errorf("can't create a delete request, try again later")
	case s.aliasesToDeleteCh <- struct {
		userID  uint64
		aliases []string
	}{userID: userID, aliases: aliases}:
		s.Logger.Infow(
			"AddAliasesToDelete: add aliases to delete",
			"userID", userID,
			"aliases", aliases,
		)
	}
	return nil
}

// --------------------------------------------------
//
//	Make short alias from URL
func (s *AliasMakerServise) createAliasKey() (string, error) {

	var charset = []string{
		"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
		"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
		"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
	}

	if s.lastKey == "" {
		s.lastKey = strings.Repeat("0", aliasKeyLen) //"000000000"
		return s.lastKey, nil
	}

	newKey := strings.Split(s.lastKey, "")
	if len(newKey) != aliasKeyLen {
		return "", fmt.Errorf("a non-valid key was received from the repository: %s", s.lastKey)
	}

	for i := aliasKeyLen - 1; i > 0; i-- {
		for n, char := range charset {
			if newKey[i] == char {
				if n == len(charset)-1 {
					newKey[i] = charset[0]
					break
				} else {
					newKey[i] = charset[n+1]
					s.lastKey = strings.Join(newKey, "")
					return s.lastKey, nil
				}
			}
		}
	}
	return "", fmt.Errorf("it is impossible to generate a new string because the storage is full")
}

// --------------------------------------------------
//
//	Stop service and full release
func (s *AliasMakerServise) Stop() {

	s.deleter.stop()
	s.Storage.Close()
	s.Logger.Close()
}
