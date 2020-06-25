package endpoints

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	ece "github.com/crow-misia/http-ece"
	"github.com/gorilla/mux"
	"github.com/potproject/ikuradon-salmon/dataAccess"
)

func PostWebPush(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subscribeID := vars["subscribeID"]

	headerEncryption := r.Header.Get("Encryption")
	headerContentEncoding := r.Header.Get("Content-Encoding")
	headerTtl := r.Header.Get("Ttl")
	headerContentType := r.Header.Get("Content-Type")
	headerCryptoKey := r.Header.Get("Crypto-Key")
	headerAuthorization := r.Header.Get("Authorization")
	b, _ := ioutil.ReadAll(r.Body)
	bodyB64 := base64.StdEncoding.EncodeToString(b)

	salt := ""
	if strings.HasPrefix(headerEncryption, "salt=") {
		salt = headerEncryption[5:]
	}
	sliceHeaderCryptoKey := strings.Split(headerCryptoKey, ";")
	dh := ""
	p256ecdsa := ""
	for _, ckey := range sliceHeaderCryptoKey {
		if strings.HasPrefix(ckey, "dh=") {
			dh = ckey[3:]
		}
		if strings.HasPrefix(ckey, "p256ecdsa=") {
			p256ecdsa = ckey[10:]
		}
	}
	jwt := ""
	if strings.HasPrefix(headerAuthorization, "WebPush ") {
		jwt = headerAuthorization[8:]
	}
	fmt.Println("Sid        :", subscribeID)
	fmt.Println("C-E        :", headerContentEncoding) // "aesgcm"
	fmt.Println("ttl        :", headerTtl)             // 172800
	fmt.Println("c-t        :", headerContentType)     // application/octet-stream
	fmt.Println("Encryption :", headerEncryption)      //
	fmt.Println("salt       :", salt)                  // EhJrnT2cqiZXXXXXX
	fmt.Println("dh         :", dh)                    // BCC42wgRWCcMIquAAyegXXXXXXXXXXXAhzIc61XXXXXXXXXPL5r2Ndh9RRGYvpaH2_BU
	fmt.Println("p256ecdsa  :", p256ecdsa)             // BPf7TFNX-XXXXXX_XXXXXXX-pyAI8sJyYYt62Dus0Mxpy8OF9kbG5gIxxxxxxXXkwsKvcnTA
	fmt.Println("jwt        :", jwt)                   // XXX.xxxx.xxxx
	fmt.Println("bodyb64    :", bodyB64)

	check, err := dataAccess.DA.Has(subscribeID)
	if err != nil {
		ErrorResponse(w, r, http.StatusNotFound, err)
		return
	}
	if !check {
		ErrorResponse(w, r, http.StatusNotFound, errors.New("NotFound"))
		return
	}
	ds, err := dataAccess.DA.Get(subscribeID)
	if err != nil {
		ErrorResponse(w, r, http.StatusInternalServerError, err)
		return
	}

	plaintext, err := ece.Decrypt(b,
		ece.WithEncoding(ece.AESGCM),
		ece.WithSalt(d(salt)),
		ece.WithAuthSecret(d(ds.PushAuth)),
		ece.WithPrivate(d(ds.PushPrivateKey)),
		ece.WithDh(d(dh)),
	)
	if err != nil {
		fmt.Println("error Decrypt.", err.Error())
		ErrorResponse(w, r, http.StatusInternalServerError, err)
		return
	}
	fmt.Println("plaintext  :", string(plaintext))
}

func d(text string) []byte {
	b, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		panic(err)
	}
	return b
}
