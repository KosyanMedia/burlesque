package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"strconv"
	"strings"
)

func startServer() {
	port := fmt.Sprintf(":%d", config.port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		alert(err, "Error starting server on port %d", config.port)
	}
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	info := make(map[string]map[string]uint)

	for _, q := range queues {
		info[q.name] = map[string]uint{
			"messages":      q.counter.distance(),
			"subscriptions": 0,
		}
	}

	for _, r := range pool.requests {
		for _, q := range r.queues {
			info[q]["subscriptions"]++
		}
	}

	jsn, _ := json.Marshal(info)
	w.Write(jsn)
}

func debugHandler(w http.ResponseWriter, r *http.Request) {
	info := make(map[string]interface{})
	info["version"] = version
	info["goroutines"] = runtime.NumGoroutine()

	s, err := storage.Status()
	if err != nil {
		alert(err, "Failed to get Kyoto Cabinet status")
	}
	s = s[:len(s)-1] // Removing trailing new line

	ks := make(map[string]interface{})
	tokens := strings.Split(s, "\n")
	for _, t := range tokens {
		tt := strings.Split(t, "\t")
		num, err := strconv.Atoi(tt[1])
		if err != nil {
			ks[tt[0]] = tt[1]
		} else {
			ks[tt[0]] = num
		}
	}
	info["kyoto_cabinet"] = ks

	jsn, _ := json.Marshal(info)
	w.Write(jsn)
}

func publishHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	msg, _ := ioutil.ReadAll(r.Body)
	if len(msg) == 0 {
		msg = message(r.FormValue("msg"))
	}

	qname := r.FormValue("queue")
	ok := registerPublication(qname, msg)

	if ok {
		w.Write([]byte("OK"))
	} else {
		http.Error(w, "FAIL", 500)
	}
}

func subscriptionHandler(w http.ResponseWriter, r *http.Request) {
	rch := make(chan response)
	abort := make(chan bool, 1)
	req := &request{
		queues:     strings.Split(r.FormValue("queues"), ","),
		responseCh: rch,
		abort:      abort,
	}
	go registerSubscription(req)

	disconnected := w.(http.CloseNotifier).CloseNotify()
	finished := make(chan bool)
	go func() {
		select {
		case <-disconnected:
			close(rch)
			abort <- true
		case <-finished:
		}
		req.purge()
	}()

	res, ok := <-rch
	if !ok {
		return
	}

	w.Header().Set("Queue", res.queue)
	w.Write(res.message)

	finished <- true
}

func setupServer() {
	http.HandleFunc("/status", statusHandler)
	http.HandleFunc("/debug", debugHandler)
	http.HandleFunc("/publish", publishHandler)
	http.HandleFunc("/subscribe", subscriptionHandler)
}
