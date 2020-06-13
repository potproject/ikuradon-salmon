package network

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const pushSubscribeMastodonEndpoints = "/api/v1/push/subscription"

type ResPushSubscribe struct {
	ID       int64  `json:"id"`
	Endpoint string `json:"endpoint"`
	Alerts   struct {
		Follow    bool `json:"follow"`
		Favourite bool `json:"favourite"`
		Reblog    bool `json:"reblog"`
		Mention   bool `json:"mention"`
		Poll      bool `json:"poll"`
	} `json:"alerts"`
	ServerKey string `json:"server_key"`
}

func PushSubscribeMastodon(domain string, accessToken string, subscriptionEndpoint string, subscriptionKeysP256dh string, subscriptionKeysAuth string) (ResPushSubscribe, error) {
	endpoints := fmt.Sprintf("https://%s%s", domain, pushSubscribeMastodonEndpoints)
	values := url.Values{}
	values.Set("subscription[endpoint]", subscriptionEndpoint)
	values.Add("subscription[keys][p256dh]", subscriptionKeysP256dh)
	values.Add("subscription[keys][auth]", subscriptionKeysAuth)
	values.Add("data[alerts][follow]", "true")
	values.Add("data[alerts][favourite]", "true")
	values.Add("data[alerts][reblog]", "true")
	values.Add("data[alerts][mention]", "true")
	values.Add("data[alerts][poll]", "true")
	req, _ := http.NewRequest("POST", endpoints, strings.NewReader(values.Encode()))

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	client := http.Client{
		Timeout: mastodonTimeout,
	}
	resp, err := client.Do(req)
	var rp ResPushSubscribe
	if err != nil {
		return rp, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return rp, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return rp, fmt.Errorf("Status:%d %s", resp.StatusCode, string(b))
	}
	err = json.Unmarshal(b, &rp)
	if err != nil {
		return rp, err
	}
	return rp, nil
}
