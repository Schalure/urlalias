package aliasmaker

import (
	"errors"
	"testing"

	"github.com/Schalure/urlalias/cmd/shortener/config"
	"github.com/Schalure/urlalias/internal/app/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

			s, err := NewAliasMakerServise(config.NewConfig())
			require.NoError(t, err)
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
		name string
		userID uint64
		aliases []models.AliasURLModel
	}{
		{
			name: "sympleTest",
			userID: 1,
			aliases: []models.AliasURLModel{
				{
					ID: 1,
					UserID: 1,
					LongURL: "https://no_matter1.com",
					ShortKey: "000000000",
					DeletedFlag: false,
				},
				{
					ID: 2,
					UserID: 2,
					LongURL: "https://no_matter2.com",
					ShortKey: "000000001",
					DeletedFlag: false,
				},
				{
					ID: 3,
					UserID: 1,
					LongURL: "https://no_matter3.com",
					ShortKey: "000000002",
					DeletedFlag: false,
				},
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {

			s, err := NewAliasMakerServise(config.NewConfig())
			require.NoError(t, err)
			defer s.Stop()

			s.Storage.SaveAll(test.aliases)
			require.NoError(t, err)

			aliases := make([]string, 0)
			for _, alias := range test.aliases {
				aliases = append(aliases, alias.ShortKey)
			}

			s.DeleteUserURLs(test.userID, aliases)

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
