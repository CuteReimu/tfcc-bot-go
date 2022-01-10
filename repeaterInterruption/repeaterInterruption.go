package repeaterInterruption

import (
	//"fmt"
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/config"
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"strings"
	"sync"
	"time"
)

func init() {
	bot.RegisterModule(instance)
}

var instance = &mh{data: make(map[int64]*repeaterData)}
var logger = utils.GetModuleLogger("tfcc-bot-go.repeaterInterruption")

type repeaterData struct {
	sync.Mutex
	lastMessage string
	counter     int
	lastTrigger time.Time
}

type mh struct {
	data map[int64]*repeaterData
}

func (m *mh) MiraiGoModule() bot.ModuleInfo {
	return bot.ModuleInfo{
		ID:       "tfcc-bot-go.repeaterInterruption",
		Instance: instance,
	}
}

func (m *mh) Init() {
	for _, groupCode := range config.GlobalConfig.GetIntSlice("repeater_interruption.qq_group") {
		m.data[int64(groupCode)] = &repeaterData{}
	}
}

func (m *mh) PostInit() {
}

func (m *mh) Serve(b *bot.Bot) {
	b.OnGroupMessage(func(c *client.QQClient, msg *message.GroupMessage) {
		data, ok := m.data[msg.GroupCode]
		if !ok {
			return
		}
		data.Lock()
		defer data.Unlock()
		m := msg.ToString()
		if m != data.lastMessage {
			data.counter = 1
			data.lastMessage = m
		} else {
			data.counter++
			if data.counter >= config.GlobalConfig.GetInt("repeater_interruption.allowance") {
				now := time.Now()
				coolDown := config.GlobalConfig.GetInt64("repeater_interruption.cool_down")
				if now.After(data.lastTrigger.Add(time.Duration(coolDown) * time.Second)) {
					text := "打断复读~~ (^-^)"
					if strings.Contains(data.lastMessage, text) {
						text = `(*/ω\*)`
					}
					data.counter = 1
					data.lastTrigger = now
					go func() {
						groupMsg := message.NewSendingMessage().Append(message.NewText(text))
						retGroupMsg := c.SendGroupMessage(msg.GroupCode, groupMsg)
						if retGroupMsg.Id == -1 {
							logger.Info("群聊消息被风控了")
						}
					}()
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

// func parseMessage(elements []message.IMessageElement) string {
// 	var s []string
// 	for _, e := range elements {
// 		if text, ok := e.(*message.TextElement); ok {
// 			content := strings.TrimSpace(text.Content)
// 			s = append(s, content)
// 		} else if img, ok := e.(*message.GroupImageElement); ok {
// 			s = append(s, fmt.Sprintf("[pic={%s}]", img.ImageId))
// 		} else if at, ok := e.(*message.AtElement); ok {
// 			s = append(s, fmt.Sprintf("[<@%d>]", at.Target))
// 		}
// 	}
// 	return strings.Join(s, "")
// }
