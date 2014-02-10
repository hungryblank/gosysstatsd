package disk_usage

import (
	"log"
	"os/exec"
	"strings"
	"regexp"
	"strconv"
	"math"
	"github.com/hungryblank/gosysstatsd/statsd"
)

func Pct(available int, total int) int {
	if (available == 0) {
		return 0
	} else {
		fraction := 1 - float64(available) / float64(total)
		return int(math.Ceil(100 * fraction))
	}
}

type Usage struct {
	device string
	mountPoint string
	total_blocks int
	available_blocks int
	used_blocks int
	total_inodes int
	available_inodes int
	used_inodes int
}

func (usage Usage) BlockPct() int {
	return Pct(usage.available_blocks, usage.total_blocks)
}

func (usage Usage) InodePct() int {
	return Pct(usage.available_inodes, usage.total_inodes)
}

func (usage Usage) ToMetrics() *[]statsd.Metric {
	list := []statsd.Metric{
		statsd.Gauge("disk_usage.blocks.total." + usage.device, usage.total_blocks),
		statsd.Gauge("disk_usage.blocks.used." + usage.device, usage.used_blocks),
		statsd.Gauge("disk_usage.blocks.available." + usage.device, usage.available_blocks),
		statsd.Gauge("disk_usage.blocks.usagePct." + usage.device, usage.BlockPct()),
		statsd.Gauge("disk_usage.inodes.total." + usage.device, usage.total_inodes),
		statsd.Gauge("disk_usage.inodes.used." + usage.device, usage.used_inodes),
		statsd.Gauge("disk_usage.inodes.available." + usage.device, usage.available_inodes),
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
	rowTokens := regexp.MustCompile(" +").Split(row, -1)
	totalBlocks, _ := strconv.Atoi(rowTokens[1])
	usedBlocks, _ := strconv.Atoi(rowTokens[2])
	availableBlocks, _ := strconv.Atoi(rowTokens[3])
	usage := Usage{
		rowTokens[0],
		rowTokens[5],
		totalBlocks,
		availableBlocks,
		usedBlocks,
		0,
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
	rowTokens := regexp.MustCompile(" +").Split(row, -1)
	totalInodes, _ := strconv.Atoi(rowTokens[1])
	usedInodes, _ := strconv.Atoi(rowTokens[2])
	availableInodes, _ := strconv.Atoi(rowTokens[3])
	usage.total_inodes = totalInodes
	usage.available_inodes = availableInodes
	usage.used_inodes = usedInodes
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
