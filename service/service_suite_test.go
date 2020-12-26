package service_test

import (
	"testing"

	"github.com/TaylorOno/golandreporter"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestInternal(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecsWithCustomReporters(t, "Service Suite", []Reporter{golandreporter.NewAutoGolandReporter()})
}
