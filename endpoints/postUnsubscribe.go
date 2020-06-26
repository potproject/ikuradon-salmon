package endpoints

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/potproject/ikuradon-salmon/dataAccess"
)

type unSubscribeRequest struct {
	// 64 characters
	SubscribeID string // subscribe_id
}

// PostUnsubscribe post unsubscribe
func PostUnsubscribe(w http.ResponseWriter, r *http.Request) {
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	if len(splitToken) != 2 || splitToken[1] == "" {
		ErrorResponse(w, r, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	req := unSubscribeRequest{
		SubscribeID: splitToken[1],
	}

	// 存在チェック
	check, err := dataAccess.DA.Has(req.SubscribeID)
	if err != nil {
		ErrorResponse(w, r, http.StatusInternalServerError, err)
		return
	}
	if !check {
		ErrorResponse(w, r, http.StatusNotFound, errors.New("NotFound"))
		return
	}

	// 消します
	err = dataAccess.DA.Delete(req.SubscribeID)
	if err != nil {
		ErrorResponse(w, r, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	res, _ := json.Marshal(subscribeResponse{
		Result: true,
		Data: subscribeResponseData{
			SubscribeID: req.SubscribeID,
		},
	})
	w.Write(res)
}
