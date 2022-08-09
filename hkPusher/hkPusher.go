package hkPusher

import (
	"encoding/json"
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/config"
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/Touhou-Freshman-Camp/tfcc-bot-go/db"
	"github.com/Touhou-Freshman-Camp/tfcc-bot-go/translate"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

func init() {
	bot.RegisterModule(instance)
}

var instance = &mh{}
var logger = utils.GetModuleLogger("tfcc-bot-go.hkPusher")

type mh struct {
}

func (m *mh) MiraiGoModule() bot.ModuleInfo {
	return bot.ModuleInfo{
		ID:       "tfcc-bot-go.hkPusher",
		Instance: instance,
	}
}

func (m *mh) Init() {
}

func (m *mh) PostInit() {
}

func (m *mh) Serve(b *bot.Bot) {
	re := regexp.MustCompile("<.*?>")
	delay := config.GlobalConfig.GetInt64("schedule.speedrun_push_delay")
	if delay <= 0 {
		return
	}
	qqGroups := config.GlobalConfig.GetIntSlice("schedule.qq_group")
	if len(qqGroups) == 0 {
		return
	}
	apiKey := config.GlobalConfig.GetString("schedule.speedrun_api_key")
	if len(apiKey) == 0 {
		return
	}
	go func() {
		ticker := time.NewTicker(time.Duration(delay) * time.Second)
		defer ticker.Stop()
		for {
			<-ticker.C
			resp, err := resty.New().SetTimeout(time.Second * 20).R().SetHeaders(map[string]string{
				"Accept":    "application/json",
				"X-API-Key": apiKey,
			}).Get("https://www.speedrun.com/api/v1/notifications")
			if err != nil {
				logger.WithError(err).Error("cannot access speedrun.com")
				continue
			}
			if resp.StatusCode() != 200 {
				logger.Error("speedrun.com return code: ", resp.StatusCode())
				continue
			}
			buf := resp.Body()
			if !gjson.ValidBytes(buf) {
				logger.Error("speedrun.com return invalid json: ", string(buf))
				continue
			}
			result := gjson.ParseBytes(buf)
			value := db.Get([]byte("pushed_messages"))
			if value == nil {
				value = []byte("{}")
			}
			var pushedMessages map[string]int
			err = json.Unmarshal(value, &pushedMessages)
			if err != nil {
				logger.WithError(err).Error("json unmarshal failed")
				continue
			}
			var arr []string
			for _, r := range result.Get("data").Array() {
				id := r.Get("id").String()
				if _, ok := pushedMessages[id]; !ok {
					pushedMessages[id] = 1
					s := re.ReplaceAllString(r.Get("text").String(), "")
					if strings.Contains(s, "beat the WR") || strings.Contains(s, "got a new top 3 PB") {
						arr = append(arr, translate.Translate(s))
					}
				}
			}
			value, err = json.Marshal(&pushedMessages)
			if err != nil {
				logger.WithError(err).Error("json unmarshal failed")
				continue
			}
			db.Set([]byte("pushed_messages"), value)
			if len(arr) < 0 {
				continue
			}
			for _, qqGroup := range qqGroups {
				groupCode := int64(qqGroup)
				key := []byte("unsend:" + strconv.Itoa(qqGroup))
				value = db.Get(key)
				var oldArr []string
				if len(value) > 0 {
					oldArr = append(oldArr, string(value))
				}
				str := strings.Join(append(oldArr, arr...), "\r\n")
				if len(str) == 0 {
					continue
				}
				groupMsg := message.NewSendingMessage().Append(message.NewText(str))
				retGroupMsg := b.SendGroupMessage(groupCode, groupMsg)
				if retGroupMsg == nil {
					logger.Info("群聊消息发送失败了")
				} else if retGroupMsg.Id == -1 {
					logger.Info("群聊消息被风控了")
				} else {
					db.Del(key)
					continue
				}
				db.Set(key, []byte(str))
			}
		}
	}()
}

func (m *mh) Start(*bot.Bot) {
}

func (m *mh) Stop(_ *bot.Bot, wg *sync.WaitGroup) {
	defer wg.Done()
}
