package testutil

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type HTTPTestCase struct {
	URL            string
	Method         string
	ExpectedStatus int
	Body           []byte
}

func HandlerTestCase(t *testing.T, h http.Handler, tc HTTPTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		resprec := httptest.NewRecorder()
		var body io.Reader
		if tc.Body != nil {
			body = bytes.NewReader(tc.Body)
		}

		req, err := http.NewRequest(tc.Method, tc.URL, body)
		if err != nil {
			t.Fatal(err)
		}

		h.ServeHTTP(resprec, req)

		resp := resprec.Result()
		if resp.StatusCode != tc.ExpectedStatus {
			t.Errorf("StatusCode mismatch: expected %d; got %d", tc.ExpectedStatus, resp.StatusCode)
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}
			t.Errorf("Response Body: %s", string(b))
			resp.Body.Close()
		}
	}
}
