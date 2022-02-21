package commandHandler

import (
	"fmt"
	"github.com/CuteReimu/bilibili"
	"github.com/Logiase/MiraiGo-Template/config"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/Touhou-Freshman-Camp/tfcc-bot-go/db"
	"github.com/Touhou-Freshman-Camp/tfcc-bot-go/perm"
	"github.com/ozgio/strutil"
	"strconv"
)

func init() {
	register(&getLiveState{})
	register(&startLive{})
	register(&stopLive{})
	register(&changeLiveTitle{})
}

type getLiveState struct{}

func (g *getLiveState) Name() string {
	return "直播状态"
}

func (g *getLiveState) ShowTips(int64, int64) string {
	return "直播状态"
}

func (g *getLiveState) CheckAuth(int64, int64) bool {
	return true
}

func (g *getLiveState) Execute(_ *message.GroupMessage, content string) (groupMsg *message.SendingMessage, privateMsg *message.SendingMessage) {
	if len(content) != 0 {
		return
	}
	rid := config.GlobalConfig.GetInt("bilibili.room_id")
	ret, err := bilibili.GetRoomInfo(rid)
	if err != nil {
		logger.WithError(err).Errorln("获取直播状态失败")
		return
	}
	var text string
	if ret.LiveStatus == 0 {
		text = "直播间状态：未开播"
	} else {
		text = fmt.Sprintf("直播间状态：开播\n直播标题：%s\n人气：%d\n直播间地址：%s", ret.Title, ret.Online, getLiveUrl())
	}
	groupMsg = message.NewSendingMessage().Append(message.NewText(text))
	return
}

type startLive struct{}

func (s *startLive) Name() string {
	return "开始直播"
}

func (s *startLive) ShowTips(int64, int64) string {
	return "开始直播"
}

func (s *startLive) CheckAuth(_ int64, senderId int64) bool {
	return perm.IsWhitelist(senderId)
}

func (s *startLive) Execute(msg *message.GroupMessage, content string) (groupMsg *message.SendingMessage, privateMsg *message.SendingMessage) {
	if len(content) != 0 {
		return
	}
	rid := config.GlobalConfig.GetInt("bilibili.room_id")
	area := config.GlobalConfig.GetInt("bilibili.area_v2")
	ret, err := bilibili.StartLive(rid, area)
	if err != nil {
		logger.WithError(err).Errorln("开启直播间失败")
		return
	}
	var publicText string
	if ret.Change == 0 {
		val := db.Get([]byte("bilibili_live"))
		if val != nil {
			uin, _ := strconv.ParseInt(string(val), 10, 64)
			if uin != msg.Sender.Uin {
				publicText = fmt.Sprintf("已经有人正在直播了\n直播间地址：%s\n快来围观吧！", getLiveUrl())
				groupMsg = message.NewSendingMessage().Append(message.NewText(publicText))
				return
			}
		} else {
			db.Set([]byte("bilibili_live"), []byte(strconv.FormatInt(msg.Sender.Uin, 10)))
		}
		publicText = fmt.Sprintf("直播间本来就是开启的，推流码已私聊\n直播间地址：%s\n快来围观吧！", getLiveUrl())
	} else {
		db.Set([]byte("bilibili_live"), []byte(strconv.FormatInt(msg.Sender.Uin, 10)))
		publicText = fmt.Sprintf("直播间已开启，推流码已私聊，别忘了修改直播间标题哦！\n直播间地址：%s\n快来围观吧！", getLiveUrl())
	}
	privateText := fmt.Sprintf("RTMP推流地址：%s\n密钥：%s", ret.Rtmp.Addr, ret.Rtmp.Code)
	groupMsg = message.NewSendingMessage().Append(message.NewText(publicText))
	privateMsg = message.NewSendingMessage().Append(message.NewText(privateText))
	return
}

type stopLive struct{}

func (s *stopLive) Name() string {
	return "关闭直播"
}

func (s *stopLive) ShowTips(int64, int64) string {
	return "关闭直播"
}

func (s *stopLive) CheckAuth(_ int64, senderId int64) bool {
	return perm.IsWhitelist(senderId)
}

func (s *stopLive) Execute(msg *message.GroupMessage, content string) (groupMsg *message.SendingMessage, privateMsg *message.SendingMessage) {
	if len(content) != 0 {
		return
	}
	if !perm.IsAdmin(msg.Sender.Uin) {
		val := db.Get([]byte("bilibili_live"))
		if val != nil {
			uin, _ := strconv.ParseInt(string(val), 10, 64)
			if uin != msg.Sender.Uin {
				groupMsg = message.NewSendingMessage().Append(message.NewText("谢绝唐突关闭直播"))
				return
			}
		}
	}
	rid := config.GlobalConfig.GetInt("bilibili.room_id")
	changed, err := bilibili.StopLive(rid)
	if err != nil {
		logger.WithError(err).Errorln("关闭直播间失败")
		return
	}
	db.Del([]byte("bilibili_live"))
	var text string
	if !changed {
		text = "直播间本来就是关闭的"
	} else {
		text = "直播间已关闭"
	}
	groupMsg = message.NewSendingMessage().Append(message.NewText(text))
	return
}

type changeLiveTitle struct{}

func (c *changeLiveTitle) Name() string {
	return "修改直播标题"
}

func (c *changeLiveTitle) ShowTips(int64, int64) string {
	return "修改直播标题 新标题"
}

func (c *changeLiveTitle) CheckAuth(_ int64, senderId int64) bool {
	return perm.IsWhitelist(senderId)
}

func (c *changeLiveTitle) Execute(msg *message.GroupMessage, content string) (groupMsg *message.SendingMessage, privateMsg *message.SendingMessage) {
	if len(content) == 0 {
		groupMsg = message.NewSendingMessage().Append(message.NewText("指令格式如下：\n修改直播标题 新标题"))
		return
	}
	if strutil.Len(content) > 20 {
		return
	}
	if !perm.IsAdmin(msg.Sender.Uin) {
		val := db.Get([]byte("bilibili_live"))
		if val != nil {
			uin, _ := strconv.ParseInt(string(val), 10, 64)
			if uin != msg.Sender.Uin {
				groupMsg = message.NewSendingMessage().Append(message.NewText("谢绝唐突修改直播标题"))
				return
			}
		}
	}
	rid := config.GlobalConfig.GetInt("bilibili.room_id")
	err := bilibili.UpdateLive(rid, content)
	var text string
	if err != nil {
		logger.WithError(err).Errorln("修改直播间标题失败")
		text = "修改直播间标题失败，请联系管理员"
	} else {
		text = "直播间标题已修改为：" + content
	}
	groupMsg = message.NewSendingMessage().Append(message.NewText(text))
	return
}

func getLiveUrl() string {
	return "https://live.bilibili.com/" + config.GlobalConfig.GetString("bilibili.room_id")
}
