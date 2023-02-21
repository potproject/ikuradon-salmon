package endpoints

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	ece "github.com/crow-misia/http-ece"
	"github.com/gorilla/mux"
	"github.com/potproject/ikuradon-salmon/dataaccess"
	"github.com/potproject/ikuradon-salmon/network"
	"github.com/potproject/ikuradon-salmon/notification"
	"github.com/potproject/ikuradon-salmon/rfc8188"
)

type pushPayload struct {
	SubscribeID       string
	ContentEncoding   ece.ContentEncoding
	Salt              []byte
	DH                []byte
	EncryptedBody     []byte
	AuthSecret        []byte
	PublicKey         []byte
	PrivateKey        []byte
	PlainText         string
	ExponentPushToken string
}

func setPayload(r *http.Request) (p pushPayload, sns string, err error) {
	p = pushPayload{}
	err = nil
	// SubscribeID
	vars := mux.Vars(r)
	p.SubscribeID = vars["subscribeID"]

	// ContentEncoding
	headerContentEncoding := r.Header.Get("Content-Encoding")
	if strings.ToLower(headerContentEncoding) == "aes128gcm" {
		p.ContentEncoding = ece.AES128GCM
	} else {
		p.ContentEncoding = ece.AESGCM

		// Salt
		headerEncryption := r.Header.Get("Encryption")
		if strings.HasPrefix(headerEncryption, "salt=") {
			salt := headerEncryption[5:]
			p.Salt, _ = base64.RawURLEncoding.DecodeString(salt)
		}

		// DH
		headerCryptoKey := r.Header.Get("Crypto-Key")
		sliceHeaderCryptoKey := strings.Split(headerCryptoKey, ";")
		for _, ckey := range sliceHeaderCryptoKey {
			if strings.HasPrefix(ckey, "dh=") {
				p.DH, _ = base64.RawURLEncoding.DecodeString(ckey[3:])
			}
		}
	}

	// EncryptedBody
	p.EncryptedBody, _ = ioutil.ReadAll(r.Body)

	// AuthSecret
	// PrivateKey
	// ExponentPushToken
	check, errDB := dataaccess.DA.Has(p.SubscribeID)
	if errDB != nil {
		err = errDB
		return
	}
	if !check {
		err = errors.New("NotFound")
		return
	}
	ds, errDB := dataaccess.DA.Get(p.SubscribeID)
	if errDB != nil {
		err = errDB
		return
	}
	sns = ds.Sns
	p.AuthSecret, _ = base64.RawURLEncoding.DecodeString(ds.PushAuth)
	p.PrivateKey, _ = base64.RawURLEncoding.DecodeString(ds.PushPrivateKey)
	p.PublicKey, _ = base64.RawURLEncoding.DecodeString(ds.PushPublicKey)

	p.ExponentPushToken = ds.ExponentPushToken

	var plaintextByte []byte
	if p.ContentEncoding == ece.AESGCM {
		plaintextByte, err = ece.Decrypt(p.EncryptedBody,
			ece.WithEncoding(p.ContentEncoding),
			ece.WithSalt(p.Salt),
			ece.WithAuthSecret(p.AuthSecret),
			ece.WithPrivate(p.PrivateKey),
			ece.WithDh(p.DH),
		)
	} else {
		plaintextByte, err = rfc8188.Decrypt(p.EncryptedBody, p.PublicKey, p.PrivateKey, p.AuthSecret)
	}
	if err != nil {
		return
	}
	p.PlainText = string(plaintextByte)
	return
}

// PostWebPush Webpush Endpoint
func PostWebPush(w http.ResponseWriter, r *http.Request) {
	p, sns, err := setPayload(r)
	if err != nil {
		fmt.Println("error Decrypt.", err.Error())
		ErrorResponse(w, r, http.StatusInternalServerError, err)
		return
	}
	var n notification.N
	if sns == "misskey" {
		//misskey migration from notification
		w.WriteHeader(http.StatusOK)
		return
	} else {
		err = json.Unmarshal([]byte(p.PlainText), &n)
		if err != nil {
			fmt.Println("error Decrypt.", err.Error())
			ErrorResponse(w, r, http.StatusInternalServerError, err)
			return
		}
	}

	e := network.Expo{}
	e.PushExpo(p.ExponentPushToken, n)
	w.WriteHeader(http.StatusOK)
}
