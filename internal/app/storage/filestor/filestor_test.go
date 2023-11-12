package filestor

import (
	"bufio"
	"os"
	"testing"

	"github.com/Schalure/urlalias/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileStorage_Save(t *testing.T) {

	file, err := os.CreateTemp("", "storage*.json")
	require.NoError(t, err)
	file.Close()
	defer os.Remove(file.Name())

	stor, _ := NewFileStorage(file.Name())

	testCases := []struct {
		testName string
		storNode storage.AliasURLModel
		want     struct {
			data string
			err  error
		}
	}{
		{
			testName: "simple save",
			storNode: storage.AliasURLModel{
				ID:       1,
				ShortKey: "000000000",
				LongURL:  "https://qqq.ru",
			},
			want: struct {
				data string
				err  error
			}{
				data: `{"uuid":1,"short_url":"000000000","original_url":"https://qqq.ru"}`,
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

			f, err := os.Open(file.Name())
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
