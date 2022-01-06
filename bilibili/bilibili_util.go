package bilibili

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"github.com/Logiase/MiraiGo-Template/config"
	"github.com/Logiase/MiraiGo-Template/utils"
	"time"
)

func encrypt(publicKey, data string) (string, error) {
	// pem解码
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		return "", errors.New("failed to decode public key")
	}
	// x509解码
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}
	pk := publicKeyInterface.(*rsa.PublicKey)
	// 加密
	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, pk, []byte(data))
	if err != nil {
		return "", err
	}
	// base64
	return base64.URLEncoding.EncodeToString(cipherText), nil
}

func getLiveUrl() string {
	return "https://live.bilibili.com/" + config.GlobalConfig.GetString("bilibili.room_id")
}

func getCookie(name string) string {
	now := time.Now()
	for _, cookie := range cookies {
		if cookie.Name == name && now.Before(cookie.Expires) {
			return cookie.Value
		}
	}
	return ""
}

var logger = utils.GetModuleLogger("tfcc-bot-go.bilibili")
