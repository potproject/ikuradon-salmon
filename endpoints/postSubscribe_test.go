package endpoints

import (
	"encoding/base64"
	"testing"
)

func TestMakeHMAC(t *testing.T) {
	msg := "BODY Message."
	key := "TEST_KEY"
	hmac := makeHMAC(msg, key)
	if hmac != "421e95203b54232285d7266aebe90c5a1ee45885e5e93cf446032e641c87655f" {
		t.Errorf("invalid hash: %s", hmac)
	}
}

func TestMakeEndpoints(t *testing.T) {
	id := "1234567890"
	endpoint := makeEndpoints(id)
	if endpoint != "https://0.0.0.0:8080/api/v1/webpush/1234567890" {
		t.Errorf("invalid endpoint: %s", endpoint)
	}
}

func TestGenerateAuthSecret(t *testing.T) {
	secret := generateAuthSecret()
	auth, err := base64.RawURLEncoding.DecodeString(secret)
	if err != nil {
		t.Errorf(err.Error())
	}
	if len(auth) != 16 {
		t.Errorf("invalid secret: %s", secret)
	}
}
