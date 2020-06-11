package network

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

const expoURL = "https://exp.host"
const pushExpoEndpoints = "/--/api/v2/push/send"

const expoTimeout = 10 * time.Second

func PushExpo(exponentPushToken string) error {
	url := expoURL + pushExpoEndpoints
	body := `{"ids": ["XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX", "YYYYYYYY-YYYY-YYYY-YYYY-YYYYYYYYYYYY"]}`
	req, _ := http.NewRequest("POST", url, bytes.NewReader([]byte(body)))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-encoding", "gzip, deflate")
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{
		Timeout: expoTimeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Status:%d", resp.StatusCode)
	}
	return nil
}
