package endpoints

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// ErrorResponseErrors Error Response
type ErrorResponseErrors struct {
	// Error Code
	Code string `json:"code"`
	// Error Message
	Message string `json:"message"`
}

// ErrorResponseBody Error Body
type ErrorResponseBody struct {
	// Always `false`.
	Result bool `json:"result"`

	Errors ErrorResponseErrors `json:"errors"`
}

// ErrorResponse http Error Response
func ErrorResponse(w http.ResponseWriter, r *http.Request, status int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	res, _ := json.Marshal(ErrorResponseBody{
		Result: false,
		Errors: ErrorResponseErrors{
			Code:    strconv.Itoa(status),
			Message: err.Error(),
		},
	})
	w.Write(res)
}
