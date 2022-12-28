package translate

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"io"
	"strings"
)

//go:embed translate.tsv
var transLateData []byte

var translateDict = &Trie{}
var regexpSpace = regexp.MustCompile(`(?<![()\[\]{}%'"A-Za-z-]) (?![()\[\]{}%'"A-Za-z-])`, regexp.None)

func init() {
	reader := bufio.NewReader(bytes.NewReader(transLateData))
	for {
		line, _, err := reader.ReadLine()
		if err != nil && err != io.EOF {
			panic(err)
		}
		if len(line) > 0 {
			arr := strings.Split(string(line), "\t")
			var key, val string
			key = arr[0]
			if len(arr) >= 2 {
				val = arr[1]
			}
			if !translateDict.PutIfAbsent(key, val) {
				panic(fmt.Sprint("出现重复数据：", string(line)))
			}
		}
		if err == io.EOF {
			break
		}
	}
	m := map[string]string{
		"beat the WR":                            "打破了世界纪录：",
		"in Hollow Knight Category Extensions -": "",
		"King's Pass: Level":                     "国王山道",
		"in Hollow Knight -":                     "",
		"The new WR is":                          "新的世界纪录是",
		"Their time is":                          "时间是",
		"Its time is":                            "时间是",
		"The time is":                            "时间是",
		"His time is":                            "时间是",
		"Her time is":                            "时间是",
		"got a new top 3 PB":                     "获得了前三：",
		"Pantheon of the Master: Level":          "大师万神殿",
		"Pantheon of the Artist: Level":          "艺术家万神殿",
		"Pantheon of the Sage: Level":            "贤者万神殿",
		"Pantheon of the Knight: Level":          "骑士万神殿",
		"Pantheon of Hallownest: Level":          "圣巢万神殿",
		"White Palace: Level":                    "白色宫殿",
		"Path of Pain: Level":                    "苦痛之路",
		"Trial of the Warrior: Level":            "勇士的试炼",
		"Trial of the Conqueror: Level":          "征服者的试炼",
		"Trial of the Fool: Level":               "愚人的试炼",
		"NMG.":                                   "无主要邪道.",
		"- NMG":                                  "- 无主要邪道",
		"- NMG.":                                 "- 无主要邪道.",
		"Console Runs":                           "主机速通",
		"Any Bindings":                           "任意锁",
	}
	for key, val := range m {
		if !translateDict.PutIfAbsent(key, val) {
			panic(fmt.Sprint("出现重复数据：", key))
		}
	}
}

func Translate(s string) string {
	s = translateDict.ReplaceAll(s)
	s, err := regexpSpace.Replace(s, "", -1, -1)
	if err != nil {
		panic(err)
	}
	return s
}

type trieNode struct {
	child  map[rune]*trieNode
	value  string
	exists bool
}

type Trie struct {
	root trieNode
}

func (t *Trie) PutIfAbsent(key, value string) bool {
	node := &t.root
	for _, c := range strings.ToLower(key) {
		var n *trieNode
		if node.child == nil {
			node.child = make(map[rune]*trieNode)
		}
		n = node.child[c]
		if n != nil {
			node = n
		} else {
			newNode := &trieNode{}
			node.child[c] = newNode
			node = newNode
		}
	}
	if node.exists {
		return false
	}
	node.exists = true
	node.value = value
	return true
}

func (t *Trie) getLongest(s string) (string, string) {
	var node, node2 *trieNode
	var key, key2 string
	node = &t.root
	r := []rune(strings.ToLower(s))
	for idx, c := range r {
		if node.child != nil {
			if n, ok := node.child[c]; ok {
				key += string(c)
				node = n
				if node.exists && (idx+1 >= len(s) || symbols[r[idx+1]]) {
					node2 = node
					key2 = key
				}
				continue
			}
		}
		break
	}
	if node2 != nil {
		return key2, node2.value
	}
	return "", ""
}

func (t *Trie) ReplaceAll(str string) string {
	s := []rune(str)
	var s2 []rune
	for len(s) > 0 {
		if !(len(s2) == 0 || symbols[s2[len(s2)-1]]) {
			s2 = append(s2, s[0])
			s = s[1:]
			continue
		}
		key, value := t.getLongest(string(s))
		if len(key) > 0 {
			s2 = append(s2, []rune(value)...)
			s = s[len([]rune(key)):]
		} else {
			s2 = append(s2, s[0])
			s = s[1:]
		}
	}
	return string(s2)
}

var symbols = map[rune]bool{
	' ':  true,
	'(':  true,
	')':  true,
	'[':  true,
	']':  true,
	'{':  true,
	'}':  true,
	'%':  true,
	'\'': true,
	'"':  true,
	'.':  true,
	',':  true,
}
