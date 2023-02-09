package commandHandler

import (
	"bufio"
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"github.com/Logiase/MiraiGo-Template/config"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/Touhou-Freshman-Camp/tfcc-bot-go/db"
	"github.com/ozgio/strutil"
	"io"
	"io/fs"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
)

func init() {
	register(newRandGame())
	register(newRandCharacter())
	register(newRandSpell())
}

type randGame struct {
	games []string
}

func newRandGame() *randGame {
	return &randGame{
		games: []string{"东方红魔乡", "东方妖妖梦", "东方永夜抄", "东方风神录", "东方地灵殿", "东方星莲船", "东方神灵庙", "东方辉针城", "东方绀珠传", "东方天空璋", "东方鬼形兽", "东方虹龙洞"},
	}
}

func (r *randGame) Name() string {
	return "随作品"
}

func (r *randGame) ShowTips(int64, int64) string {
	return "随作品"
}

func (r *randGame) CheckAuth(int64, int64) bool {
	return true
}

func (r *randGame) Execute(*message.GroupMessage, string) (groupMsg *message.SendingMessage, privateMsg *message.SendingMessage) {
	n := rand.Intn(len(r.games))
	groupMsg = message.NewSendingMessage().Append(message.NewText(r.games[n]))
	return
}

type randCharacter struct {
	gameMap map[string][][]string
}

func newRandCharacter() *randCharacter {
	r := &randCharacter{
		gameMap: map[string][][]string{
			"东方红魔乡": {{"灵梦", "魔理沙"}, {"A", "B"}},
			"东方妖妖梦": {{"灵梦", "魔理沙", "咲夜"}, {"A", "B"}},
			"东方永夜抄": {{"结界组", "咏唱组", "红魔组", "幽冥组", "灵梦", "紫", "魔理沙", "爱丽丝", "咲夜", "蕾米莉亚", "妖梦", "幽幽子"}},
			"东方风神录": {{"灵梦", "魔理沙"}, {"A", "B", "C"}},
			"东方地灵殿": {{"灵梦", "魔理沙"}, {"A", "B", "C"}},
			"东方星莲船": {{"灵梦", "魔理沙", "早苗"}, {"A", "B"}},
			"东方神灵庙": {{"灵梦", "魔理沙", "早苗", "妖梦"}},
			"东方辉针城": {{"灵梦", "魔理沙", "咲夜"}, {"A", "B"}},
			"东方绀珠传": {{"灵梦", "魔理沙", "早苗", "铃仙"}},
			"东方天空璋": {{"灵梦", "琪露诺", "射命丸文", "魔理沙"}, {"（春）", "（夏）", "（秋）", "（冬）"}},
			"东方鬼形兽": {{"灵梦", "魔理沙", "妖梦"}, {"（狼）", "（獭）", "（鹰）"}},
			"东方虹龙洞": {{"灵梦", "魔理沙", "早苗", "咲夜"}},
		},
	}
	games := make([]string, 0, len(r.gameMap))
	for k := range r.gameMap {
		games = append(games, k)
	}
	for _, k := range games {
		r.gameMap[strutil.MustSubstring(k, 2, 3)] = r.gameMap[k]
		r.gameMap[strutil.MustSubstring(k, 4, 5)] = r.gameMap[k]
		r.gameMap[strutil.MustSubstring(k, 2, 5)] = r.gameMap[k]
	}
	return r
}

func (r *randCharacter) Name() string {
	return "随机体"
}

func (r *randCharacter) ShowTips(int64, int64) string {
	return "随机体"
}

func (r *randCharacter) CheckAuth(int64, int64) bool {
	return true
}

func (r *randCharacter) Execute(_ *message.GroupMessage, content string) (groupMsg *message.SendingMessage, privateMsg *message.SendingMessage) {
	if len(content) == 0 {
		groupMsg = message.NewSendingMessage().Append(message.NewText(`请输入要随机的作品，例如：“随机体 红”`))
		return
	}
	if val, ok := r.gameMap[content]; ok {
		var ret string
		for _, v := range val {
			ret += v[rand.Intn(len(v))]
		}
		groupMsg = message.NewSendingMessage().Append(message.NewText(ret))
	}
	return
}

type randSpellData struct {
	LastRandTime int64
	Count        int64
}

type spells struct {
	sync.Mutex
	spells []string
}

func (s *spells) randN(count int) []string {
	s.Lock()
	defer s.Unlock()
	var text []string
	for i := 0; i < count; i++ {
		index := i + rand.Intn(len(s.spells)-i)
		if i != index {
			s.spells[i], s.spells[index] = s.spells[index], s.spells[i]
		}
		text = append(text, s.spells[i])
	}
	return text
}

type randSpell struct {
	gameMap map[string]*spells
}

//go:embed spells
var spellFs embed.FS

func newRandSpell() *randSpell {
	r := &randSpell{gameMap: make(map[string]*spells)}
	files, err := spellFs.ReadDir("spells")
	if err != nil {
		logger.WithError(err).Errorln("init spells failed")
	}
	var count int
	for _, file := range files {
		name := file.Name()
		if strings.HasSuffix(name, ".txt") {
			err = r.loadSpells("spells/" + name)
			if err != nil {
				logger.WithError(err).Errorln("load file failed: " + name)
			} else {
				count++
			}
		}
	}
	logger.Infof("load %d spell files successful\n", count)
	return r
}

func (r *randSpell) loadSpells(name string) error {
	f, err := spellFs.Open(name)
	if err != nil {
		return err
	}
	defer func(f fs.File) { _ = f.Close() }(f)
	reader := bufio.NewReader(f)
	arr := strings.Split(name[len("spells/"):len(name)-len(".txt")], " ")
	for _, s := range arr {
		r.gameMap[s] = &spells{}
	}
	for {
		line, _, err := reader.ReadLine() // 不太可能出现太长的行数，所以 isPrefix 参数可以忽略
		if err != nil && err != io.EOF {
			return err
		}
		line = bytes.TrimSpace(line)
		if len(line) > 0 {
			for _, s := range arr {
				r.gameMap[s].spells = append(r.gameMap[s].spells, string(line))
			}
		}
		if err == io.EOF {
			break
		}
	}
	return nil
}

func (r *randSpell) Name() string {
	return "随符卡"
}

func (r *randSpell) ShowTips(int64, int64) string {
	return "随符卡"
}

func (r *randSpell) CheckAuth(int64, int64) bool {
	return true
}

func (r *randSpell) Execute(msg *message.GroupMessage, content string) (groupMsg *message.SendingMessage, privateMsg *message.SendingMessage) {
	oneTimeLimit := config.GlobalConfig.GetInt("qq.rand_one_time_limit")
	if len(content) == 0 {
		groupMsg = message.NewSendingMessage().Append(message.NewText(fmt.Sprintf(`请输入要随机的作品与符卡数量，例如：“随符卡 红”或“随符卡 全部 %d”`, oneTimeLimit)))
		return
	}
	cmds := strings.Split(content, " ")
	content = cmds[0]
	var count int
	if len(cmds) <= 1 {
		count = 1 // 默认抽取一张符卡
	} else {
		countStr := cmds[1]
		var err error
		count, err = strconv.Atoi(countStr)
		if err != nil || count == 0 || count > oneTimeLimit {
			groupMsg = message.NewSendingMessage().Append(message.NewText(fmt.Sprintf(`请输入%d以内数字，例如：“随符卡 红 %d”或“随符卡 全部 %d”`, oneTimeLimit, oneTimeLimit, oneTimeLimit)))
			return
		}
	}
	if val, ok := r.gameMap[content]; ok {
		if count > len(val.spells) {
			groupMsg = message.NewSendingMessage().Append(message.NewText(fmt.Sprintf(`请输入小于或等于该作符卡数量%d的数字`, len(val.spells))))
			return
		}
		db.UpdateWithTtl([]byte("rand_spell:"+strconv.FormatInt(msg.Sender.Uin, 10)), func(oldValue []byte) ([]byte, time.Duration) {
			var d *randSpellData
			if oldValue == nil {
				d = &randSpellData{}
			} else {
				err := json.Unmarshal(oldValue, &d)
				if err != nil {
					logger.WithError(err).Errorln("unmarshal json failed")
					return nil, 0
				}
			}
			now := time.Now()
			yy, mm, dd := now.Date()
			yy2, mm2, dd2 := time.Unix(d.LastRandTime, 0).Date()
			if !(yy == yy2 && mm == mm2 && dd == dd2) {
				d.Count = 0
			}
			d.Count++
			limitCount := config.GlobalConfig.GetInt64("qq.rand_count")
			if d.Count <= limitCount {
				text := val.randN(count)
				groupMsg = message.NewSendingMessage().Append(message.NewText(strings.Join(text, "\n")))
			} else if d.Count == limitCount+1 {
				relatedUrl := config.GlobalConfig.GetString("qq.related_url")
				s := fmt.Sprintf("随符卡一天只能使用%d次", limitCount)
				if len(relatedUrl) > 0 {
					s += "\n你可以前往 " + relatedUrl + "继续使用"
				}
				groupMsg = message.NewSendingMessage().Append(message.NewText(s))
			}
			d.LastRandTime = now.Unix()
			newValue, err := json.Marshal(d)
			if err != nil {
				logger.WithError(err).Errorln("unmarshal json failed")
				return nil, 0
			}
			return newValue, time.Hour * 24
		})
	}
	return
}
