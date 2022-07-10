package network

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const pushSubscribeMisskeyEndpoints = "/api/sw/register"
const pushUnsubscribeMisskeyEndpoints = "/api/sw/unregister"

// MisskeyPush Misskey Push Backend Server JSON-API Access Interface
type MisskeyPush struct {
	SNSPushInterface
}

// ReqPushSubscribeMisskey Misskey POST:/api/sw/register Request
type ReqPushSubscribeMisskey struct {
	I         string `json:"i"`
	Endpoint  string `json:"endpoint"`
	Publickey string `json:"publickey"`
	Auth      string `json:"auth"`
}

// ReqPushUnsubscribeMisskey Misskey POST:/api/sw/unregister Request
type ReqPushUnsubscribeMisskey struct {
	I        string `json:"i"`
	Endpoint string `json:"endpoint"`
}

// ResPushSubscribeMisskey Misskey POST:/api/sw/register Response
type ResPushSubscribeMisskey struct {
	State     string `json:"state"`
	ServerKey string `json:"key"`
}

// MisskeyPush Subscribe to push notifications
// See: https://misskey.io/api-doc#operation/sw/register
func (mp MisskeyPush) PushSubscribe(domain string, accessToken string, subscriptionEndpoint string, subscriptionKeysP256dh string, subscriptionKeysAuth string) (string, error) {
	endpoints := fmt.Sprintf("https://%s%s", domain, pushSubscribeMisskeyEndpoints)
	reqJson := ReqPushSubscribeMisskey{
		I:         accessToken,
		Endpoint:  subscriptionEndpoint,
		Publickey: subscriptionKeysP256dh,
		Auth:      subscriptionKeysAuth,
	}
	reqjsonString, err := json.Marshal(reqJson)
	if err != nil {
		return "", err
	}
	req, _ := http.NewRequest("POST", endpoints, bytes.NewBuffer(reqjsonString))

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{
		Timeout: misskeyTimeout,
	}
	resp, err := client.Do(req)
	var rp ResPushSubscribeMisskey
	if err != nil {
		return rp.ServerKey, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return rp.ServerKey, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return rp.ServerKey, fmt.Errorf("Status:%d %s", resp.StatusCode, string(b))
	}
	err = json.Unmarshal(b, &rp)
	if err != nil {
		return rp.ServerKey, err
	}
	return rp.ServerKey, nil
}

// MisskeyPush Subscribe to push notifications
// See: https://misskey.io/api-doc#operation/sw/unregister
func (mp MisskeyPush) PushUnsubscribe(domain string, accessToken string, subscriptionEndpoint string) error {
	endpoints := fmt.Sprintf("https://%s%s", domain, pushUnsubscribeMisskeyEndpoints)
	reqJson := ReqPushUnsubscribeMisskey{
		I:        accessToken,
		Endpoint: subscriptionEndpoint,
	}
	reqjsonString, err := json.Marshal(reqJson)
	if err != nil {
		return err
	}
	req, _ := http.NewRequest("POST", endpoints, bytes.NewBuffer(reqjsonString))

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{
		Timeout: misskeyTimeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 204 {
		return fmt.Errorf("Status:%d Error", resp.StatusCode)
	}
	return nil
}
