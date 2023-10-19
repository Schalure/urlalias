package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Schalure/urlalias/cmd/shortener/config"
	"github.com/Schalure/urlalias/internal/app/storage"
	"github.com/Schalure/urlalias/internal/app/storage/memstor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ------------------------------------------------------------
//
//	Test request
//	Input:
//		t *testing.T
//		ts *httptest.Server - test server
//		method string - request method
//		path string - path of request
//	Output:
//		*http.Response - response object
//		string - response body
func testRequest(t *testing.T, ts *httptest.Server, method, contentType, path string) (*http.Response, string) {

	req, err := http.NewRequest(method, ts.URL+"/"+path, nil)
	require.NoError(t, err)
	req.Header.Add("Content-type", contentType)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

// ------------------------------------------------------------
//
//	Test mainHandlerMethodGet: "/{shortKey}"
func Test_mainHandlerMethodGet(t *testing.T) {

	var listOfURL = []storage.AliasURLModel{
		{ID: 0, LongURL: "https://ya.ru", ShortKey: "123456789"},
		{ID: 1, LongURL: "https://google.com", ShortKey: "987654321"},
	}
	testStor := memstor.NewMemStorage()
	for i, nodeURL := range listOfURL {
		if err := testStor.Save(&storage.AliasURLModel{ID: uint64(i), LongURL: nodeURL.LongURL, ShortKey: nodeURL.ShortKey}); err != nil {
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
				contentType: "text/plain",
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
			h := NewHandlers(testStor, config.NewConfig()).mainHandlerGet
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

	listOfURL := []storage.AliasURLModel{
		{ID: 0, LongURL: "https://ya.ru", ShortKey: "123456789"},
		{ID: 1, LongURL: "https://google.com", ShortKey: "987654321"},
		{ID: 2, LongURL: "https://go.dev", ShortKey: ""},
	}

	testStor := memstor.NewMemStorage()
	for i, nodeURL := range listOfURL {
		if err := testStor.Save(&storage.AliasURLModel{ID: uint64(i), LongURL: nodeURL.LongURL, ShortKey: nodeURL.ShortKey}); err != nil {
			require.NotNil(t, err)
		}
	}

	testConfig := config.NewConfig()

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
				contentType: "text/plain",
				requestURI:  listOfURL[0].LongURL,
			},
			want: struct {
				code        int
				contentType string
				response    string
			}{
				code:        http.StatusCreated,
				contentType: "text/plain",
				response:    testConfig.BaseURL() + "/" + listOfURL[0].ShortKey,
			},
		},
	}

	//	Start test cases
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.request.requestURI))
			request.Header.Add("Content-type", tt.request.contentType)

			recorder := httptest.NewRecorder()
			h := NewHandlers(testStor, testConfig).mainHandlerPost
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
