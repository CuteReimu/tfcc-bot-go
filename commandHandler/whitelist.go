package commandHandler

import (
	"fmt"
	"github.com/CuteReimu/tfcc-bot-go/perm"
	"github.com/Mrs4s/MiraiGo/message"
	"strconv"
	"strings"
)

func init() {
	register(&delWhitelist{})
	register(&addWhitelist{})
	register(&listAllWhitelist{})
	register(&checkWhitelist{})
	register(&enableAllWhitelist{})
	register(&disableAllWhitelist{})
}

type delWhitelist struct{}

func (d *delWhitelist) Name() string {
	return "删除白名单"
}

func (d *delWhitelist) ShowTips(int64, int64) string {
	return "删除白名单 对方QQ号"
}

func (d *delWhitelist) CheckAuth(_ int64, senderId int64) bool {
	return perm.IsAdmin(senderId)
}

func (d *delWhitelist) Execute(_ *message.GroupMessage, content string) (groupMsg *message.SendingMessage, privateMsg *message.SendingMessage) {
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
		if !perm.IsWhitelist(qq) {
			groupMsg = message.NewSendingMessage().Append(message.NewText(s + "并不是白名单"))
			return
		}
		qqNumbers = append(qqNumbers, qq)
	}
	if len(qqNumbers) == 0 {
		groupMsg = message.NewSendingMessage().Append(message.NewText("指令格式如下：\n删除白名单 对方QQ号"))
		return
	}
	for _, qq := range qqNumbers {
		perm.DelWhitelist(qq)
	}
	ret := "已删除白名单"
	if len(qqNumbers) == 1 {
		ret += "：" + strconv.FormatInt(qqNumbers[0], 10)
	}
	groupMsg = message.NewSendingMessage().Append(message.NewText(ret))
	return
}

type addWhitelist struct{}

func (a *addWhitelist) Name() string {
	return "增加白名单"
}

func (a *addWhitelist) ShowTips(int64, int64) string {
	return "增加白名单 对方QQ号"
}

func (a *addWhitelist) CheckAuth(_ int64, senderId int64) bool {
	return perm.IsAdmin(senderId)
}

func (a *addWhitelist) Execute(_ *message.GroupMessage, content string) (groupMsg *message.SendingMessage, privateMsg *message.SendingMessage) {
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
		if perm.IsWhitelist(qq) {
			groupMsg = message.NewSendingMessage().Append(message.NewText(s + "已经是白名单了"))
			return
		}
		qqNumbers = append(qqNumbers, qq)
	}
	if len(qqNumbers) == 0 {
		groupMsg = message.NewSendingMessage().Append(message.NewText("指令格式如下：\n增加白名单 对方QQ号"))
		return
	}
	for _, qq := range qqNumbers {
		perm.AddWhitelist(qq)
	}
	ret := "已增加白名单"
	if len(qqNumbers) == 1 {
		ret += "：" + strconv.FormatInt(qqNumbers[0], 10)
	}
	groupMsg = message.NewSendingMessage().Append(message.NewText(ret))
	return
}

type listAllWhitelist struct{}

func (g *listAllWhitelist) Name() string {
	return "列出所有白名单"
}

func (g *listAllWhitelist) ShowTips(int64, int64) string {
	return "列出所有白名单"
}

func (g *listAllWhitelist) CheckAuth(_ int64, senderId int64) bool {
	return perm.IsAdmin(senderId)
}

func (g *listAllWhitelist) Execute(_ *message.GroupMessage, content string) (groupMsg *message.SendingMessage, privateMsg *message.SendingMessage) {
	if len(content) != 0 {
		return
	}
	list := perm.ListWhitelist()
	if len(list) > 0 {
		groupMsg = message.NewSendingMessage().Append(message.NewText("白名单列表：\n" + strings.Join(list, "\n")))
	}
	return
}

type checkWhitelist struct{}

func (e *checkWhitelist) Name() string {
	return "查看白名单"
}

func (e *checkWhitelist) ShowTips(int64, int64) string {
	return "查看白名单 对方QQ号"
}

func (e *checkWhitelist) CheckAuth(int64, int64) bool {
	return true
}

func (e *checkWhitelist) Execute(_ *message.GroupMessage, content string) (groupMsg *message.SendingMessage, privateMsg *message.SendingMessage) {
	qq, err := strconv.ParseInt(content, 10, 64)
	if err != nil {
		groupMsg = message.NewSendingMessage().Append(message.NewText("指令格式如下：\n查看白名单 对方QQ号"))
		return
	}
	if perm.IsWhitelist(qq) {
		groupMsg = message.NewSendingMessage().Append(message.NewText(content + "是白名单"))
	} else {
		groupMsg = message.NewSendingMessage().Append(message.NewText(content + "不是白名单"))
	}
	return
}

type enableAllWhitelist struct{}

func (e *enableAllWhitelist) Name() string {
	return "启用所有白名单"
}

func (e *enableAllWhitelist) ShowTips(int64, int64) string {
	return "启用所有白名单"
}

func (e *enableAllWhitelist) CheckAuth(_ int64, senderId int64) bool {
	return perm.IsSuperAdmin(senderId)
}

func (e *enableAllWhitelist) Execute(*message.GroupMessage, string) (groupMsg *message.SendingMessage, privateMsg *message.SendingMessage) {
	count := perm.EnableAllWhitelist()
	groupMsg = message.NewSendingMessage().Append(message.NewText(fmt.Sprintf("已启用%d个白名单", count)))
	return
}

type disableAllWhitelist struct{}

func (d *disableAllWhitelist) Name() string {
	return "禁用所有白名单"
}

func (d *disableAllWhitelist) ShowTips(int64, int64) string {
	return "禁用所有白名单"
}

func (d *disableAllWhitelist) CheckAuth(_ int64, senderId int64) bool {
	return perm.IsSuperAdmin(senderId)
}

func (d *disableAllWhitelist) Execute(*message.GroupMessage, string) (groupMsg *message.SendingMessage, privateMsg *message.SendingMessage) {
	count := perm.DisableAllWhitelist()
	groupMsg = message.NewSendingMessage().Append(message.NewText(fmt.Sprintf("已禁用%d个白名单", count)))
	return
}
