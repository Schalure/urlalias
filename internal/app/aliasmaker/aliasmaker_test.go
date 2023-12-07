package aliasmaker

import (
	"errors"
	"testing"

	"github.com/Schalure/urlalias/cmd/shortener/config"
	"github.com/Schalure/urlalias/internal/app/aliaslogger/zaplogger"
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
