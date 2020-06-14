package metric

import (
	"github.com/sandbreaker/goservice/log"

	gometrics "github.com/rcrowley/go-metrics"
	"github.com/rcrowley/go-metrics/exp"
)

// TODO add datadog

type Metric struct {
	Registry gometrics.Registry
}

var StaticClient *Metric

func init() {
	var err error
	StaticClient, err = NewMetric("")
	if err != nil {
		log.Error(err)
	}
}

func NewMetric(appName string) (*Metric, error) {
	// enable expvar
	exp.Exp(gometrics.DefaultRegistry)

	return &Metric{}, nil
}

func (m *Metric) RegisterGCDebug() {
	gometrics.CaptureDebugGCStatsOnce(m.Registry)
}

func (m *Metric) Inc(name string, value int64, tags []string) {
	gometrics.GetOrRegisterCounter(name, m.Registry).Inc(value)
}

func (m *Metric) Guage(name string, value float64, tags []string) {
	registryVal := int64(value)
	gometrics.GetOrRegisterGauge(name, m.Registry).Update(registryVal)
}

func (m *Metric) Histogram(name string, value float64, tags []string) {
	registryVal := int64(value)

	h := gometrics.GetOrRegisterHistogram(
		name,
		m.Registry,
		gometrics.NewExpDecaySample(1028, 0.015),
	)
	h.Update(registryVal)
}
