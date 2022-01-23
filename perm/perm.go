package perm

import (
	"github.com/Logiase/MiraiGo-Template/config"
	"github.com/Touhou-Freshman-Camp/tfcc-bot-go/db"
	"strconv"
)

const adminPrefix = "admin:"

func IsSuperAdmin(qq int64) bool {
	return qq == config.GlobalConfig.GetInt64("qq.super_admin_qq")
}

func IsAdmin(qq int64) bool {
	if IsSuperAdmin(qq) {
		return true
	}
	buf := db.Get([]byte(adminPrefix + strconv.FormatInt(qq, 10)))
	return buf != nil
}

func AddAdmin(qq int64) {
	db.Set([]byte(adminPrefix+strconv.FormatInt(qq, 10)), []byte{'1'})
}

func DelAdmin(qq int64) {
	db.Del([]byte(adminPrefix + strconv.FormatInt(qq, 10)))
}

const whitelistPrefix = "whitelist:"

func IsWhitelist(qq int64) bool {
	buf := db.Get([]byte(whitelistPrefix + strconv.FormatInt(qq, 10)))
	return buf != nil
}

func AddWhitelist(qq int64) {
	db.Set([]byte(whitelistPrefix+strconv.FormatInt(qq, 10)), []byte{'1'})
}

func DelWhitelist(qq int64) {
	db.Del([]byte(whitelistPrefix + strconv.FormatInt(qq, 10)))
}
