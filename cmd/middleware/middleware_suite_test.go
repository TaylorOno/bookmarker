package middleware_test

import (
	"testing"

	"github.com/TaylorOno/golandreporter"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMiddleware(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecsWithCustomReporters(t, "Middleware Suite", []Reporter{golandreporter.NewAutoGolandReporter()})
}
