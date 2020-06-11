package endpoints

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/potproject/ikuradon-salmon/dataAccess"
)

type UnSubscribeRequest struct {
	// 64 characters
	SubscribeId string // subscribe_id
}

func PostUnsubscribe(w http.ResponseWriter, r *http.Request) {
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	if len(splitToken) != 2 || splitToken[1] == "" {
		ErrorResponse(w, r, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	req := UnSubscribeRequest{
		SubscribeId: splitToken[1],
	}

	// 存在チェック
	check, err := dataAccess.DA.Has(req.SubscribeId)
	if err != nil {
		ErrorResponse(w, r, http.StatusInternalServerError, err)
		return
	}
	if !check {
		ErrorResponse(w, r, http.StatusNotFound, errors.New("NotFound"))
		return
	}

	// 消します
	err = dataAccess.DA.Delete(req.SubscribeId)
	if err != nil {
		ErrorResponse(w, r, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	res, _ := json.Marshal(SubscribeResponse{
		Result: true,
		Data: SubscribeResponseData{
			SubscribeId: req.SubscribeId,
		},
	})
	w.Write(res)
}
