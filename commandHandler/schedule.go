package commandHandler

import (
	"encoding/json"
	"fmt"
	"github.com/CuteReimu/tfcc-bot-go/bot"
	"github.com/CuteReimu/tfcc-bot-go/config"
	"github.com/CuteReimu/tfcc-bot-go/db"
	"github.com/CuteReimu/tfcc-bot-go/perm"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/araddon/dateparse"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func initSchedule(b *bot.Bot) {
	go func() {
		ticker := time.NewTicker(time.Second * 30)
		defer ticker.Stop()
		for {
			<-ticker.C
			now := time.Now().Unix()
			fatalErr := false
			var tiggeredData []*scheduleData
			db.Update([]byte("schedule"), func(oldValue []byte) []byte {
				if oldValue == nil {
					return nil
				}
				var data []*scheduleData
				err := json.Unmarshal(oldValue, &data)
				if err != nil {
					fatalErr = true
					logger.WithError(err).Errorln("unmarshal json failed")
					return nil
				}
				var updated bool
				// 删除过期的事件
				var newData []*scheduleData
				for _, d := range data {
					if now < d.EndTime {
						newData = append(newData, d)
					} else {
						updated = true
					}
				}
				data = newData
				// 判断有哪些事件该提醒了
				for _, d := range data {
					var newNotifyTime []int64
					var triggered bool
					for _, t := range d.NotifyTime {
						if now >= t {
							triggered = true
						} else {
							newNotifyTime = append(newNotifyTime, t)
						}
					}
					if triggered {
						updated = true
						d.NotifyTime = newNotifyTime
						tiggeredData = append(tiggeredData, d)
					}
				}
				// 更新数据库里的值
				if updated {
					buf, err := json.Marshal(data)
					if err != nil {
						logger.WithError(err).Errorln("json marshal failed")
						return nil
					}
					return buf
				}
				return nil
			})
			if len(tiggeredData) > 0 {
				text := "温馨提醒："
				for _, d := range tiggeredData {
					text += fmt.Sprintf("\n%s 将于%s开始", d.Tips, time.Unix(d.EndTime, 0).Format("2006/01/02 15:04:05"))
				}
				for _, groupCode := range config.GlobalConfig.GetIntSlice("schedule.qq_group") {
					b.SendGroupMessage(int64(groupCode), message.NewSendingMessage().Append(message.NewText(text)))
				}
			}
			if fatalErr {
				logger.Error("定时任务功能出现严重异常，已停止")
				return
			}
		}
	}()
}

func init() {
	register(newAddSchedule())
	register(&delSchedule{})
	register(&listAllSchedule{})
}

type scheduleData struct {
	NotifyTime []int64 `json:"notify_time,omitempty"`
	EndTime    int64   `json:"end_time,omitempty"`
	Tips       string  `json:"tips,omitempty"`
}

type addSchedule struct {
	reg *regexp.Regexp
}

func newAddSchedule() *addSchedule {
	return &addSchedule{reg: regexp.MustCompile(`[^0-9:\\/-]`)}
}

func (a *addSchedule) Name() string {
	return "增加预约"
}

func (a *addSchedule) ShowTips(int64, int64) string {
	return "增加预约"
}

func (a *addSchedule) CheckAuth(_ int64, senderId int64) bool {
	return perm.IsWhitelist(senderId)
}

func (a *addSchedule) Execute(_ *message.GroupMessage, content string) (groupMsg *message.SendingMessage, privateMsg *message.SendingMessage) {
	if len(content) == 0 {
		groupMsg = message.NewSendingMessage().Append(message.NewText("指令格式：增加预约 年月日时分秒 预约文字\n例如：增加预约 2020-12-25 12:23:00 风神录L避弹\n（时间可以不用分隔符）"))
		return
	}
	arr := strings.Split(content, " ")
	i := 0
	for _, s := range arr {
		if a.reg.MatchString(s) {
			break
		}
		i++
	}
	timeStr := strings.TrimSpace(strings.Join(arr[:i], " "))
	contentStr := strings.TrimSpace(strings.Join(arr[i:], " "))
	if len(contentStr) == 0 {
		groupMsg = message.NewSendingMessage().Append(message.NewText("指令格式：增加预约 年月日时分秒 预约文字\n例如：增加预约 2020-12-25 12:23:00 风神录L避弹\n（时间可以不用分隔符）"))
		return
	}
	t, err := dateparse.ParseLocal(timeStr)
	if err != nil {
		groupMsg = message.NewSendingMessage().Append(message.NewText("日期或时间格式错误"))
		return
	}
	if !time.Now().Add(time.Minute).Before(t) {
		groupMsg = message.NewSendingMessage().Append(message.NewText("请预约一个将来的时间"))
		return
	}
	var success bool
	db.Update([]byte("schedule"), func(oldValue []byte) []byte {
		var data []*scheduleData
		if oldValue != nil {
			err := json.Unmarshal(oldValue, &data)
			if err != nil {
				logger.WithError(err).Errorln("unmarshal json failed")
				return nil
			}
		}
		endTime := t.Unix()
		newScheduleData := &scheduleData{EndTime: endTime, Tips: contentStr}
		for _, before := range config.GlobalConfig.GetIntSlice("schedule.before") {
			newScheduleData.NotifyTime = append(newScheduleData.NotifyTime, endTime-int64(before))
		}
		data = append(data, newScheduleData)
		newValue, err := json.Marshal(data)
		if err != nil {
			logger.WithError(err).Errorln("marshal json failed")
			return nil
		}
		success = true
		return newValue
	})
	if success {
		groupMsg = message.NewSendingMessage().Append(message.NewText(fmt.Sprintf("已于%s增加预约：%s", t.Format("2006/01/02 15:04:05"), contentStr)))
	}
	return
}

type delSchedule struct{}

func (d *delSchedule) Name() string {
	return "删除预约"
}

func (d *delSchedule) ShowTips(int64, int64) string {
	buf := db.Get([]byte("schedule"))
	if buf == nil {
		return ""
	}
	var data []*scheduleData
	err := json.Unmarshal(buf, &data)
	if err != nil {
		logger.WithError(err).Errorln("unmarshal json failed")
		return ""
	}
	if len(data) == 0 {
		return ""
	}
	return "删除预约 序号"
}

func (d *delSchedule) CheckAuth(_ int64, senderId int64) bool {
	return perm.IsWhitelist(senderId)
}

func (d *delSchedule) Execute(_ *message.GroupMessage, content string) (groupMsg *message.SendingMessage, privateMsg *message.SendingMessage) {
	if len(content) == 0 {
		groupMsg = message.NewSendingMessage().Append(message.NewText("指令格式如下：\\n删除预约 序号（请先用“预约列表”查询序号）"))
		return
	}
	i, err := strconv.Atoi(content)
	if err != nil {
		return
	}
	i--
	var success bool
	db.Update([]byte("schedule"), func(oldValue []byte) []byte {
		var data []*scheduleData
		if oldValue != nil {
			err := json.Unmarshal(oldValue, &data)
			if err != nil {
				logger.WithError(err).Errorln("unmarshal json failed")
				return nil
			}
		}
		if i < 0 || i >= len(data) {
			return nil
		}
		newValue, err := json.Marshal(append(data[:i], data[i+1:]...))
		if err != nil {
			logger.WithError(err).Errorln("marshal json failed")
			return nil
		}
		success = true
		return newValue
	})
	if success {
		groupMsg = message.NewSendingMessage().Append(message.NewText("删除预约成功"))
	} else {
		groupMsg = message.NewSendingMessage().Append(message.NewText("找不到这条预约，请再次确认序号是否正确"))
	}
	return
}

type listAllSchedule struct{}

func (l *listAllSchedule) Name() string {
	return "预约列表"
}

func (l *listAllSchedule) ShowTips(int64, int64) string {
	buf := db.Get([]byte("schedule"))
	if buf == nil {
		return ""
	}
	var data []*scheduleData
	err := json.Unmarshal(buf, &data)
	if err != nil {
		logger.WithError(err).Errorln("unmarshal json failed")
		return ""
	}
	if len(data) == 0 {
		return ""
	}
	return "预约列表 行数"
}

func (l *listAllSchedule) CheckAuth(int64, int64) bool {
	return true
}

func (l *listAllSchedule) Execute(_ *message.GroupMessage, content string) (groupMsg *message.SendingMessage, privateMsg *message.SendingMessage) {
	count := 5
	if len(content) > 0 {
		var err error
		count, err = strconv.Atoi(content)
		if err != nil {
			return
		}
	}
	buf := db.Get([]byte("schedule"))
	if buf != nil {
		var data []*scheduleData
		err := json.Unmarshal(buf, &data)
		if err != nil {
			logger.WithError(err).Errorln("unmarshal json failed")
			return
		}
		if len(data) > 0 {
			var text []string
			for i, d := range data {
				if i >= count {
					break
				}
				text = append(text, fmt.Sprintf("%d %s %s", i+1, time.Unix(d.EndTime, 0).Format("2006/01/02 15:04:05"), d.Tips))
			}
			groupMsg = message.NewSendingMessage().Append(message.NewText(strings.Join(text, "\n")))
			return
		}
	}
	groupMsg = message.NewSendingMessage().Append(message.NewText("目前没有预约"))
	return
}
