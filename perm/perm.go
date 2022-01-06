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

// ListAdmin 因为这个接口一般用来展示，所以返回[]string
func ListAdmin() (list []string) {
	list = append(list, strconv.FormatInt(config.GlobalConfig.GetInt64("qq.super_admin_qq"), 10))
	db.PrefixScanKey([]byte(adminPrefix), func(key []byte) error {
		if len(key) > len(adminPrefix) {
			list = append(list, string(key)[len(adminPrefix):])
		}
		return nil
	})
	return
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

// ListWhitelist 因为这个接口一般用来展示，所以返回[]string
func ListWhitelist() (list []string) {
	db.PrefixScanKey([]byte(whitelistPrefix), func(key []byte) error {
		if len(key) > len(whitelistPrefix) {
			list = append(list, string(key)[len(whitelistPrefix):])
		}
		return nil
	})
	return
}
