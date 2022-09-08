package factory

import (
	"Frieza/master"
)

var sl master.Slaves

func New() (master.Slaves, error) {
	if sl != nil {
		return sl, nil
	} else {
		return nil, master.ErrNotFound
	}
}

func Register(slave master.Slaves) {
	if slave == nil {
		panic("streamsList: Register ch is nil")
	}
	if sl != nil {
		panic("streamsList: Register called twice for slaves")
	}
	sl = slave
}
