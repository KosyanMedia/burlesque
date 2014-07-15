package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"runtime"
	"strings"
)

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	info := make(map[string]map[string]uint)

	for _, q := range queues {
		info[q.Name] = map[string]uint{
			"messages":      q.Counter.Distance(),
			"subscriptions": 0,
		}
	}

	for _, r := range pool {
		for _, q := range r.Queues {
			info[q]["subscriptions"]++
		}
	}

	infoJson, _ := json.Marshal(info)
	w.Write(infoJson)
}

func DebugHandler(w http.ResponseWriter, r *http.Request) {
	info := make(map[string]int)
	info["goroutines"] = runtime.NumGoroutine()

	infoJson, _ := json.Marshal(info)
	w.Write(infoJson)
}

func PublishHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	msg, _ := ioutil.ReadAll(r.Body)
	if len(msg) == 0 {
		msg = Message(r.FormValue("msg"))
	}

	queueName := r.FormValue("queue")
	go Register(queueName, msg)

	Debug("Published message of %d bytes to queue %s", len(msg), queueName)
	w.Write([]byte("OK"))
}

func SubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	rch := make(chan *Response)
	req := &Request{
		Queues: strings.Split(r.FormValue("queues"), ","),
		Callback: func(r *Response) {
			rch <- r
		},
	}

	go Process(req)

	disconnected := w.(http.CloseNotifier).CloseNotify()
	finished := make(chan bool)
	go func() {
		select {
		case <-disconnected:
			rch <- nil
		case <-finished:
			break
		}
		Purge(req)
	}()

	res := <-rch
	if res == nil {
		return
	}

	w.Header().Set("Queue", res.Queue)
	w.Write(res.Message)

	Debug("Recieved message of %d bytes from queue %s", len(res.Message), res.Queue)
	finished <- true
}

func SetupServer() {
	http.HandleFunc("/status", StatusHandler)
	http.HandleFunc("/debug", DebugHandler)
	http.HandleFunc("/publish", PublishHandler)
	http.HandleFunc("/subscribe", SubscriptionHandler)
}
