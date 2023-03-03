package tfcc

import (
	"github.com/CuteReimu/tfcc-bot-go/utils"
	"gopkg.in/yaml.v3"
	"os"
)

type JfNNData struct {
	Work      string  `yaml:"work"`
	Rank      string  `yaml:"rank"`
	Route     string  `yaml:"route"`
	Character string  `yaml:"character"`
	CType     string  `yaml:"ctype"`
	AllSpell  bool    `yaml:"allspell"`
	Jf        float64 `yaml:"jf"`
}

var jfNN map[string][]*JfNNData

func init() {
	m := make(map[string][]*JfNNData)
	buf, err := os.ReadFile("assets/score.yaml")
	if err != nil {
		logger.WithError(err).Errorln("load score.yaml failed")
		return
	}
	err = yaml.Unmarshal(buf, &m)
	if err != nil {
		logger.WithError(err).Errorln("unmarshal json failed")
		return
	}
	jfNN = m
}

func IsOK() bool {
	return jfNN != nil
}

func GetJf(work string) []*JfNNData {
	l, ok := jfNN[work]
	if ok {
		return l
	}
	return nil
}

var logger = utils.GetModuleLogger("tfcc-bot-go.tfcc")
