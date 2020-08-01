package metrics

import (
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
func init() {
	MyMetricGauge.Set(0)
	prometheus.MustRegister(MyMetricGauge)
	MyMetricGaugeVec.WithLabelValues("l1", "l2").Set(1)
	MyMetricGaugeVec.WithLabelValues("l2", "l3").Set(2)
	prometheus.MustRegister(MyMetricGaugeVec)
	prometheus.MustRegister(ElasticMetricCountVec)
}
