package repository_test

import (
	"testing"

	"github.com/TaylorOno/golandreporter"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRepository(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecsWithCustomReporters(t, "Repository Suite", []Reporter{golandreporter.NewAutoGolandReporter()})
}
