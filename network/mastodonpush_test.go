package network

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
)

const mastodonPushTestDomain = "test.mastodon.com"
const mastodonPushTestAccessToken = "ACCESSTOKEN"
const mastodonPushTestEndpoints = "https://salmon.mastodon.net/api/v1/webpush/1F2D3E4D5C6B7A80"
const mastodonPushTestPrivateKey = "PRIVATE_KEY"
const mastodonPushTestAuth = "AuthRandamData"

func TestPushSubscribeMastodonSuccess(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	url := fmt.Sprintf("https://%s%s", mastodonPushTestDomain, pushSubscribeMastodonEndpoints)
	rps := ResPushSubscribe{
		ID:       100200300,
		Endpoint: mastodonPushTestEndpoints,
		Alerts: struct {
			Follow    bool `json:"follow"`
			Favourite bool `json:"favourite"`
			Reblog    bool `json:"reblog"`
			Mention   bool `json:"mention"`
			Poll      bool `json:"poll"`
		}{
			Follow:    true,
			Favourite: true,
			Reblog:    true,
			Mention:   true,
			Poll:      true,
		},
		ServerKey: "SERVER_KEY",
	}
	httpmock.RegisterResponder("POST", url,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, rps)
		},
	)
	mp := MastodonPush{}
	resp, err := mp.PushSubscribeMastodon(
		mastodonPushTestDomain,
		mastodonPushTestAccessToken,
		mastodonPushTestEndpoints,
		mastodonPushTestPrivateKey,
		mastodonPushTestAuth,
	)
	if resp != rps {
		t.Error("Invalid Response")
	}
	if err != nil {
		t.Error(err)
	}
}

func TestPushSubscribeMastodonServerError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	url := fmt.Sprintf("https://%s%s", mastodonPushTestDomain, pushSubscribeMastodonEndpoints)
	httpmock.RegisterResponder("POST", url,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(500, ""), nil
		},
	)
	mp := MastodonPush{}
	_, err := mp.PushSubscribeMastodon(
		mastodonPushTestDomain,
		mastodonPushTestAccessToken,
		mastodonPushTestEndpoints,
		mastodonPushTestPrivateKey,
		mastodonPushTestAuth,
	)
	if err == nil {
		t.Error("invaild status")
	}
}

func TestPushSubscribeMastodonClientError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	url := fmt.Sprintf("https://%s%s", mastodonPushTestDomain, pushSubscribeMastodonEndpoints)
	httpmock.RegisterResponder("POST", url,
		func(req *http.Request) (*http.Response, error) {
			return nil, fmt.Errorf("Client Err")
		},
	)
	mp := MastodonPush{}
	_, err := mp.PushSubscribeMastodon(
		mastodonPushTestDomain,
		mastodonPushTestAccessToken,
		mastodonPushTestEndpoints,
		mastodonPushTestPrivateKey,
		mastodonPushTestAuth,
	)
	if err == nil {
		t.Error("invaild status")
	}
}

func TestPushSubscribeMastodonJSONParseError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	url := fmt.Sprintf("https://%s%s", mastodonPushTestDomain, pushSubscribeMastodonEndpoints)
	httpmock.RegisterResponder("POST", url,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, "a")
		},
	)
	mp := MastodonPush{}
	_, err := mp.PushSubscribeMastodon(
		mastodonPushTestDomain,
		mastodonPushTestAccessToken,
		mastodonPushTestEndpoints,
		mastodonPushTestPrivateKey,
		mastodonPushTestAuth,
	)
	if err == nil || err.Error() != "json: cannot unmarshal string into Go value of type network.ResPushSubscribe" {
		t.Error("invaild json Parse")
	}
}

func TestPushUnsubscribeMastodonSuccess(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	url := fmt.Sprintf("https://%s%s", mastodonPushTestDomain, pushSubscribeMastodonEndpoints)
	httpmock.RegisterResponder("DELETE", url,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, map[string]interface{}{
				"status": "ok",
			})
		},
	)
	mp := MastodonPush{}
	err := mp.PushUnsubscribeMastodon(mastodonPushTestDomain, mastodonPushTestAccessToken)
	if err != nil {
		t.Error(err)
	}
}

func TestPushUnsubscribeMastodonServerError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	url := fmt.Sprintf("https://%s%s", mastodonPushTestDomain, pushSubscribeMastodonEndpoints)
	httpmock.RegisterResponder("DELETE", url,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(500, "")
		},
	)
	mp := MastodonPush{}
	err := mp.PushUnsubscribeMastodon(mastodonPushTestDomain, mastodonPushTestAccessToken)
	if err == nil {
		t.Error("invaild status")
	}
}

func TestPushUnsubscribeMastodonClientError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	url := fmt.Sprintf("https://%s%s", mastodonPushTestDomain, pushSubscribeMastodonEndpoints)
	httpmock.RegisterResponder("DELETE", url,
		func(req *http.Request) (*http.Response, error) {
			return nil, fmt.Errorf("Client Err")
		},
	)
	mp := MastodonPush{}
	err := mp.PushUnsubscribeMastodon(mastodonPushTestDomain, mastodonPushTestAccessToken)
	if err == nil {
		t.Error("invaild status")
	}
}
