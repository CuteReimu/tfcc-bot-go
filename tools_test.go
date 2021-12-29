package main_test

import (
	"testing"

	"github.com/Logiase/MiraiGo-Template/bot"
)

// 用于生成device.json
func TestGenDevice(t *testing.T) {
	bot.GenRandomDevice()
}
