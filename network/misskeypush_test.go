package network

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
)

const misskeyPushTestDomain = "test.misskey.com"
const misskeyPushTestAccessToken = "ACCESSTOKEN"
const misskeyPushTestEndpoints = "https://salmon.misskey.net/api/v1/webpush/1F2D3E4D5C6B7A80"
const misskeyPushTestPrivateKey = "PRIVATE_KEY"
const misskeyPushTestAuth = "AuthRandamData"

func TestPushSubscribeMisskeySuccess(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	url := fmt.Sprintf("https://%s%s", misskeyPushTestDomain, pushSubscribeMisskeyEndpoints)
	rps := ResPushSubscribeMisskey{
		State:     "subscribed",
		ServerKey: "SERVER_KEY",
	}
	httpmock.RegisterResponder("POST", url,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, rps)
		},
	)
	mp := MisskeyPush{}
	resp, err := mp.PushSubscribe(
		misskeyPushTestDomain,
		misskeyPushTestAccessToken,
		misskeyPushTestEndpoints,
		misskeyPushTestPrivateKey,
		misskeyPushTestAuth,
	)
	if resp != rps.ServerKey {
		t.Error("Invalid Response")
	}
	if err != nil {
		t.Error(err)
	}
}

func TestPushSubscribeMisskeyServerError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	url := fmt.Sprintf("https://%s%s", misskeyPushTestDomain, pushSubscribeMisskeyEndpoints)
	httpmock.RegisterResponder("POST", url,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(500, ""), nil
		},
	)
	mp := MisskeyPush{}
	_, err := mp.PushSubscribe(
		misskeyPushTestDomain,
		misskeyPushTestAccessToken,
		misskeyPushTestEndpoints,
		misskeyPushTestPrivateKey,
		misskeyPushTestAuth,
	)
	if err == nil {
		t.Error("invaild status")
	}
}

func TestPushSubscribemisskeyClientError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	url := fmt.Sprintf("https://%s%s", misskeyPushTestDomain, pushSubscribeMisskeyEndpoints)
	httpmock.RegisterResponder("POST", url,
		func(req *http.Request) (*http.Response, error) {
			return nil, fmt.Errorf("Client Err")
		},
	)
	mp := MisskeyPush{}
	_, err := mp.PushSubscribe(
		misskeyPushTestDomain,
		misskeyPushTestAccessToken,
		misskeyPushTestEndpoints,
		misskeyPushTestPrivateKey,
		misskeyPushTestAuth,
	)
	if err == nil {
		t.Error("invaild status")
	}
}

func TestPushSubscribemisskeyJSONParseError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	url := fmt.Sprintf("https://%s%s", misskeyPushTestDomain, pushSubscribeMisskeyEndpoints)
	httpmock.RegisterResponder("POST", url,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, "a")
		},
	)
	mp := MisskeyPush{}
	_, err := mp.PushSubscribe(
		misskeyPushTestDomain,
		misskeyPushTestAccessToken,
		misskeyPushTestEndpoints,
		misskeyPushTestPrivateKey,
		misskeyPushTestAuth,
	)
	if err == nil {
		t.Error("invaild json Parse")
	}
}

func TestPushUnsubscribemisskeySuccess(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	url := fmt.Sprintf("https://%s%s", misskeyPushTestDomain, pushUnsubscribeMisskeyEndpoints)
	httpmock.RegisterResponder("POST", url,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(204, map[string]interface{}{})
		},
	)
	mp := MisskeyPush{}
	err := mp.PushUnsubscribe(misskeyPushTestDomain, misskeyPushTestAccessToken, misskeyPushTestEndpoints)
	if err != nil {
		t.Error(err)
	}
}

func TestPushUnsubscribemisskeyServerError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	url := fmt.Sprintf("https://%s%s", misskeyPushTestDomain, pushUnsubscribeMisskeyEndpoints)
	httpmock.RegisterResponder("DELETE", url,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(500, "")
		},
	)
	mp := MisskeyPush{}
	err := mp.PushUnsubscribe(misskeyPushTestDomain, misskeyPushTestAccessToken, misskeyPushTestEndpoints)
	if err == nil {
		t.Error("invaild status")
	}
}

func TestPushUnsubscribemisskeyClientError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	url := fmt.Sprintf("https://%s%s", misskeyPushTestDomain, pushUnsubscribeMisskeyEndpoints)
	httpmock.RegisterResponder("DELETE", url,
		func(req *http.Request) (*http.Response, error) {
			return nil, fmt.Errorf("Client Err")
		},
	)
	mp := MisskeyPush{}
	err := mp.PushUnsubscribe(misskeyPushTestDomain, misskeyPushTestAccessToken, misskeyPushTestEndpoints)
	if err == nil {
		t.Error("invaild status")
	}
}
