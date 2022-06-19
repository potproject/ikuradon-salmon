package endpoints

import (
	"encoding/json"
	"net/http"

	"github.com/potproject/ikuradon-salmon/dataaccess"
)

type unsubscribeRequest struct {
	SubscribeID string // subscribe_id
}

// PostUnsubscribe post unsubscribe
func PostUnsubscribe(w http.ResponseWriter, r *http.Request) {
	s := unsubscribeRequest{
		SubscribeID: r.FormValue("subscribe_id"),
	}

	// Deleting DA
	err := dataaccess.DA.Delete(s.SubscribeID)
	if err != nil {
		ErrorResponse(w, r, http.StatusInternalServerError, err)
		return
	}

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
