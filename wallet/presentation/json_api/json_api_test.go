package jsonapi_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ageeknamedslickback/wallet-API/wallet/dto"
	"github.com/ageeknamedslickback/wallet-API/wallet/presentation"
	"github.com/shopspring/decimal"
)

func accessToken(t *testing.T) string {
	router := presentation.Router()
	url := "/access_token"
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		t.Fatal(err)
	}
	router.ServeHTTP(w, req)

	// Response ..
	type Response struct {
		Resp dto.AccessToken `json:"response"`
	}

	var resp Response

	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}

	return resp.Resp.AccessToken

}
func TestWalletJsonAPI_WalletBalance(t *testing.T) {
	router := presentation.Router()
	type args struct {
		url    string
		method string
	}
	tests := []struct {
		name           string
		args           args
		wantStatusCode int
	}{
		{
			name: "happy case",
			args: args{
				url:    "/api/v1/1/balance",
				method: http.MethodGet,
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "sad case - not found",
			args: args{
				url:    "/1/balance",
				method: http.MethodGet,
			},
			wantStatusCode: http.StatusNotFound,
		},
		{
			name: "sad case - bad request",
			args: args{
				url:    "/api/v1/0/balance",
				method: http.MethodGet,
			},
			wantStatusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			req, _ := http.NewRequest(tt.args.method, tt.args.url, nil)
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken(t)))

			router.ServeHTTP(w, req)

			if tt.wantStatusCode != w.Code {
				t.Fatalf(
					"expected status code %v, but got %v",
					tt.wantStatusCode,
					w.Code,
				)
			}

			if tt.wantStatusCode == http.StatusOK {
				if !strings.Contains(w.Body.String(), "wallet") {
					t.Fatalf("expected balance to be found in response")
				}
			}

			if tt.wantStatusCode == http.StatusBadRequest {
				if !strings.Contains(w.Body.String(), "error") {
					t.Fatalf("expected error to be found in response")
				}
			}
		})
	}
}

func TestWalletJsonAPI_CreditWallet(t *testing.T) {
	router := presentation.Router()

	crAmount := dto.AmountInput{
		Amount: decimal.NewFromFloat(2.98),
	}
	crAmountBs, err := json.Marshal(crAmount)
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		url    string
		method string
		body   io.Reader
	}
	tests := []struct {
		name           string
		args           args
		wantStatusCode int
	}{
		{
			name: "happy case",
			args: args{
				url:    "/api/v1/3/credit",
				method: http.MethodPost,
				body:   bytes.NewBuffer(crAmountBs),
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "sad case - not found",
			args: args{
				url:    "/1/credit",
				method: http.MethodPost,
				body:   nil,
			},
			wantStatusCode: http.StatusNotFound,
		},
		{
			name: "sad case - bad request",
			args: args{
				url:    "/api/v1/0/credit",
				method: http.MethodPost,
				body:   nil,
			},
			wantStatusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			req, _ := http.NewRequest(tt.args.method, tt.args.url, tt.args.body)
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken(t)))

			router.ServeHTTP(w, req)

			if tt.wantStatusCode != w.Code {
				t.Fatalf(
					"expected status code %v, but got %v",
					tt.wantStatusCode,
					w.Code,
				)
			}

			if tt.wantStatusCode == http.StatusOK {
				if !strings.Contains(w.Body.String(), "wallet") {
					t.Fatalf("expected wallet to be found in response")
				}
			}

			if tt.wantStatusCode == http.StatusBadRequest {
				if !strings.Contains(w.Body.String(), "error") {
					t.Fatalf("expected error to be found in response")
				}
			}
		})
	}
}

func TestWalletJsonAPI_DebitWallet(t *testing.T) {
	router := presentation.Router()

	drAmount := dto.AmountInput{
		Amount: decimal.NewFromFloat(2.98),
	}
	drAmountBs, err := json.Marshal(drAmount)
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		url    string
		method string
		body   io.Reader
	}
	tests := []struct {
		name           string
		args           args
		wantStatusCode int
	}{
		{
			name: "happy case",
			args: args{
				url:    "/api/v1/3/debit",
				method: http.MethodPost,
				body:   bytes.NewBuffer(drAmountBs),
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "sad case - not found",
			args: args{
				url:    "/1/debit",
				method: http.MethodPost,
				body:   nil,
			},
			wantStatusCode: http.StatusNotFound,
		},
		{
			name: "sad case - bad request",
			args: args{
				url:    "/api/v1/0/debit",
				method: http.MethodPost,
				body:   nil,
			},
			wantStatusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			req, _ := http.NewRequest(tt.args.method, tt.args.url, tt.args.body)
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken(t)))

			router.ServeHTTP(w, req)

			if tt.wantStatusCode != w.Code {
				t.Fatalf(
					"expected status code %v, but got %v",
					tt.wantStatusCode,
					w.Code,
				)
			}

			if tt.wantStatusCode == http.StatusOK {
				if !strings.Contains(w.Body.String(), "wallet") {
					t.Fatalf("expected wallet to be found in response")
				}
			}

			if tt.wantStatusCode == http.StatusBadRequest {
				if !strings.Contains(w.Body.String(), "error") {
					t.Fatalf("expected error to be found in response")
				}
			}
		})
	}
}

func TestWalletJsonAPI_Authenticate(t *testing.T) {
	router := presentation.Router()
	type args struct {
		url    string
		method string
	}
	tests := []struct {
		name           string
		args           args
		wantStatusCode int
	}{
		{
			name: "happy case",
			args: args{
				url:    "/access_token",
				method: http.MethodPost,
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "sad case - not found",
			args: args{
				url:    "/not-found",
				method: http.MethodGet,
			},
			wantStatusCode: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			req, _ := http.NewRequest(tt.args.method, tt.args.url, nil)
			router.ServeHTTP(w, req)

			if tt.wantStatusCode != w.Code {
				t.Fatalf(
					"expected status code %v, but got %v",
					tt.wantStatusCode,
					w.Code,
				)
			}

			if tt.wantStatusCode == http.StatusOK {
				if !strings.Contains(w.Body.String(), "response") {
					t.Fatalf("expected response to be found in response")
				}
			}
		})
	}
}
