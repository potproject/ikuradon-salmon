package endpoints

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/potproject/ikuradon-salmon/dataAccess"
	"github.com/potproject/ikuradon-salmon/network"
	"github.com/potproject/ikuradon-salmon/setting"
	"github.com/sethvargo/go-password/password"
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

	subscribeID, err := password.Generate(24, 10, 0, false, true)
	if err != nil {
		ErrorResponse(w, r, http.StatusInternalServerError, err)
		return
	}

	// Mastodonに信頼チェック
	id, username, err := network.VerifyMastodon(req.Domain, req.AccessToken)
	if err != nil {
		ErrorResponse(w, r, http.StatusServiceUnavailable, errors.New("Mastodon Server Unavailable: "+err.Error()))
		return
	}
	// AuthSecret = Base64 encoded string of 16 bytes of random data.
	authSecret := make([]byte, 16)
	rand.Read(authSecret)
	auth := base64.RawURLEncoding.EncodeToString([]byte(authSecret))

	// Web Push APIに登録
	// VAPID Keyを生成
	privateKey, publicKey, err := webpush.GenerateVAPIDKeys()

	endpoints := setting.S.BaseURL + "api/v1/webpush/" + subscribeID
	if setting.S.BaseURL == "" {
		endpoints = fmt.Sprintf("https://%s:%d/api/v1/webpush/%s", setting.S.ApiHost, setting.S.ApiPort, subscribeID)
	}

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

	// データ登録
	now := time.Now().Unix()
	err = dataAccess.DA.Set(subscribeID, dataAccess.DataSet{
		SubscribeId:        subscribeID,
		UserID:             id,
		Username:           username,
		Domain:             req.Domain,
		AccessToken:        req.AccessToken,
		ExponentPushToken:  req.ExponentPushToken,
		PushPrivateKey:     privateKey,
		PushPublicKey:      publicKey,
		PushAuth:           auth,
		ServerKey:          rps.ServerKey,
		CreatedAt:          now,
		ExpiredAt:          0,
		LastUpdatedAt:      now,
		ServerLastId:       0,
		NotificationsCount: 0,
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
