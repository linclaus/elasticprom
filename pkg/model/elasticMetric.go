package model

import (
	"sync"
	"time"
)

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

type ElasticMetricMap struct {
	elasticMetricMap map[string]*StrategyMetic
	lock             sync.RWMutex
}

func (m ElasticMetricMap) Get(k string) *StrategyMetic {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if v, exit := m.elasticMetricMap[k]; exit {
		return v
	}
	return nil
}

func (m *ElasticMetricMap) Set(k string, v *StrategyMetic) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.elasticMetricMap == nil {
		m.elasticMetricMap = make(map[string]*StrategyMetic)
	}
	m.elasticMetricMap[k] = v
}

func (m *ElasticMetricMap) Delete(k string) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.elasticMetricMap == nil {
		return
	}
	delete(m.elasticMetricMap, k)
}
