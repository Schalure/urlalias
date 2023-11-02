package filestor

import (
	"bufio"
	"fmt"
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

	stor := NewFileStorage(file.Name())


	testCases := []struct {
		testName string
		storNode storage.AliasURLModel
		want struct {
			data string
			err error
		}
	}{
		{
			testName: "simple save",
			storNode: storage.AliasURLModel{
				ID: 1,
				ShortKey: "555555555",
				LongURL: "https://qqq.ru",
			},
			want: struct{data string; err error}{
				data: `{"uuid":1,"short_url":"555555555","original_url":"https://qqq.ru"}`,
				err: nil,
			},
		},
		{
			testName: "dublicate key save",
			storNode: storage.AliasURLModel{
				ID: 2,
				ShortKey: "555555555",
				LongURL: "https://eee.ru",
			},
			want: struct{data string; err error}{
				data: ``,
				err: fmt.Errorf("the key \"%s\" is already in the database", "555555555"),
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.testName, func(t *testing.T) {

			err := stor.Save(&test.storNode)
			if err != nil{
				assert.Equal(t, test.want.err, err)
				return
			}

			f, err := os.Open(file.Name())
			require.NoError(t, err)
			scanner := bufio.NewScanner(f)

			var resalt = make([]string, 0)
			for scanner.Scan(){
				resalt = append(resalt, scanner.Text())
			}

			assert.Equal(t, resalt[len(resalt) - 1], test.want.data)
		})
	}	
}
