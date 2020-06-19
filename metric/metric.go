package metric

import (
	gometrics "github.com/rcrowley/go-metrics"
	"github.com/rcrowley/go-metrics/exp"
	datadog "gopkg.in/zorkian/go-datadog-api.v2"

	"github.com/sandbreaker/goservice/config"
	"github.com/sandbreaker/goservice/lib/timelib"
	"github.com/sandbreaker/goservice/log"
)

type Metric struct {
	Registry gometrics.Registry

	// note that this is not a default agent based ddog client.
	// however, this allows unlimited datadog api metrics for free version
	DDogClient *datadog.Client
}

var (
	StaticClient        *Metric
	DdogMetricTypeGauge = "gauge"
	DdogMetricTypeRate  = "rate"
	DdogMetricTypeCount = "count"
)

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

	var ddogAPI *datadog.Client

	if config.StaticConfig.GetBool("metric.datadog.is_enable", false) {
		ddogAPI = datadog.NewClient(
			config.StaticConfig.GetString("datadog.api.key", ""),
			config.StaticConfig.GetString("datadog.app.key", ""),
		)

	}

	return &Metric{
		DDogClient: ddogAPI,
		Registry:   gometrics.DefaultRegistry,
	}, nil
}

func (m *Metric) RegisterGCDebug() {
	gometrics.CaptureDebugGCStatsOnce(m.Registry)
}

func (m *Metric) postMetrics(metricSeries []datadog.Metric) {
	if err := m.DDogClient.PostMetrics(metricSeries); err != nil {
		log.Error("Metric post problem: ", err)
	}
}

func (m *Metric) Guage(name string, value float64, tags []string) {
	registryVal := int64(value)
	gometrics.GetOrRegisterGauge(name, m.Registry).Update(registryVal)

	if m.DDogClient != nil {
		now := timelib.GetEpoch()
		val := float64(value)
		dp := [2]*float64{&now, &val}
		host := config.StaticConfig.GetRole()
		tags = append(tags, "env:"+config.StaticConfig.GetEnv())

		ddogMetric := datadog.Metric{
			Metric: &name,
			Type:   &DdogMetricTypeGauge,
			Points: []datadog.DataPoint{dp},
			Tags:   tags,
			Host:   &host,
		}
		metricSeries := []datadog.Metric{ddogMetric}

		go m.postMetrics(metricSeries)
	}
}

// ddog api doesn't have histogram API
func (m *Metric) Histogram(name string, value float64, tags []string) {
	registryVal := int64(value)

	h := gometrics.GetOrRegisterHistogram(
		name,
		m.Registry,
		gometrics.NewExpDecaySample(1028, 0.015),
	)
	h.Update(registryVal)
}

func (m *Metric) Inc(name string, value int64, tags []string) {
	gometrics.GetOrRegisterCounter(name, m.Registry).Inc(value)

	if m.DDogClient != nil {
		now := timelib.GetEpoch()
		val := float64(value)
		dp := [2]*float64{&now, &val}
		host := config.StaticConfig.GetRole()
		tags = append(tags, "env:"+config.StaticConfig.GetEnv())

		ddogMetric := datadog.Metric{
			Metric: &name,
			Type:   &DdogMetricTypeCount,
			Points: []datadog.DataPoint{dp},
			Tags:   tags,
			Host:   &host,
		}
		metricSeries := []datadog.Metric{ddogMetric}

		go m.postMetrics(metricSeries)
	}
}
