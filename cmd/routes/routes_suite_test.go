package routes_test

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"testing"

	"github.com/TaylorOno/golandreporter"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRoutes(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecsWithCustomReporters(t, "Routes Suite", []Reporter{golandreporter.NewAutoGolandReporter()})
}

func bodyFromFile(s string) io.Reader {
	body, err := ioutil.ReadFile(fmt.Sprintf("test_data/%v", s))
	if err != nil {
		Fail(fmt.Sprintf("failed to read file %v: %v", s, err.Error()))
	}
	return bytes.NewReader(body)
}
