package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEtherscanHandle(t *testing.T) {
	tests := []struct {
		name       string
		query      string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "Invalid module parameter",
			query:      "module=invalid&action=getLogs",
			wantStatus: http.StatusOK,
			wantBody:   `"code":0`,
		},
		{
			name:       "Missing required parameters",
			query:      "module=logs",
			wantStatus: http.StatusOK,
			wantBody:   `"code":0`,
		},
		{
			name:       "Valid logs-getLogs request",
			query:      "module=logs&action=getLogs&offset=10&page=1",
			wantStatus: http.StatusOK,
			wantBody:   `"status":1`,
		},
		{
			name:       "Valid transaction-getStatus request",
			query:      "module=transaction&action=getstatus&txhash=0xa9972c6f84de1e56d924d7cdcfbfc7ba06eba92f1e3823bf6cd7147c4e277621",
			wantStatus: http.StatusOK,
			wantBody:   `"status":1`,
		},
		{
			name:       "Valid transaction-getxReceiptStatus request",
			query:      "module=transaction&action=gettxreceiptstatus&txhash=0xa9972c6f84de1e56d924d7cdcfbfc7ba06eba92f1e3823bf6cd7147c4e277621",
			wantStatus: http.StatusOK,
			wantBody:   `"status":1`,
		},
		{
			name:       "Valid account-balance request",
			query:      "module=account&action=balance&address=0x66b8c60c79dfad02fc04f1f13aab0f6feff8615b",
			wantStatus: http.StatusOK,
			wantBody:   `"status":1`,
		},
		{
			name:       "Valid account-balanceMulti request",
			query:      "module=account&action=balancemulti&address=0x66b8c60c79dfad02fc04f1f13aab0f6feff8615b,0xe22d73f5dcccb31a994ad4e7ad265cf69b4e725a",
			wantStatus: http.StatusOK,
			wantBody:   `"status":1`,
		},
		{
			name:       "Valid account-txlist request",
			query:      "module=account&action=txlist&address=0x66b8c60c79dfad02fc04f1f13aab0f6feff8615b&startblock=0&endblock=99999999&sort=asc&offset=100",
			wantStatus: http.StatusOK,
			wantBody:   `"status":1`,
		},
		{
			name:       "Valid account-tokentx request", // erc20
			query:      "module=account&action=tokentx&address=0x66b8c60c79dfad02fc04f1f13aab0f6feff8615b&offset=100&page=1",
			wantStatus: http.StatusOK,
			wantBody:   `"status":1`,
		},
		{
			name:       "Valid account-tokennfttx request", // erc721
			query:      "module=account&action=tokennfttx&address=0x66b8c60c79dfad02fc04f1f13aab0f6feff8615b&offset=100&page=1",
			wantStatus: http.StatusOK,
			wantBody:   `"status":1`,
		},
		{
			name:       "Valid account-token1155tx request", // erc1155
			query:      "module=account&action=token1155tx&address=0x66b8c60c79dfad02fc04f1f13aab0f6feff8615b&offset=100&page=1",
			wantStatus: http.StatusOK,
			wantBody:   `"status":1`,
		},
		{
			name:       "Valid contract-getABI request",
			query:      "module=contract&action=getabi&address=0x66b8c60c79dfad02fc04f1f13aab0f6feff8615b",
			wantStatus: http.StatusOK,
			wantBody:   `"status":1`,
		},
		{
			name:       "Valid contract-getContractCreation request",
			query:      "module=contract&action=getcontractcreation&contractaddresses=0x66b8c60c79dfad02fc04f1f13aab0f6feff8615b",
			wantStatus: http.StatusOK,
			wantBody:   `"status":1`,
		},
		{
			name:       "Valid contract-checkVerifyStatus request",
			query:      "module=contract&action=checkverifystatus&guid=0x66b8c60c79dfad02fc04f1f13aab0f6feff8615b",
			wantStatus: http.StatusOK,
			wantBody:   `"status":1`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/etherscan?"+tt.query, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_ = etherscanHandle(w, r)
			})

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.wantStatus {
				t.Errorf("Handler returned wrong status code: got %v want %v", status, tt.wantStatus)
			}

			if !strings.Contains(rr.Body.String(), tt.wantBody) {
				t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), tt.wantBody)
			}
		})
	}
}
