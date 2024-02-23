package server

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Schalure/urlalias/internal/app/aliaslogger/zaplogger"
	"github.com/Schalure/urlalias/internal/app/aliasmaker"
	"github.com/Schalure/urlalias/internal/app/mocks"
	"github.com/Schalure/urlalias/internal/app/models/aliasentity"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_redirect(t *testing.T) {

	mockController := gomock.NewController(t)
	defer mockController.Finish()

	userManager := mocks.NewMockUserManager(mockController)
	shortner := mocks.NewMockShortner(mockController)
	logger, err := zaplogger.NewZapLogger("")
	require.NoError(t, err)

	//	originalURL, err := h.shortner.GetOriginalURL(r.Context(), shortKey)
	testCases := []struct {
		name string
		requesURI string
		getOriginalURLParams struct {
			inpURI string
			outURL string
			outErr error
		}
		want struct {
			statusCode int
			responseURL string
		}
	}{
		{
			name: "simple test",
			requesURI: "/000000000",
			getOriginalURLParams: struct{inpURI string; outURL string; outErr error}{
				inpURI: "000000000",
				outURL: "https://ya.ru",
				outErr: nil,
			},
			want: struct{statusCode int; responseURL string}{
				statusCode: http.StatusTemporaryRedirect,
				responseURL: "https://ya.ru",
			},
		},
		{
			name: "deleted test",
			requesURI: "/000000000",
			getOriginalURLParams: struct{inpURI string; outURL string; outErr error}{
				inpURI: "000000000",
				outURL: "",
				outErr: aliasmaker.ErrURLWasDeleted,
			},
			want: struct{statusCode int; responseURL string}{
				statusCode: http.StatusGone,
				responseURL: "",
			},
		},
		{
			name: "not found test",
			requesURI: "/000000000",
			getOriginalURLParams: struct{inpURI string; outURL string; outErr error}{
				inpURI: "000000000",
				outURL: "",
				outErr: aliasmaker.ErrURLNotFound,
			},
			want: struct{statusCode int; responseURL string}{
				statusCode: http.StatusBadRequest,
				responseURL: "",
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {

			shortner.EXPECT().GetOriginalURL(gomock.Any(), test.getOriginalURLParams.inpURI).Return(test.getOriginalURLParams.outURL, test.getOriginalURLParams.outErr)

			request := httptest.NewRequest(http.MethodGet, test.requesURI, nil)

			recorder := httptest.NewRecorder()
			h := New(userManager, shortner, logger, "http://localhost/").redirect
			h(recorder, request)

			resp := recorder.Result()
			err := resp.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, test.want.statusCode, resp.StatusCode)
			assert.Equal(t, test.want.responseURL, resp.Header.Get("Location"))
		})
	}
}

func Benchmark_redirect(b *testing.B) {

	testMethod := "GET"
	testURL := "/000000002"
	testLocalHost := "http://localhost"
	userID := uint64(1)

	//b.StopTimer()
	mockController := gomock.NewController(b)
	defer mockController.Finish()

	storage := mocks.NewMockStorager(mockController)
	storage.EXPECT().GetLastShortKey().Return("000000001").AnyTimes()
	storage.EXPECT().CreateUser().Return(userID, nil).AnyTimes()
	storage.EXPECT().FindByShortKey(gomock.Any(), "000000002").Return(&aliasentity.AliasURLModel{
		ID: 1,
		UserID: userID,
		ShortKey: "000000002",
		LongURL: "https://ya.ru",
		DeletedFlag: false,
	}, nil).AnyTimes()

	logger, err := zaplogger.NewZapLogger("")
	require.NoError(b, err)

	service, err := aliasmaker.New(storage, logger)
	require.NoError(b, err)

	request := httptest.NewRequest(testMethod, testURL, nil)
	request.Header.Add("Content-type", "text/plain")

	recorder := httptest.NewRecorder()
	h := New(service, service, logger, testLocalHost).redirect

	for i := 0; i < b.N; i++ {

		h(recorder, request)
	}
}

func Test_getShortURL(t *testing.T) {

	testLocalHost := "http://localhost"
	testMethod := "POST"
	testURL := "/"
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
		name string
		requestURL string
		getShortKeyOut struct {
			shortKey string
			err error
		}
		want struct {
			statusCode int
			responseURL string
		}
	}{
		{
			name: "simple test",
			requestURL: "https://ya.ru",
			getShortKeyOut: struct{shortKey string; err error}{
				shortKey: "000000001",
				err: nil,
			},
			want: struct{statusCode int; responseURL string}{
				statusCode: http.StatusCreated,
				responseURL: testLocalHost + "/000000001",
			},
		},
		{
			name: "conflict test",
			requestURL: "https://ya.ru",
			getShortKeyOut: struct{shortKey string; err error}{
				shortKey: "000000001",
				err: aliasmaker.ErrConflictURL,
			},
			want: struct{statusCode int; responseURL string}{
				statusCode: http.StatusConflict,
				responseURL: testLocalHost + "/000000001",
			},
		},
		{
			name: "not found test",
			requestURL: "https://ya.ru",
			getShortKeyOut: struct{shortKey string; err error}{
				shortKey: "",
				err: aliasmaker.ErrInternal,
			},
			want: struct{statusCode int; responseURL string}{
				statusCode: http.StatusBadRequest,
				responseURL: aliasmaker.ErrInternal.Error(),
			},
		},
	}


	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {

			userManager.EXPECT().CreateUser().Return(userID, nil)
			shortner.EXPECT().GetShortKey(gomock.Any(), userID, test.requestURL).Return(test.getShortKeyOut.shortKey, test.getShortKeyOut.err)


			request, err := http.NewRequest(testMethod, testServer.URL + testURL, strings.NewReader(test.requestURL))
			require.NoError(t, err)
			request.Header.Add("Content-type", "text/plain")

			client := testServer.Client()
			transport := &http.Transport{
				Proxy: http.ProxyFromEnvironment,
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
				assert.Equal(t, test.want.responseURL, string(data))
			}
		})
	}
}

func Benchmark_getShortURL(b *testing.B) {

	testLocalHost := "http://localhost"
	testMethod := "POST"
	testURL := "/"
	userID := uint64(1)

	mockController := gomock.NewController(b)
	defer mockController.Finish()

	storage := mocks.NewMockStorager(mockController)
	storage.EXPECT().GetLastShortKey().Return("000000001").AnyTimes()
	storage.EXPECT().CreateUser().Return(userID, nil).AnyTimes()
	storage.EXPECT().FindByLongURL(gomock.Any(), "https://ya.ru").Return(&aliasentity.AliasURLModel{
		ID: 1,
		UserID: userID,
		ShortKey: "000000002",
		LongURL: "https://ya.ru",
		DeletedFlag: false,
	}, nil).AnyTimes()

	logger, err := zaplogger.NewZapLogger("")
	require.NoError(b, err)

	service, err := aliasmaker.New(storage, logger)
	require.NoError(b, err)

	// testServer := httptest.NewServer(NewRouter(New(service, service, logger, testLocalHost)))
	// defer testServer.Close()

	for i := 0; i < b.N; i++ {

		request := httptest.NewRequest(testMethod, testURL, strings.NewReader("https://ya.ru"))
		request.Header.Add("Content-type", "text/plain")

		recorder := httptest.NewRecorder()
		h := New(service, service, logger, testLocalHost).getShortURL

		h(recorder, request.WithContext(context.WithValue(request.Context(), UserID, userID)))
	}
}

// import (
// 	"io"
// 	"net/http"
// 	"net/http/httptest"
// 	"strings"
// 	"testing"

// 	"github.com/Schalure/urlalias/cmd/shortener/config"
// 	"github.com/Schalure/urlalias/internal/app/aliaslogger"
// 	"github.com/Schalure/urlalias/internal/app/aliasmaker"
// 	"github.com/Schalure/urlalias/internal/app/models/aliasentity"
// 	"github.com/Schalure/urlalias/internal/app/storage/memstor"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// )

// func newService(t *testing.T) *aliasmaker.AliasMakerServise {

// 	logger, err := aliaslogger.NewLogger(aliaslogger.LoggerTypeZap)
// 	require.NoError(t, err)
// 	stor, err := memstor.NewStorage()
// 	require.NoError(t, err)
// 	s, err := aliasmaker.New(config.NewConfig(), stor, logger)
// 	require.NoError(t, err)
// 	return s
// }

// // ------------------------------------------------------------
// //
// //	Test mainHandlerMethodGet: "/{shortKey}"
// func Test_mainHandlerMethodGet(t *testing.T) {

// 	service := newService(t)
// 	defer service.Stop()

// 	var listOfURL = []aliasentity.AliasURLModel{
// 		{ID: 0, LongURL: "https://ya.ru", ShortKey: "123456789"},
// 		{ID: 1, LongURL: "https://google.com", ShortKey: "987654321"},
// 	}

// 	for i, nodeURL := range listOfURL {
// 		if err := service.Storage.Save(&aliasentity.AliasURLModel{ID: uint64(i), LongURL: nodeURL.LongURL, ShortKey: nodeURL.ShortKey}); err != nil {
// 			require.NotNil(t, err)
// 		}
// 	}

// 	//	Test cases
// 	testCases := []struct {
// 		name    string
// 		request struct {
// 			contentType string
// 			requestURI  string
// 		}
// 		want struct {
// 			code     int
// 			response string
// 		}
// 	}{
// 		//----------------------------------
// 		//	1. simple test
// 		{
// 			name: "simple test",
// 			request: struct {
// 				contentType string
// 				requestURI  string
// 			}{
// 				contentType: textPlain,
// 				requestURI:  listOfURL[0].ShortKey,
// 			},
// 			want: struct {
// 				code     int
// 				response string
// 			}{
// 				code:     http.StatusTemporaryRedirect,
// 				response: listOfURL[0].LongURL,
// 			},
// 		},
// 	}

// 	//	Start test cases
// 	for _, testCase := range testCases {
// 		t.Run(testCase.name, func(t *testing.T) {

// 			request := httptest.NewRequest(http.MethodGet, "/"+testCase.request.requestURI, nil)
// 			request.Header.Add("Content-type", testCase.request.contentType)

// 			recorder := httptest.NewRecorder()
// 			h := New(service).mainHandlerGet
// 			h(recorder, request)

// 			result := recorder.Result()
// 			err := result.Body.Close()
// 			require.NoError(t, err)

// 			//	check status code
// 			assert.Equal(t, testCase.want.code, result.StatusCode)

// 			//	check response
// 			assert.Equal(t, testCase.want.response, result.Header.Get("Location"))
// 		})
// 	}
// }

// func Test_mainHandlerMethodPost(t *testing.T) {

// 	service := newService(t)
// 	defer service.Stop()

// 	listOfURL := []aliasentity.AliasURLModel{
// 		{ID: 0, LongURL: "https://ya.ru", ShortKey: "123456789"},
// 		{ID: 1, LongURL: "https://google.com", ShortKey: "987654321"},
// 		{ID: 2, LongURL: "https://go.dev", ShortKey: ""},
// 	}

// 	for i, nodeURL := range listOfURL {
// 		if err := service.Storage.Save(&aliasentity.AliasURLModel{ID: uint64(i), LongURL: nodeURL.LongURL, ShortKey: nodeURL.ShortKey}); err != nil {
// 			require.NotNil(t, err)
// 		}
// 	}

// 	testCases := []struct {
// 		name    string
// 		request struct {
// 			contentType string
// 			requestURI  string
// 		}
// 		want struct {
// 			code        int
// 			contentType string
// 			response    string
// 		}
// 	}{
// 		//----------------------------------
// 		//	1. simple test
// 		{
// 			name: "simple test",
// 			request: struct {
// 				contentType string
// 				requestURI  string
// 			}{
// 				contentType: textPlain,
// 				requestURI:  listOfURL[0].LongURL,
// 			},
// 			want: struct {
// 				code        int
// 				contentType string
// 				response    string
// 			}{
// 				code:        http.StatusConflict,
// 				contentType: textPlain,
// 				response:    service.Config.BaseURL() + "/" + listOfURL[0].ShortKey,
// 			},
// 		},
// 	}

// 	//	Start test cases
// 	for _, tt := range testCases {
// 		t.Run(tt.name, func(t *testing.T) {
// 			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.request.requestURI))
// 			request.Header.Add("Content-type", tt.request.contentType)

// 			recorder := httptest.NewRecorder()
// 			h := New(service).mainHandlerPost
// 			h(recorder, request)

// 			result := recorder.Result()

// 			//	check status code
// 			assert.Equal(t, tt.want.code, result.StatusCode)

// 			//	check response
// 			data, err := io.ReadAll(recorder.Body)
// 			require.NoError(t, err)
// 			err = result.Body.Close()
// 			require.NoError(t, err)

// 			assert.Equal(t, tt.want.response, string(data))

// 			assert.Contains(t, recorder.Header().Get("Content-type"), tt.want.contentType)
// 		})
// 	}
// }
