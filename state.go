package main

import (
	"encoding/json"
	"time"
)

type (
	QueueState  map[string]uint
	ServerState map[string]QueueState
)

const (
	StateMetaKey      = "state"
	StateSaveInterval = 1 // seconds
)

func SaveState() {
	state := make(ServerState)
	for _, q := range queues {
		state[q.Name] = QueueState{
			"wi": q.Counter.WriteIndex,
			"ri": q.Counter.ReadIndex,
		}
	}

	jsn, _ := json.Marshal(state)
	key := Key(StateMetaKey)
	if err := storage.Set(key, jsn); err != nil {
		Error(err, "Failed to persist state")
		return
	}
}

func LoadState() {
	state := make(ServerState)
	key := Key(StateMetaKey)

	jsn, err := storage.Get(key)
	if err != nil {
		Log("State not found")
		return
	}

	err = json.Unmarshal(jsn, &state)
	if err != nil {
		Log("Failed to load state")
		return
	}

	for qname, meta := range state {
		RegisterQueue(qname, meta["wi"], meta["ri"])
	}

	Log("State successfully loaded")
}

func KeepStatePersisted() {
	t := time.NewTicker(StateSaveInterval * time.Second)
	for {
		<-t.C
		SaveState()
		err := storage.Sync(false)
		if err != nil {
			Error(err, "Failed to sync storage")
		}
	}
}
