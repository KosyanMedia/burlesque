package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	publishEndpoint   = "/publish"
	subscribeEndpoint = "/subscribe"
	flushEndpoint     = "/flush"
	statusEndpoit     = "/status"
	debugEndpoint     = "/debug"
)

type (
	Config struct {
		Host    string
		Port    int
		Timeout time.Duration
	}
	Client struct {
		Config     *Config
		httpClient *http.Client
	}
	Message struct {
		Queue string
		Body  []byte
	}
	QueueInfo struct {
		Name        string
		Messages    int
		Subscribers int
	}
	DebugInfo struct {
		Version      string                 `json:"version"`
		Gomaxprocs   int                    `json:"gomaxprocs"`
		Goroutines   int                    `json:"goroutines"`
		KyotoCabinet map[string]interface{} `json:"kyoto_cabinet"`
	}
)

func NewConfig() *Config {
	return &Config{
		Host:    "127.0.0.1",
		Port:    4401,
		Timeout: 60 * time.Second,
	}
}

func NewClient(c *Config) *Client {
	if c == nil {
		c = NewConfig()
	}

	transport := http.Transport{
		Dial: func(network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, c.Timeout)
		},
	}

	client := http.Client{
		Transport: &transport,
	}

	return &Client{
		Config:     c,
		httpClient: &client,
	}
}

func (c *Client) Publish(m *Message) bool {
	rdr := bytes.NewReader(m.Body)
	url := c.url(publishEndpoint, "?queue=", m.Queue)

	res, err := http.Post(url, "text/plain", rdr)
	if err != nil {
		return false
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return false
	}

	return (string(body) == "OK")
}

func (c *Client) Subscribe(queues ...string) *Message {
	url := c.url(subscribeEndpoint, "?queues=", strings.Join(queues, ","))
	head, body := c.get(url)

	return &Message{
		Queue: head.Get("Queue"),
		Body:  body,
	}
}

func (c *Client) Flush(queues ...string) (messages []*Message) {
	url := c.url(flushEndpoint, "?queues=", strings.Join(queues, ","))
	_, body := c.get(url)

	var tmp []map[string]string
	if err := json.Unmarshal(body, &tmp); err != nil {
		return
	}

	messages = []*Message{}
	for _, msg := range tmp {
		messages = append(messages, &Message{msg["queue"], []byte(msg["message"])})
	}

	return
}

func (c *Client) Status() (stat []*QueueInfo) {
	url := c.url(statusEndpoit)
	_, body := c.get(url)

	tmp := make(map[string]map[string]int)
	if err := json.Unmarshal(body, &tmp); err != nil {
		return
	}

	for queue, info := range tmp {
		qi := &QueueInfo{
			Name:        queue,
			Messages:    info["messages"],
			Subscribers: info["subscribers"],
		}
		stat = append(stat, qi)
	}

	return
}

func (c *Client) Debug() *DebugInfo {
	url := c.url(debugEndpoint)
	_, body := c.get(url)

	var dbg DebugInfo
	if err := json.Unmarshal(body, &dbg); err != nil {
		return nil
	}

	return &dbg
}

func (c *Client) get(url string) (http.Header, []byte) {
	res, err := http.Get(url)
	if err != nil {
		return nil, nil
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, nil
	}
	res.Body.Close()

	return res.Header, body
}

func (c *Client) url(path ...string) string {
	parts := []string{"http://", c.Config.Host, ":", strconv.Itoa(c.Config.Port)}
	parts = append(parts, path...)

	return strings.Join(parts, "")
}
