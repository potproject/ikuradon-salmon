package endpoints

import (
	"errors"
	"io/ioutil"
	"net/http/httptest"
	"testing"
)

func TestErrorResponse(t *testing.T) {
	r := httptest.NewRequest("GET", "http://0.0.0.0:8080/", nil)
	w := httptest.NewRecorder()
	ErrorResponse(w, r, 500, errors.New("Err"))
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 500 {
		t.Errorf("Status Code invalid: %d", resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Content-Type invalid: %s", resp.Header.Get("Content-Type"))
	}
	if string(body) != `{"result":false,"errors":{"code":"500","message":"Err"}}` {
		t.Errorf("Status Code invalid: %s", string(body))
	}
}
