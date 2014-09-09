package main

import (
	"encoding/json"
	"time"
)

type (
	queueState  map[string]uint
	serverState map[string]queueState
)

const (
	stateMetaKey      = "state"
	stateSaveInterval = 1 // seconds
)

func saveState() {
	state := make(serverState)
	for _, q := range queues {
		state[q.name] = queueState{
			"wi": q.counter.writeIndex,
			"ri": q.counter.readIndex,
		}
	}

	jsn, _ := json.Marshal(state)
	k := key(stateMetaKey)
	if err := storage.Set(k, jsn); err != nil {
		alert(err, "Failed to persist state")
		return
	}
}

func loadState() {
	state := make(serverState)
	k := key(stateMetaKey)

	jsn, err := storage.Get(k)
	if err != nil {
		log("State not found")
		return
	}

	err = json.Unmarshal(jsn, &state)
	if err != nil {
		log("Failed to load state")
		return
	}

	for qname, meta := range state {
		registerQueue(qname, meta["wi"], meta["ri"])
	}

	log("State successfully loaded")
}

func keepStatePersisted() {
	t := time.NewTicker(stateSaveInterval * time.Second)

	for {
		<-t.C
		saveState()
		err := storage.Sync(false)
		if err != nil {
			alert(err, "Failed to sync storage")
		}
	}
}
