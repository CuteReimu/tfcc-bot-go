package db

import (
	"github.com/dgraph-io/badger/v3"
	"time"
)

// Set 设置键值对，ttl是超时时间（可选）
func Set(key, value []byte, ttl ...time.Duration) {
	err := DB.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry(key, value)
		if len(ttl) > 0 {
			e = e.WithTTL(ttl[0])
		}
		err := txn.SetEntry(e)
		return err
	})
	if err != nil {
		logger.WithError(err).Error("set key value failed: {%s, %s}", string(key), string(value))
	}
}

// Del 删除Key
func Del(key []byte) {
	err := DB.Update(func(txn *badger.Txn) error {
		err := txn.Delete(key)
		return err
	})
	if err != nil {
		logger.WithError(err).Error("delete key failed: %s", string(key))
	}
	return
}

// Get 根据Key获取Value
func Get(key []byte) (value []byte) {
	err := DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err == badger.ErrKeyNotFound {
			return nil
		}
		logger.WithError(err).Error("get failed, key: %s", string(key))
		value, err = item.ValueCopy(nil)
		return err
	})
	if err != nil {
		logger.WithError(err).Error("get failed, key: %s", string(key))
	}
	return
}

func PrefixScan(prefix []byte, f func(key, value []byte) error) {
	err := DB.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			k := item.Key()
			err := item.Value(func(v []byte) error {
				return f(k, v)
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		logger.WithError(err).Error("prefix scan failed")
	}
}
