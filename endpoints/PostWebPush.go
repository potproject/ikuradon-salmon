package endpoints

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	ece "github.com/crow-misia/http-ece"
	"github.com/gorilla/mux"
	"github.com/potproject/ikuradon-salmon/dataAccess"
)

type pushPayload struct {
	SubscribeID     string
	ContentEncoding ece.ContentEncoding
	TTL             int
	ContentType     string
	Salt            []byte
	DH              []byte
	P256ECDSA       []byte
	JWT             string
	EncryptedBody   []byte
	AuthSecret      []byte
	PrivateKey      []byte
	PlainText       string
}

func setPayload(r *http.Request) (p pushPayload, err error) {
	p = pushPayload{}
	err = nil
	// SubscribeID
	vars := mux.Vars(r)
	p.SubscribeID = vars["subscribeID"]

	// ContentEncoding
	headerContentEncoding := r.Header.Get("Content-Encoding")
	if headerContentEncoding == "aes128gcm" {
		p.ContentEncoding = ece.AES128GCM
	} else {
		p.ContentEncoding = ece.AESGCM
	}

	// TTL
	p.TTL, _ = strconv.Atoi(r.Header.Get("Ttl"))

	// ContentType
	p.ContentType = r.Header.Get("Content-Type")

	// Salt
	headerEncryption := r.Header.Get("Encryption")
	if strings.HasPrefix(headerEncryption, "salt=") {
		salt := headerEncryption[5:]
		p.Salt, _ = base64.RawURLEncoding.DecodeString(salt)
	}

	// DH
	// P256ECDSA
	headerCryptoKey := r.Header.Get("Crypto-Key")
	sliceHeaderCryptoKey := strings.Split(headerCryptoKey, ";")
	for _, ckey := range sliceHeaderCryptoKey {
		if strings.HasPrefix(ckey, "dh=") {
			p.DH, _ = base64.RawURLEncoding.DecodeString(ckey[3:])
		}
		if strings.HasPrefix(ckey, "p256ecdsa=") {
			p.P256ECDSA, _ = base64.RawURLEncoding.DecodeString(ckey[10:])
		}
	}
	// JWT
	headerAuthorization := r.Header.Get("Authorization")
	if strings.HasPrefix(headerAuthorization, "WebPush ") {
		p.JWT = headerAuthorization[8:]
	}

	// EncryptedBody
	p.EncryptedBody, _ = ioutil.ReadAll(r.Body)

	// AuthSecret
	// PrivateKey
	check, errDB := dataAccess.DA.Has(p.SubscribeID)
	if errDB != nil {
		err = errDB
		return
	}
	if !check {
		err = errors.New("NotFound")
		return
	}
	ds, errDB := dataAccess.DA.Get(p.SubscribeID)
	if errDB != nil {
		err = errDB
		return
	}
	p.AuthSecret, _ = base64.RawURLEncoding.DecodeString(ds.PushAuth)
	p.PrivateKey, _ = base64.RawURLEncoding.DecodeString(ds.PushPrivateKey)
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
		plaintextByte, err = ece.Decrypt(p.EncryptedBody,
			ece.WithEncoding(p.ContentEncoding),
			ece.WithAuthSecret(p.AuthSecret),
			ece.WithPrivate(p.PrivateKey),
			ece.WithDh(p.DH),
		)
	}
	if err != nil {
		return
	}
	p.PlainText = string(plaintextByte)
	return
}

func PostWebPush(w http.ResponseWriter, r *http.Request) {
	p, err := setPayload(r)
	if err != nil {
		fmt.Println("error Decrypt.", err.Error())
		ErrorResponse(w, r, http.StatusInternalServerError, err)
		return
	}
	fmt.Printf("%+v", p)
	fmt.Println()
}
