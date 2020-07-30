package model

import "time"

type ElasticMetric struct {
	Keyword    string
	StrategyId string
	Count      float64
}

type StrategyMetic struct {
	StrategyId   string
	container    string
	Keyword      string
	tickInterval time.Duration
	esDuration   time.Duration
	quit         chan struct{}
}

type StrategyMetricMap map[string]StrategyMetic
