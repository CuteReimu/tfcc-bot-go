package bilibili

import (
	"fmt"
	"github.com/Logiase/MiraiGo-Template/config"
	"github.com/Touhou-Freshman-Camp/tfcc-bot-go/db"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"net/http"
	"strings"
	"time"
)

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

func GetLiveStatus() string {
	rid := config.GlobalConfig.GetString("bilibili.room_id")
	resp, err := resty.New().SetTimeout(20*time.Second).SetHeader("Content-Type", "application/x-www-form-urlencoded").SetLogger(logger).SetCookies(cookies).
		R().SetQueryParam("id", rid).Get("https://api.live.bilibili.com/room/v1/Room/get_info")
	if err != nil {
		logger.WithError(err).Errorln("请求直播间信息失败")
		return ""
	}
	if resp.StatusCode() != 200 {
		logger.Errorf("请求直播间信息失败，错误码：%d，返回内容：%s\n", resp.StatusCode(), resp.String())
		return ""
	}
	if !gjson.ValidBytes(resp.Body()) {
		logger.Errorf("解析json失败：%s\n", resp.String())
		return ""
	}
	liveStatusResp := gjson.ParseBytes(resp.Body())
	if liveStatusResp.Get("code").Int() != 0 {
		logger.Errorf("请求直播间状态失败，错误码：%d，错误信息：%s\n", liveStatusResp.Get("code").Int(), liveStatusResp.Get("message").String())
		return ""
	}
	if liveStatusResp.Get("data.live_status").Int() == 0 {
		return "直播间状态：未开播"
	}
	return fmt.Sprintf("直播间状态：开播\n直播标题：%s\n人气：%d\n直播间地址：%s",
		liveStatusResp.Get("data.title").String(),
		liveStatusResp.Get("data.online").Int(),
		getLiveUrl())
}

func StartLive() (string, string) {
	biliJct := getCookie("bili_jct")
	if len(biliJct) == 0 {
		return "B站登录过期", ""
	}
	rid := config.GlobalConfig.GetString("bilibili.room_id")
	area := config.GlobalConfig.GetString("bilibili.area_v2")
	resp, err := resty.New().SetTimeout(20*time.Second).SetHeader("Content-Type", "application/x-www-form-urlencoded").SetLogger(logger).SetCookies(cookies).R().SetQueryParams(map[string]string{
		"room_id":    rid,
		"platform":   "pc",
		"area_v2":    area,
		"csrf_token": biliJct,
		"csrf":       biliJct,
	}).Post("https://api.live.bilibili.com/room/v1/Room/startLive")
	if err != nil {
		logger.WithError(err).Errorln("开启直播间失败")
		return "", ""
	}
	if resp.StatusCode() != 200 {
		logger.Errorf("开启直播间失败，错误码：%d，返回内容：%s\n", resp.StatusCode(), resp.String())
		return "", ""
	}
	if !gjson.ValidBytes(resp.Body()) {
		logger.Errorf("解析json失败：%s\n", resp.String())
		return "", ""
	}
	startLiveResp := gjson.ParseBytes(resp.Body())
	if startLiveResp.Get("code").Int() != 0 {
		logger.Errorf("开启直播间失败，错误码：%d，错误信息1：%s，错误信息2：%s\n", startLiveResp.Get("code").Int(), startLiveResp.Get("message").String(), startLiveResp.Get("msg").String())
		return "", ""
	}
	var ret string
	if startLiveResp.Get("data.change").Int() == 0 {
		ret = fmt.Sprintf("直播间本来就是开启的，推流码已私聊\n直播间地址：%s\n快来围观吧！", getLiveUrl())
	} else {
		ret = fmt.Sprintf("直播间已开启，推流码已私聊，别忘了修改直播间标题哦！\n直播间地址：%s\n快来围观吧！", getLiveUrl())
	}
	rtmpAddr := startLiveResp.Get("data.rtmp.addr").String()
	rtmpCode := startLiveResp.Get("data.rtmp.code").String()
	return ret, fmt.Sprintf("RTMP推流地址：%s\n秘钥：%s", rtmpAddr, rtmpCode)
}

func StopLive() string {
	biliJct := getCookie("bili_jct")
	if len(biliJct) == 0 {
		return "B站登录过期"
	}
	rid := config.GlobalConfig.GetString("bilibili.room_id")
	resp, err := resty.New().SetTimeout(20*time.Second).SetHeader("Content-Type", "application/x-www-form-urlencoded").SetLogger(logger).SetCookies(cookies).R().SetQueryParams(map[string]string{
		"room_id": rid,
		"csrf":    biliJct,
	}).Post("https://api.live.bilibili.com/room/v1/Room/stopLive")
	if err != nil {
		logger.WithError(err).Errorln("关闭直播间失败")
		return ""
	}
	if resp.StatusCode() != 200 {
		logger.Errorf("关闭直播间失败，错误码：%d，返回内容：%s\n", resp.StatusCode(), resp.String())
		return ""
	}
	if !gjson.ValidBytes(resp.Body()) {
		logger.Errorf("解析json失败：%s\n", resp.String())
		return ""
	}
	stopLiveResp := gjson.ParseBytes(resp.Body())
	if stopLiveResp.Get("code").Int() != 0 {
		logger.Errorf("关闭直播间失败，错误码：%d，错误信息1：%s，错误信息2：%s\n", stopLiveResp.Get("code").Int(), stopLiveResp.Get("message").String(), stopLiveResp.Get("msg").String())
		return ""
	}
	if stopLiveResp.Get("data.change").Int() == 0 {
		return "直播间本来就是关闭的"
	}
	return "直播间已关闭"
}

func ChangeLiveTitle(title string) string {
	biliJct := getCookie("bili_jct")
	if len(biliJct) == 0 {
		return "B站登录过期"
	}
	rid := config.GlobalConfig.GetString("bilibili.room_id")
	resp, err := resty.New().SetTimeout(20*time.Second).SetHeader("Content-Type", "application/x-www-form-urlencoded").SetLogger(logger).SetCookies(cookies).R().SetQueryParams(map[string]string{
		"room_id": rid,
		"title":   title,
		"csrf":    biliJct,
	}).Post("https://api.live.bilibili.com/room/v1/Room/update")
	if err != nil {
		logger.WithError(err).Errorln("修改直播间标题失败")
		return ""
	}
	if resp.StatusCode() != 200 {
		logger.Errorf("修改直播间标题失败，错误码：%d，返回内容：%s\n", resp.StatusCode(), resp.String())
		return ""
	}
	if !gjson.ValidBytes(resp.Body()) {
		logger.Errorf("解析json失败：%s\n", resp.String())
		return ""
	}
	changeLiveTitleResp := gjson.ParseBytes(resp.Body())
	if changeLiveTitleResp.Get("code").Int() != 0 {
		logger.Errorf("修改直播间标题失败，错误码：%d，错误信息1：%s，错误信息2：%s\n", changeLiveTitleResp.Get("code").Int(), changeLiveTitleResp.Get("message").String(), changeLiveTitleResp.Get("msg").String())
		return "修改直播间标题失败，请联系管理员"
	}
	return "直播间标题已修改为：" + title
}
