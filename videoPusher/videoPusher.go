package videoPusher

import (
	"bytes"
	"fmt"
	"github.com/CuteReimu/bilibili"
	"github.com/CuteReimu/tfcc-bot-go/bot"
	"github.com/CuteReimu/tfcc-bot-go/config"
	"github.com/CuteReimu/tfcc-bot-go/db"
	"github.com/CuteReimu/tfcc-bot-go/repeaterInterruption"
	"github.com/CuteReimu/tfcc-bot-go/utils"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/go-resty/resty/v2"
	"github.com/ozgio/strutil"
	"sync"
	"time"
)

func init() {
	bot.RegisterModule(instance)
}

var instance = &mh{}
var logger = utils.GetModuleLogger("tfcc-bot-go.videoPusher")

type mh struct {
}

func (m *mh) MiraiGoModule() bot.ModuleInfo {
	return bot.ModuleInfo{
		ID:       "tfcc-bot-go.videoPusher",
		Instance: instance,
	}
}

func (m *mh) Init() {
}

func (m *mh) PostInit() {
}

func (m *mh) Serve(b *bot.Bot) {
	delay := config.GlobalConfig.GetInt64("schedule.video_push_delay")
	if delay <= 0 {
		return
	}
	qqGroups := config.GlobalConfig.GetIntSlice("schedule.qq_group")
	if len(qqGroups) == 0 {
		return
	}
	go func() {
		ticker := time.NewTicker(time.Duration(delay) * time.Second)
		defer ticker.Stop()
		for {
			<-ticker.C
			video := getNewVideo()
			if video != nil {
				for _, qqGroup := range qqGroups {
					groupCode := int64(qqGroup)
					groupMsg := message.NewSendingMessage()
					var text string
					if len(video.Pic) > 0 {
						resp, err := resty.New().SetTimeout(20 * time.Second).SetLogger(logger).R().Get(video.Pic)
						if err != nil {
							logger.WithError(err).Errorln("获取视频封面失败")
						} else {
							elem, err := b.UploadImage(message.Source{SourceType: message.SourceGroup, PrimaryID: groupCode}, bytes.NewReader(resp.Body()))
							if err != nil {
								logger.WithError(err).Errorln("上传封面失败")
							} else {
								groupMsg.Append(elem)
								text = "\n"
							}
						}
					}
					if newStr, err := strutil.Substring(video.Description, 0, 100); err == nil {
						video.Description = newStr + "。。。"
					}
					groupMsg.Append(message.NewText(fmt.Sprintf(text+"%s\nhttps://www.bilibili.com/video/%s\nUP主：%s\n视频简介：%s", video.Title, video.Bvid, video.Author, video.Description)))
					retGroupMsg := b.SendGroupMessage(groupCode, groupMsg)
					if retGroupMsg == nil {
						logger.Info("群聊消息发送失败了")
					} else if retGroupMsg.Id == -1 {
						logger.Info("群聊消息被风控了")
					} else {
						repeaterInterruption.Clean(groupCode)
					}
				}
			}
		}
	}()
}

func (m *mh) Start(*bot.Bot) {
}

func (m *mh) Stop(_ *bot.Bot, wg *sync.WaitGroup) {
	defer wg.Done()
}

func getNewVideo() *bilibili.Video {
	mid := config.GlobalConfig.GetInt("bilibili.mid")
	videoList, err := bilibili.GetUserVideos(mid, bilibili.OrderPubDate, 0, "", 1, 1)
	if err != nil {
		logger.WithError(err).Errorln("获取用户视频失败")
		return nil
	}
	var newVideo *bilibili.Video
	db.Update([]byte("latest_video_id"), func(oldValue []byte) []byte {
		var latestId string
		if oldValue != nil {
			latestId = string(oldValue)
		}
		if len(videoList.List.Vlist) > 0 {
			latestVideo := videoList.List.Vlist[0]
			if latestId != latestVideo.Bvid {
				newVideo = &latestVideo
				return []byte(latestVideo.Bvid)
			}
		}
		return nil
	})
	return newVideo
}
