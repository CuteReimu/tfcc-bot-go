package commandHandler

import (
	"fmt"
	"github.com/CuteReimu/tfcc-bot-go/tfcc"
	"github.com/Mrs4s/MiraiGo/message"
	"strings"
)

func init() {
	if tfcc.IsOK() {
		register(&jfNN{})
	}
}

type jfNN struct{}

func (j *jfNN) Name() string {
	return "分数表"
}

func (j *jfNN) ShowTips(int64, int64) string {
	return "分数表"
}

func (j *jfNN) CheckAuth(int64, int64) bool {
	return true
}

func (j *jfNN) Execute(_ *message.GroupMessage, content string) (groupMsg *message.SendingMessage, privateMsg *message.SendingMessage) {
	if len(content) <= 0 || len(content) > 50 {
		groupMsg = message.NewSendingMessage().Append(message.NewText(`请输入要查找的内容，例如“分数表 庙”`))
		return
	}
	pr := tfcc.ParseMsg(content)
	if len(pr.Work) == 0 {
		groupMsg = message.NewSendingMessage().Append(message.NewText("请至少限制一部作品"))
		return
	}
	if pr.Work == "8" && len(pr.Character) == 0 {
		pr.Character["RY"] = struct{}{}
		pr.Character["MA"] = struct{}{}
		pr.Character["SR"] = struct{}{}
		pr.Character["YY"] = struct{}{}
	}
	if len(pr.Rank) == 0 {
		pr.Rank = "Lunatic"
	}
	var rank string
	if len(pr.Rank) > 0 {
		rank = pr.Rank[:1]
		if pr.Rank == "Extra" {
			rank = "Ex"
		}
	}
	data := tfcc.GetJf(pr.Work)
	if len(data) == 0 {
		groupMsg = message.NewSendingMessage().Append(message.NewText("这一作品的" + rank + "难度目前没有数据，也许使用了交叉避弹或者其它分数表，请换一部作品试试"))
		return
	}
	text := []string{"作品：TH" + pr.Work + rank}
	for _, d := range data {
		if len(pr.Rank) > 0 && len(d.Rank) > 0 && pr.Rank != d.Rank {
			continue
		}
		if _, ok := pr.Route[d.Route]; !ok && len(pr.Route) > 0 && len(d.Route) > 0 {
			continue
		}
		if _, ok := pr.Character[d.Character]; !ok && len(pr.Character) > 0 && len(d.Character) > 0 {
			continue
		}
		if _, ok := pr.CType[d.CType]; !ok && len(pr.CType) > 0 && len(d.CType) > 0 {
			continue
		}
		if pr.AllSpell != d.AllSpell {
			continue
		}
		if d.Work == "8" {
			text = append(text, fmt.Sprintf("%s(%s)：%.2f", tfcc.Translate(d.Character), d.Route, d.Jf))
		} else {
			text = append(text, fmt.Sprintf("%s%s：%.2f", tfcc.Translate(d.Character), tfcc.Translate(d.CType), d.Jf))
		}
		if len(text) > 11 {
			groupMsg = message.NewSendingMessage().Append(message.NewText("请缩小查询范围"))
			return
		}
	}
	groupMsg = message.NewSendingMessage().Append(message.NewText(strings.Join(text, "\n")))
	return
}
