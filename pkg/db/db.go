package db

import "github.com/linclaus/elasticprom/pkg/model"

var (
	dateTemplate      = "2006-01-02T15:04:05"
	indexDateTemplate = "2006.01.02"
	indexPrefix       = "filebeat-6.8.3-"
)

type Storer interface {
	GetVersion() error
	GetMetric(metricChan chan<- model.ElasticMetric, sm *model.StrategyMetic)
}
