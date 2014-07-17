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

	for _, r := range pool.Requests {
		for _, q := range r.Queues {
			info[q]["subscriptions"]++
		}
	}

	jsn, _ := json.Marshal(info)
	w.Write(jsn)
}

func DebugHandler(w http.ResponseWriter, r *http.Request) {
	info := make(map[string]interface{})
	info["goroutines"] = runtime.NumGoroutine()

	s, err := storage.Status()
	if err != nil {
		Error(err, "Failed to get Kyoto Cabinet status")
	}
	info["kyoto_cabinet"] = s

	jsn, _ := json.Marshal(info)
	w.Write(jsn)
}

func PublishHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	msg, _ := ioutil.ReadAll(r.Body)
	if len(msg) == 0 {
		msg = Message(r.FormValue("msg"))
	}

	qname := r.FormValue("queue")
	ok := RegisterPublication(qname, msg)

	if ok {
		Debug("Published message of %d bytes to queue %s", len(msg), qname)
		w.Write([]byte("OK"))
	} else {
		Debug("Failed to publish message of %d bytes to queue %s", len(msg), qname)
		http.Error(w, "FAIL", 500)
	}
}

func SubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	rch := make(chan Response)
	abort := make(chan bool, 1)
	req := &Request{
		Queues:     strings.Split(r.FormValue("queues"), ","),
		ResponseCh: rch,
		Abort:      abort,
	}
	go RegisterSubscription(req)

	disconnected := w.(http.CloseNotifier).CloseNotify()
	finished := make(chan bool)
	go func() {
		select {
		case <-disconnected:
			close(rch)
			abort <- true
		case <-finished:
		}
		req.Purge()
	}()

	res, ok := <-rch
	if !ok {
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
