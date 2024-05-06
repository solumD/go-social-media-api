package authorization

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestLogin(t *testing.T) {
	testTable := []struct {
		name               string
		reqBody            []byte
		expectedStatusCode int
		expextedBody       string
	}{
		{
			name:               "Login OK",
			reqBody:            []byte(`{"login":"test","password":"password"}`),
			expectedStatusCode: http.StatusOK,
			expextedBody:       `{"login":"test"}`,
		},
		{
			name:               "Login in the same account",
			reqBody:            []byte(`{"login":"test","password":"password"}`),
			expectedStatusCode: http.StatusBadRequest,
			expextedBody:       `{"error": "user test already logged in!"}`,
		},
		{
			name:               "Invalid password",
			reqBody:            []byte(`{"login":"user","password":"123"}`),
			expectedStatusCode: http.StatusBadRequest,
			expextedBody:       `{"error":"invalid password"}`,
		},
		{
			name:               "Invalid json login input",
			reqBody:            []byte(`{"login":"test,"password":"password"}`),
			expectedStatusCode: http.StatusBadRequest,
			expextedBody:       `{"error":"invalid json User Input"}`,
		},
		{
			name:               "Invalid json password input",
			reqBody:            []byte(`{"login":"test","password":"password}`),
			expectedStatusCode: http.StatusBadRequest,
			expextedBody:       `{"error":"invalid json User Input"}`,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			t.Logf("Calling request: %s", testCase.reqBody)
			req, _ := http.NewRequest(http.MethodPost, "http://localhost:8080/login", bytes.NewBuffer(testCase.reqBody))
			resp, _ := http.DefaultClient.Do(req)
			body, _ := io.ReadAll(resp.Body)
			if strings.TrimSpace(string(body)) != testCase.expextedBody {
				t.Errorf("Incorrect result. Expexted %s, got %s", testCase.expextedBody, strings.TrimSpace(string(body)))
			}
			if resp.StatusCode != testCase.expectedStatusCode {
				t.Errorf("Incorrect result. Expexted %d, got %d", testCase.expectedStatusCode, resp.StatusCode)
			}
		})
	}
}

