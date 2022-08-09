package main

import (
	"fmt"
	"github.com/CuteReimu/bilibili"
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/config"
	"github.com/Logiase/MiraiGo-Template/utils"
	_ "github.com/Touhou-Freshman-Camp/tfcc-bot-go/commandHandler"
	"github.com/Touhou-Freshman-Camp/tfcc-bot-go/db"
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
	_, err = os.Stat("device.json")
	if err != nil {
		bot.GenRandomDevice()
	}
	utils.WriteLogToFS(utils.LogWithStack)
	config.Init()
}

func main() {
	// 初始化
	db.Init()
	initBilibili()
	bot.Init()
	bot.StartService()

	// 使用协议
	// 不同协议可能会有部分功能无法使用
	// 在登陆前切换协议
	bot.UseProtocol(bot.IPad)

	// 登录
	if err := bot.Login(); err != nil {
		logrus.Errorf("%+v", err)
		bot.Stop()
		db.Stop()
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
	db.Stop()
}

func writeConfig() {
	config.GlobalConfig = &config.Config{Viper: viper.New()}
	config.GlobalConfig.Set("bot.loginmethod", "qrcode")
	config.GlobalConfig.Set("bot.account", int64(0))
	config.GlobalConfig.Set("bot.password", "")
	config.GlobalConfig.Set("qq.super_admin_qq", int64(12345678))
	config.GlobalConfig.Set("qq.qq_group", []int64{12345678})
	config.GlobalConfig.Set("bilibili.username", "13888888888")
	config.GlobalConfig.Set("bilibili.password", "12345678")
	config.GlobalConfig.Set("bilibili.mid", "12345678")
	config.GlobalConfig.Set("bilibili.room_id", "12345678")
	config.GlobalConfig.Set("bilibili.area_v2", "236")
	err := config.GlobalConfig.WriteConfigAs("application.yaml")
	if err != nil {
		fmt.Println("生成application.yaml失败，请检查")
	} else {
		fmt.Println("已生成application.yaml，请修改配置后重新启动")
	}
}

func initBilibili() {
	savedCookies := db.Get([]byte("cookies"))
	if savedCookies != nil {
		bilibili.SetCookiesString(string(savedCookies))
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
		logrus.Fatalf("%+v", qrCode)
		return
	}
	qrCode.Print()
	fmt.Println("请扫码后按回车")
	var line string
	_, _ = fmt.Scanln(&line)
	if err = bilibili.LoginWithQRCode(qrCode); err != nil {
		logrus.Fatalf("%+v", err)
	}
	logrus.Infoln("登录bilibili成功")
	db.Set([]byte("cookies"), []byte(bilibili.GetCookiesString()))
}
