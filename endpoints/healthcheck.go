package endpoints

import (
	"encoding/json"
	"net/http"

	"github.com/potproject/ikuradon-salmon/setting"
)

type healthcheckResponse struct {
	// Always `true`.
	Result bool `json:"result"`
	// env.APP_VERSION
	Data healthcheckResponseData `json:"data"`
}

type healthcheckResponseData struct {
	// env.APP_NAME
	Name string `json:"name"`
	// env.APP_VERSION
	Version string `json:"version"`
}

// Healthcheck Endpoint
func Healthcheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	res, _ := json.Marshal(healthcheckResponse{
		Result: true,
		Data: healthcheckResponseData{
			Name:    setting.S.AppName,
			Version: setting.S.AppVersion,
		},
	})
	w.Write(res)
}
