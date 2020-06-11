package endpoints

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/potproject/ikuradon-salmon/dataAccess"
)

type IdRequest struct {
	// 64 characters
	SubscribeId string // subscribe_id
}

type IdResponse struct {
	Result bool `json:"result"`

	Data dataAccess.DataSet `json:"data"`
}

func GetId(w http.ResponseWriter, r *http.Request) {
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	if len(splitToken) != 2 || splitToken[1] == "" {
		ErrorResponse(w, r, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	reqParam := IdRequest{
		SubscribeId: splitToken[1],
	}
	check, err := dataAccess.DA.Has(reqParam.SubscribeId)
	if err != nil {
		ErrorResponse(w, r, http.StatusInternalServerError, err)
		return
	}
	if !check {
		ErrorResponse(w, r, http.StatusNotFound, errors.New("NotFound"))
		return
	}
	ds, err := dataAccess.DA.Get(reqParam.SubscribeId)
	if err != nil {
		ErrorResponse(w, r, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	res, _ := json.Marshal(IdResponse{
		Result: true,
		Data:   ds,
	})
	w.Write(res)
}
