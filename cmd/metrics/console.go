package metrics

import (
	"fmt"
	"log"
	"strings"
)

type console struct{}

func NewConsoleReporter() *console {
	return &console{}
}

func (c *console) ObserverHistogram(name string, value float64, labels ...string) {
	writeLog(name, value, labels)
}

func (c *console) ObserverSummary(name string, value float64, labels ...string) {
	writeLog(name, value, labels)
}

func (c *console) ObserverCount(name string, value float64, labels ...string) {
	writeLog(name, value, labels)
}

func writeLog(name string, value float64, labels []string) {
	metric := fmt.Sprintf("%v \t %v \t", name, value)
	labelString := strings.Join(labels, "\t")
	log.Print(metric + labelString)
}
