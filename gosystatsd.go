package main

import (
	"github.com/hungryblank/gosysstatsd/statsd"
	"github.com/hungryblank/gosysstatsd/disk_usage"
	"github.com/hungryblank/gosysstatsd/memory"
	"github.com/hungryblank/gosysstatsd/loadavg"
	"flag"
	"fmt"
	"os"
)

func Usage() {
        fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
        flag.PrintDefaults()
}

func main() {
	host := flag.String("h", "localhost", "statasd host")
	port := flag.Int("p", 8125, "statasd port")
	help := flag.Bool("help", false, "print this help message")
	flag.Parse()
	if (*help) {
		Usage()
	}
	client := statsd.New(*host, *port)
	dataPoint := memory.Poll()
	statsd.Report(client, dataPoint.ToMetrics())
	diskDataPoint := disk_usage.Poll()
	statsd.Report(client, diskDataPoint.ToMetrics())
	loadAvgDataPoint := loadavg.Poll()
	statsd.Report(client, loadAvgDataPoint.ToMetrics())
}
