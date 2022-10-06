package memory

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

func (db *Db) Store(p entity.Profile) {
	if _, ok := db.memory[p.Path]; !ok {
		db.memory[p.Path] = p
	}
}

func (db *Db) Find(key string) entity.Profile {
	return db.memory[key]
}

func (db *Db) Delete(p entity.Profile) {
	delete(db.memory, p.Path)
}
