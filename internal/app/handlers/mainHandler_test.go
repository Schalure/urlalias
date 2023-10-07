package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Schalure/urlalias/cmd/shortener/config"
	"github.com/Schalure/urlalias/models"
	"github.com/Schalure/urlalias/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)


func Test_mainHandlerMethodGet(t *testing.T) {

	listOfURL := []models.AliasURLModel{
		{ID: 0, LongURL: "https://ya.ru", ShortKey: "123456789"},
		{ID: 1, LongURL: "https://google.com", ShortKey: "987654321"},
	}

	//	create storage
	stor := repositories.NewStorageURL()
	if _, err := stor.Save(models.AliasURLModel{ID: 0, LongURL: listOfURL[0].LongURL, ShortKey: listOfURL[0].ShortKey}); err != nil{
		require.NotNil(t, err)
	}
	if _, err := stor.Save(models.AliasURLModel{ID: 1, LongURL: listOfURL[1].LongURL, ShortKey: listOfURL[1].ShortKey}); err != nil{
		require.NotNil(t, err)
	}


	testCases := []struct{
		name string
		request struct{
			contentType string
			requestURI string
		}
		want struct{
			code int
			response string
		}
	}{
		//----------------------------------
		//	1. simple test
		{
			name: "simple test",
			request: struct{contentType string; requestURI string}{
				contentType: "text/plain",
				requestURI: listOfURL[0].ShortKey,
			},
			want: struct{code int; response string}{
				code: http.StatusTemporaryRedirect,
				response: listOfURL[0].LongURL,
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/" + tt.request.requestURI, nil)
			request.Header.Add("Content-type", tt.request.contentType)

			recorder := httptest.NewRecorder()
			h := http.HandlerFunc(MainHandlerMethodGet(stor))
            h(recorder, request)

			result := recorder.Result()
            err := result.Body.Close()
            require.NoError(t, err)
			
			//	check status code
			assert.Equal(t, tt.want.code, result.StatusCode)

			//	check response
			assert.Equal(t, tt.want.response, result.Header.Get("Location"))
		})
	}
}



func Test_mainHandlerMethodPost(t *testing.T) {

	listOfURL := []models.AliasURLModel{
		{ID: 0, LongURL: "https://ya.ru", ShortKey: "123456789"},
		{ID: 1, LongURL: "https://google.com", ShortKey: "987654321"},
		{ID: 2, LongURL: "https://go.dev", ShortKey: ""},
	}

	//	create storage
	stor := repositories.NewStorageURL()
	if _, err := stor.Save(models.AliasURLModel{ID: 0, LongURL: listOfURL[0].LongURL, ShortKey: listOfURL[0].ShortKey}); err != nil{
		require.NotNil(t, err)
	}
	if _, err := stor.Save(models.AliasURLModel{ID: 1, LongURL: listOfURL[1].LongURL, ShortKey: listOfURL[1].ShortKey}); err != nil{
		require.NotNil(t, err)
	}


	testCases := []struct{
		name string
		request struct{
			contentType string
			requestURI string
		}
		want struct{
			code int
			contentType string
			response string
		}
	}{
		//----------------------------------
		//	1. simple test
		{
			name: "simple test",
			request: struct{contentType string; requestURI string}{
				contentType: "text/plain",
				requestURI: listOfURL[0].LongURL,
			},
			want: struct{code int; contentType string; response string}{
				code: http.StatusCreated,
				contentType: "text/plain",
				response: "http://" + *config.ResponseBaseAddres + "/" + listOfURL[0].ShortKey,
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, *config.Host, strings.NewReader(tt.request.requestURI))
			request.Header.Add("Content-type", tt.request.contentType)

			recorder := httptest.NewRecorder()
			h := http.HandlerFunc(MainHandlerMethodPost(stor))
            h(recorder, request)

			result := recorder.Result()

			//	check status code
			assert.Equal(t, tt.want.code, result.StatusCode)

			//	check response
			data, err := io.ReadAll(recorder.Body)
            require.NoError(t, err)
            err = result.Body.Close()
            require.NoError(t, err)

			assert.Equal(t, tt.want.response, string(data))

			assert.Contains(t, recorder.Header().Get("Content-type"), tt.want.contentType)
		})
	}
}

