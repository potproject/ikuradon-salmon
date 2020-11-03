package network

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/potproject/ikuradon-salmon/notification"
)

// ExpoInterface Expo Backend Server JSON-API Access Interface
type ExpoInterface interface {
	PushExpo(exponentPushToken string, data notification.N) error
}

// Expo  Expo Backend Server JSON-API Access Interface
type Expo struct {
	ExpoInterface
}

const expoURL = "https://exp.host"
const pushExpoEndpoints = "/--/api/v2/push/send"

const expoTimeout = 10 * time.Second

type message struct {
	To        string         `json:"to"`
	Data      notification.N `json:"data"`
	Title     string         `json:"title"`
	Body      string         `json:"body"`
	TTL       string         `json:"ttl,omitempty"`
	Priority  string         `json:"priority,omitempty"`  // ('default' | 'normal' | 'high')
	Subtitle  string         `json:"subtitle,omitempty"`  // iOS Only
	Sound     string         `json:"sound,omitempty"`     // iOS Only ('default' | null)
	Badge     int            `json:"badge,omitempty"`     // iOS Only
	ChannelID string         `json:"channelId,omitempty"` // Android Only
}

// PushExpo Sending Expo Backend Server
func (m Expo) PushExpo(exponentPushToken string, data notification.N) error {
	url := expoURL + pushExpoEndpoints
	message := message{
		To:    exponentPushToken,
		Data:  data,
		Title: data.Title,
		Body:  data.Body,
		Sound: "default",
	}
	body, _ := json.Marshal(message)
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
