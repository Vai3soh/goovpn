package cache

import (
	"github.com/Vai3soh/goovpn/entity"
)

type Db struct {
	memory map[string]entity.Profile
}

type Option func(*Db)

func NewDb(opts ...Option) *Db {
	c := &Db{}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func WithMapMemory(m map[string]entity.Profile) Option {
	return func(d *Db) {
		d.memory = m
	}
}

func (db *Db) Save(cfgPath, body string) {
	if _, ok := db.memory[cfgPath]; !ok {
		db.memory[cfgPath] = entity.Profile{Body: body}
	}
}

func (db *Db) GetProfile(cfgPath string) entity.Profile {
	return db.memory[cfgPath]
}
