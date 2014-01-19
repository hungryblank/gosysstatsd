package memory

import (
	"log"
	"os/exec"
	"strings"
	"regexp"
	"strconv"
	"../statsd"
)

type DataPoint struct {
	total int
	used int
	free int
	shared int
	buffers int
	cached int
}

func Poll() DataPoint {
	out, err := exec.Command("free", "-mt").Output()
	if err != nil {
		log.Fatal(err)
	}
	split := strings.Split(string(out), "\n")
	mem := split[1]
	mem_details := regexp.MustCompile(" +").Split(mem, -1)[1:]
	total, _ := strconv.Atoi(mem_details[0])
	used, _ := strconv.Atoi(mem_details[1])
	free, _ := strconv.Atoi(mem_details[2])
	shared, _ := strconv.Atoi(mem_details[3])
	buffers, _ := strconv.Atoi(mem_details[4])
	cached, _ := strconv.Atoi(mem_details[5])
	point := DataPoint{
		total,
		used,
		free,
		shared,
		buffers,
		cached,
	}
	return point
}

func (point DataPoint) Available() int {
	return point.free + point.buffers + point.cached
}

func (point DataPoint) UsagePct() int {
	return 100 - int(float32(point.Available()) / float32(point.total) * 100.0)
}

func (point DataPoint) ToMetrics() *[]statsd.Metric {
	metrics := []statsd.Metric{
		statsd.Gauge("system.memory.total", point.total),
		statsd.Gauge("system.memory.used", point.used),
		statsd.Gauge("system.memory.free", point.free),
		statsd.Gauge("system.memory.shared", point.shared),
		statsd.Gauge("system.memory.buffers", point.buffers),
		statsd.Gauge("system.memory.cached", point.cached),
		statsd.Gauge("system.memory.available", point.Available()),
		statsd.Gauge("system.memory.usagePct", point.UsagePct()),
	}
	return &metrics
}
