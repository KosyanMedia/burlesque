package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strings"
	"text/template"
  // "time"
	"../hub"
  "expvar"
)

var (
    counts = expvar.NewMap("counters")
)

func init() {
	counts.Add("PubCount", 0)
	counts.Add("PubBs", 0)
	counts.Add("SubCount", 0)
	counts.Add("SubBs", 0)
}

type (
	Server struct {
    server        *http.Server
		port          int
		hub           *hub.Hub
		dashboardTmpl string
	}
)

const (
	Version = "1.2.0"
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
	http.HandleFunc("/dashboard", s.dashboardHandler)

	return &s
}

func (s *Server) Start() {
  srv := &http.Server{
    Addr:           fmt.Sprintf(":%d", s.port),
    // ReadTimeout:    5 * time.Second,
    // WriteTimeout:   5 * time.Second,
    MaxHeaderBytes: 1 << 20,
  }
  srv.SetKeepAlivesEnabled(true)
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}

func (s *Server) statusHandler(w http.ResponseWriter, r *http.Request) {
	var (
		res       = map[string]map[string]interface{}{}
		info      = s.hub.Info()
		withRates = (r.FormValue("rates") != "")
	)

	for queue, meta := range info {
		res[queue] = map[string]interface{}{}

		for key, val := range meta {
			res[queue][key] = val
		}
		if withRates {
			inRate, outRate := s.hub.Rates(queue)
			inHist, outHist := s.hub.RateHistory(queue)
			res[queue]["in_rate"] = inRate
			res[queue]["out_rate"] = outRate
			res[queue]["in_rate_history"] = inHist
			res[queue]["out_rate_history"] = outHist
		}
	}

	jsn, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(jsn)
}

func (s *Server) debugHandler(w http.ResponseWriter, r *http.Request) {
	info := make(map[string]interface{})
	info["version"] = Version
	info["gomaxprocs"] = runtime.GOMAXPROCS(-1)
	info["goroutines"] = runtime.NumGoroutine()
	jsn, _ := json.Marshal(info)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
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
  counts.Add("PubCount", 1)
 	counts.Add("PubBs", int64(len(msg)))
}

func (s *Server) subHandler(w http.ResponseWriter, r *http.Request) {
  queues_param := r.FormValue("queues")
  var queues []string
  if len(queues_param) > 0 {
	  queues = strings.Split(queues_param, ",")
  }

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

    counts.Add("SubCount", 1)
 		counts.Add("SubBs", int64(len(res.Message)))
	}
}

func (s *Server) flushHandler(w http.ResponseWriter, r *http.Request) {
	queues := strings.Split(r.FormValue("queues"), ",")
	messages := s.hub.Flush(queues)
	jsn, _ := json.Marshal(messages)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(jsn)
}

func (s *Server) dashboardHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.New("dashboard")
	tmpl, _ = tmpl.Parse(s.dashboardTmpl)

	w.Header().Set("Content-Type", "text/html; charset=utf8")
	hostname, _ := os.Hostname()
	if hostname == "" {
		hostname = "Unknown Host"
	}

	tmpl.ExecuteTemplate(w, "dashboard", map[string]interface{}{
		"version":  Version,
		"hostname": hostname,
		"port":     s.port,
	})
}
