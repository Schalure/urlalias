package filestor

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/Schalure/urlalias/internal/app/storage"
	"github.com/stretchr/testify/require"
)

func TestFileStorage_Save(t *testing.T) {

	file, err := os.CreateTemp("", "storage*.json")
	require.NoError(t, err)
	defer os.Remove(file.Name())

	fmt.Println(file.Name())

	fileContent := storage.AliasURLModel{
		ID: 1,
		ShortKey: "4rSPg8ap",
		LongURL: "http://yandex.ru",
	}

	data, err := json.Marshal(fileContent)
	require.NoError(t, err)
	_, err = file.Write(append(data, '\n'))
	require.NoError(t, err)

	

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {



		})
	}	
	// type args struct {
	// 	urlAliasNode *storage.AliasURLModel
	// }
	// tests := []struct {
	// 	name    string
	// 	s       *FileStorage
	// 	args    args
	// 	wantErr bool
	// }{
	// 	// TODO: Add test cases.
	// }

	// 		if err := tt.s.Save(tt.args.urlAliasNode); (err != nil) != tt.wantErr {
	// 			t.Errorf("FileStorage.Save() error = %v, wantErr %v", err, tt.wantErr)
	// 		}
	// 	})
	// }
}

func createStorFile()(*os.File, string, error){

}

