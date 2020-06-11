package endpoints

import (
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

type SubscribeRequest struct {
	Domain            string // domain
	AccessToken       string // access_token
	ExponentPushToken string // exponent_push_token
}

type SubscribeResponse struct {
	Result bool `json:"result"`

	Data SubscribeResponseData `json:"data"`
}

type SubscribeResponseData struct {
	// 64 characters
	SubscribeId string `json:"subscribe_id"`
}

func PostSubscribe(w http.ResponseWriter, r *http.Request) {
	req := SubscribeRequest{
		Domain:            r.FormValue("domain"),
		AccessToken:       r.FormValue("access_token"),
		ExponentPushToken: r.FormValue("exponent_push_token"),
	}
	if req.Domain == "" || req.AccessToken == "" || req.ExponentPushToken == "" {
		ErrorResponse(w, r, http.StatusBadRequest, errors.New("Bad Request"))
		return
	}

	SubscribeId, err := password.Generate(24, 10, 0, false, true)
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
	auth := base64.StdEncoding.EncodeToString([]byte(SubscribeId))
	fmt.Println(auth)
	// Web Push APIに登録
	// VAPID Keyを生成
	privateKey, publicKey, err := webpush.GenerateVAPIDKeys()

	endpoints := setting.S.BaseURL + "api/v1/webpush"
	if setting.S.BaseURL == "" {
		endpoints = fmt.Sprintf("https://%s:%d/api/v1/webpush", setting.S.ApiHost, setting.S.ApiPort)
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
	fmt.Printf("%+v/r/n", rps)

	// データ登録
	now := time.Now().Unix()
	err = dataAccess.DA.Set(SubscribeId, dataAccess.DataSet{
		SubscribeId:        SubscribeId,
		UserID:             id,
		Username:           username,
		Domain:             req.Domain,
		AccessToken:        req.AccessToken,
		ExponentPushToken:  req.ExponentPushToken,
		PushPrivateKey:     privateKey,
		PushPublicKey:      publicKey,
		PushAuth:           auth,
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
	res, _ := json.Marshal(SubscribeResponse{
		Result: true,
		Data: SubscribeResponseData{
			SubscribeId: SubscribeId,
		},
	})
	w.Write(res)
}
