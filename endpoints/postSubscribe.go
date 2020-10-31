package endpoints

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/potproject/ikuradon-salmon/dataaccess"
	"github.com/potproject/ikuradon-salmon/network"
	"github.com/potproject/ikuradon-salmon/setting"
)

type subscribeRequest struct {
	Domain            string // domain
	AccessToken       string // access_token
	ExponentPushToken string // exponent_push_token
}

type subscribeResponse struct {
	Result bool `json:"result"`

	Data subscribeResponseData `json:"data"`
}

type subscribeResponseData struct {
	// 64 characters
	SubscribeID string `json:"subscribe_id"`
}

// PostSubscribe post subscribe
func PostSubscribe(w http.ResponseWriter, r *http.Request) {
	req := subscribeRequest{
		Domain:            r.FormValue("domain"),
		AccessToken:       r.FormValue("access_token"),
		ExponentPushToken: r.FormValue("exponent_push_token"),
	}
	if req.Domain == "" || req.AccessToken == "" || req.ExponentPushToken == "" {
		ErrorResponse(w, r, http.StatusBadRequest, errors.New("Bad Request"))
		return
	}

	// Unique HMAC-SHA256
	uniq := req.Domain + ":" + req.AccessToken + ":" + req.ExponentPushToken
	subscribeID := makeHMAC(uniq, setting.S.Salt)
	exist, _ := dataaccess.DA.Has(subscribeID)
	if exist {
		updateSubscribe(w, r, subscribeID, req)
	} else {
		newSubscribe(w, r, subscribeID, req)
	}
}

func updateSubscribe(w http.ResponseWriter, r *http.Request, subscribeID string, req subscribeRequest) {
	ds, err := dataaccess.DA.Get(subscribeID)
	if err != nil {
		ErrorResponse(w, r, http.StatusInternalServerError, err)
		return
	}
	endpoints := makeEndpoints(subscribeID)
	rps, err := network.PushSubscribeMastodon(
		req.Domain,
		req.AccessToken,
		endpoints,
		ds.PushPublicKey,
		ds.PushAuth,
	)

	if err != nil {
		ErrorResponse(w, r, http.StatusServiceUnavailable, err)
		return
	}

	// Data Set
	now := time.Now().Unix()
	err = dataaccess.DA.Set(subscribeID, dataaccess.DataSet{
		SubscribeID:       subscribeID,
		UserID:            ds.UserID,
		Username:          ds.Username,
		Domain:            req.Domain,
		AccessToken:       req.AccessToken,
		ExponentPushToken: req.ExponentPushToken,
		PushPrivateKey:    ds.PushPrivateKey,
		PushPublicKey:     ds.PushPublicKey,
		PushAuth:          ds.PushAuth,
		ServerKey:         rps.ServerKey,
		CreatedAt:         ds.CreatedAt,
		LastUpdatedAt:     now,
	})
	if err != nil {
		ErrorResponse(w, r, http.StatusInternalServerError, err)
		return
	}
	// OK!
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	res, _ := json.Marshal(subscribeResponse{
		Result: true,
		Data: subscribeResponseData{
			SubscribeID: subscribeID,
		},
	})
	w.Write(res)
}

func newSubscribe(w http.ResponseWriter, r *http.Request, subscribeID string, req subscribeRequest) {
	// Mastodon Vertify
	id, username, err := network.VerifyMastodon(req.Domain, req.AccessToken)
	if err != nil {
		ErrorResponse(w, r, http.StatusServiceUnavailable, errors.New("Mastodon Server Unavailable: "+err.Error()))
		return
	}

	// AuthSecret = Base64 encoded string of 16 bytes of random data.
	auth := generateAuthSecret()

	// Web Push API Subscribing
	// generate VAPID Key
	privateKey, publicKey, err := webpush.GenerateVAPIDKeys()

	endpoints := makeEndpoints(subscribeID)

	rps, err := network.PushSubscribeMastodon(
		req.Domain,
		req.AccessToken,
		endpoints,
		publicKey,
		auth,
	)
	if err != nil {
		ErrorResponse(w, r, http.StatusServiceUnavailable, errors.New("Mastodon Server Unavailable: "+err.Error()))
		return
	}

	// Data Set
	now := time.Now().Unix()
	err = dataaccess.DA.Set(subscribeID, dataaccess.DataSet{
		SubscribeID:       subscribeID,
		UserID:            id,
		Username:          username,
		Domain:            req.Domain,
		AccessToken:       req.AccessToken,
		ExponentPushToken: req.ExponentPushToken,
		PushPrivateKey:    privateKey,
		PushPublicKey:     publicKey,
		PushAuth:          auth,
		ServerKey:         rps.ServerKey,
		CreatedAt:         now,
		LastUpdatedAt:     now,
	})
	if err != nil {
		ErrorResponse(w, r, http.StatusInternalServerError, err)
		return
	}
	// OK!
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	res, _ := json.Marshal(subscribeResponse{
		Result: true,
		Data: subscribeResponseData{
			SubscribeID: subscribeID,
		},
	})
	w.Write(res)
}

// makeHMAC generating HMAC (Digest Algorithm: SHA-256)
func makeHMAC(msg, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(msg))
	return hex.EncodeToString(mac.Sum(nil))
}

// makeEndpoints generaing WebPush Endpoint
func makeEndpoints(subscribeID string) string {
	endpoints := setting.S.BaseURL + "api/v1/webpush/" + subscribeID
	if setting.S.BaseURL == "" {
		endpoints = fmt.Sprintf("https://%s:%d/api/v1/webpush/%s", setting.S.APIHost, setting.S.APIPort, subscribeID)
	}
	return endpoints
}

// generateAuthSecret generating Base64 encoded string of 16 bytes of random data.
func generateAuthSecret() string {
	authSecret := make([]byte, 16)
	rand.Read(authSecret)
	auth := base64.RawURLEncoding.EncodeToString([]byte(authSecret))
	return auth
}
