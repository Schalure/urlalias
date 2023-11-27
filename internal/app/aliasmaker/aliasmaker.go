package aliasmaker

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/Schalure/urlalias/cmd/shortener/config"
	"github.com/Schalure/urlalias/internal/app/aliaslogger/zaplogger"
	"github.com/Schalure/urlalias/internal/app/models"
	"github.com/Schalure/urlalias/internal/app/storage/filestor"
	"github.com/Schalure/urlalias/internal/app/storage/memstor"
	"github.com/Schalure/urlalias/internal/app/storage/postgrestor"
)

const (
	aliasKeyLen int = 9
)

// Type of service
type AliasMakerServise struct {
	Config  *config.Configuration
	Logger  Loggerer
	Storage Storager
	lastKey string
}

// --------------------------------------------------
//
//	Constructor
func NewAliasMakerServise(c *config.Configuration) (*AliasMakerServise, error) {

	var errs []error

	logger, loggerErr := chooseLogger(LoggerTypeZap)
	errs = append(errs, loggerErr)

	storage, storageErr := chooseStorage(c)
	errs = append(errs, storageErr)

	if errors.Join(errs...) != nil {
		return nil, errors.Join(errs...)
	}

	lastKey := storage.GetLastShortKey()

	return &AliasMakerServise{
		Config:  c,
		Logger:  logger,
		Storage: storage,
		lastKey: lastKey,
	}, nil
}

// --------------------------------------------------
//
//	Choose logger for service
func chooseLogger(loggerType LoggerType) (Loggerer, error) {
	switch loggerType {
	case LoggerTypeZap:
		return zaplogger.NewZapLogger("")
	default:
		return nil, fmt.Errorf("logger type is not supported: %s", loggerType.String())
	}
}

// --------------------------------------------------
//
//	Choose storage for service
func chooseStorage(c *config.Configuration) (Storager, error) {

	switch c.StorageType() {
	case config.DataBaseStor:
		return postgrestor.NewStorage(c.DBConnection())
	case config.FileStor:
		return filestor.NewStorage(c.AliasesFile(), c.UsersFile())
	default:
		return memstor.NewStorage()
	}
}

// --------------------------------------------------
//
//	Create new URL pair
func (s *AliasMakerServise) NewPairURL(longURL string) (*models.AliasURLModel, error) {

	newAliasKey, err := s.createAliasKey()
	if err != nil {
		return nil, err
	}

	return &models.AliasURLModel{
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
func (s *AliasMakerServise) CreateAlias(userID uint64, originalURL string) (*models.AliasURLModel, int, error) {

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
//	Delete users URLs
func (s *AliasMakerServise) DeleteUserURLs(userID uint64, shortKeys []string) {

	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	
	inputCh := func () chan string {
		inputCh := make(chan string)
		go func () {
			defer close(inputCh)
			for i, shortKey := range shortKeys {
				select {
				case <- ctx.Done():
					s.Logger.Errorw("func DeleteUserURLs: context deadline", "nums ellements added to inputCh", i)
					return
				case inputCh <- shortKey:
				}
			}
		}()
		return inputCh
	}()

	//	get nodes from DB
	resultChannels := func () []chan models.AliasURLModel {
		
		numWorkers := runtime.NumCPU()
		resultChannels := make([]chan models.AliasURLModel, numWorkers)

		for i := 0; i < numWorkers; i++ {
			resultChannels[i] = func () chan models.AliasURLModel {

				resultCh := make(chan models.AliasURLModel)

				go func (resultCh chan models.AliasURLModel) {

					defer close(resultCh)			
					for shortKey := range inputCh {
						node := s.Storage.FindByShortKey(shortKey)
						if node == nil {
							s.Logger.Infow("func DeleteUserURLs: can't Storage.FindByShortKey", "shortKey", shortKey)
							return
						}
						s.Logger.Info(node)
						select {
						case <-ctx.Done():
							s.Logger.Errorw("func DeleteUserURLs: context deadline", "nums ellements added to work", i)
							return
						case resultCh <- *node:
							s.Logger.Infow("func DeleteUserURLs: write to resultCh", "shortKey", shortKey)
						}
					}
				}(resultCh)
				return resultCh

			}()
		}
		return resultChannels
	}()

	//	get aliases id to mark deleted
	outCh := func () chan models.AliasURLModel {

		var wg sync.WaitGroup
		outCh := make(chan models.AliasURLModel)

		for _, result := range resultChannels {
			wg.Add(1)
			go func (result chan models.AliasURLModel) {
				defer wg.Done()
				for aliasNode := range result {
					select {
					case <-ctx.Done():
						s.Logger.Errorw("func DeleteUserURLs: context deadline")
						return
					case outCh <- aliasNode:
					}
				}
			}(result)
		}

		//	wait all gorutins
		go func ()  {
			wg.Wait()
			close(outCh)
		}()
		return outCh
	}()

	//	mark deleted
	aliasesID := make([]uint64, 0)
	for aliasNode := range outCh {
		if aliasNode.UserID == userID{
			aliasesID = append(aliasesID, aliasNode.ID)
			s.Logger.Infow(
				"DeleteUserURLs choose to delete", 
				"user ID", aliasNode.UserID,
				"alias ID", aliasNode.ID,
				"original URL", aliasNode.LongURL,
			)
		}
	}

	ctx, cancel = context.WithTimeout(ctx, 10 * time.Second)
	defer cancel()
	go func() {
		<-ctx.Done()
		if ctx.Err() == context.DeadlineExceeded {
			s.Logger.Info("DeleteUserURLs context deadline while updating DB")
		}
	}()

	err := s.Storage.MarkDeleted(ctx, aliasesID)
	if err != nil{
		s.Logger.Info(err)
	}
}


// --------------------------------------------------
//
//	Make short alias from URL
func (s *AliasMakerServise) createAliasKey() (string, error) {

	var charset = []string{
		"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
		"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
		"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
	}

	if s.lastKey == "" {
		s.lastKey = "000000000"
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

	s.Storage.Close()
	s.Logger.Close()
}
