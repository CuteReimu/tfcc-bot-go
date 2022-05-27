package chatPipeline

import (
	"bytes"
	"fmt"
	"github.com/CuteReimu/bilibili"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	regexp "github.com/dlclark/regexp2"
	"github.com/go-resty/resty/v2"
	"github.com/ozgio/strutil"
	"github.com/pkg/errors"
	"strconv"
	"time"
)

func init() {
	register(newBilibiliVideoAnalysis())
}

type bilibiliVideoAnalysis struct {
	avReg, bvReg, shortReg *regexp.Regexp
}

func newBilibiliVideoAnalysis() *bilibiliVideoAnalysis {
	return &bilibiliVideoAnalysis{
		avReg:    regexp.MustCompile(`(?<![A-Za-z0-9])(?:https?://www\.bilibili\.com/video/)?av(\d+)`, regexp.IgnoreCase),
		bvReg:    regexp.MustCompile(`(?<![A-Za-z0-9])(?:https?://www\.bilibili\.com/video/|https?://b23\.tv)?bv([0-9A-Za-z]{10})`, regexp.IgnoreCase),
		shortReg: regexp.MustCompile(`(?<![A-Za-z0-9])https?://b23\.tv/[0-9A-Za-z]{7}`, regexp.IgnoreCase),
	}
}

func (b *bilibiliVideoAnalysis) Execute(c *client.QQClient, msg *message.GroupMessage, content string) (groupMsg *message.SendingMessage) {
	result, found, err := b.getVideoInfo(content)
	if found {
		if err != nil {
			logger.WithError(err).Errorln("获取视频信息失败")
			return
		}
		var text string
		groupMsg = message.NewSendingMessage()
		if len(result.Pic) > 0 {
			resp, err := resty.New().SetTimeout(20 * time.Second).SetLogger(logger).R().Get(result.Pic)
			if err != nil {
				logger.WithError(err).Errorln("获取视频封面失败")
			} else {
				elem, err := c.UploadImage(message.Source{SourceType: message.SourceGroup, PrimaryID: msg.GroupCode}, bytes.NewReader(resp.Body()))
				if err != nil {
					logger.WithError(err).Errorln("上传封面失败")
				} else {
					groupMsg.Append(elem)
					text = "\n"
				}
			}
		}
		if newStr, err := strutil.Substring(result.Desc, 0, 100); err == nil {
			result.Desc = newStr + "。。。"
		}
		groupMsg.Append(message.NewText(fmt.Sprintf(text+"%s\nhttps://www.bilibili.com/video/%s\nUP主：%s\n视频简介：%s",
			result.Title, result.Bvid, result.Owner.Name, result.Desc)))
	}
	return
}

func (b *bilibiliVideoAnalysis) getVideoInfo(content string) (*bilibili.VideoInfo, bool, error) {
	if avRes, _ := b.avReg.FindStringMatch(content); avRes != nil {
		avid, err := strconv.Atoi(avRes.GroupByNumber(1).String())
		if err != nil {
			return nil, true, errors.Wrap(err, "解析avid失败："+avRes.GroupByNumber(1).String())
		}
		result, err := bilibili.GetVideoInfoByAvid(avid)
		return result, true, err
	}
	if bvRes, _ := b.bvReg.FindStringMatch(content); bvRes != nil {
		result, err := bilibili.GetVideoInfoByBvid(bvRes.GroupByNumber(1).String())
		return result, true, err
	}
	if shortRes, _ := b.shortReg.FindStringMatch(content); shortRes != nil {
		result, err := bilibili.GetVideoInfoByShortUrl(shortRes.String())
		return result, true, err
	}
	return nil, false, nil
}
