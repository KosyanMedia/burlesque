package main

import (
	"strconv"
	"strings"

	"bitbucket.org/ww/cabinet"
)

type (
	message []byte
	key     []byte
)

var (
	storage = cabinet.New()
)

func newKey(queue string, index uint) key {
	istr := strconv.FormatUint(uint64(index), 10)
	k := strings.Join([]string{queue, istr}, "_")

	return key(k)
}

func setupStorage() {
	err := storage.Open(config.storage, cabinet.KCOWRITER|cabinet.KCOCREATE)
	if err != nil {
		alert(err, "Failed to open database '%s'", config.storage)
	}
}

func closeStorage() {
	var err error

	err = storage.Sync(true)
	if err != nil {
		alert(err, "Failed to sync storage (hard)")
	} else {
		log("Storage synchronized")
	}

	err = storage.Close()
	if err != nil {
		alert(err, "Failed to close storage")
	} else {
		log("Storage closed")
	}
}
