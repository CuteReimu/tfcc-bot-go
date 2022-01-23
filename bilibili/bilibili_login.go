package bilibili

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"github.com/Logiase/MiraiGo-Template/config"
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/Touhou-Freshman-Camp/tfcc-bot-go/db"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"net/http"
	"strings"
	"time"
)

var logger = utils.GetModuleLogger("tfcc-bot-go.bilibili")

var cookies []*http.Cookie

func Init() {
	savedCookies := db.Get([]byte("cookies"))
	if savedCookies != nil {
		cookies = (&resty.Response{RawResponse: &http.Response{Header: http.Header{
			"Set-Cookie": strings.Split(string(savedCookies), "\n"),
		}}}).Cookies()
		now := time.Now()
		upToDate := true
		for _, cookie := range cookies {
			if now.After(cookie.Expires) {
				upToDate = false
				break
			}
		}
		if upToDate {
			return
		}
	}
	client := resty.New().SetTimeout(20 * time.Second).SetLogger(logger)
	resp, err := client.R().SetQueryParam("plat", "6").Get("https://passport.bilibili.com/web/captcha/combine")
	if err != nil {
		logger.WithError(err).Fatalln("login failed")
	}
	if resp.StatusCode() != 200 {
		logger.Fatalf("login failed, status code: %d\n", resp.StatusCode())
	}
	if !gjson.ValidBytes(resp.Body()) {
		logger.Fatalf("json invalid: %s\n", resp.String())
	}
	loginResp := gjson.ParseBytes(resp.Body())
	if loginResp.Get("code").Int() != 0 {
		logger.Fatalf("登录bilibili获取人机校验失败, code: %d\n", loginResp.Get("code").Int())
	}
	gt := loginResp.Get("data.result.gt").String()
	challenge := loginResp.Get("data.result.challenge").String()
	key := loginResp.Get("data.result.key").String()
	fmt.Println("gt:", gt)
	fmt.Println("challenge:", challenge)
	fmt.Println("请前往以下链接进行人机验证：")
	fmt.Println("https://kuresaru.github.io/geetest-validator/")
	fmt.Println("验证后请输入validate：")
	var line string
	_, err = fmt.Scanln(&line)
	if err != nil {
		logger.WithError(err).Fatalln("读取stdin失败")
	}
	validate := strings.TrimSpace(line)
	seccode := validate
	resp, err = client.R().SetQueryParam("act", "getkey").Get("https://passport.bilibili.com/login")
	if err != nil {
		logger.WithError(err).Fatalln("登录bilibili失败")
	}
	if resp.StatusCode() != 200 {
		logger.Fatalf("登录bilibili失败, status code: %d\n", resp.StatusCode())
	}
	if !gjson.ValidBytes(resp.Body()) {
		logger.Fatalf("json invalid: %s\n", resp.String())
	}
	getKeyResp := gjson.ParseBytes(resp.Body())
	userName := config.GlobalConfig.GetString("bilibili.username")
	pwd := config.GlobalConfig.GetString("bilibili.password")
	encryptPwd, err := encrypt(getKeyResp.Get("key").String(), getKeyResp.Get("hash").String()+pwd)
	if err != nil {
		logger.WithError(err).Fatalln("encrypt failed")
	}
	resp, err = client.R().SetQueryParams(map[string]string{
		"captchaType": "6",
		"username":    userName,
		"password":    encryptPwd,
		"keep":        "true",
		"key":         key,
		"challenge":   challenge,
		"validate":    validate,
		"seccode":     seccode,
	}).Post("https://passport.bilibili.com/web/login/v2")
	if err != nil {
		logger.WithError(err).Fatalln("登录bilibili失败")
	}
	if resp.StatusCode() != 200 {
		logger.Fatalf("登录bilibili失败, status code: %d\n", resp.StatusCode())
	}
	if !gjson.ValidBytes(resp.Body()) {
		logger.Fatalf("json invalid: %s\n", resp.String())
	}
	loginSuccessResp := gjson.ParseBytes(resp.Body())
	if loginSuccessResp.Get("code").Int() != 0 {
		logger.Fatalf("登录bilibili失败，错误码：%d, 错误信息：%s\n", loginSuccessResp.Get("code").Int(), loginSuccessResp.Get("message").String())
	}
	logger.Infoln("登录bilibili成功")
	cookies = resp.Cookies()
	var cookieStrings []string
	for _, cookie := range cookies {
		cookieStrings = append(cookieStrings, cookie.String())
	}
	db.Set([]byte("cookies"), []byte(strings.Join(cookieStrings, "\n")))
}

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

func getCookie(name string) string {
	now := time.Now()
	for _, cookie := range cookies {
		if cookie.Name == name && now.Before(cookie.Expires) {
			return cookie.Value
		}
	}
	return ""
}
