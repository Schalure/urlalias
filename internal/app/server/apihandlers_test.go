package server

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Schalure/urlalias/internal/app/aliaslogger/zaplogger"
	"github.com/Schalure/urlalias/internal/app/aliasmaker"
	"github.com/Schalure/urlalias/internal/app/mocks"
	"github.com/Schalure/urlalias/internal/app/models/aliasentity"
)

func Test_apiGetShortURL(t *testing.T) {

	testLocalHost := "http://localhost"
	testMethod := "POST"
	testURL := "/api/shorten"
	userID := uint64(1)

	mockController := gomock.NewController(t)
	defer mockController.Finish()

	userManager := mocks.NewMockUserManager(mockController)
	shortner := mocks.NewMockShortner(mockController)
	logger, err := zaplogger.NewZapLogger("")
	require.NoError(t, err)

	testServer := httptest.NewServer(NewRouter(New(userManager, shortner, logger, testLocalHost)))
	defer testServer.Close()

	testCases := []struct {
		name           string
		requestBody    string
		getShortKeyOut struct {
			requestURL string
			shortKey   string
			err        error
		}
		want struct {
			statusCode   int
			responseBody string
		}
	}{
		{
			name:        "simple test",
			requestBody: `{"url": "https://ya.ru"}`,
			getShortKeyOut: struct {
				requestURL string
				shortKey   string
				err        error
			}{
				requestURL: "https://ya.ru",
				shortKey:   "000000001",
				err:        nil,
			},
			want: struct {
				statusCode   int
				responseBody string
			}{
				statusCode:   http.StatusCreated,
				responseBody: `{"result":"` + testLocalHost + `/000000001"}`,
			},
		},
		{
			name:        "conflict test",
			requestBody: `{"url": "https://ya.ru"}`,
			getShortKeyOut: struct {
				requestURL string
				shortKey   string
				err        error
			}{
				requestURL: "https://ya.ru",
				shortKey:   "000000001",
				err:        aliasmaker.ErrConflictURL,
			},
			want: struct {
				statusCode   int
				responseBody string
			}{
				statusCode:   http.StatusConflict,
				responseBody: `{"result":"` + testLocalHost + `/000000001"}`,
			},
		},
		{
			name:        "not found test",
			requestBody: `{"url": "https://ya.ru"}`,
			getShortKeyOut: struct {
				requestURL string
				shortKey   string
				err        error
			}{
				requestURL: "https://ya.ru",
				shortKey:   "",
				err:        aliasmaker.ErrInternal,
			},
			want: struct {
				statusCode   int
				responseBody string
			}{
				statusCode:   http.StatusBadRequest,
				responseBody: aliasmaker.ErrInternal.Error(),
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {

			userManager.EXPECT().CreateUser().Return(userID, nil)
			shortner.EXPECT().GetShortKey(gomock.Any(), userID, test.getShortKeyOut.requestURL).Return(test.getShortKeyOut.shortKey, test.getShortKeyOut.err)

			request, err := http.NewRequest(testMethod, testServer.URL+testURL, strings.NewReader(test.requestBody))
			require.NoError(t, err)
			request.Header.Add("Content-type", "application/json")

			client := testServer.Client()
			transport := &http.Transport{
				Proxy:              http.ProxyFromEnvironment,
				DisableCompression: true,
			}
			client.Transport = transport

			response, err := client.Do(request)
			require.NoError(t, err)

			//	check status code
			assert.Equal(t, test.want.statusCode, response.StatusCode)

			data, err := io.ReadAll(response.Body)
			require.NoError(t, err)
			err = response.Body.Close()
			require.NoError(t, err)

			if response.StatusCode != http.StatusBadRequest {
				assert.Equal(t, test.want.responseBody, string(data))
			}
		})
	}
}

func Benchmark_apiGetShortURL(b *testing.B) {

	testLocalHost := "http://localhost"
	testMethod := "POST"
	testURL := "/api/shorten"
	userID := uint64(1)

	mockController := gomock.NewController(b)
	defer mockController.Finish()

	storage := mocks.NewMockStorager(mockController)
	storage.EXPECT().GetLastShortKey().Return("000000001").AnyTimes()
	storage.EXPECT().CreateUser().Return(userID, nil).AnyTimes()
	storage.EXPECT().FindByLongURL(gomock.Any(), "https://ya.ru").Return(nil, errors.New("")).AnyTimes()
	storage.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	logger, err := zaplogger.NewZapLogger("")
	require.NoError(b, err)

	service, err := aliasmaker.New(storage, logger)
	require.NoError(b, err)

	for i := 0; i < b.N; i++ {

		request, err := http.NewRequest(testMethod, testURL, strings.NewReader(`{"url": "https://ya.ru"}`))
		require.NoError(b, err)
		request.Header.Add("Content-type", "application/json")

		recorder := httptest.NewRecorder()
		h := New(service, service, logger, testLocalHost).apiGetShortURL

		h(recorder, request.WithContext(context.WithValue(request.Context(), UserID, userID)))
	}
}

func Test_apiGetBatchShortURL(t *testing.T) {

	testLocalHost := "http://localhost"
	testMethod := "POST"
	testURL := "/api/shorten/batch"
	userID := uint64(1)

	mockController := gomock.NewController(t)
	defer mockController.Finish()

	userManager := mocks.NewMockUserManager(mockController)
	shortner := mocks.NewMockShortner(mockController)
	logger, err := zaplogger.NewZapLogger("")
	require.NoError(t, err)

	testServer := httptest.NewServer(NewRouter(New(userManager, shortner, logger, testLocalHost)))
	defer testServer.Close()

	testCases := []struct {
		name                string
		requestBody         string
		getBatchShortURLOut struct {
			batchRequestURL []string
			batchRsponseKey []string
			err             error
		}
		want struct {
			statusCode   int
			responseBody string
		}
	}{
		{
			name:        "simple test",
			requestBody: `[{"correlation_id": "1","original_url": "https://ya.ru"},{"correlation_id": "2","original_url": "https://google.com"}]`,
			getBatchShortURLOut: struct {
				batchRequestURL []string
				batchRsponseKey []string
				err             error
			}{
				batchRequestURL: []string{"https://ya.ru", "https://google.com"},
				batchRsponseKey: []string{"000000001", "000000002"},
				err:             nil,
			},
			want: struct {
				statusCode   int
				responseBody string
			}{
				statusCode:   http.StatusCreated,
				responseBody: `[{"correlation_id":"1","short_url":"http://localhost/000000001"},{"correlation_id":"2","short_url":"http://localhost/000000002"}]`,
			},
		},
		{
			name:        "bad request test",
			requestBody: `[{"correlation_id": "1","original_url": "https://ya.ru"},{"correlation_id": "2","original_url": "https://google.com"}]`,
			getBatchShortURLOut: struct {
				batchRequestURL []string
				batchRsponseKey []string
				err             error
			}{
				batchRequestURL: []string{"https://ya.ru", "https://google.com"},
				batchRsponseKey: []string{"000000001", "000000002"},
				err:             errors.New(""),
			},
			want: struct {
				statusCode   int
				responseBody string
			}{
				statusCode:   http.StatusBadRequest,
				responseBody: ``,
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {

			userManager.EXPECT().CreateUser().Return(userID, nil)
			shortner.EXPECT().GetBatchShortURL(gomock.Any(), userID, test.getBatchShortURLOut.batchRequestURL).Return(test.getBatchShortURLOut.batchRsponseKey, test.getBatchShortURLOut.err)

			request, err := http.NewRequest(testMethod, testServer.URL+testURL, strings.NewReader(test.requestBody))
			require.NoError(t, err)
			request.Header.Add("Content-type", "application/json")

			client := testServer.Client()
			transport := &http.Transport{
				Proxy:              http.ProxyFromEnvironment,
				DisableCompression: true,
			}
			client.Transport = transport

			response, err := client.Do(request)
			require.NoError(t, err)

			//	check status code
			assert.Equal(t, test.want.statusCode, response.StatusCode)

			data, err := io.ReadAll(response.Body)
			require.NoError(t, err)
			err = response.Body.Close()
			require.NoError(t, err)

			if response.StatusCode != http.StatusBadRequest {
				assert.Equal(t, test.want.responseBody, string(data))
			}
		})
	}
}

func Benchmark_apiGetBatchShortURL(b *testing.B) {

	testLocalHost := "http://localhost"
	testMethod := "POST"
	testURL := "/api/shorten/batch"
	userID := uint64(1)

	mockController := gomock.NewController(b)
	defer mockController.Finish()

	storage := mocks.NewMockStorager(mockController)
	storage.EXPECT().GetLastShortKey().Return("000000001").AnyTimes()
	storage.EXPECT().CreateUser().Return(userID, nil).AnyTimes()
	storage.EXPECT().FindAllByLongURLs(gomock.Any(), gomock.Any()).Return(map[string]*aliasentity.AliasURLModel{}, nil).AnyTimes()
	storage.EXPECT().SaveAll(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	logger, err := zaplogger.NewZapLogger("")
	require.NoError(b, err)

	service, err := aliasmaker.New(storage, logger)
	require.NoError(b, err)

	for i := 0; i < b.N; i++ {

		request, err := http.NewRequest(testMethod, testURL, strings.NewReader(`[{"correlation_id": "1","original_url": "https://ya.ru"},{"correlation_id": "2","original_url": "https://google.com"}]`))
		require.NoError(b, err)
		request.Header.Add("Content-type", "application/json")

		recorder := httptest.NewRecorder()
		h := New(service, service, logger, testLocalHost).apiGetBatchShortURL

		h(recorder, request.WithContext(context.WithValue(request.Context(), UserID, userID)))
	}
}

func Test_apiGetUserAliases(t *testing.T) {

	testLocalHost := "http://localhost"
	testMethod := "GET"
	testURL := "/api/user/urls"
	userID := uint64(1)

	mockController := gomock.NewController(t)
	defer mockController.Finish()

	userManager := mocks.NewMockUserManager(mockController)
	shortner := mocks.NewMockShortner(mockController)
	logger, err := zaplogger.NewZapLogger("")
	require.NoError(t, err)

	testServer := httptest.NewServer(NewRouter(New(userManager, shortner, logger, testLocalHost)))
	defer testServer.Close()

	testCases := []struct {
		name              string
		getUserAliasesOut struct {
			nodesOut []aliasentity.AliasURLModel
			err      error
		}
		want struct {
			statusCode   int
			responseBody string
		}
	}{
		{
			name: "simple test",
			getUserAliasesOut: struct {
				nodesOut []aliasentity.AliasURLModel
				err      error
			}{
				nodesOut: []aliasentity.AliasURLModel{
					{
						ShortKey: "000000001",
						LongURL:  "https://ya.ru",
					},
					{
						ShortKey: "000000002",
						LongURL:  "https://goo.com",
					},
				},
				err: nil,
			},
			want: struct {
				statusCode   int
				responseBody string
			}{
				statusCode:   http.StatusOK,
				responseBody: `[{"short_url":"http://localhost/000000001","original_url":"https://ya.ru"},{"short_url":"http://localhost/000000002","original_url":"https://goo.com"}]`,
			},
		},
		{
			name: "not found test",
			getUserAliasesOut: struct {
				nodesOut []aliasentity.AliasURLModel
				err      error
			}{
				nodesOut: []aliasentity.AliasURLModel{},
				err:      nil,
			},
			want: struct {
				statusCode   int
				responseBody string
			}{
				statusCode:   http.StatusNoContent,
				responseBody: ``,
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {

			userManager.EXPECT().GetUserAliases(gomock.Any(), userID).Return(test.getUserAliasesOut.nodesOut, test.getUserAliasesOut.err)

			request, err := http.NewRequest(testMethod, testServer.URL+testURL, nil)
			require.NoError(t, err)
			request.Header.Add("Content-type", "application/json")
			tokenString, err := createTokenJWT(userID)
			require.NoError(t, err)
			request.AddCookie(&http.Cookie{
				Name:  authorization,
				Value: tokenString,
			})

			client := testServer.Client()
			transport := &http.Transport{
				Proxy:              http.ProxyFromEnvironment,
				DisableCompression: true,
			}
			client.Transport = transport

			response, err := client.Do(request)
			require.NoError(t, err)

			//	check status code
			assert.Equal(t, test.want.statusCode, response.StatusCode)

			data, err := io.ReadAll(response.Body)
			require.NoError(t, err)
			err = response.Body.Close()
			require.NoError(t, err)

			if response.StatusCode != http.StatusBadRequest {
				assert.Equal(t, test.want.responseBody, string(data))
			}
		})
	}
}

func Benchmark_apiGetUserAliases(b *testing.B) {

	testLocalHost := "http://localhost"
	testMethod := "GET"
	testURL := "/api/user/urls"
	userID := uint64(1)

	mockController := gomock.NewController(b)
	defer mockController.Finish()

	storage := mocks.NewMockStorager(mockController)
	storage.EXPECT().GetLastShortKey().Return("000000001").AnyTimes()
	storage.EXPECT().FindByUserID(gomock.Any(), userID).Return([]aliasentity.AliasURLModel{{ShortKey: "000000001", LongURL: "https://ya.ru"}, {ShortKey: "000000002", LongURL: "https://goo.com"}}, nil).AnyTimes()

	logger, err := zaplogger.NewZapLogger("")
	require.NoError(b, err)

	service, err := aliasmaker.New(storage, logger)
	require.NoError(b, err)

	for i := 0; i < b.N; i++ {

		request, err := http.NewRequest(testMethod, testURL, nil)
		require.NoError(b, err)
		request.Header.Add("Content-type", "application/json")
		tokenString, err := createTokenJWT(userID)
		require.NoError(b, err)
		request.AddCookie(&http.Cookie{
			Name:  authorization,
			Value: tokenString,
		})

		recorder := httptest.NewRecorder()
		h := New(service, service, logger, testLocalHost).apiGetUserAliases

		h(recorder, request.WithContext(context.WithValue(request.Context(), UserID, userID)))
	}
}

func Test_aipDeleteUserAliases(t *testing.T) {

	testLocalHost := "http://localhost"
	testMethod := "DELETE"
	testURL := "/api/user/urls"
	userID := uint64(1)

	mockController := gomock.NewController(t)
	defer mockController.Finish()

	userManager := mocks.NewMockUserManager(mockController)
	shortner := mocks.NewMockShortner(mockController)
	logger, err := zaplogger.NewZapLogger("")
	require.NoError(t, err)

	testServer := httptest.NewServer(NewRouter(New(userManager, shortner, logger, testLocalHost)))
	defer testServer.Close()

	testCases := []struct {
		name                     string
		requestBody              string
		addAliasesToDeleteParams struct {
			aliases []string
			err     error
		}
		want struct {
			statusCode int
		}
	}{
		{
			name:        "simple test",
			requestBody: `["6qxTVvsy","RTfd56hn","Jlfd67ds"]`,
			addAliasesToDeleteParams: struct {
				aliases []string
				err     error
			}{
				aliases: []string{"6qxTVvsy", "RTfd56hn", "Jlfd67ds"},
				err:     nil,
			},
			want: struct{ statusCode int }{
				statusCode: http.StatusAccepted,
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {

			shortner.EXPECT().AddAliasesToDelete(gomock.Any(), userID, test.addAliasesToDeleteParams.aliases).Return(test.addAliasesToDeleteParams.err)

			request, err := http.NewRequest(testMethod, testServer.URL+testURL, strings.NewReader(test.requestBody))
			require.NoError(t, err)
			request.Header.Add("Content-type", "application/json")
			tokenString, err := createTokenJWT(userID)
			require.NoError(t, err)
			request.AddCookie(&http.Cookie{
				Name:  authorization,
				Value: tokenString,
			})

			client := testServer.Client()
			transport := &http.Transport{
				Proxy:              http.ProxyFromEnvironment,
				DisableCompression: true,
			}
			client.Transport = transport

			response, err := client.Do(request)
			require.NoError(t, err)
			response.Body.Close()

			//	check status code
			assert.Equal(t, test.want.statusCode, response.StatusCode)
		})
	}
}
