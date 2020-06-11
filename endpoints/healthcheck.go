package endpoints

import (
	"encoding/json"
	"net/http"

	"github.com/potproject/ikuradon-salmon/setting"
)

type HealthcheckResponse struct {
	// Always `true`.
	Result bool `json:"result"`
	// env.APP_VERSION
	Data HealthcheckResponseData `json:"data"`
}

type HealthcheckResponseData struct {
	// env.APP_NAME
	Name string `json:"name"`
	// env.APP_VERSION
	Version string `json:"version"`
}

func Healthcheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	res, _ := json.Marshal(HealthcheckResponse{
		Result: true,
		Data: HealthcheckResponseData{
			Name:    setting.S.AppName,
			Version: setting.S.AppVersion,
		},
	})
	w.Write(res)
}
