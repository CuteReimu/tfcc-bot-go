package commandHandler

import (
	"github.com/Mrs4s/MiraiGo/message"
	"math/rand"
)

func init() {
	register(newRandGame())
	register(newRandCharacter())
}

type randGame struct {
	games []string
}

func newRandGame() cmdHandler {
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

func newRandCharacter() cmdHandler {
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
		r.gameMap[string([]rune(k)[2:3])] = r.gameMap[k]
		r.gameMap[string([]rune(k)[4:5])] = r.gameMap[k]
		r.gameMap[string([]rune(k)[2:5])] = r.gameMap[k]
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
