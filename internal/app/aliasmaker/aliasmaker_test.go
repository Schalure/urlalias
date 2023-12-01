package aliasmaker

import (
	"errors"
	"testing"

	"github.com/Schalure/urlalias/cmd/shortener/config"
	"github.com/Schalure/urlalias/internal/app/aliaslogger/zaplogger"
	"github.com/Schalure/urlalias/internal/app/models"
	"github.com/Schalure/urlalias/internal/app/storage/memstor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)


func newService(t *testing.T) *AliasMakerServise {

	logger, err := zaplogger.NewZapLogger("")
	require.NoError(t, err)
	stor, err := memstor.NewStorage()
	require.NoError(t, err)
	s, err := NewAliasMakerServise(config.NewConfig(), stor, logger)
	require.NoError(t, err)
	return s
}

func Test_createAliasKey(t *testing.T) {

	testCases := []struct {
		name    string
		lastKey string
		want    struct {
			newKey string
			err    error
		}
	}{
		{
			name:    "symple test",
			lastKey: "000000000",
			want: struct {
				newKey string
				err    error
			}{
				newKey: "000000001",
				err:    nil,
			},
		},
		{
			name:    "overload test",
			lastKey: "00000000Z",
			want: struct {
				newKey string
				err    error
			}{
				newKey: "000000010",
				err:    nil,
			},
		},
		{
			name:    "storage full test",
			lastKey: "ZZZZZZZZZ",
			want: struct {
				newKey string
				err    error
			}{
				newKey: "",
				err:    errors.New("it is impossible to generate a new string because the storage is full"),
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {


			s := newService(t)
			defer s.Stop()

			s.lastKey = test.lastKey

			aliasKey, err := s.createAliasKey()

			assert.Equal(t, aliasKey, test.want.newKey)
			assert.Equal(t, err, test.want.err)
		})
	}
}

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

			s := newService(t)
			defer s.Stop()

			err := s.Storage.SaveAll(test.aliases)
			require.NoError(t, err)

			s.DeleteUserURLs(test.userID, test.shortKeys)

			for _, alias := range test.aliases {

				node := s.Storage.FindByShortKey(alias.ShortKey)

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
