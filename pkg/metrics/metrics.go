package metrics

import (
	"github.com/linclaus/elasticprom/pkg/model"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	MyMetricGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "my_metric_gauge",
		Help: "metric test",
	})
	MyMetricGaugeVec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "my_metric_gauge_vec",
		Help: "metric_vec test",
	}, []string{"label1", "label2"})
	ElasticMetricCountVec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "elastic_metric_gauge_vec",
		Help: "elastic count",
	}, []string{"keyword", "strategy_id"})
)

//Init metric
func Init(metricChan chan model.ElasticMetric) {
	// Register the summary and the histogram with Prometheus's default registry.
	MyMetricGauge.Set(0)
	prometheus.MustRegister(MyMetricGauge)
	MyMetricGaugeVec.WithLabelValues("l1", "l2").Set(1)
	MyMetricGaugeVec.WithLabelValues("l2", "l3").Set(2)
	prometheus.MustRegister(MyMetricGaugeVec)
	prometheus.MustRegister(ElasticMetricCountVec)
	// Add Go module build info.
	// prometheus.MustRegister(prometheus.NewBuildInfoCollector())

	for em := range metricChan {
		MyMetricGauge.Inc()
		MyMetricGaugeVec.WithLabelValues("l1", "l2").Inc()
		MyMetricGaugeVec.WithLabelValues("l2", "l3").Inc()
		ElasticMetricCountVec.WithLabelValues(em.Keyword, em.StrategyId).Set(em.Count)
	}
}
