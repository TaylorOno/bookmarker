package routes_test

import (
	"bytes"
	"fmt"
	"github.com/TaylorOno/bookmarker/cmd/routes"
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

var _ = Describe("Router", func() {
	Context("SetRoutes", func() {
		It("Returns a router", func() {
			server := routes.Server{}
			router := server.SetRoutes()
			Expect(fmt.Sprintf("%T", router)).To(Equal("*mux.Router"))
		})
	})
})

func bodyFromFile(s string) io.Reader {
	body, err := ioutil.ReadFile(fmt.Sprintf("test_data/%v", s))
	if err != nil {
		Fail(fmt.Sprintf("failed to read file %v: %v", s, err.Error()))
	}
	return bytes.NewReader(body)
}
