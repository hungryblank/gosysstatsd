package disk_usage

import (
	"log"
	"os/exec"
	"strings"
	"regexp"
	"strconv"
	"../statsd"
)

type Usage struct {
	device string
	mountPoint string
	total_blocks int
	used_blocks int
	total_inodes int
	used_inodes int
}

func (usage Usage) BlockPct() int {
	if (usage.used_blocks == 0) {
		return 0
	} else {
		return int(float32(usage.used_blocks) / float32(usage.total_blocks) * 100)
	}
}

func (usage Usage) InodePct() int {
	if (usage.used_inodes == 0) {
		return 0
	} else {
		return int(float32(usage.used_inodes) / float32(usage.total_inodes) * 100)
	}
}

func (usage Usage) ToMetrics() *[]statsd.Metric {
	list := []statsd.Metric{
		statsd.Gauge("disk_usage.blocks.total." + usage.device, usage.total_blocks),
		statsd.Gauge("disk_usage.blocks.used." + usage.device, usage.used_blocks),
		statsd.Gauge("disk_usage.blocks.usagePct." + usage.device, usage.BlockPct()),
		statsd.Gauge("disk_usage.inodes.total." + usage.device, usage.total_blocks),
		statsd.Gauge("disk_usage.inodes.used." + usage.device, usage.used_blocks),
		statsd.Gauge("disk_usage.inodes.usagePct." + usage.device, usage.InodePct()),
	}
	return &list
}

func (usage Usage) AppendMetrics(list *[]statsd.Metric) {
	for _, metric := range *usage.ToMetrics() {
		*list = append(*list, metric)
	}
}

func (usage Usage) IsNormalDevice() bool {
	return regexp.MustCompile("/dev").MatchString(usage.device)
}

type DataPoint struct {
	usages []*Usage
}

func rowToUsage(row string) *Usage {
	row_tokens := regexp.MustCompile(" +").Split(row, -1)
	total_blocks, _ := strconv.Atoi(row_tokens[1])
	used_blocks, _ := strconv.Atoi(row_tokens[2])
	usage := Usage{
		row_tokens[0],
		row_tokens[5],
		total_blocks,
		used_blocks,
		0,
		0,
	}
	return &usage
}

func addInodeToUsages(usages []*Usage, out string) {
	rows := strings.Split(string(out), "\n")
	for index, row := range rows[:len(rows) - 1] {
		if (index > 0) {
			addInodeToUsage(row, usages[index - 1])
		}
	}
}

func addInodeToUsage(row string, usage *Usage) {
	row_tokens := regexp.MustCompile(" +").Split(row, -1)
	total_inodes, _ := strconv.Atoi(row_tokens[3])
	used_inodes, _ := strconv.Atoi(row_tokens[4])
	usage.total_inodes = total_inodes
	usage.used_inodes = used_inodes
}

func parseOutput(out string) []*Usage {
	rows := strings.Split(string(out), "\n")
	list := []*Usage{}
	for index, row := range rows[:len(rows) - 1] {
		if (index > 0) {
			list = append(list, rowToUsage(row))
		}
	}
	return list
}

func (point DataPoint) ToMetrics() *[]statsd.Metric {
	list := []statsd.Metric{}
	for _, usage := range point.usages {
		if (usage.IsNormalDevice()) {
			usage.AppendMetrics(&list)
		}
	}
	return &list
}

func Poll() DataPoint {
	blocksOut, err := exec.Command("df").Output()
	if err != nil {
		log.Fatal(err)
	}
	inodeOut, err := exec.Command("df", "-i").Output()
	if err != nil {
		log.Fatal(err)
	}
	usages := parseOutput(string(blocksOut))
	addInodeToUsages(usages, string(inodeOut))
	return DataPoint{usages}
}
