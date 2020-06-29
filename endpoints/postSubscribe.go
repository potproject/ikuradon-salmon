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
	"github.com/potproject/ikuradon-salmon/dataAccess"
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

	// SubscribeID ExponentPushToken HMAC-SHA256
	subscribeID := makeHMAC(req.ExponentPushToken, setting.S.Salt)
	exist, _ := dataAccess.DA.Has(subscribeID)
	if exist {
		updateSubscribe(w, r, subscribeID, req)
	} else {
		newSubscribe(w, r, subscribeID, req)
	}
}

func updateSubscribe(w http.ResponseWriter, r *http.Request, subscribeID string, req subscribeRequest) {
	ds, err := dataAccess.DA.Get(subscribeID)
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

	// Data Set
	now := time.Now().Unix()
	err = dataAccess.DA.Set(subscribeID, dataAccess.DataSet{
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
	authSecret := make([]byte, 16)
	rand.Read(authSecret)
	auth := base64.RawURLEncoding.EncodeToString([]byte(authSecret))

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
	err = dataAccess.DA.Set(subscribeID, dataAccess.DataSet{
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

func makeHMAC(msg, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(msg))
	return hex.EncodeToString(mac.Sum(nil))
}

func makeEndpoints(subscribeID string) string {
	endpoints := setting.S.BaseURL + "api/v1/webpush/" + subscribeID
	if setting.S.BaseURL == "" {
		endpoints = fmt.Sprintf("https://%s:%d/api/v1/webpush/%s", setting.S.ApiHost, setting.S.ApiPort, subscribeID)
	}
	return endpoints
}
