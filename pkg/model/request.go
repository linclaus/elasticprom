package model

type StrategyMetricRequest struct {
	StrategyId   string `json:strategyId`
	Container    string `json:container`
	Keyword      string `json:keyword`
	TickInterval int64  `json:tickInterval`
	ESDuration   int64  `json:esDuration`
}
