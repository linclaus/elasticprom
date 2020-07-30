package model

type ElasticMetric struct {
	Keyword    string
	StrategyId string
	Count      float64
}

type StrategyMetic struct {
	ElasticMetric
	quit chan bool
}

type StrategyMetricMap map[string]StrategyMetic
