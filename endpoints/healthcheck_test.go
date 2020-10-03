package endpoints

import (
	"io/ioutil"
	"net/http/httptest"
	"testing"
)

func TestHealthcheck(t *testing.T) {
	r := httptest.NewRequest("GET", "http://0.0.0.0:8080/", nil)
	w := httptest.NewRecorder()
	Healthcheck(w, r)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		t.Errorf("Status Code invalid: %d", resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Content-Type invalid: %s", resp.Header.Get("Content-Type"))
	}
	if string(body) != `{"result":true,"data":{"name":"ikuradon-salmon","version":"unversioned"}}` {
		t.Errorf("Status Code invalid: %s", string(body))
	}
}
