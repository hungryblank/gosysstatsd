package statsd

import (
	"net"
	"fmt"
	"log"
)

type StatsdClient struct {
	Host string
	Port int
	conn net.Conn
}

/**
 * Method to open udp connection, called by default client factory
 **/
func (client *StatsdClient) Open() {
	connectionString := fmt.Sprintf("%s:%d", client.Host, client.Port)
	conn, err := net.Dial("udp", connectionString)
	if err != nil {
		log.Println(err)
	}
	client.conn = conn
}

/**
 * Method to close udp connection
 **/
func (client *StatsdClient) Close() {
	client.conn.Close()
}

func New(host string, port int) *StatsdClient {
	client := StatsdClient{Host: host, Port: port}
	client.Open()
	return &client
}

func (client *StatsdClient) UpdateMetrics(metrics *[]Metric) {
	for _, metric := range *metrics {
		_,err := fmt.Fprintf(client.conn, metric.ToMessage())
		if err != nil {
			log.Println(err)
		}
	}
}

type Metric interface {
	ToMessage() (string)
}

type CounterMetric struct {
	label string
	value int
}

func (metric CounterMetric) ToMessage() string {
	return fmt.Sprintf("%s:%d|c", metric.label, metric.value)
}

type GaugeMetric struct {
	label string
	value int
}

func (metric GaugeMetric) ToMessage() string {
	return fmt.Sprintf("%s:%d|g", metric.label, metric.value)
}

func Counter(label string, value int) Metric {
	return Metric(CounterMetric{label, value})
}

func Gauge(label string, value int) Metric {
	return Metric(GaugeMetric{label, value})
}

func Report(client *StatsdClient, metrics *[]Metric) {
	client.UpdateMetrics(metrics)
}
