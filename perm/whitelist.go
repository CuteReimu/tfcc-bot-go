package perm

import (
	"github.com/Touhou-Freshman-Camp/tfcc-bot-go/db"
	"strconv"
)

const whitelistPrefix = "whitelist:"

func IsWhitelist(qq int64) bool {
	buf := db.Get([]byte(whitelistPrefix + strconv.FormatInt(qq, 10)))
	return buf != nil
}

func AddWhitelist(qq int64) bool {
	db.Set([]byte(whitelistPrefix+strconv.FormatInt(qq, 10)), []byte{'1'})
	return true
}

func DelWhitelist(qq int64) {
	db.Del([]byte(whitelistPrefix + strconv.FormatInt(qq, 10)))
}

// ListWhitelist 因为这个接口一般用来展示，所以返回[]string
func ListWhitelist() (list []string) {
	db.PrefixScan([]byte(whitelistPrefix), func(key, value []byte) error {
		if len(key) > len(whitelistPrefix) {
			list = append(list, string(key)[len(whitelistPrefix):])
		}
		return nil
	})
	return
}
