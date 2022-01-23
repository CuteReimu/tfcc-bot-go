package db

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/pkg/errors"
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
		logger.WithError(err).Errorf("set key value failed: {%s, %s}\n", string(key), string(value))
	}
}

// Del 删除Key
func Del(key []byte) {
	err := DB.Update(func(txn *badger.Txn) error {
		err := txn.Delete(key)
		return err
	})
	if err != nil {
		logger.WithError(err).Errorf("delete key failed: %s\n", string(key))
	}
}

// Get 根据Key获取Value
func Get(key []byte) (value []byte) {
	err := DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err == badger.ErrKeyNotFound {
			return nil
		} else if err != nil {
			return errors.Wrapf(err, "get failed, key: %s", string(key))
		}
		value, err = item.ValueCopy(nil)
		return err
	})
	if err != nil {
		logger.WithError(err).Errorf("get failed, key: %s\n", string(key))
	}
	return
}

func PrefixScanKey(prefix []byte, f func(key []byte) error) {
	err := DB.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.IteratorOptions{PrefetchSize: 100})
		defer it.Close()
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			k := item.Key()
			err := f(k)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		logger.WithError(err).Errorln("prefix scan failed")
	}
}
