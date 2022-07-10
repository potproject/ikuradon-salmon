package network

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
)

const misskeyTestDomain = "test.misskey.net"
const misskeyTestAccessToken = "ACCESS_TOKEN"

func TestVerifyMisskeySuccess(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	url := fmt.Sprintf("https://%s%s", misskeyTestDomain, verifyMisskeyEndpoints)
	httpmock.RegisterResponder("POST", url,
		func(req *http.Request) (*http.Response, error) {
			resv := ResVerifyMisskey{
				ID:       "100200300",
				Username: "testuser",
			}
			return httpmock.NewJsonResponse(200, resv)
		},
	)
	m := Misskey{}
	id, uName, err := m.Verify(misskeyTestDomain, misskeyTestAccessToken)
	if id != "100200300" || uName != "testuser" {
		t.Error("invalid ID/User")
	}
	if err != nil {
		t.Error(err)
	}
}

func TestVerifyMisskeyServerError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	url := fmt.Sprintf("https://%s%s", misskeyTestDomain, verifyMisskeyEndpoints)
	httpmock.RegisterResponder("POST", url,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(500, ""), nil
		},
	)
	m := Misskey{}
	_, _, err := m.Verify(misskeyTestDomain, misskeyTestAccessToken)
	if err == nil {
		t.Error("invaild status")
	}
}

func TestVerifyMisskeyClientError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	url := fmt.Sprintf("https://%s%s", misskeyTestDomain, verifyMisskeyEndpoints)
	httpmock.RegisterResponder("POST", url,
		func(req *http.Request) (*http.Response, error) {
			return nil, fmt.Errorf("Client Err")
		},
	)
	m := Misskey{}
	_, _, err := m.Verify(misskeyTestDomain, misskeyTestAccessToken)
	if err == nil {
		t.Error("invaild status")
	}
}

func TestVerifyMisskeyJSONParseError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	url := fmt.Sprintf("https://%s%s", misskeyTestDomain, verifyMisskeyEndpoints)
	httpmock.RegisterResponder("POST", url,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, "a")
		},
	)
	m := Misskey{}
	_, _, err := m.Verify(misskeyTestDomain, misskeyTestAccessToken)
	if err == nil {
		t.Error("invaild json Parse")
	}
}
