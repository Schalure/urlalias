package aliasmaker

import (
	"errors"
	"testing"

	"github.com/Schalure/urlalias/internal/app/models"
	"github.com/Schalure/urlalias/internal/app/storage/memstor"
	"github.com/stretchr/testify/assert"
)

func Test_createAliasKey(t *testing.T) {

	testCases := []struct{
		name string
		lastKey string
		want struct{
			newKey string
			err error
		}
	}{
		{
			name: "symple test",
			lastKey: "000000000",
			want: struct{newKey string; err error}{
				newKey: "000000001",
				err: nil,
			},
		},
		{
			name: "overload test",
			lastKey: "00000000Z",
			want: struct{newKey string; err error}{
				newKey: "000000010",
				err: nil,
			},
		},
		{
			name: "storage full test",
			lastKey: "ZZZZZZZZZ",
			want: struct{newKey string; err error}{
				newKey: "",
				err: errors.New("it is impossible to generate a new string because the storage is full"),
			},
		},
	}

	for _, test := range testCases{
		t.Run(test.name, func(t *testing.T) {

			stor, _ := memstor.NewMemStorage()
			stor.Save(&models.AliasURLModel{
				ID: 0,
				ShortKey: test.lastKey,
				LongURL: "",
			})

			s := NewAliasMakerServise(stor)
			aliasKey, err := s.createAliasKey()

			assert.Equal(t, aliasKey, test.want.newKey)
			assert.Equal(t, err, test.want.err)
		})
	}
}