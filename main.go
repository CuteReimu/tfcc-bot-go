package main

import (
	"fmt"
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/config"
	"github.com/Logiase/MiraiGo-Template/utils"
	_ "github.com/Touhou-Freshman-Camp/tfcc-bot-go/commandHandler"
	"github.com/spf13/viper"
	"os"
	"os/signal"
)

func init() {
	_, err := os.Stat("application.yaml")
	if err != nil {
		writeConfig()
		b := make([]byte, 1)
		_, _ = os.Stdin.Read(b)
		os.Exit(0)
	}
	_, err = os.Stat("device.json")
	if err != nil {
		bot.GenRandomDevice()
	}
	utils.WriteLogToFS()
	config.Init()
}

func main() {
	// 快速初始化
	bot.Init()

	// 初始化 Modules
	bot.StartService()

	// 使用协议
	// 不同协议可能会有部分功能无法使用
	// 在登陆前切换协议
	bot.UseProtocol(bot.IPad)

	// 登录
	if err := bot.Login(); err != nil {
		bot.Stop()
		return
	}

	// 刷新好友列表，群列表
	bot.RefreshList()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill)
	<-ch
	bot.Stop()
}

func writeConfig() {
	config.GlobalConfig = &config.Config{Viper: viper.New()}
	config.GlobalConfig.Set("bot.account", int64(0))
	config.GlobalConfig.Set("bot.password", "")
	config.GlobalConfig.Set("qq.super_admin_qq", int64(12345678))
	config.GlobalConfig.Set("schedule.qq_group", []int64{12345678})
	config.GlobalConfig.Set("schedule.before", []int64{3 * 3600, 6 * 3600})
	config.GlobalConfig.Set("schedule.video_push_delay", int64(600))
	config.GlobalConfig.Set("repeater_interruption.qq_group", []int64{12345678})
	config.GlobalConfig.Set("repeater_interruption.allowance", 5)
	config.GlobalConfig.Set("repeater_interruption.cool_down", int64(3))
	err := config.GlobalConfig.WriteConfigAs("application.yaml")
	if err != nil {
		fmt.Println("生成application.yaml失败，请检查")
	} else {
		fmt.Println("已生成application.yaml，请修改配置后重新启动")
	}
}
