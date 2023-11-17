package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Schalure/urlalias/cmd/shortener/config"
	"github.com/Schalure/urlalias/internal/app/aliasmaker"
	"github.com/Schalure/urlalias/internal/app/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ------------------------------------------------------------
//
//	Test mainHandlerMethodGet: "/{shortKey}"
func Test_mainHandlerMethodGet(t *testing.T) {

	service, err := aliasmaker.NewAliasMakerServise(config.NewConfig())
	require.NoError(t, err)
	defer service.Stop()

	var listOfURL = []models.AliasURLModel{
		{ID: 0, LongURL: "https://ya.ru", ShortKey: "123456789"},
		{ID: 1, LongURL: "https://google.com", ShortKey: "987654321"},
	}

	for i, nodeURL := range listOfURL {
		if err := service.Storage.Save(&models.AliasURLModel{ID: uint64(i), LongURL: nodeURL.LongURL, ShortKey: nodeURL.ShortKey}); err != nil {
			require.NotNil(t, err)
		}
	}

	//	Test cases
	testCases := []struct {
		name    string
		request struct {
			contentType string
			requestURI  string
		}
		want struct {
			code     int
			response string
		}
	}{
		//----------------------------------
		//	1. simple test
		{
			name: "simple test",
			request: struct {
				contentType string
				requestURI  string
			}{
				contentType: textPlain,
				requestURI:  listOfURL[0].ShortKey,
			},
			want: struct {
				code     int
				response string
			}{
				code:     http.StatusTemporaryRedirect,
				response: listOfURL[0].LongURL,
			},
		},
	}

	//	Start test cases
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			request := httptest.NewRequest(http.MethodGet, "/"+testCase.request.requestURI, nil)
			request.Header.Add("Content-type", testCase.request.contentType)

			recorder := httptest.NewRecorder()
			h := NewHandlers(service).mainHandlerGet
			h(recorder, request)

			result := recorder.Result()
			err := result.Body.Close()
			require.NoError(t, err)

			//	check status code
			assert.Equal(t, testCase.want.code, result.StatusCode)

			//	check response
			assert.Equal(t, testCase.want.response, result.Header.Get("Location"))
		})
	}
}

func Test_mainHandlerMethodPost(t *testing.T) {

	service, err := aliasmaker.NewAliasMakerServise(config.NewConfig())
	require.NoError(t, err)
	defer service.Stop()

	listOfURL := []models.AliasURLModel{
		{ID: 0, LongURL: "https://ya.ru", ShortKey: "123456789"},
		{ID: 1, LongURL: "https://google.com", ShortKey: "987654321"},
		{ID: 2, LongURL: "https://go.dev", ShortKey: ""},
	}

	for i, nodeURL := range listOfURL {
		if err := service.Storage.Save(&models.AliasURLModel{ID: uint64(i), LongURL: nodeURL.LongURL, ShortKey: nodeURL.ShortKey}); err != nil {
			require.NotNil(t, err)
		}
	}


	testCases := []struct {
		name    string
		request struct {
			contentType string
			requestURI  string
		}
		want struct {
			code        int
			contentType string
			response    string
		}
	}{
		//----------------------------------
		//	1. simple test
		{
			name: "simple test",
			request: struct {
				contentType string
				requestURI  string
			}{
				contentType: textPlain,
				requestURI:  listOfURL[0].LongURL,
			},
			want: struct {
				code        int
				contentType string
				response    string
			}{
				code:        http.StatusConflict,
				contentType: textPlain,
				response:    service.Config.BaseURL() + "/" + listOfURL[0].ShortKey,
			},
		},
	}

	//	Start test cases
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.request.requestURI))
			request.Header.Add("Content-type", tt.request.contentType)

			recorder := httptest.NewRecorder()
			h := NewHandlers(service).mainHandlerPost
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
