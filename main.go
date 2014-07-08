package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ezotrank/cabinetgo"
	"github.com/stvp/rollbar"
	"io/ioutil"
	logpkg "log"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type (
	Message []byte
	Key     []byte
	Counter struct {
		Write  uint
		Read   uint
		stream chan uint
		inLoop bool
	}
	Queue struct {
		Name    string
		Counter *Counter
	}
	Request struct {
		Queues   []string
		Callback func(*Response)
	}
	Response struct {
		Queue   string
		Message Message
	}
	Payload struct {
		Queue   *Queue
		Message Message
	}
	QueueState  map[string]uint
	ServerState map[string]QueueState
)

const (
	MaxIndex          = ^uint(0)
	StateMetaKey      = "state"
	StateSaveInterval = 10
)

//
// Key
//

func NewKey(queue string, index uint) Key {
	istr := strconv.FormatUint(uint64(index), 10)
	key := strings.Join([]string{queue, istr}, "_")
	return Key(key)
}

//
// Counter
//

func NewCounter(wi, ri uint) *Counter {
	c := &Counter{Write: wi, Read: ri}
	c.stream = make(chan uint)
	go c.Loop()
	return c
}

func (c *Counter) Incr() {
	c.Write++
	if !c.inLoop {
		c.inLoop = true
		go c.Loop()
	}
}

func (c *Counter) Next() uint {
	return <-c.stream
}

func (c *Counter) Distance() uint {
	return c.Write - c.Read
}

func (c *Counter) Loop() {
	for c.Write > c.Read {
		c.stream <- c.Read + 1
		c.Read++
	}
	c.inLoop = false
}

//
// Queue
//

func PersistMessages() {
	for {
		select {
		case payload := <-saver:
			i := payload.Queue.Counter.Write + 1
			key := NewKey(payload.Queue.Name, i)

			if err := storage.Set(key, payload.Message); err != nil {
				rollbar.Error("error", err)
			} else {
				payload.Queue.Counter.Incr()
			}
		}
	}
}

func (q *Queue) Push(msg Message) {
	saver <- Payload{Queue: q, Message: msg}
}

func (q *Queue) Fetch() (Message, error) {
	i := q.Counter.Next()
	key := NewKey(q.Name, i)

	msg, err := storage.Get(key)
	if err != nil {
		rollbar.Error("error", err)
		return msg, err
	}

	defer func() {
		if err := storage.Remove(key); err != nil {
			rollbar.Error("error", err)
		}
	}()

	return msg, nil
}

func (q *Queue) Size() uint {
	size := q.Counter.Distance()
	if size < 0 {
		size += MaxIndex
	}
	return size
}

func GetQueue(name string) *Queue {
	if _, ok := queues[name]; !ok {
		RegisterQueue(name, 0, 0)
	}
	return queues[name]
}

func RegisterQueue(name string, wi, ri uint) {
	queues[name] = &Queue{Name: name, Counter: NewCounter(wi, ri)}
}

//
// Request
//

func Register(q string, msg Message) {
	for i, r := range pool {
		for _, queueName := range r.Queues {
			if queueName == q {
				go r.Callback(&Response{Queue: queueName, Message: msg})
				pool = append(pool[:i], pool[i+1:]...)
				return
			}
		}
	}
	GetQueue(q).Push(msg)
}

func Process(r *Request) {
	for _, queueName := range r.Queues {
		q := GetQueue(queueName)
		if q.Size() > 0 {
			if msg, err := q.Fetch(); err != nil {
				go r.Callback(nil)
			} else {
				go r.Callback(&Response{Queue: queueName, Message: msg})
			}
			return
		}
	}
	pool = append(pool, r)
}

func Purge(r *Request) {
	for i, req := range pool {
		if reflect.ValueOf(r).Pointer() == reflect.ValueOf(req).Pointer() {
			pool = append(pool[:i], pool[i+1:]...)
			return
		}
	}
}

//
// State
//

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

//
// HTTP handlers
//

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	info := make(map[string]map[string]uint)
	for _, q := range queues {
		info[q.Name] = map[string]uint{
			"messages":      q.Size(),
			"subscriptions": 0,
		}
	}
	for _, r := range pool {
		for _, q := range r.Queues {
			info[q]["subscriptions"]++
		}
	}
	infoJson, _ := json.Marshal(info)
	fmt.Fprintf(w, string(infoJson))
}

func DebugHandler(w http.ResponseWriter, r *http.Request) {
	info := make(map[string]int)
	info["goroutines"] = runtime.NumGoroutine()
	infoJson, _ := json.Marshal(info)
	fmt.Fprintf(w, string(infoJson))
}

func PublishHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	msg, _ := ioutil.ReadAll(r.Body)
	if len(msg) == 0 {
		msg = Message(r.FormValue("msg"))
	}
	queueName := r.FormValue("queue")

	go Register(queueName, msg)

	log.Println("Published message of", len(msg), "bytes to queue", queueName)
	fmt.Fprintf(w, "OK")
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
	fmt.Fprintf(w, string(res.Message))

	log.Println("Recieved message of", len(res.Message), "bytes from queue", res.Queue)
	finished <- true
}

//
// main
//

var (
	log     *logpkg.Logger
	storage = cabinet.New()
	queues  = make(map[string]*Queue)
	pool    = []*Request{}
	saver   = make(chan Payload, 1000)
)

func main() {
	log = logpkg.New(os.Stdout, "", logpkg.Ldate|logpkg.Lmicroseconds)

	storagep := flag.String("storage", "-", "Kyoto Cabinet storage path (e.g. storage.kch#zcomp=gz#capsiz=524288000)")
	env := flag.String("environment", "development", "Process environment: development or production")
	port := flag.Int("port", 4401, "HTTP port to listen")
	flag.Parse()

	rollbar.Token = "c91028beb8434b66882f59f55f22659d" // klit access token
	rollbar.Environment = *env

	// Init storage
	err := storage.Open(*storagep, cabinet.KCOWRITER|cabinet.KCOCREATE)
	if err != nil {
		panic(err)
	}

	// Handle SIGTERM
	shutdown := make(chan os.Signal)
	signal.Notify(shutdown, os.Interrupt, os.Kill, syscall.SIGTERM)
	go func() {
		<-shutdown
		SaveState()
		log.Printf("State successfully persisted")
		storage.Close()
		rollbar.Wait()
		log.Println("Storage closed")
		log.Printf("Server stopped")
		os.Exit(0)
	}()

	LoadState()

	go KeepStatePersisted()
	go PersistMessages()

	log.Printf("GOMAXPROCS = %d", runtime.GOMAXPROCS(-1))
	log.Printf("Starting HTTP server on port %d", *port)

	http.HandleFunc("/status", StatusHandler)
	http.HandleFunc("/debug", DebugHandler)
	http.HandleFunc("/publish", PublishHandler)
	http.HandleFunc("/subscribe", SubscriptionHandler)
	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}
