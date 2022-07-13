package rfc8188

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/elliptic"
	"crypto/hmac"
	"crypto/sha256"

	"github.com/aead/ecdh"
)

// 参考: https://github.com/pantasystem/Milktea/blob/develop/PushToFCM/webPushDecipher.js
func Decrypt(body []byte, public []byte, private []byte, authSecret []byte) ([]byte, error) {
	salt := body[0:16]
	//rs := body[16:20]
	idlen := body[20:21]
	keyid := body[21 : 21+uint64(idlen[0])]
	content := body[21+uint64(idlen[0]):]

	// Create ECDH Public
	p256 := ecdh.Generic(elliptic.P256())
	sharedSecret := p256.ComputeSecret(private, keyid)

	/*
	  # HKDF-Extract(salt=auth_secret, IKM=ecdh_secret)
	  PRK_key = HMAC-SHA-256(auth_secret, ecdh_secret)
	  # HKDF-Expand(PRK_key, key_info, L_key=32)
	  key_info = "WebPush: info" || 0x00 || ua_public || as_public
	  IKM = HMAC-SHA-256(PRK_key, key_info || 0x01)
	  ## HKDF calculations from RFC 8188
	  # HKDF-Extract(salt, IKM)
	  PRK = HMAC-SHA-256(salt, IKM)
	  # HKDF-Expand(PRK, cek_info, L_cek=16)
	  cek_info = "Content-Encoding: aes128gcm" || 0x00
	  CEK = HMAC-SHA-256(PRK, cek_info || 0x01)[0..15]
	  # HKDF-Expand(PRK, nonce_info, L_nonce=12)
	  nonce_info = "Content-Encoding: nonce" || 0x00
	  NONCE = HMAC-SHA-256(PRK, nonce_info || 0x01)[0..11]
	*/

	// prkKey
	prkKey := HMACSHA256(authSecret, sharedSecret)

	// KeyInfo
	keyInfo := []byte{}
	keyInfo = append(keyInfo, []byte("WebPush: info")...)
	keyInfo = append(keyInfo, 0x00)
	keyInfo = append(keyInfo, public...)
	keyInfo = append(keyInfo, keyid...)

	// IKM
	ikm := HMACSHA256(prkKey, append(keyInfo, 0x01))
	prk := HMACSHA256(salt, ikm)

	// cekInfo
	cekInfo := []byte{}
	cekInfo = append(cekInfo, []byte("Content-Encoding: aes128gcm")...)
	cekInfo = append(cekInfo, 0x00)

	// cek
	cek := HMACSHA256(prk, append(cekInfo, 0x01))[0:16]

	// nonceInfo
	nonceInfo := []byte{}
	nonceInfo = append(nonceInfo, []byte("Content-Encoding: nonce")...)
	nonceInfo = append(nonceInfo, 0x00)

	// nonce
	nonce := HMACSHA256(prk, append(nonceInfo, 0x01))[0:12]

	// aes-128-gcm
	block, err := aes.NewCipher(cek)
	if err != nil {
		return []byte{}, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return []byte{}, err
	}
	plain, err := aesgcm.Open(nil, nonce, content, nil)
	if err != nil {
		return []byte{}, err
	}

	return plain, nil
}

func HMACSHA256(key, data []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}
