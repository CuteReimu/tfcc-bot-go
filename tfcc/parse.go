package tfcc

import "strings"

type trieNode struct {
	child map[rune]*trieNode
	value func(*ParseResult)
}

func newTrieNode() *trieNode {
	return &trieNode{child: make(map[rune]*trieNode)}
}

type trie struct {
	root *trieNode
}

func (t *trie) putIfAbsent(key string, value func(*ParseResult)) bool {
	if value == nil {
		panic("cannot put a nil value")
	}
	if t.root == nil {
		t.root = newTrieNode()
	}
	node := t.root
	for _, c := range key {
		n, ok := node.child[c]
		if ok {
			node = n
		} else {
			newNode := newTrieNode()
			node.child[c] = newNode
			node = newNode
		}
	}
	if node.value != nil {
		return false
	}
	node.value = value
	return true
}

func (t *trie) get(key string) func(*ParseResult) {
	if t.root == nil {
		return nil
	}
	node := t.root
	for _, c := range key {
		n, ok := node.child[c]
		if ok {
			node = n
		} else {
			return nil
		}
	}
	return node.value
}

type ParseResult struct {
	Work, Rank              string
	Route, Character, CType map[string]struct{}
	AllSpell                bool
}

func newParseResult() *ParseResult {
	return &ParseResult{
		Route:     make(map[string]struct{}),
		Character: make(map[string]struct{}),
		CType:     make(map[string]struct{}),
	}
}

var workDict = &trie{}
var otherDict = &trie{}

func addWorkMap(result string, represent ...string) {
	for _, s := range represent {
		if !workDict.putIfAbsent(strings.ToLower(s), func(res *ParseResult) { parseWork(res, result) }) {
			panic("repeated trie keys: " + s)
		}
	}
}

func addOtherMap(f func(*ParseResult, string), result string, represent ...string) {
	for _, s := range represent {
		if !otherDict.putIfAbsent(strings.ToLower(s), func(res *ParseResult) { f(res, result) }) {
			panic("repeated trie keys: " + s)
		}
	}
}

func parseWork(res *ParseResult, result string) {
	if len(res.Work) == 0 {
		res.Work = result
	}
}

func parseRank(res *ParseResult, result string) {
	if len(res.Rank) == 0 {
		res.Rank = result
	}
}

func parseRoute(res *ParseResult, result string) {
	res.Route[result] = struct{}{}
}

func parseCharacter(res *ParseResult, result string) {
	res.Character[result] = struct{}{}
}

func parseCType(res *ParseResult, result string) {
	res.CType[result] = struct{}{}
}

func parseCharacterCType(res *ParseResult, result string) {
	switch result[len(result)-2:] {
	case "SA":
		res.CType["Spring"] = struct{}{}
		parseCharacter(res, result[:len(result)-2])
	case "SB":
		res.CType["Summer"] = struct{}{}
		parseCharacter(res, result[:len(result)-2])
	case "SC":
		res.CType["Autumn"] = struct{}{}
		parseCharacter(res, result[:len(result)-2])
	case "SD":
		res.CType["Winter"] = struct{}{}
		parseCharacter(res, result[:len(result)-2])
	default:
		switch result[len(result)-1] {
		case 'W':
			res.CType["Wolf"] = struct{}{}
			parseCharacter(res, result[:len(result)-1])
		case 'O':
			res.CType["Otter"] = struct{}{}
			parseCharacter(res, result[:len(result)-1])
		case 'E':
			res.CType["Eagle"] = struct{}{}
			parseCharacter(res, result[:len(result)-1])
		}
	}
}

func init() {
	addWorkMap("6", "红", "红魔乡", "hmx", "th6", "th06", "EoSD")
	addWorkMap("7", "妖", "妖妖梦", "yym", "th7", "th07", "PCB")
	addWorkMap("8", "永", "永夜抄", "yyc", "th8", "th08", "IN")
	addWorkMap("9", "花", "花映冢", "hyz", "th9", "th09", "PoFV")
	addWorkMap("10", "风", "风神录", "fsl", "th10", "MoF")
	addWorkMap("11", "地", "殿", "地灵殿", "dld", "th11", "SA")
	addWorkMap("12", "星", "船", "星莲船", "xlc", "th12", "UFO")
	addWorkMap("128", "大", "大战争", "dzz", "th128", "128")
	addWorkMap("13", "神", "庙", "神灵庙", "slm", "th13", "TD")
	addWorkMap("14", "辉", "城", "辉针城", "hzc", "th14", "DDC")
	addWorkMap("15", "绀", "绀珠传", "gzz", "th15", "LoLK")
	addWorkMap("16", "天", "璋", "天空璋", "tkz", "th16", "HSiFS")
	addWorkMap("17", "鬼", "鬼形兽", "gxs", "th17", "WBaWC")
	addWorkMap("18", "虹", "洞", "虹龙洞", "hld", "th18", "UM")
	addOtherMap(parseRank, "Easy", "e")
	addOtherMap(parseRank, "Normal", "n")
	addOtherMap(parseRank, "Hard", "h")
	addOtherMap(parseRank, "Lunatic", "l")
	addOtherMap(parseRank, "Extra", "ex", "Phantasm", "ph")
	addOtherMap(parseCharacter, "Reimu", "灵", "梦", "灵梦", "单灵梦", "单灵", "单梦", "Reimu", "博麗霊夢")
	addOtherMap(parseCharacter, "Marisa", "魔", "魔理沙", "m", "单魔", "单魔理沙", "Marisa")
	addOtherMap(parseCharacter, "Sakuya", "咲", "咲夜", "s", "16", "单咲", "单咲夜", "单16", "Sakuya")
	addOtherMap(parseCharacter, "Sanae", "苗", "早苗", "Sanae")
	addOtherMap(parseCharacter, "Youmu", "妖", "妖梦", "单妖", "单妖梦", "Youmu")
	addOtherMap(parseCharacter, "RY", "结界组", "RY", "霊夢＆紫")
	addOtherMap(parseCharacter, "MA", "咏唱组", "MA", "魔理沙＆アリス")
	addOtherMap(parseCharacter, "SR", "红魔组", "SR", "咲夜＆レミリア")
	addOtherMap(parseCharacter, "YY", "幽冥组", "YY", "妖夢＆幽々子")
	addOtherMap(parseCharacter, "Yukari", "紫", "八云紫", "单紫", "Yukari")
	addOtherMap(parseCharacter, "Alice", "小爱", "爱丽丝", "单小爱", "单爱", "单爱丽丝", "Alice", "アリス・Ｍ")
	addOtherMap(parseCharacter, "Remilia", "蕾米", "蕾米莉亚", "单蕾米", "单蕾米莉亚", "Remilia", "レミリア・Ｓ")
	addOtherMap(parseCharacter, "Yuyuko", "幽幽子", "uuz", "单幽幽子", "Yuyuko", "西行寺幽々子")
	addOtherMap(parseCharacter, "Reisen", "铃仙", "灵仙", "兔子", "Reisen")
	addOtherMap(parseCharacter, "Cirno", "琪露诺", "⑨", "Cirno")
	addOtherMap(parseCharacter, "Aya", "文", "文文", "射命丸文", "Aya")
	addOtherMap(parseCType, "A", "A")
	addOtherMap(parseCType, "B", "B")
	addOtherMap(parseCType, "C", "C")
	addOtherMap(parseCType, "Spring", "春")
	addOtherMap(parseCType, "Summer", "夏")
	addOtherMap(parseCType, "Autumn", "秋")
	addOtherMap(parseCType, "Winter", "冬")
	addOtherMap(parseCType, "Wolf", "狼")
	addOtherMap(parseCType, "Otter", "獭")
	addOtherMap(parseCType, "Eagle", "鹰")
	addOtherMap(parseRoute, "6A", "6A")
	addOtherMap(parseRoute, "6B", "6B")
	for _, cType := range []string{"SA", "SB", "SC", "SD"} {
		for _, ch := range []string{"Reimu", "Marisa", "Cirno", "Aya"} {
			addOtherMap(parseCharacterCType, ch+cType, ch+cType)
		}
	}
	for _, cType := range []string{"W", "O", "E"} {
		for _, ch := range []string{"Reimu", "Marisa", "Youmu"} {
			addOtherMap(parseCharacterCType, ch+cType, ch+cType)
		}
	}
}

func tryParse(res *ParseResult, t *trie, s *string, nLen int) bool {
	ref := []rune(*s)
	length := len(ref)
	for i := range ref {
		m := nLen
		if i+m > length {
			m = length - i
		}
		for n := m; n > 0; n-- {
			f := t.get(string(ref[i : i+n]))
			if f != nil {
				f(res)
				*s = string(append(ref[:i], ref[i+n:]...))
				return true
			}
		}
	}
	return false
}

func ParseMsg(s string) *ParseResult {
	if len(s) == 0 {
		panic("s is empty")
	}
	length := len(s)
	s = strings.ToLower(s)
	res := newParseResult()
	tryParse(res, workDict, &s, length)
	for i := 0; i < 10; i++ {
		if !tryParse(res, otherDict, &s, length) {
			break
		}
	}
	for i := 0; i < 10; i++ {
		if !tryParse(res, workDict, &s, length) {
			break
		}
	}
	if strings.Contains(s, "全卡") {
		res.AllSpell = true
		s = strings.Replace(s, "全卡", "", 1)
	}
	return res
}
