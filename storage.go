package main

import (
	"github.com/ezotrank/cabinetgo"
	"strconv"
	"strings"
)

type (
	Message []byte
	Key     []byte
)

var (
	storage = cabinet.New()
)

func NewKey(queue string, index uint) Key {
	istr := strconv.FormatUint(uint64(index), 10)
	key := strings.Join([]string{queue, istr}, "_")
	return Key(key)
}

func SetupStorage() {
	err := storage.Open(Config.Storage, cabinet.KCOWRITER|cabinet.KCOCREATE)
	if err != nil {
		Error(err, "Failed to open database '%s'", Config.Storage)
	}
}

func CloseStorage() {
	var err error

	err = storage.Sync(true)
	if err != nil {
		Error(err, "Failed to sync storage (hard)")
	} else {
		Log("Storage synchronized")
	}

	err = storage.Close()
	if err != nil {
		Error(err, "Failed to close storage")
	} else {
		Log("Storage closed")
	}
}
