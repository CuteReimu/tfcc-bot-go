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
	"strings"
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
	db.Stop()
}

func writeConfig() {
	config.GlobalConfig = &config.Config{Viper: viper.New()}
	config.GlobalConfig.Set("bot.account", int64(0))
	config.GlobalConfig.Set("bot.password", "")
	config.GlobalConfig.Set("qq.super_admin_qq", int64(12345678))
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
	captchaResult, err := bilibili.Captcha()
	if err != nil {
		logrus.Fatalf("%+v", err)
	}
	if captchaResult.Code != 0 {
		logrus.Fatalf("登录bilibili获取人机校验失败, code: %d\n", captchaResult.Code)
	}
	fmt.Println("gt:", captchaResult.Data.Geetest.Gt)
	fmt.Println("challenge:", captchaResult.Data.Geetest.Challenge)
	fmt.Println("请前往以下链接进行人机验证：")
	fmt.Println("https://kuresaru.github.io/geetest-validator/")
	fmt.Println("验证后请输入validate：")
	var line string
	_, err = fmt.Scanln(&line)
	if err != nil {
		logrus.WithError(err).Fatalln("读取stdin失败")
	}
	validate := strings.TrimSpace(line)
	seccode := validate
	userName := config.GlobalConfig.GetString("bilibili.username")
	pwd := config.GlobalConfig.GetString("bilibili.password")
	err = bilibili.LoginWithPassword(userName, pwd, captchaResult, validate, seccode)
	if err != nil {
		logrus.Fatalf("%+v", err)
	}
	logrus.Infoln("登录bilibili成功")
	db.Set([]byte("cookies"), []byte(bilibili.GetCookiesString()))
}
