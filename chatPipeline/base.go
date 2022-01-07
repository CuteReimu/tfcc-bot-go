package chatPipeline

import (
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"sync"
)

func init() {
	bot.RegisterModule(instance)
}

var instance = &mh{}
var logger = utils.GetModuleLogger("tfcc-bot-go.chatPipeline")

type pipelineHandler interface {
	Execute(msg *message.GroupMessage, content string) (groupMsg *message.SendingMessage)
}

var handlers []pipelineHandler

func register(handler pipelineHandler) {
	handlers = append(handlers, handler)
}

type mh struct {
}

func (m *mh) MiraiGoModule() bot.ModuleInfo {
	return bot.ModuleInfo{
		ID:       "tfcc-bot-go.chatPipeline",
		Instance: instance,
	}
}

func (m *mh) Init() {
}

func (m *mh) PostInit() {
}

func (m *mh) Serve(b *bot.Bot) {
	b.OnGroupMessage(func(c *client.QQClient, msg *message.GroupMessage) {
		elem := msg.Elements
		if len(elem) != 1 {
			return
		}
		if text, ok := elem[0].(*message.TextElement); ok {
			for _, handler := range handlers {
				groupMsg := handler.Execute(msg, text.Content)
				if groupMsg != nil {
					retGroupMsg := c.SendGroupMessage(msg.GroupCode, groupMsg)
					if retGroupMsg.Id == -1 {
						logger.Info("群聊消息被风控了")
					}
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
