package network

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/potproject/ikuradon-salmon/notification"
)

func TestPushExpoSuccess(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	url := expoURL + pushExpoEndpoints
	httpmock.RegisterResponder("POST", url,
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"data": map[string]interface{}{
					"XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX": map[string]interface{}{
						"status": "ok",
					},
				},
			})
			return resp, err
		},
	)
	n := notification.N{}
	err := PushExpo("Expo[xxxxxxxxx]", n)
	if err != nil {
		t.Error(err)
	}
}

func TestPushExpoServerError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	url := expoURL + pushExpoEndpoints
	httpmock.RegisterResponder("POST", url,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(500, ""), nil
		},
	)
	n := notification.N{}
	err := PushExpo("Expo[xxxxxxxxx]", n)
	if err == nil {
		t.Error("invaild status")
	}
}

func TestPushExpoClientError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	url := expoURL + pushExpoEndpoints
	httpmock.RegisterResponder("POST", url,
		func(req *http.Request) (*http.Response, error) {
			return nil, fmt.Errorf("Client Err")
		},
	)
	n := notification.N{}
	err := PushExpo("Expo[xxxxxxxxx]", n)
	if err == nil {
		t.Error("invaild status")
	}
}
