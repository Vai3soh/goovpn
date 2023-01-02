package boltdb

import (
	"fmt"
	"sync"

	"github.com/Vai3soh/goovpn/entity"
	bolt "go.etcd.io/bbolt" //https://github.com/etcd-io/bbolt
)

type BoltDB struct {
	db         *bolt.DB
	path       string
	nameBucket string
	m          []entity.Message
	mu         sync.Mutex
}

func openBolt(filePath string) (*bolt.DB, error) {
	db, err := bolt.Open(filePath, 0600, nil)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func NewBoltDB(filePath string) (*BoltDB, error) {
	db, err := openBolt(filePath)
	if err != nil {
		return nil, err
	}
	return &BoltDB{db: db, path: filePath}, nil
}

func (b *BoltDB) SetDB(db *bolt.DB) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.db = db
}

func (b *BoltDB) ReOpen() error {
	b.Close()
	db, err := openBolt(b.path)
	if err != nil {
		return err
	}
	b.SetDB(db)
	return nil
}

func (b *BoltDB) SetMessage(m []entity.Message) {
	b.m = m
}

func (b *BoltDB) Message() []entity.Message {
	return b.m
}

func (b *BoltDB) SetNameBucket(name string) {
	b.nameBucket = name
}

func (b *BoltDB) NameBucket() string {
	return b.nameBucket
}

func (b *BoltDB) Path() string {
	return b.db.Path()
}

func (b *BoltDB) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.db.Close()
}

func (b *BoltDB) CreateBucket(name string) error {
	defer b.Close()
	err := b.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(name))
		if err != nil {
			return fmt.Errorf(`create bucket: [%s]`, err)
		}
		b.SetNameBucket(name)
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (b *BoltDB) DeleteBucket(name string) error {
	defer b.Close()
	err := b.db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket([]byte(name))
		if err != nil {
			return fmt.Errorf(`delete bucket: [%s]`, err)
		}
		b.SetNameBucket("")
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (b *BoltDB) Store(key, value string) error {
	defer b.Close()
	err := b.db.Update(func(tx *bolt.Tx) error {
		buck := tx.Bucket([]byte(b.NameBucket()))
		err := buck.Put([]byte(key), []byte(value))
		if err != nil {
			return fmt.Errorf(`put value,key: [%s]`, err)
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

func (b *BoltDB) StoreBulk(result []entity.Message) error {
	defer b.Close()
	err := b.db.Update(func(tx *bolt.Tx) error {
		buck := tx.Bucket([]byte(b.NameBucket()))
		for i := 0; i < len(result); i++ {
			err := buck.Put([]byte(result[i].AtrId), []byte(result[i].Value))
			if err != nil {
				return fmt.Errorf(`put value,key: [%s]`, err)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (b *BoltDB) DeleteKey(key string) error {
	defer b.Close()
	err := b.db.Update(func(tx *bolt.Tx) error {
		buck := tx.Bucket([]byte(b.NameBucket()))
		err := buck.Delete([]byte(key))
		if err != nil {
			return fmt.Errorf(`delete key: [%s]`, err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (b *BoltDB) BucketIsCreate() bool {
	defer b.Close()
	b.db.View(func(tx *bolt.Tx) error {
		buck := tx.Bucket([]byte(b.NameBucket()))
		if buck == nil {
			b.SetNameBucket("")
		}
		return nil
	})
	return b.nameBucket != ""
}

func (b *BoltDB) GetAllValue() {
	defer b.Close()
	b.db.View(func(tx *bolt.Tx) error {
		buck := tx.Bucket([]byte(b.NameBucket()))
		res := make([]entity.Message, 0)
		buck.ForEach(func(k, v []byte) error {
			m := entity.Message{AtrId: string(k), Value: string(v)}
			res = append(res, m)
			return nil
		})
		b.SetMessage(res)
		return nil
	})
}

func (b *BoltDB) GetValueFromBucket(key string) error {
	defer b.Close()
	err := b.db.View(func(tx *bolt.Tx) error {
		buck := tx.Bucket([]byte(b.NameBucket()))
		value := buck.Get([]byte(key))
		if value == nil {
			return fmt.Errorf("key [%s] not found in bucket [%s]", key, b.NameBucket())
		}
		m := entity.Message{
			AtrId: string(key),
			Value: string(value),
		}
		b.SetMessage([]entity.Message{m})
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
