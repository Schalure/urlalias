package handlers

import (
	"context"
	"fmt"
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

func Test_ApiShortenHandlerPost(t *testing.T) {

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
	//service := aliasmaker.NewAliasMakerServise(testStor)

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
				code:        http.StatusConflict,
				contentType: appJSON,
				response:    fmt.Sprintf("{\"result\":\"%s/%s\"}", testConfig.BaseURL(), listOfURL[0].ShortKey),
			},
		},
	}

	//	Start test cases
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {

			service.Logger.Infow(
				"Response URL",
				"URL", tt.requestURI,
				"test", fmt.Sprintf("\"url\": \"%v\"", listOfURL[0].LongURL),
			)

			request := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(tt.requestURI))
			request.Header.Add("Content-type", appJSON)

			recorder := httptest.NewRecorder()
			h := NewHandlers(service).APIShortenHandlerPost
			ctx := context.WithValue(request.Context(), UserID, uint64(0))

			h(recorder, request.WithContext(ctx))

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
			testName: "memstor_simple_test",
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
				response: fmt.Sprintf(`[{"correlation_id":"1","short_url":"%s/123456789"},{"correlation_id":"2","short_url":"%s/987654321"}]`,
					service.Config.BaseURL(),
					service.Config.BaseURL(),
				),
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.testName, func(t *testing.T) {

			request := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", strings.NewReader(test.data))
			request.Header.Add(contentType, appJSON)

			recorder := httptest.NewRecorder()
			h := NewHandlers(service).APIShortenBatchHandlerPost

			ctx := context.WithValue(request.Context(), UserID, uint64(0))
			h(recorder, request.WithContext(ctx))

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

func Test_APIUserURLsHandlerDelete(t *testing.T) {

	service, err := aliasmaker.NewAliasMakerServise(config.NewConfig())
	require.NoError(t, err)
	defer service.Stop()

	testCases := []struct {
		name   string
		userID uint64
		data   string
		want   struct {
			statusCode int
		}
	}{
		{
			name:   "sympleTest",
			userID: 1,
			data:   `["6qxTVvsy","RTfd56hn","Jlfd67ds"]`,
			want: struct{ statusCode int }{
				statusCode: http.StatusAccepted,
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {

			request := httptest.NewRequest(http.MethodDelete, "/api/user/urls", strings.NewReader(test.data))
			request.Header.Add(contentType, appJSON)

			recorder := httptest.NewRecorder()
			h := NewHandlers(service).APIUserURLsHandlerDelete

			ctx := context.WithValue(request.Context(), UserID, test.userID)
			h(recorder, request.WithContext(ctx))

			result := recorder.Result()
			err = result.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, result.StatusCode, test.want.statusCode)
		})
	}
}
