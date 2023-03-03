package commandHandler

import (
	"github.com/CuteReimu/tfcc-bot-go/perm"
	"github.com/Mrs4s/MiraiGo/message"
	"strconv"
	"strings"
)

func init() {
	register(&delAdmin{})
	register(&addAdmin{})
	register(&listAllAdmin{})
}

type delAdmin struct{}

func (d *delAdmin) Name() string {
	return "删除管理员"
}

func (d *delAdmin) ShowTips(int64, int64) string {
	return "删除管理员 对方QQ号"
}

func (d *delAdmin) CheckAuth(_ int64, senderId int64) bool {
	return perm.IsSuperAdmin(senderId)
}

func (d *delAdmin) Execute(_ *message.GroupMessage, content string) (groupMsg *message.SendingMessage, privateMsg *message.SendingMessage) {
	var qqNumbers []int64
	for _, s := range strings.Split(content, " ") {
		s = strings.TrimSpace(s)
		if len(s) == 0 {
			continue
		}
		qq, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			logger.WithError(err).Errorf("parse failed: %s", s)
			return
		}
		if perm.IsSuperAdmin(qq) {
			groupMsg = message.NewSendingMessage().Append(message.NewText("你不能删除自己"))
			return
		}
		if !perm.IsAdmin(qq) {
			groupMsg = message.NewSendingMessage().Append(message.NewText(s + "并不是管理员"))
			return
		}
		qqNumbers = append(qqNumbers, qq)
	}
	if len(qqNumbers) == 0 {
		return
	}
	for _, qq := range qqNumbers {
		perm.DelAdmin(qq)
	}
	ret := "已删除管理员"
	if len(qqNumbers) == 1 {
		ret += "：" + strconv.FormatInt(qqNumbers[0], 10)
	}
	groupMsg = message.NewSendingMessage().Append(message.NewText(ret))
	return
}

type addAdmin struct{}

func (a *addAdmin) Name() string {
	return "增加管理员"
}

func (a *addAdmin) ShowTips(int64, int64) string {
	return "增加管理员 对方QQ号"
}

func (a *addAdmin) CheckAuth(_ int64, senderId int64) bool {
	return perm.IsSuperAdmin(senderId)
}

func (a *addAdmin) Execute(_ *message.GroupMessage, content string) (groupMsg *message.SendingMessage, privateMsg *message.SendingMessage) {
	var qqNumbers []int64
	for _, s := range strings.Split(content, " ") {
		s = strings.TrimSpace(s)
		if len(s) == 0 {
			continue
		}
		qq, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			logger.WithError(err).Errorf("parse failed: %s", s)
			return
		}
		if perm.IsSuperAdmin(qq) || perm.IsAdmin(qq) {
			groupMsg = message.NewSendingMessage().Append(message.NewText(s + "已经是管理员了"))
			return
		}
		qqNumbers = append(qqNumbers, qq)
	}
	if len(qqNumbers) == 0 {
		return
	}
	for _, qq := range qqNumbers {
		perm.AddAdmin(qq)
	}
	ret := "已增加管理员"
	if len(qqNumbers) == 1 {
		ret += "：" + strconv.FormatInt(qqNumbers[0], 10)
	}
	groupMsg = message.NewSendingMessage().Append(message.NewText(ret))
	return
}

type listAllAdmin struct{}

func (g *listAllAdmin) Name() string {
	return "查看管理员"
}

func (g *listAllAdmin) ShowTips(int64, int64) string {
	return "查看管理员"
}

func (g *listAllAdmin) CheckAuth(int64, int64) bool {
	return true
}

func (g *listAllAdmin) Execute(_ *message.GroupMessage, content string) (groupMsg *message.SendingMessage, privateMsg *message.SendingMessage) {
	if len(content) != 0 {
		return
	}
	list := perm.ListAdmin()
	if len(list) > 0 {
		groupMsg = message.NewSendingMessage().Append(message.NewText("管理员列表：\n" + strings.Join(list, "\n")))
	}
	return
}
