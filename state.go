package main

import (
	"encoding/json"
	"github.com/stvp/rollbar"
	"time"
)

type (
	QueueState  map[string]uint
	ServerState map[string]QueueState
)

const (
	StateMetaKey      = "state"
	StateSaveInterval = 10
)

func SaveState() {
	state := make(ServerState)
	for _, q := range queues {
		state[q.Name] = QueueState{
			"wi": q.Counter.Write,
			"ri": q.Counter.Read,
		}
	}

	stateJson, _ := json.Marshal(state)
	key := Key(StateMetaKey)
	if err := storage.Set(key, stateJson); err != nil {
		rollbar.Error("error", err)
		log.Printf("Failed to persist state")
		return
	}
}

func LoadState() {
	state := make(ServerState)
	key := Key(StateMetaKey)

	stateJson, err := storage.Get(key)
	if err != nil {
		log.Printf("State not found")
		return
	}

	if err := json.Unmarshal(stateJson, &state); err != nil {
		rollbar.Error("error", err)
		log.Printf("Failed to load state")
		return
	}

	for queueName, meta := range state {
		RegisterQueue(queueName, meta["wi"], meta["ri"])
	}

	log.Printf("State successfully loaded")
}

func KeepStatePersisted() {
	t := time.NewTicker(time.Second)
	for {
		<-t.C
		SaveState()
	}
}
