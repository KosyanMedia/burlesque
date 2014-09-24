package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"strings"

	"github.com/KosyanMedia/burlesque/hub"
)

type (
	Server struct {
		port int
		hub  *hub.Hub
	}
)

const (
	Version = "0.2.0"
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
	http.HandleFunc("/flush", s.flushHandler)

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
	info := make(map[string]interface{})
	info["version"] = Version
	info["gomaxprocs"] = runtime.GOMAXPROCS(-1)
	info["goroutines"] = runtime.NumGoroutine()
	info["kyoto_cabinet"] = s.hub.StorageInfo()

	jsn, _ := json.Marshal(info)
	w.Write(jsn)
}

func (s *Server) pubHandler(w http.ResponseWriter, r *http.Request) {
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
	queues := strings.Split(r.FormValue("queues"), ",")
	sub := hub.NewSubscription(queues)

	finished := make(chan struct{})
	defer close(finished)

	disconnected := w.(http.CloseNotifier).CloseNotify()
	go func() {
		select {
		case <-disconnected:
		case <-finished:
		}
		sub.Close()
	}()

	go s.hub.Sub(sub)

	if res, ok := <-sub.Result(); ok {
		w.Header().Set("Queue", res.Queue)
		w.Write(res.Message)
	}
}

func (s *Server) flushHandler(w http.ResponseWriter, r *http.Request) {
	queues := strings.Split(r.FormValue("queues"), ",")
	messages := s.hub.Flush(queues)

	jsn, _ := json.Marshal(messages)
	w.Write(jsn)
}
