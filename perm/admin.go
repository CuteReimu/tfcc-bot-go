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

func AddAdmin(qq int64) bool {
	db.Set([]byte(adminPrefix+strconv.FormatInt(qq, 10)), []byte{'1'})
	return true
}

func DelAdmin(qq int64) {
	db.Del([]byte(adminPrefix + strconv.FormatInt(qq, 10)))
}

// ListAdmin 因为这个接口一般用来展示，所以返回[]string
func ListAdmin() (list []string) {
	db.PrefixScan([]byte(adminPrefix), func(key, value []byte) error {
		if len(key) > len(adminPrefix) {
			list = append(list, string(key)[len(adminPrefix):])
		}
		return nil
	})
	return
}
