package model

type ElasticMetric struct {
	Keyword    string
	StrategyId string
	Count      float64
}

type StrategyMetic struct {
	StrategyId      string
	Keyword         string
	intervalMinutes int
	quit            chan struct{}
}

type StrategyMetricMap map[string]StrategyMetic
