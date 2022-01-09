package tfcc

import "strings"

var translateMap = map[string]string{
	"Reimu":   "灵梦",
	"Marisa":  "魔理沙",
	"Sakuya":  "咲夜",
	"Sanae":   "早苗",
	"Youmu":   "妖梦",
	"RY":      "结界组",
	"MA":      "咏唱组",
	"SR":      "红魔组",
	"YY":      "幽冥组",
	"Yukari":  "紫",
	"Alice":   "爱丽丝",
	"Remilia": "蕾米莉亚",
	"Yuyuko":  "幽幽子",
	"Reisen":  "铃仙",
	"Cirno":   "琪露诺",
	"Aya":     "射命丸文",
	"Spring":  "(春)",
	"Summer":  "(夏)",
	"Autumn":  "(秋)",
	"Winter":  "(冬)",
	"Wolf":    "(狼)",
	"Otter":   "(獭)",
	"Eagle":   "(鹰)",
	"Border":  "结界",
	"Magic":   "咏唱",
	"Scarlet": "红魔",
	"Ghost":   "幽冥",
	"Team":    "组",
	"FinalA":  "(6A)",
	"FinalB":  "(6B)",
}

func Translate(s string) string {
	if s2, ok := translateMap[s]; ok {
		return s2
	}
	arr := []rune(s)
	for i := 0; i < len(arr); i++ {
		n := len(arr) - 1
		if n > 10 {
			n = 7
		}
		for ; n >= 1; n-- {
			if i+n > len(arr) {
				continue
			}
			key := string(arr[i : i+n])
			if val, ok := translateMap[key]; ok {
				return Translate(strings.Replace(string(arr[:i+n]), key, val, 1) + string(arr[i+n:]))
			}
		}
	}
	return s
}
