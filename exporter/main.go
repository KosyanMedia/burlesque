package main

import (
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	namespace = "burlesque"
)

var (
	queue_label = []string{"queue"}
)

type Exporter struct {
	url     string
	timeout time.Duration

	up            *prometheus.Desc
	queues        *prometheus.Desc
	messages      *prometheus.Desc
	subscriptions *prometheus.Desc
}

type Status map[string]Queue

type Queue struct {
	Messages      int `json:"messages"`
	Subscriptions int `json:"subscriptions"`
}

func NewExporter(url string, timeout time.Duration) *Exporter {
	return &Exporter{
		url:     url,
		timeout: timeout,

		up: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "up"),
			"Could the burlesque server be reached.",
			nil,
			nil,
		),
		queues: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "queues"),
			"Burlesque queues count.",
			nil,
			nil,
		),
		messages: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "messages"),
			"Burlesque queue messages.",
			queue_label,
			nil,
		),
		subscriptions: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "subscriptions"),
			"Burlesque queue subscriptions.",
			queue_label,
			nil,
		),
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.up
	ch <- e.queues
	ch <- e.messages
	ch <- e.subscriptions
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	data := Status{}

	client := http.Client{
		Timeout: e.timeout,
	}
	res, err := client.Get(e.url)
	if err != nil {
		ch <- prometheus.MustNewConstMetric(e.up, prometheus.GaugeValue, 0)
		log.Error(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	ch <- prometheus.MustNewConstMetric(e.up, prometheus.GaugeValue, 1)

	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Error(err)
	}

	ch <- prometheus.MustNewConstMetric(e.queues, prometheus.CounterValue, float64(len(data)))
	for queue, info := range data {
		ch <- prometheus.MustNewConstMetric(e.messages, prometheus.CounterValue, float64(info.Messages), queue)
		ch <- prometheus.MustNewConstMetric(e.subscriptions, prometheus.CounterValue, float64(info.Subscriptions), queue)
	}
}

func main() {
	var (
		url           = kingpin.Flag("url", "Burlesque server url.").Default("http://localhost:4401/status").String()
		timeout       = kingpin.Flag("timeout", "Burlesque server connect timeout.").Default("30s").Duration()
		listenAddress = kingpin.Flag("web.listen-address", "Address to listen on for web interface and telemetry.").Default(":9118").String()
		metricsPath   = kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").String()
	)

	log.AddFlags(kingpin.CommandLine)
	kingpin.Version(version.Print("burlesque_exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	log.Infoln("Starting burlesque_exporter", version.Info())
	log.Infoln("Build context", version.BuildContext())

	prometheus.MustRegister(NewExporter(*url, *timeout))

	http.Handle(*metricsPath, prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
      <head><title>Burlesque Exporter</title></head>
      <body>
      <h1>Burlesque Exporter</h1>
      <p><a href='` + *metricsPath + `'>Metrics</a></p>
      </body>
      </html>`))
	})
	log.Infoln("Starting HTTP server on", *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
