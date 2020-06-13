package endpoints

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func PostWebPush(w http.ResponseWriter, r *http.Request) {
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
	fmt.Println("C-E        :", headerContentEncoding) // "aesgcm"
	fmt.Println("ttl        :", headerTtl)             // 172800
	fmt.Println("c-t        :", headerContentType)     // application/octet-stream
	fmt.Println("salt       :", salt)                  // EhJrnT2cqiZXXXXXX
	fmt.Println("dh         :", dh)                    // BCC42wgRWCcMIquAAyegXXXXXXXXXXXAhzIc61XXXXXXXXXPL5r2Ndh9RRGYvpaH2_BU
	fmt.Println("p256ecdsa  :", p256ecdsa)             // BPf7TFNX-XXXXXX_XXXXXXX-pyAI8sJyYYt62Dus0Mxpy8OF9kbG5gIxxxxxxXXkwsKvcnTA
	fmt.Println("jwt        :", jwt)                   // XXX.xxxx.xxxx
	fmt.Println("bodyb64    :", bodyB64)
}
