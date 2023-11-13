package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Schalure/urlalias/cmd/shortener/config"
	"github.com/Schalure/urlalias/internal/app/aliasmaker"
	"github.com/Schalure/urlalias/internal/app/storage"
	"github.com/Schalure/urlalias/internal/app/storage/memstor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ApiShortenHandlerPost(t *testing.T) {

	logger, err := NewLogger(LoggerTypeZap)
	require.NoError(t, err)
	defer logger.Close()

	listOfURL := []storage.AliasURLModel{
		{ID: 0, LongURL: "https://ya.ru", ShortKey: "123456789"},
		{ID: 1, LongURL: "https://google.com", ShortKey: "987654321"},
		{ID: 2, LongURL: "https://go.dev", ShortKey: ""},
	}

	testStor, _ := memstor.NewMemStorage()
	for i, nodeURL := range listOfURL {
		if err := testStor.Save(&storage.AliasURLModel{ID: uint64(i), LongURL: nodeURL.LongURL, ShortKey: nodeURL.ShortKey}); err != nil {
			require.NotNil(t, err)
		}
	}
	service := aliasmaker.NewAliasMakerServise(testStor)

	testConfig := config.NewConfig()

	testCases := []struct {
		name       string
		requestURI string
		want       struct {
			code        int
			contentType string
			response    string
		}
	}{
		//----------------------------------
		//	1. simple test
		{
			name:       "simple test",
			requestURI: fmt.Sprintf("{\"url\":\"%s\"}", listOfURL[0].LongURL),
			want: struct {
				code        int
				contentType string
				response    string
			}{
				code:        http.StatusCreated,
				contentType: appJSON,
				response:    fmt.Sprintf("{\"result\":\"%s/%s\"}", testConfig.BaseURL(), listOfURL[0].ShortKey),
			},
		},
	}

	//	Start test cases
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {

			logger.Infow(
				"Response URL",
				"URL", tt.requestURI,
				"test", fmt.Sprintf("\"url\": \"%v\"", listOfURL[0].LongURL),
			)

			request := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(tt.requestURI))
			request.Header.Add("Content-type", appJSON)

			recorder := httptest.NewRecorder()
			h := NewHandlers(service, testConfig, logger).APIShortenHandlerPost
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

func Test_ApiShortenBatchHandlerPost(t *testing.T) {

	logger, err := NewLogger(LoggerTypeZap)
	require.NoError(t, err)
	defer logger.Close()

	listOfURL := []storage.AliasURLModel{
		{ID: 0, LongURL: "https://ya.ru", ShortKey: "123456789"},
		{ID: 1, LongURL: "https://google.com", ShortKey: "987654321"},
		{ID: 2, LongURL: "https://go.dev", ShortKey: ""},
	}

	testStor, _ := memstor.NewMemStorage()
	for i, nodeURL := range listOfURL {
		if err := testStor.Save(&storage.AliasURLModel{ID: uint64(i), LongURL: nodeURL.LongURL, ShortKey: nodeURL.ShortKey}); err != nil {
			require.NotNil(t, err)
		}
	}
	service := aliasmaker.NewAliasMakerServise(testStor)
	testConfig := config.NewConfig()

	testCases := []struct {
		testName string
		data     string
		want     struct {
			code        int
			contentType string
			response    string
		}
	}{
		//	memstor simple test
		{
			testName: "memstor simple test",
			data: `[
				{
					"correlation_id": "1",
					"original_url": "https://ya.ru"
				},
				{
					"correlation_id": "2",
					"original_url": "https://google.com"
				}
			]`,
			want: struct {
				code        int
				contentType string
				response    string
			}{
				code:        http.StatusCreated,
				contentType: appJSON,
				response:    `[{"correlation_id":"1","short_url":"http://localhost:8080/123456789"},{"correlation_id":"2","short_url":"http://localhost:8080/987654321"}]`,
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.testName, func(t *testing.T) {

			request := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", strings.NewReader(test.data))
			request.Header.Add(contentType, appJSON)

			recorder := httptest.NewRecorder()
			h := NewHandlers(service, testConfig, logger).APIShortenBatchHandlerPost
			h(recorder, request)

			result := recorder.Result()

			//	check status code
			assert.Equal(t, test.want.code, result.StatusCode)

			//	check contentType
			assert.Contains(t, recorder.Header().Get("Content-type"), test.want.contentType)

			//	check response
			data, err := io.ReadAll(recorder.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, test.want.response, string(data))
		})
	}
}
