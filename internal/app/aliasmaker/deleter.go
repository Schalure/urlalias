package aliasmaker

import (
	"context"
	"runtime"
	"sync"
	"time"

	"github.com/Schalure/urlalias/internal/app/models/aliasentity"
)

type deleter struct {
	storage Storager
	logger  Loggerer

	aliasesToDelete chan struct {
		userID  uint64
		aliases []string
	}
	cancel context.CancelFunc
}

func newDeleter(cansel context.CancelFunc, s Storager, l Loggerer, aliasesToDelete chan struct {
	userID  uint64
	aliases []string
}) *deleter {

	return &deleter{
		cancel:          cansel,
		storage:         s,
		logger:          l,
		aliasesToDelete: aliasesToDelete,
	}
}

func (d *deleter) run(ctx context.Context) {

	go func() {
		for {
			select {
			case <-ctx.Done():
				d.logger.Infow(
					"func (d *Deleter) RunDeleter(ctx context.Context)",
					"error", "function stopped by ctx.Done()",
				)
			case aliasesToDelete := <-d.aliasesToDelete:
				d.deleteUserURLs(ctx, aliasesToDelete.userID, aliasesToDelete.aliases)
			}
		}
	}()
}

func (d *deleter) stop() {

	d.logger.Infow(
		"func (d *deleter) stop()",
		"message", "deleter stop",
	)
	d.cancel()
}

// --------------------------------------------------
//
//	Delete users URLs
func (d *deleter) deleteUserURLs(ctx context.Context, userID uint64, shortKeys []string) {

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	inputCh := func() chan string {
		inputCh := make(chan string)
		go func() {
			defer close(inputCh)
			for i, shortKey := range shortKeys {
				select {
				case <-ctx.Done():
					d.logger.Errorw("func DeleteUserURLs: context deadline", "nums ellements added to inputCh", i)
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
						node := d.storage.FindByShortKey(shortKey)
						if node == nil {
							d.logger.Infow("func DeleteUserURLs: can't Storage.FindByShortKey", "shortKey", shortKey)
							return
						}
						d.logger.Info(node)
						select {
						case <-ctx.Done():
							d.logger.Errorw("func DeleteUserURLs: context deadline", "nums ellements added to work", i)
							return
						case resultCh <- *node:
							d.logger.Infow("func DeleteUserURLs: write to resultCh", "shortKey", shortKey)
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
						d.logger.Errorw("func DeleteUserURLs: context deadline")
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
	for aliasNode := range outCh {
		if aliasNode.UserID == userID {
			aliasesID = append(aliasesID, aliasNode.ID)
			d.logger.Infow(
				"DeleteUserURLs choose to delete",
				"user ID", aliasNode.UserID,
				"alias ID", aliasNode.ID,
				"original URL", aliasNode.LongURL,
			)
		}
	}

	ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	go func() {
		<-ctx.Done()
		if ctx.Err() == context.DeadlineExceeded {
			d.logger.Info("DeleteUserURLs context deadline while updating DB")
		}
	}()

	err := d.storage.MarkDeleted(ctx, aliasesID)
	if err != nil {
		d.logger.Info(err)
	}
}
