package infoblox

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	logging "mlmtool/pkg/util/logger"
)

func TestClientGet(t *testing.T) {
	logger := logging.NewTestingLogger(t.Name())
	tests := []struct {
		name           string
		reqURL         string
		expectedStatus int
		expectedError  bool
	}{
		{"ValidRequest", "valid", http.StatusOK, false},
		{"NotFound", "notfound", http.StatusNotFound, true},
		{"ServerError", "error", http.StatusInternalServerError, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/wapi/v1/valid" {
					w.WriteHeader(http.StatusOK)
					_, _ = w.Write([]byte(`{"status":"success"}`))
				} else if r.URL.Path == "/wapi/v1/notfound" {
					w.WriteHeader(http.StatusNotFound)
				} else {
					w.WriteHeader(http.StatusInternalServerError)
				}
			}))
			defer server.Close()

			client := &Client{
				baseURL: server.URL,
				version: "v1",
				client:  server.Client(),
				log:     logger,
			}

			var response map[string]interface{}
			err := client.get(tt.reqURL, nil, &response)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, "success", response["status"])
			}
		})
	}
}
