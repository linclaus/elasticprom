package model

import "time"

type ElasticMetric struct {
	Keyword    string
	StrategyId string
	Count      float64
}

type StrategyMetic struct {
	StrategyId   string
	Container    string
	Keyword      string
	TickInterval time.Duration
	ESDuration   time.Duration
	Quit         chan struct{}
}
