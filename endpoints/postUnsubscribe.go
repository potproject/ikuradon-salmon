package endpoints

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/potproject/ikuradon-salmon/dataAccess"
	"github.com/potproject/ikuradon-salmon/network"
	"github.com/potproject/ikuradon-salmon/setting"
)

// PostUnsubscribe post unsubscribe
func PostUnsubscribe(w http.ResponseWriter, r *http.Request) {
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

	// if exists?
	check, err := dataAccess.DA.Has(subscribeID)
	if err != nil {
		ErrorResponse(w, r, http.StatusInternalServerError, err)
		return
	}
	if !check {
		ErrorResponse(w, r, http.StatusNotFound, errors.New("NotFound"))
		return
	}

	// Deleting Mastodon Server
	_ = network.PushUnsubscribeMastodon(
		req.Domain,
		req.AccessToken,
	)

	// Deleting DA
	err = dataAccess.DA.Delete(subscribeID)
	if err != nil {
		ErrorResponse(w, r, http.StatusInternalServerError, err)
		return
	}
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
