package commandHandler

import (
	"encoding/json"
	"fmt"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/Touhou-Freshman-Camp/tfcc-bot-go/db"
	"github.com/Touhou-Freshman-Camp/tfcc-bot-go/perm"
	"sort"
	"strings"
)

type voteData struct {
	GroupCode int64            `json:"group_code,omitempty"`
	Content   string           `json:"content,omitempty"`
	Cache     map[int64]string `json:"cache,omitempty"`
	Forbidden []string         `json:"forbidden,omitempty"`
}

func init() {
	register(&addVote{})
	register(&delVote{})
	register(&showVote{})
	register(&doVote{})
	register(&addVoteForbiddenWords{})
	register(&delVoteForbiddenWords{})
}

type addVote struct{}

func (a *addVote) Name() string {
	return "发起投票"
}

func (a *addVote) ShowTips(int64, int64) string {
	return "发起投票 投票内容"
}

func (a *addVote) CheckAuth(_ int64, senderId int64) bool {
	return perm.IsAdmin(senderId)
}

func (a *addVote) Execute(msg *message.GroupMessage, content string) (groupMsg *message.SendingMessage, privateMsg *message.SendingMessage) {
	if len(content) == 0 {
		groupMsg = message.NewSendingMessage().Append(message.NewText("指令格式如下：\n发起投票 投票内容"))
		return
	}
	db.Update([]byte("vote"), func(oldValue []byte) []byte {
		if oldValue != nil {
			groupMsg = message.NewSendingMessage().Append(message.NewText("目前只支持同时存在一个投票，请先停止当前的投票"))
			return nil
		}
		data := &voteData{
			GroupCode: msg.GroupCode,
			Content:   content,
		}
		newValue, err := json.Marshal(data)
		if err != nil {
			logger.WithError(err).Errorln("marshal json failed")
			return nil
		}
		groupMsg = message.NewSendingMessage().Append(message.NewText(fmt.Sprintf("发起“%s”的投票成功", content)))
		return newValue
	})
	return
}

type delVote struct{}

func (d *delVote) Name() string {
	return "确定清除投票"
}

func (d *delVote) ShowTips(int64, int64) string {
	if db.Get([]byte("vote")) == nil {
		return ""
	}
	return "确定清除投票"
}

func (d *delVote) CheckAuth(_ int64, senderId int64) bool {
	return perm.IsAdmin(senderId)
}

func (d *delVote) Execute(_ *message.GroupMessage, content string) (groupMsg *message.SendingMessage, privateMsg *message.SendingMessage) {
	if len(content) != 0 {
		return
	}
	db.Update([]byte("vote"), func(oldValue []byte) []byte {
		if oldValue == nil {
			groupMsg = message.NewSendingMessage().Append(message.NewText("目前没有正在进行的投票"))
			return nil
		}
		groupMsg = message.NewSendingMessage().Append(message.NewText("清除投票成功"))
		return []byte{}
	})
	return
}

type showVote struct{}

func (s *showVote) Name() string {
	return "查看投票"
}

func (s *showVote) ShowTips(int64, int64) string {
	if db.Get([]byte("vote")) == nil {
		return ""
	}
	return "查看投票"
}

func (s *showVote) CheckAuth(int64, int64) bool {
	return true
}

func (s *showVote) Execute(*message.GroupMessage, string) (groupMsg *message.SendingMessage, privateMsg *message.SendingMessage) {
	buf := db.Get([]byte("vote"))
	if buf == nil {
		groupMsg = message.NewSendingMessage().Append(message.NewText("目前没有正在进行的投票"))
		return
	}
	var data *voteData
	err := json.Unmarshal(buf, &data)
	if err != nil {
		logger.WithError(err).Errorln("unmarshal json failed")
		return
	}
	cache := make(map[string]int)
	reverseCache := make(map[int][]string)
	if data.Cache != nil {
		for _, v := range data.Cache {
			count, ok := cache[v]
			if ok {
				cache[v] = count + 1
			} else {
				cache[v] = 0
			}
		}
	}
	for k, v := range cache {
		c, ok := reverseCache[v]
		if ok {
			reverseCache[v] = append(c, k)
		} else {
			reverseCache[v] = []string{k}
		}
	}
	keys := make([]int, 0, len(reverseCache))
	for k := range reverseCache {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	count := 0
	text := fmt.Sprintf("投票“%s”正在火热进行中，目前排名前三的是：", data.Content)
	for i := len(keys) - 1; i >= 0; i-- {
		for _, v := range reverseCache[keys[i]] {
			text += fmt.Sprintf("\n%s %d票", v, keys[i])
			count++
			if count >= 3 {
				goto end
			}
		}
	}
end:
	groupMsg = message.NewSendingMessage().Append(message.NewText(text))
	return
}

type doVote struct{}

func (d *doVote) Name() string {
	return "投票"
}

func (d *doVote) ShowTips(int64, int64) string {
	if db.Get([]byte("vote")) == nil {
		return ""
	}
	return "投票 投票答案"
}

func (d *doVote) CheckAuth(int64, int64) bool {
	return db.Get([]byte("vote")) != nil
}

func (d *doVote) Execute(msg *message.GroupMessage, content string) (groupMsg *message.SendingMessage, privateMsg *message.SendingMessage) {
	if len(content) == 0 {
		groupMsg = message.NewSendingMessage().Append(message.NewText("指令格式如下：\n投票 投票答案"))
		return
	}
	if len(content) > 18 {
		return
	}
	db.Update([]byte("vote"), func(oldValue []byte) []byte {
		if oldValue == nil {
			groupMsg = message.NewSendingMessage().Append(message.NewText("目前没有正在进行的投票"))
			return nil
		}
		var data *voteData
		err := json.Unmarshal(oldValue, &data)
		if err != nil {
			logger.WithError(err).Errorln("unmarshal json failed")
			return nil
		}
		for _, forbidden := range data.Forbidden {
			if strings.Contains(content, forbidden) {
				groupMsg = message.NewSendingMessage().Append(message.NewText(fmt.Sprintf("“%s”被禁止了", forbidden)))
				return nil
			}
		}
		if data.Cache == nil {
			data.Cache = make(map[int64]string)
		}
		_, ok := data.Cache[msg.Sender.Uin]
		data.Cache[msg.Sender.Uin] = content
		newValue, err := json.Marshal(data)
		if err != nil {
			logger.WithError(err).Errorln("marshal json failed")
			return nil
		}
		if ok {
			groupMsg = message.NewSendingMessage().Append(message.NewText("你将投票结果改为：" + content))
		} else {
			groupMsg = message.NewSendingMessage().Append(message.NewText("你进行了投票：" + content))
		}
		return newValue
	})
	return
}

type addVoteForbiddenWords struct{}

func (a *addVoteForbiddenWords) Name() string {
	return "增加投票屏蔽词"
}

func (a *addVoteForbiddenWords) ShowTips(int64, int64) string {
	return ""
}

func (a *addVoteForbiddenWords) CheckAuth(_ int64, senderId int64) bool {
	return perm.IsAdmin(senderId)
}

func (a *addVoteForbiddenWords) Execute(_ *message.GroupMessage, content string) (groupMsg *message.SendingMessage, privateMsg *message.SendingMessage) {
	if len(content) == 0 {
		groupMsg = message.NewSendingMessage().Append(message.NewText("指令格式如下：\n增加投票屏蔽词 词语"))
		return
	}
	db.Update([]byte("vote"), func(oldValue []byte) []byte {
		if oldValue == nil {
			groupMsg = message.NewSendingMessage().Append(message.NewText("目前没有正在进行的投票"))
			return nil
		}
		var data *voteData
		err := json.Unmarshal(oldValue, &data)
		if err != nil {
			logger.WithError(err).Errorln("unmarshal json failed")
			return nil
		}
		data.Forbidden = append(data.Forbidden, content)
		if data.Cache != nil {
			newCache := make(map[int64]string)
			for k, v := range data.Cache {
				if !strings.Contains(v, content) {
					newCache[k] = v
				}
			}
			data.Cache = newCache
		}
		newValue, err := json.Marshal(data)
		if err != nil {
			logger.WithError(err).Errorln("marshal json failed")
			return nil
		}
		groupMsg = message.NewSendingMessage().Append(message.NewText("投票屏蔽词增加成功，现在的屏蔽词有：\n" + strings.Join(data.Forbidden, "\n")))
		return newValue
	})
	return
}

type delVoteForbiddenWords struct{}

func (d *delVoteForbiddenWords) Name() string {
	return "删除投票屏蔽词"
}

func (d *delVoteForbiddenWords) ShowTips(int64, int64) string {
	return ""
}

func (d *delVoteForbiddenWords) CheckAuth(_ int64, senderId int64) bool {
	return perm.IsAdmin(senderId)
}

func (d *delVoteForbiddenWords) Execute(_ *message.GroupMessage, content string) (groupMsg *message.SendingMessage, privateMsg *message.SendingMessage) {
	if len(content) == 0 {
		groupMsg = message.NewSendingMessage().Append(message.NewText("指令格式如下：\n删除投票屏蔽词 词语"))
		return
	}
	db.Update([]byte("vote"), func(oldValue []byte) []byte {
		if oldValue == nil {
			groupMsg = message.NewSendingMessage().Append(message.NewText("目前没有正在进行的投票"))
			return nil
		}
		var data *voteData
		err := json.Unmarshal(oldValue, &data)
		if err != nil {
			logger.WithError(err).Errorln("unmarshal json failed")
			return nil
		}
		var ok bool
		var newForbidden []string
		for _, forbidden := range data.Forbidden {
			if forbidden == content {
				ok = true
			} else {
				newForbidden = append(newForbidden, forbidden)
			}
		}
		data.Forbidden = newForbidden
		newValue, err := json.Marshal(data)
		if err != nil {
			logger.WithError(err).Errorln("marshal json failed")
			return nil
		}
		if ok {
			groupMsg = message.NewSendingMessage().Append(message.NewText("投票屏蔽词删除成功，现在的屏蔽词有：\n" + strings.Join(data.Forbidden, "\n")))
		} else {
			groupMsg = message.NewSendingMessage().Append(message.NewText("没有这个屏蔽词"))
		}
		return newValue
	})
	return
}
