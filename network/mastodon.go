package network

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const verifyMastodonEndpoints = "/api/v1/accounts/verify_credentials"

const notificationsMastodonEndpoints = "/api/v1/notifications"

const mastodonTimeout = 30 * time.Second

// MastodonInterface Mastodon Backend Server JSON-API Access Interface
type MastodonInterface interface {
	VerifyMastodon(domain string, accessToken string) (string, string, error)
}

// Mastodon Mastodon Backend Server JSON-API Access Interface
type Mastodon struct {
	MastodonInterface
}

// ResVerify Mastodon id and username Response
type ResVerify struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

// VerifyMastodon GET:/api/v1/verify_credentials Mastodon Server
func (m Mastodon) VerifyMastodon(domain string, accessToken string) (string, string, error) {
	url := fmt.Sprintf("https://%s%s", domain, verifyMastodonEndpoints)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	client := http.Client{
		Timeout: mastodonTimeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	if resp.StatusCode != 200 {
		return "", "", fmt.Errorf("Status:%d", resp.StatusCode)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	var rv ResVerify
	err = json.Unmarshal(b, &rv)
	if err != nil {
		return "", "", err
	}
	return rv.ID, rv.Username, nil
}
