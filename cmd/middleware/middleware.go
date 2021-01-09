//go:generate mockgen -destination=../../test/mocks/mock_middleware.go -package=mocks -source middleware.go

package middleware

type reporter interface {
	ObserverHistogram(name string, value float64, labels ...string)
	ObserverSummary(name string, value float64, labels ...string)
}
