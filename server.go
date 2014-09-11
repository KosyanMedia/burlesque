package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/KosyanMedia/burlesque/hub"
)

func startServer() {
	http.HandleFunc("/status", statusHandler)
	http.HandleFunc("/debug", debugHandler)
	http.HandleFunc("/publish", pubHandler)
	http.HandleFunc("/subscribe", subHandler)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", config.port), nil); err != nil {
		panic(err)
	}
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	// info := make(map[string]map[string]uint)

	// for _, q := range queues {
	// 	info[q.name] = map[string]uint{
	// 		// "messages":      q.counter.distance(),
	// 		"subscriptions": 0,
	// 	}
	// }

	// for _, r := range pool.requests {
	// 	for _, q := range r.queues {
	// 		info[q]["subscriptions"]++
	// 	}
	// }

	// jsn, _ := json.Marshal(info)
	// w.Write(jsn)
}

func debugHandler(w http.ResponseWriter, r *http.Request) {
	// info := make(map[string]interface{})
	// info["version"] = version
	// info["goroutines"] = runtime.NumGoroutine()

	// s, err := storage.Status()
	// if err != nil {
	// 	alert(err, "Failed to get Kyoto Cabinet status")
	// }
	// s = s[:len(s)-1] // Removing trailing new line

	// ks := make(map[string]interface{})
	// tokens := strings.Split(s, "\n")
	// for _, t := range tokens {
	// 	tt := strings.Split(t, "\t")
	// 	num, err := strconv.Atoi(tt[1])
	// 	if err != nil {
	// 		ks[tt[0]] = tt[1]
	// 	} else {
	// 		ks[tt[0]] = num
	// 	}
	// }
	// info["kyoto_cabinet"] = ks

	// jsn, _ := json.Marshal(info)
	// w.Write(jsn)
}

func pubHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	msg, _ := ioutil.ReadAll(r.Body)
	if len(msg) == 0 {
		msg = []byte(r.FormValue("msg"))
	}
	queue := r.FormValue("queue")

	if ok := theHub.Pub(queue, msg); ok {
		w.Write([]byte("OK"))
	} else {
		http.Error(w, "FAIL", 500)
	}
}

func subHandler(w http.ResponseWriter, r *http.Request) {
	result := make(chan hub.Result)
	queues := strings.Split(r.FormValue("queues"), ",")
	sub := hub.NewSubscription(queues, result)

	disconnected := w.(http.CloseNotifier).CloseNotify()
	finished := make(chan struct{})
	go func() {
		select {
		case <-disconnected:
			sub.Close()
			close(finished)
		case <-finished:
		}
	}()
	defer sub.Close()

	theHub.Sub(sub)
	res := <-result

	w.Header().Set("Queue", res.Queue)
	w.Write(res.Message)

	finished <- struct{}{}
}
