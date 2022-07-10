package network

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const verifyMisskeyEndpoints = "/api/i"

const misskeyTimeout = 30 * time.Second

// Misskey Misskey Backend Server JSON-API Access Interface
type Misskey struct {
	SnsInterface
}

// ReqVerifyMisskey Misskey POST:/api/i Request
type ReqVerifyMisskey struct {
	I string `json:"i"`
}

// ResVerifyMisskey Misskey POST:/api/sw/register Response
type ResVerifyMisskey struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

// VerifyMastodon POST:/api/is Misskey Server
func (m Misskey) Verify(domain string, accessToken string) (string, string, error) {
	url := fmt.Sprintf("https://%s%s", domain, verifyMisskeyEndpoints)
	reqJson := ReqVerifyMisskey{
		I: accessToken,
	}
	reqjsonString, err := json.Marshal(reqJson)
	if err != nil {
		return "", "", err
	}
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(reqjsonString))

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{
		Timeout: misskeyTimeout,
	}
	resp, err := client.Do(req)
	var rp ResVerifyMisskey
	if err != nil {
		return "", "", err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", "", fmt.Errorf("Status:%d Error", resp.StatusCode)
	}
	err = json.Unmarshal(b, &rp)
	return rp.ID, rp.Username, err
}
