package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/KosyanMedia/burlesque/hub"
)

type (
	Server struct {
		port int
		hub  *hub.Hub
	}
)

func New(port int, h *hub.Hub) *Server {
	s := Server{
		port: port,
		hub:  h,
	}

	http.HandleFunc("/status", s.statusHandler)
	http.HandleFunc("/debug", s.debugHandler)
	http.HandleFunc("/publish", s.pubHandler)
	http.HandleFunc("/subscribe", s.subHandler)

	return &s
}

func (s *Server) Start() {
	port := fmt.Sprintf(":%d", s.port)
	if err := http.ListenAndServe(port, nil); err != nil {
		panic(err)
	}
}

func (s *Server) statusHandler(w http.ResponseWriter, r *http.Request) {
	info := s.hub.Info()
	jsn, _ := json.Marshal(info)
	w.Write(jsn)
}

func (s *Server) debugHandler(w http.ResponseWriter, r *http.Request) {
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

func (s *Server) pubHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	msg, _ := ioutil.ReadAll(r.Body)
	if len(msg) == 0 {
		msg = []byte(r.FormValue("msg"))
	}
	queue := r.FormValue("queue")

	if ok := s.hub.Pub(queue, msg); ok {
		w.Write([]byte("OK"))
	} else {
		http.Error(w, "FAIL", 500)
	}
}

func (s *Server) subHandler(w http.ResponseWriter, r *http.Request) {
	result := make(chan hub.Result)
	queues := strings.Split(r.FormValue("queues"), ",")

	sub := hub.NewSubscription(queues, result)
	defer sub.Close()

	finished := make(chan struct{})
	defer close(finished)

	disconnected := w.(http.CloseNotifier).CloseNotify()
	go func() {
		select {
		case <-disconnected:
			sub.Close()
		case <-finished:
		}
	}()

	go s.hub.Sub(sub)
	res := <-result

	w.Header().Set("Queue", res.Queue)
	w.Write(res.Message)
}
