package network

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
)

const mastodonTestDomain = "test.mastodon.net"
const mastodonTestAccessToken = "ACCESS_TOKEN"

func TestVerifyMastodonSuccess(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	url := fmt.Sprintf("https://%s%s", mastodonTestDomain, verifyMastodonEndpoints)
	httpmock.RegisterResponder("GET", url,
		func(req *http.Request) (*http.Response, error) {
			resv := ResVerify{
				ID:       "100200300",
				Username: "testuser",
			}
			return httpmock.NewJsonResponse(200, resv)
		},
	)
	m := Mastodon{}
	id, uName, err := m.VerifyMastodon(mastodonTestDomain, mastodonTestAccessToken)
	if id != "100200300" || uName != "testuser" {
		t.Error("invalid ID/User")
	}
	if err != nil {
		t.Error(err)
	}
}

func TestVerifyMastodonServerError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	url := fmt.Sprintf("https://%s%s", mastodonTestDomain, verifyMastodonEndpoints)
	httpmock.RegisterResponder("GET", url,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(500, ""), nil
		},
	)
	m := Mastodon{}
	_, _, err := m.VerifyMastodon(mastodonTestDomain, mastodonTestAccessToken)
	if err == nil {
		t.Error("invaild status")
	}
}

func TestVerifyMastodonClientError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	url := fmt.Sprintf("https://%s%s", mastodonTestDomain, verifyMastodonEndpoints)
	httpmock.RegisterResponder("GET", url,
		func(req *http.Request) (*http.Response, error) {
			return nil, fmt.Errorf("Client Err")
		},
	)
	m := Mastodon{}
	_, _, err := m.VerifyMastodon(mastodonTestDomain, mastodonTestAccessToken)
	if err == nil {
		t.Error("invaild status")
	}
}

func TestVerifyMastodonJSONParseError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	url := fmt.Sprintf("https://%s%s", mastodonTestDomain, verifyMastodonEndpoints)
	httpmock.RegisterResponder("GET", url,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, "a")
		},
	)
	m := Mastodon{}
	_, _, err := m.VerifyMastodon(mastodonTestDomain, mastodonTestAccessToken)
	if err == nil || err.Error() != "json: cannot unmarshal string into Go value of type network.ResVerify" {
		t.Error("invaild json Parse")
	}
}
