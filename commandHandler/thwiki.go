package commandHandler

import (
	"encoding/json"
	"fmt"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/go-resty/resty/v2"
	"sort"
	"strings"
	"sync"
	"time"
)

func init() {
	register(newGetThwikiEvent())
}

type thwikiEvent struct {
	Results thwikiEventResultList `json:"results"`
	Version string                `json:"version"`
	Meta    struct {
		Hash   string `json:"hash"`
		Count  int    `json:"count"`
		Offset int    `json:"offset"`
		Source string `json:"source"`
		Time   string `json:"time"`
	} `json:"meta"`
}

func (t *thwikiEvent) String() string {
	var result []string
	for _, r := range t.Results {
		var res string
		if len(r.Type) != 0 {
			res = fmt.Sprintf("%s【%s】%s", r.StartStr, r.Type[0], r.Desc)
		} else {
			res = fmt.Sprintf("%s %s", r.StartStr, r.Desc)
		}
		result = append(result, res)
	}
	return strings.Join(result, "\n")
}

type thwikiEventResult struct {
	Id       string   `json:"id"`
	Start    int      `json:"start"`
	End      int      `json:"end"`
	StartStr string   `json:"startStr"`
	EndStr   string   `json:"endStr"`
	Title    string   `json:"title"`
	Desc     string   `json:"desc"`
	Url      string   `json:"url"`
	Icon     string   `json:"icon"`
	Type     []string `json:"type"`
	Color    string   `json:"color"`
}

type thwikiEventResultList []*thwikiEventResult

func (t thwikiEventResultList) Len() int {
	return len(t)
}

func (t thwikiEventResultList) Less(i, j int) bool {
	return t[i].Start < t[j].End
}

func (t thwikiEventResultList) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

type getThwikiEvent struct {
	sync.Mutex
	event           *thwikiEvent
	lastFetchTime   time.Time
	lastRequestTime map[int64]time.Time
}

func newGetThwikiEvent() *getThwikiEvent {
	return &getThwikiEvent{
		lastRequestTime: make(map[int64]time.Time),
	}
}

func (e *getThwikiEvent) Name() string {
	return "看新闻"
}

func (e *getThwikiEvent) ShowTips(int64, int64) string {
	return "看新闻"
}

func (e *getThwikiEvent) CheckAuth(int64, int64) bool {
	return true
}

func (e *getThwikiEvent) Execute(msg *message.GroupMessage, _ string) (groupMsg *message.SendingMessage, privateMsg *message.SendingMessage) {
	e.Lock()
	defer e.Unlock()
	now := time.Now()
	lastRequestTime, ok := e.lastRequestTime[msg.GroupCode]
	if ok && now.Before(lastRequestTime.Add(5*time.Minute)) {
		groupMsg = message.NewSendingMessage().Append(message.NewText("这个功能每5分钟才能使用一次"))
		return
	}
	if now.After(e.lastFetchTime.Add(6 * time.Hour)) {
		e.getEvents(now)
		e.lastFetchTime = now
	}
	if e.event == nil {
		return nil, nil
	}
	e.lastRequestTime[msg.GroupCode] = now
	text := e.event.String()
	if len(text) > 0 {
		groupMsg = message.NewSendingMessage().Append(message.NewText(text))
	}
	return
}

func (e *getThwikiEvent) getEvents(now time.Time) {
	resp, err := resty.New().SetTimeout(20 * time.Second).R().SetQueryParams(map[string]string{
		"start": now.Add(-3 * 24 * time.Hour).Format("2006-01-02"),
		"end":   now.Add(4 * 24 * time.Hour).Format("2006-01-02"),
	}).Get("https://calendar.thwiki.cc/events/")
	if err != nil {
		logger.WithError(err).Error("failed to access thwiki")
		return
	}
	if resp.StatusCode() != 200 {
		logger.WithField("StatusCode", resp.StatusCode()).Error("failed to access thwiki")
		return
	}
	err = json.Unmarshal(resp.Body(), &e.event)
	if err != nil {
		logger.WithError(err).Error("failed to unmarshal json")
		return
	}
	sort.Sort(e.event.Results)
}
