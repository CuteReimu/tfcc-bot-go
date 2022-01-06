package commandHandler

import (
	"strings"
	"sync"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
)

func init() {
	bot.RegisterModule(instance)
}

var instance = &mh{}
var logger = utils.GetModuleLogger("tfcc-bot-go.cmdHandler")

// 这是聊天指令处理器的接口，当你想要新增自己的聊天指令处理器时，继承这个类即可。最后，不要忘记在init里调用register
type cmdHandler interface {
	Name() string
	ShowTips(groupCode int64, senderId int64) string
	CheckAuth(groupCode int64, senderId int64) bool
	Execute(msg *message.GroupMessage, content string) (groupMsg *message.SendingMessage, privateMsg *message.SendingMessage)
}

var handlers = make(map[string]cmdHandler)

func register(handler cmdHandler) {
	handlers[handler.Name()] = handler
}

type mh struct {
}

func (m *mh) MiraiGoModule() bot.ModuleInfo {
	return bot.ModuleInfo{
		ID:       "tfcc-bot-go.cmdHandler",
		Instance: instance,
	}
}

func (m *mh) Init() {
}

func (m *mh) PostInit() {
}

func (m *mh) Serve(b *bot.Bot) {
	b.OnGroupMessage(func(c *client.QQClient, msg *message.GroupMessage) {
		var isAt bool
		elem := msg.Elements
		if len(elem) > 0 {
			if at, ok := elem[0].(*message.AtElement); ok && at.Target == b.Uin {
				elem = elem[1:]
				isAt = true
			}
		}
		var cmd, content string
		if len(elem) > 1 {
			return
		}
		if len(elem) == 1 {
			if text, ok := elem[0].(*message.TextElement); ok {
				arr := strings.SplitN(strings.TrimSpace(text.Content), " ", 2)
				cmd = strings.TrimSpace(arr[0])
				if len(arr) > 1 {
					content = strings.TrimSpace(arr[1])
				}
			} else {
				return
			}
		}
		if len(cmd) == 0 {
			if isAt {
				tips(c, msg)
			}
			return
		}
		if handler, ok := handlers[cmd]; ok {
			if handler.CheckAuth(msg.GroupCode, msg.Sender.Uin) {
				if len(content) > 0 {
					logger.WithField("SenderID", msg.Sender.Uin).Info(cmd, " ", content)
				} else {
					logger.WithField("SenderID", msg.Sender.Uin).Info(cmd)
				}
				groupMsg, privateMsg := handler.Execute(msg, content)
				if groupMsg != nil {
					retGroupMsg := c.SendGroupMessage(msg.GroupCode, groupMsg)
					if retGroupMsg.Id == -1 {
						logger.Info("消息被风控了")
					}
				}
				if privateMsg != nil {
					c.SendPrivateMessage(msg.Sender.Uin, privateMsg)
				}
			}
		}
	})
}

func (m *mh) Start(*bot.Bot) {
}

func (m *mh) Stop(_ *bot.Bot, wg *sync.WaitGroup) {
	defer wg.Done()
}

func tips(c *client.QQClient, msg *message.GroupMessage) {
	var ret []string
	for _, handler := range handlers {
		if handler.CheckAuth(msg.GroupCode, msg.Sender.Uin) {
			tip := handler.ShowTips(msg.GroupCode, msg.Sender.Uin)
			ret = append(ret, tip)
		}
	}
	c.SendGroupMessage(msg.GroupCode, message.NewSendingMessage().Append(message.NewText("你可以使用以下功能：\n"+strings.Join(ret, "\n"))))
}
