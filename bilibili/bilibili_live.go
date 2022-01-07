package bilibili

import (
	"encoding/json"
	"github.com/Logiase/MiraiGo-Template/config"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"time"
)

type LiveStatus struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Data    struct {
		LiveStatus int    `json:"live_status,omitempty"`
		Title      string `json:"title,omitempty"`
		Online     int    `json:"online,omitempty"`
	} `json:"data,omitempty"`
}

func GetLiveStatus() (*LiveStatus, error) {
	rid := config.GlobalConfig.GetString("bilibili.room_id")
	resp, err := resty.New().SetTimeout(20*time.Second).SetHeader("Content-Type", "application/x-www-form-urlencoded").SetLogger(logger).SetCookies(cookies).
		R().SetQueryParam("id", rid).Get("https://api.live.bilibili.com/room/v1/Room/get_info")
	if err != nil {
		return nil, errors.Wrap(err, "请求直播间信息失败")
	}
	if resp.StatusCode() != 200 {
		return nil, errors.Errorf("请求直播间信息失败，错误码：%d，返回内容：%s", resp.StatusCode(), resp.String())
	}
	var ret *LiveStatus
	err = json.Unmarshal(resp.Body(), &ret)
	if err != nil {
		return nil, errors.Wrapf(err, "解析json失败：%s", resp.String())
	}
	return ret, nil
}

type StartLiveResp struct {
	Code    int    `json:"code,omitempty"`
	Msg     string `json:"msg,omitempty"`
	Message string `json:"message,omitempty"`
	Data    struct {
		Change int `json:"change,omitempty"`
		Rtmp   struct {
			Addr string `json:"addr,omitempty"`
			Code string `json:"code,omitempty"`
		} `json:"rtmp"`
	} `json:"data"`
}

func StartLive() (*StartLiveResp, error) {
	rid := config.GlobalConfig.GetString("bilibili.room_id")
	area := config.GlobalConfig.GetString("bilibili.area_v2")
	biliJct := getCookie("bili_jct")
	if len(biliJct) == 0 {
		return nil, errors.New("B站登录过期")
	}
	resp, err := resty.New().SetTimeout(20*time.Second).SetHeader("Content-Type", "application/x-www-form-urlencoded").SetLogger(logger).SetCookies(cookies).R().SetQueryParams(map[string]string{
		"room_id":    rid,
		"platform":   "pc",
		"area_v2":    area,
		"csrf_token": biliJct,
		"csrf":       biliJct,
	}).Post("https://api.live.bilibili.com/room/v1/Room/startLive")
	if err != nil {
		return nil, errors.Wrap(err, "开启直播间失败")
	}
	if resp.StatusCode() != 200 {
		return nil, errors.Wrapf(err, "开启直播间失败，错误码：%d，返回内容：%s", resp.StatusCode(), resp.String())
	}
	var ret *StartLiveResp
	err = json.Unmarshal(resp.Body(), &ret)
	if err != nil {
		return nil, errors.Wrapf(err, "解析json失败：%s", resp.String())
	}
	return ret, nil
}

type StopLiveResp struct {
	Code    int    `json:"code,omitempty"`
	Msg     string `json:"msg,omitempty"`
	Message string `json:"message,omitempty"`
	Data    struct {
		Change int `json:"change,omitempty"`
	} `json:"data,omitempty"`
}

func StopLive() (*StopLiveResp, error) {
	biliJct := getCookie("bili_jct")
	if len(biliJct) == 0 {
		return nil, errors.New("B站登录过期")
	}
	rid := config.GlobalConfig.GetString("bilibili.room_id")
	resp, err := resty.New().SetTimeout(20*time.Second).SetHeader("Content-Type", "application/x-www-form-urlencoded").SetLogger(logger).SetCookies(cookies).R().SetQueryParams(map[string]string{
		"room_id": rid,
		"csrf":    biliJct,
	}).Post("https://api.live.bilibili.com/room/v1/Room/stopLive")
	if err != nil {
		return nil, errors.Wrap(err, "关闭直播间失败")
	}
	if resp.StatusCode() != 200 {
		return nil, errors.Errorf("关闭直播间失败，错误码：%d，返回内容：%s", resp.StatusCode(), resp.String())
	}
	var ret *StopLiveResp
	err = json.Unmarshal(resp.Body(), &ret)
	if err != nil {
		return nil, errors.Wrapf(err, "解析json失败：%s", resp.String())
	}
	return ret, nil
}

type ChangeLiveTitleResp struct {
	Code    int    `json:"code,omitempty"`
	Msg     string `json:"msg,omitempty"`
	Message string `json:"message,omitempty"`
}

func ChangeLiveTitle(title string) (*ChangeLiveTitleResp, error) {
	biliJct := getCookie("bili_jct")
	if len(biliJct) == 0 {
		return nil, errors.New("B站登录过期")
	}
	rid := config.GlobalConfig.GetString("bilibili.room_id")
	resp, err := resty.New().SetTimeout(20*time.Second).SetHeader("Content-Type", "application/x-www-form-urlencoded").SetLogger(logger).SetCookies(cookies).R().SetQueryParams(map[string]string{
		"room_id": rid,
		"title":   title,
		"csrf":    biliJct,
	}).Post("https://api.live.bilibili.com/room/v1/Room/update")
	if err != nil {
		return nil, errors.Wrap(err, "修改直播间标题失败")
	}
	if resp.StatusCode() != 200 {
		return nil, errors.Errorf("修改直播间标题失败，错误码：%d，返回内容：%s\n", resp.StatusCode(), resp.String())
	}
	var ret *ChangeLiveTitleResp
	err = json.Unmarshal(resp.Body(), &ret)
	if err != nil {
		return nil, errors.Wrapf(err, "解析json失败：%s", resp.String())
	}
	return ret, nil
}

func GetLiveUrl() string {
	return "https://live.bilibili.com/" + config.GlobalConfig.GetString("bilibili.room_id")
}
