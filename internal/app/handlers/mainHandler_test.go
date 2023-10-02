package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)


func Test_mainHandlerMethodGet(t *testing.T) {



	testCases := []struct{
		name string
		request struct{
			contentType string
			request string
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
			request: struct{contentType string; request string}{
				contentType: "text/plain",
				request: "/123456789",
			},
			want: struct{code int; contentType string; response string}{
				code: http.StatusCreated,
				contentType: "text/plain",
				response: "https://ya.ru",
			},
		},
		//----------------------------------
		//	2.  test

	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tt.request.request, nil)
			request.Header.Add("Content-type", tt.request.contentType)

			recorder := httptest.NewRecorder()
			//mainHandlerMethodGet(recorder, request)

			result := recorder.Result()
			//	check status code
			assert.Equal(t, tt.want.code, result.StatusCode)
			//	check body
			//body, err := io.ReadAll(result.Body)
			//require.NotNil(err)
			//assert.Equal(t, tt.want.response, body)
			//	check content-type
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-type"))
		})
	}
}

// func Test_mainHandlerMethodPost(t *testing.T) {
// 	type args struct {
// 		w http.ResponseWriter
// 		r *http.Request
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			mainHandlerMethodPost(tt.args.w, tt.args.r)
// 		})
// 	}
// }
