package db

import (
	"github.com/CuteReimu/dets"
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/dgraph-io/badger/v3"
	"time"
)

var DB *badger.DB

var logger = utils.GetModuleLogger("db")

func Init() {
	var err error
	DB, err = badger.Open(badger.DefaultOptions("assets/database"))
	if err != nil {
		logger.WithError(err).Fatal("init database failed")
	}
	dets.SetDB(DB, logger)
	go gc()
}

func gc() {
	ticker := time.NewTicker(time.Hour)
	for range ticker.C {
	again:
		err := DB.RunValueLogGC(0.7)
		if err == nil {
			goto again
		}
	}
}

func Stop() {
	err := DB.Close()
	if err != nil {
		logger.WithError(err).Error("close database failed")
	}
}
