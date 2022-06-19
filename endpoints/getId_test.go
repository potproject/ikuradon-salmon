package endpoints

import (
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/potproject/ikuradon-salmon/dataaccess"
)

func TestGetID404(t *testing.T) {
	dataaccess.SetMemory()
	defer dataaccess.DA.Close()
	r := httptest.NewRequest("GET", "http://0.0.0.0:8080/", nil)
	r.Header.Set("Authorization", "Bearer SUBSCRIBEID_NOT_FOUND")
	w := httptest.NewRecorder()
	GetID(w, r)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 404 {
		t.Errorf("Status Code invalid: %d", resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Content-Type invalid: %s", resp.Header.Get("Content-Type"))
	}
	if string(body) != `{"result":false,"errors":{"code":"404","message":"NotFound"}}` {
		t.Errorf("Status Code invalid: %s", string(body))
	}
}

func TestGetID401(t *testing.T) {
	dataaccess.SetMemory()
	defer dataaccess.DA.Close()
	r := httptest.NewRequest("GET", "http://0.0.0.0:8080/", nil)
	w := httptest.NewRecorder()
	GetID(w, r)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 401 {
		t.Errorf("Status Code invalid: %d", resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Content-Type invalid: %s", resp.Header.Get("Content-Type"))
	}
	if string(body) != `{"result":false,"errors":{"code":"401","message":"Unauthorized"}}` {
		t.Errorf("Status Code invalid: %s", string(body))
	}
}

func TestGetID200(t *testing.T) {
	dataaccess.SetMemory()
	defer dataaccess.DA.Close()
	dataaccess.DA.Set("12F34D56C78B90A", dataaccess.DataSet{
		SubscribeID:       "12F34D56C78B90A",
		ExponentPushToken: "Expo[xxxxxx]",
		PushPrivateKey:    "PushPrivateKey",
		PushAuth:          "PushAuth",
		ServerKey:         "ServerKey",
		CreatedAt:         1600000000,
		LastUpdatedAt:     1600000000,
	})
	r := httptest.NewRequest("GET", "http://0.0.0.0:8080/", nil)
	r.Header.Set("Authorization", "Bearer 12F34D56C78B90A")
	w := httptest.NewRecorder()
	GetID(w, r)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		t.Errorf("Status Code invalid: %d", resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Content-Type invalid: %s", resp.Header.Get("Content-Type"))
	}
	if string(body) != `{"result":true,"data":{"subscribe_id":"12F34D56C78B90A","user_id":"1234567890","username":"UserName","domain":"server.mastodon.net","access_token":"AccessToken","exponent_push_token":"Expo[xxxxxx]","push_private_key":"PushPrivateKey","push_public_key":"PushPublicKey","push_auth":"PushAuth","server_key":"ServerKey","created_at":1600000000,"last_updated_at":1600000000}}` {
		t.Errorf("Status Code invalid: %s", string(body))
	}
}
