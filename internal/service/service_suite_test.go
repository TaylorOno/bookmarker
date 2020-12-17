package service_test

import (
	"github.com/TaylorOno/golandreporter"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestInternal(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecsWithCustomReporters(t, "Service Suite", []Reporter{golandreporter.NewAutoGolandReporter()})
}
