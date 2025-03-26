package link

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateSuccess(t *testing.T) {
	type want struct {
		code     int
		response createRes
	}
	tests := []struct {
		name  string
		value createReq
		want  want
	}{
		{
			name:  "valid value request",
			value: createReq{URL: "https://example.com"},
			want: want{
				code:     http.StatusCreated,
				response: createRes{URL: "https://example.com", Hash: "https://example.com"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &handler{}
			data, err := json.Marshal(&tt.value)
			if err != nil {
				t.Fatal("Error marshalling data")
			}
			reader := bytes.NewReader(data)
			w := httptest.NewRecorder()
			r, err := http.NewRequest(http.MethodPost, "/link", reader)
			if err != nil {
				t.Fatal("Error creating request")
			}

			handler.create()(w, r)

			if w.Result().StatusCode != tt.want.code {
				t.Fatalf("Error status code, expected %d, got %d", tt.want.code, w.Result().StatusCode)
			}

			body, err := io.ReadAll(w.Result().Body)
			if err != nil {
				t.Fatal("Error reading body")
			}
			defer w.Result().Body.Close()

			var respData createRes
			if err := json.Unmarshal(body, &respData); err != nil {
				t.Fatal("Failed to unmarshal response body")
			}
			if respData != tt.want.response {
				t.Fatalf("Error response body, expected %v, got %v", tt.want.response, respData)
			}
		})
	}
}

func TestCreateFail(t *testing.T) {
	type want struct {
		code     int
		response string
	}
	tests := []struct {
		name  string
		value any
		want  want
	}{
		{
			name:  "URL data request is not URL format",
			value: createReq{URL: "text"},
			want: want{
				code:     http.StatusBadRequest,
				response: "Error validating request",
			},
		},
		{
			name:  "URL data request is int type",
			value: struct{ URL int }{URL: 123},
			want: want{
				code:     http.StatusBadRequest,
				response: "Error decoding request",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &handler{}
			data, err := json.Marshal(&tt.value)
			if err != nil {
				t.Fatal("Error marshalling data")
			}
			reader := bytes.NewReader(data)
			w := httptest.NewRecorder()
			r, err := http.NewRequest(http.MethodPost, "/link", reader)
			if err != nil {
				t.Fatal("Error creating request")
			}

			handler.create()(w, r)

			if w.Result().StatusCode != tt.want.code {
				t.Fatalf("Error status code, expected %d, got %d", tt.want.code, w.Result().StatusCode)
			}

			body, err := io.ReadAll(w.Result().Body)
			if err != nil {
				t.Fatal("Error reading body")
			}
			defer w.Result().Body.Close()

			var respData string
			if err := json.Unmarshal(body, &respData); err != nil {
				t.Fatal("Failed to unmarshal response body")
			}
			if respData != tt.want.response {
				t.Fatalf("Error response body, expected %v, got %v", tt.want.response, respData)
			}
		})
	}
}

func TestGoToSuccess(t *testing.T) {
	type want struct {
		code int
	}
	tests := []struct {
		name  string
		value string
		want  want
	}{
		{
			name:  "valid value request",
			value: "https://example.com",
			want: want{
				code: http.StatusTemporaryRedirect,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &handler{}
			r, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/link/%s", tt.value), nil)
			if err != nil {
				t.Fatal("Error creating request")
			}
			w := httptest.NewRecorder()

			handler.goTo()(w, r)

			if w.Result().StatusCode != tt.want.code {
				t.Fatalf("Error status code, expected %d, got %d", tt.want.code, w.Result().StatusCode)
			}
		})
	}
}
