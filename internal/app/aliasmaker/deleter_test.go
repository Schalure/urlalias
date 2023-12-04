package aliasmaker

import (
	"context"
	"testing"

	"github.com/Schalure/urlalias/internal/app/aliaslogger/zaplogger"
	"github.com/Schalure/urlalias/internal/app/models"
	"github.com/Schalure/urlalias/internal/app/storage/memstor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_deleteUserURLs(t *testing.T) {

	testCases := []struct {
		name      string
		userID    uint64
		aliases   []models.AliasURLModel
		shortKeys []string
	}{
		{
			name:   "sympleTest",
			userID: 1,
			aliases: []models.AliasURLModel{
				{
					ID:          1,
					UserID:      1,
					LongURL:     "https://no_matter1.com",
					ShortKey:    "000000000",
					DeletedFlag: false,
				},
				{
					ID:          2,
					UserID:      2,
					LongURL:     "https://no_matter2.com",
					ShortKey:    "000000001",
					DeletedFlag: false,
				},
				{
					ID:          3,
					UserID:      1,
					LongURL:     "https://no_matter3.com",
					ShortKey:    "000000002",
					DeletedFlag: false,
				},
			},
			shortKeys: []string{"000000000", "000000001", "000000002"},
		},
		{
			name:   "nil test",
			userID: 1,
			aliases: []models.AliasURLModel{
				{
					ID:          1,
					UserID:      1,
					LongURL:     "https://no_matter1.com",
					ShortKey:    "000000000",
					DeletedFlag: false,
				},
			},
			shortKeys: []string{"000000000", "000000001", "000000002"},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {

			ctx, cancel := context.WithCancel(context.Background())
			
			logger, err := zaplogger.NewZapLogger("")
			require.NoError(t, err)
			stor, err := memstor.NewStorage()
			require.NoError(t, err)

			aliasesToDeleteCh := make(chan struct{userID uint64; aliases []string}, 1)

			d := newDeleter(cancel, stor, logger, aliasesToDeleteCh)

			err = stor.SaveAll(test.aliases)
			require.NoError(t, err)

			d.deleteUserURLs(ctx, test.userID, test.shortKeys)

			for _, alias := range test.aliases {

				node := d.storage.FindByShortKey(alias.ShortKey)

				assert.Equal(t, node.ID, alias.ID)
				assert.Equal(t, node.UserID, alias.UserID)
				assert.Equal(t, node.LongURL, alias.LongURL)
				assert.Equal(t, node.ShortKey, alias.ShortKey)

				if test.userID == node.UserID {
					assert.Equal(t, node.DeletedFlag, true)
				} else {
					assert.Equal(t, node.DeletedFlag, false)
				}
			}
		})
	}
}
