package db

import (
	"log"
	"time"

	"github.com/linclaus/elasticprom/pkg/model"
)

type NullDB struct{}

func (db NullDB) GetVersion() error {
	log.Println("this is null db")
	return nil
}

func (db NullDB) GetMetric(metricChan chan<- model.ElasticMetric, sm *model.StrategyMetic) {
	tick := time.NewTicker(sm.TickInterval)
	defer tick.Stop()
LOOP:
	for {
		select {
		case <-tick.C:
			count := float64(111)
			log.Printf("count : %f", count)
			em := model.ElasticMetric{
				Keyword:    sm.Keyword,
				StrategyId: sm.StrategyId,
				Count:      count,
			}
			select {
			case metricChan <- em:
				log.Println("send message successful")
			default:
				log.Println("send message timeout")
			}
		case <-sm.Quit:
			log.Println("stop strategy: %s", sm.StrategyId)
			break LOOP
		}
	}
}
