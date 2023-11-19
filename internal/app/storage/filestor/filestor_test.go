package filestor

import (
	"bufio"
	"os"
	"testing"

	"github.com/Schalure/urlalias/internal/app/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileStorage_Save(t *testing.T) {

	aliasesFile, err := os.CreateTemp("", "storage*.json")
	require.NoError(t, err)
	aliasesFile.Close()
	defer os.Remove(aliasesFile.Name())
	
	usersFile, err := os.CreateTemp("", "storage*.json")
	require.NoError(t, err)
	aliasesFile.Close()
	defer os.Remove(usersFile.Name())

	stor, _ := NewStorage(aliasesFile.Name(), usersFile.Name())

	testCases := []struct {
		testName string
		storNode models.AliasURLModel
		want     struct {
			data string
			err  error
		}
	}{
		{
			testName: "simple save",
			storNode: models.AliasURLModel{
				ID:       1,
				UserID: 0,
				ShortKey: "000000000",
				LongURL:  "https://qqq.ru",
			},
			want: struct {
				data string
				err  error
			}{
				data: `{"uuid":1,"user_id":0,"short_url":"000000000","original_url":"https://qqq.ru"}`,
				err:  nil,
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.testName, func(t *testing.T) {

			err := stor.Save(&test.storNode)
			if err != nil {
				assert.Equal(t, test.want.err, err)
				return
			}

			f, err := os.Open(aliasesFile.Name())
			require.NoError(t, err)
			scanner := bufio.NewScanner(f)

			var resalt = make([]string, 0)
			for scanner.Scan() {
				resalt = append(resalt, scanner.Text())
			}

			assert.Equal(t, resalt[len(resalt)-1], test.want.data)
		})
	}
}
