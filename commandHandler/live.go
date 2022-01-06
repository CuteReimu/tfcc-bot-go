package commandHandler

import (
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/Touhou-Freshman-Camp/tfcc-bot-go/bilibili"
	"github.com/Touhou-Freshman-Camp/tfcc-bot-go/perm"
	"github.com/ozgio/strutil"
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
	ret := bilibili.GetLiveStatus()
	if len(ret) != 0 {
		groupMsg = message.NewSendingMessage().Append(message.NewText(ret))
	}
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

func (s *startLive) Execute(_ *message.GroupMessage, content string) (groupMsg *message.SendingMessage, privateMsg *message.SendingMessage) {
	if len(content) != 0 {
		return
	}
	if len(content) != 0 {
		return
	}
	publicRet, privateRet := bilibili.StartLive()
	if len(publicRet) != 0 {
		groupMsg = message.NewSendingMessage().Append(message.NewText(publicRet))
	}
	if len(privateRet) != 0 {
		privateMsg = message.NewSendingMessage().Append(message.NewText(privateRet))
	}
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

func (s *stopLive) Execute(_ *message.GroupMessage, content string) (groupMsg *message.SendingMessage, privateMsg *message.SendingMessage) {
	if len(content) != 0 {
		return
	}
	ret := bilibili.StopLive()
	if len(ret) != 0 {
		groupMsg = message.NewSendingMessage().Append(message.NewText(ret))
	}
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

func (c *changeLiveTitle) Execute(_ *message.GroupMessage, content string) (groupMsg *message.SendingMessage, privateMsg *message.SendingMessage) {
	if len(content) == 0 {
		groupMsg = message.NewSendingMessage().Append(message.NewText("指令格式如下：\n修改直播标题 新标题"))
		return
	}
	if strutil.Len(content) > 20 {
		return
	}
	ret := bilibili.ChangeLiveTitle(content)
	if len(ret) != 0 {
		groupMsg = message.NewSendingMessage().Append(message.NewText(ret))
	}
	return
}
