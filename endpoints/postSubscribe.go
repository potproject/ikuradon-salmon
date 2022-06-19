package endpoints

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/potproject/ikuradon-salmon/dataaccess"
)

type subscribeRequest struct {
	SubscribeID       string // subscribe_id
	PushPrivateKey    string // push_private_key
	PushAuth          string // push_auth
	ServerKey         string // server_key
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

// PostSubscribeV2 post subscribeV2
func PostSubscribe(w http.ResponseWriter, r *http.Request) {
	s := subscribeRequest{
		SubscribeID:       r.FormValue("subscribe_id"),
		PushPrivateKey:    r.FormValue("push_private_key"),
		PushAuth:          r.FormValue("push_auth"),
		ServerKey:         r.FormValue("server_key"),
		ExponentPushToken: r.FormValue("exponent_push_token"),
	}

	dataaccess.DA.Set(s.SubscribeID, dataaccess.DataSet{
		SubscribeID:       s.SubscribeID,
		ExponentPushToken: s.ExponentPushToken,
		PushPrivateKey:    s.PushPrivateKey,
		PushAuth:          s.PushAuth,
		ServerKey:         s.ServerKey,
		LastUpdatedAt:     time.Now().Unix(),
	})

	// OK!
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	res, _ := json.Marshal(subscribeResponse{
		Result: true,
		Data: subscribeResponseData{
			SubscribeID: s.SubscribeID,
		},
	})
	w.Write(res)
}
