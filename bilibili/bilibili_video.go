package bilibili

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"regexp"
	"strconv"
	"time"
)

type VideoInfo struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Data    struct {
		Bvid  string `json:"bvid,omitempty"`
		Pic   string `json:"pic,omitempty"`
		Title string `json:"title,omitempty"`
		Desc  string `json:"desc,omitempty"`
		Owner struct {
			Name string `json:"name,omitempty"`
		} `json:"owner"`
	} `json:"data"`
}

func GetVideoInfoByAvid(avid uint64) (*VideoInfo, error) {
	resp, err := resty.New().SetTimeout(20*time.Second).SetHeader("Content-Type", "application/x-www-form-urlencoded").SetLogger(logger).SetCookies(cookies).R().
		SetQueryParam("aid", strconv.FormatUint(avid, 10)).Get("https://api.bilibili.com/x/web-interface/view")
	if err != nil {
		return nil, errors.Wrap(err, "获取视频详细信息失败")
	}
	if resp.StatusCode() != 200 {
		return nil, errors.Errorf("获取视频详细信息失败，错误码：%d", resp.StatusCode())
	}
	var ret *VideoInfo
	err = json.Unmarshal(resp.Body(), &ret)
	if err != nil {
		return nil, errors.Wrap(err, "解析json失败")
	}
	return ret, nil
}

func GetVideoInfoByBvid(bvid string) (*VideoInfo, error) {
	resp, err := resty.New().SetTimeout(20*time.Second).SetHeader("Content-Type", "application/x-www-form-urlencoded").SetLogger(logger).SetCookies(cookies).R().
		SetQueryParam("bvid", bvid).Get("https://api.bilibili.com/x/web-interface/view")
	if err != nil {
		return nil, errors.Wrap(err, "获取视频详细信息失败")
	}
	if resp.StatusCode() != 200 {
		return nil, errors.Errorf("获取视频详细信息失败，错误码：%d", resp.StatusCode())
	}
	var ret *VideoInfo
	err = json.Unmarshal(resp.Body(), &ret)
	if err != nil {
		return nil, errors.Wrap(err, "解析json失败")
	}
	return ret, nil
}

var regBv = regexp.MustCompile("(?i)bv([0-9A-Za-z]{10})")

func GetVideoInfoByShortUrl(shortUrl string) (*VideoInfo, error) {
	resp, _ := resty.New().SetRedirectPolicy(resty.NoRedirectPolicy()).SetTimeout(20 * time.Second).SetLogger(logger).SetCookies(cookies).R().Get(shortUrl)
	if resp == nil {
		return nil, errors.New("获取视频详细信息失败")
	}
	if resp.StatusCode() != 302 {
		return nil, errors.Errorf("获取视频详细信息失败，返回码：%d", resp.StatusCode())
	}
	url := resp.Header().Get("Location")
	ret := regBv.FindAllStringSubmatch(url, 1)
	if len(ret) != 1 {
		return nil, errors.New("通过短链接获取视频信息失败：" + url)
	}
	return GetVideoInfoByBvid(ret[0][0])
}
