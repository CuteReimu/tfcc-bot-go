package main

import (
	"fmt"
	"github.com/CuteReimu/bilibili"
	"github.com/CuteReimu/dets"
	"github.com/CuteReimu/tfcc-bot-go/bot"
	_ "github.com/CuteReimu/tfcc-bot-go/chatPipeline"
	_ "github.com/CuteReimu/tfcc-bot-go/commandHandler"
	"github.com/CuteReimu/tfcc-bot-go/config"
	"github.com/CuteReimu/tfcc-bot-go/db"
	_ "github.com/CuteReimu/tfcc-bot-go/repAnalyze"
	_ "github.com/CuteReimu/tfcc-bot-go/repeaterInterruption"
	"github.com/CuteReimu/tfcc-bot-go/utils"
	_ "github.com/CuteReimu/tfcc-bot-go/videoPusher"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"time"
)

func init() {
	_, err := os.Stat("application.yaml")
	if err != nil {
		writeConfig()
		b := make([]byte, 1)
		_, _ = os.Stdin.Read(b)
		os.Exit(0)
	}
	logrus.SetReportCaller(true)
	utils.WriteLogToFS(utils.LogWithStack)
	config.Init()
}

func main() {
	// 初始化
	db.Init()
	defer db.Stop()
	initBilibili()
	bot.Init()
	bot.StartService()

	// 使用协议
	// 不同协议可能会有部分功能无法使用
	// 在登陆前切换协议
	bot.UseProtocol(bot.AndroidWatch)

	// 登录
	if err := bot.Login(); err != nil {
		logrus.Errorf("%+v", err)
		bot.Stop()
		return
	}

	// 刷新好友列表，群列表
	bot.RefreshList()
	bot.SaveToken()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
	bot.SaveToken()
	bot.Stop()
}

func writeConfig() {
	config.GlobalConfig = &config.Config{Viper: viper.New()}
	config.GlobalConfig.Set("bot.login-method", "qrcode")
	config.GlobalConfig.Set("bot.account", int64(0))
	config.GlobalConfig.Set("bot.password", "")
	config.GlobalConfig.Set("qq.rand_count", int64(10))
	config.GlobalConfig.Set("qq.rand_one_time_limit", 2)
	config.GlobalConfig.Set("qq.related_url", "")
	config.GlobalConfig.Set("qq.super_admin_qq", int64(12345678))
	config.GlobalConfig.Set("qq.qq_group", []int64{12345678})
	config.GlobalConfig.Set("schedule.qq_group", []int64{12345678})
	config.GlobalConfig.Set("schedule.before", []int64{3 * 3600, 6 * 3600})
	config.GlobalConfig.Set("schedule.video_push_delay", int64(600))
	config.GlobalConfig.Set("repeater_interruption.qq_group", []int64{12345678})
	config.GlobalConfig.Set("repeater_interruption.allowance", 5)
	config.GlobalConfig.Set("repeater_interruption.cool_down", int64(3))
	config.GlobalConfig.Set("bilibili.username", "13888888888")
	config.GlobalConfig.Set("bilibili.password", "12345678")
	config.GlobalConfig.Set("bilibili.mid", "12345678")
	config.GlobalConfig.Set("bilibili.room_id", "12345678")
	config.GlobalConfig.Set("bilibili.area_v2", "236")
	config.GlobalConfig.Set("thwiki.enable", false)
	err := config.GlobalConfig.WriteConfigAs("application.yaml")
	if err != nil {
		fmt.Println("生成application.yaml失败，请检查")
	} else {
		fmt.Println("已生成application.yaml，请修改配置后重新启动")
	}
}

func initBilibili() {
	savedCookies := dets.GetString([]byte("cookies"))
	if len(savedCookies) > 0 {
		bilibili.SetCookiesString(savedCookies)
		cookies := bilibili.GetCookies()
		now := time.Now()
		upToDate := true
		for _, cookie := range cookies {
			if now.After(cookie.Expires) {
				upToDate = false
				break
			}
		}
		if upToDate {
			return
		}
	}
	qrCode, err := bilibili.GetQRCode()
	if err != nil {
		logrus.Fatalf("%+v", err)
		return
	}
	qrCode.Print()
	fmt.Println("B站登录过期，请扫码登录B站后按回车")
	var line string
	_, _ = fmt.Scanln(&line)
	if err = bilibili.LoginWithQRCode(qrCode); err != nil {
		logrus.Fatalf("%+v", err)
	}
	logrus.Infoln("登录bilibili成功")
	dets.Put([]byte("cookies"), bilibili.GetCookiesString())
}
