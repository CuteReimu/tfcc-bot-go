package repAnalyze

import (
	"github.com/CuteReimu/threp"
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/Touhou-Freshman-Camp/tfcc-bot-go/repeaterInterruption"
	"github.com/go-resty/resty/v2"
	"strings"
	"sync"
	"time"
)

func init() {
	bot.RegisterModule(instance)
}

var instance = &mh{}
var logger = utils.GetModuleLogger("tfcc-bot-go.repAnalyze")

type mh struct {
}

func (m *mh) MiraiGoModule() bot.ModuleInfo {
	return bot.ModuleInfo{
		ID:       "tfcc-bot-go.repAnalyze",
		Instance: instance,
	}
}

func (m *mh) Init() {
}

func (m *mh) PostInit() {
}

func (m *mh) Serve(b *bot.Bot) {
	b.GroupMessageEvent.Subscribe(func(c *client.QQClient, msg *message.GroupMessage) {
		for _, elem := range msg.Elements {
			if e, ok := elem.(*message.GroupFileElement); ok {
				if strings.HasSuffix(e.Name, ".rpy") {
					url := c.GetGroupFileUrl(msg.GroupCode, e.Path, e.Busid)
					fileInfo := fetchRepFileInfo(url)
					if len(fileInfo) > 0 {
						retGroupMsg := c.SendGroupMessage(msg.GroupCode, message.NewSendingMessage().Append(message.NewText(fileInfo)))
						if retGroupMsg == nil {
							logger.Info("群聊消息发送失败了")
						} else if retGroupMsg.Id == -1 {
							logger.Info("群聊消息被风控了")
						} else {
							repeaterInterruption.Clean(msg.GroupCode)
						}
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

func fetchRepFileInfo(url string) string {
	resp, err := resty.New().SetTimeout(20 * time.Second).SetLogger(logger).R().SetDoNotParseResponse(true).Get(url)
	if err != nil {
		logger.WithError(err).Errorln("fetch file failed")
		return ""
	}
	body := resp.RawBody()
	defer func() {
		if err := body.Close(); err != nil {
			logger.WithError(err).Errorln("close body failed")
		}
	}()
	rep, err := threp.DecodeReplay(body)
	if err != nil {
		logger.WithError(err).Errorln("decode replay failed")
		return ""
	}
	return rep.String()
}
