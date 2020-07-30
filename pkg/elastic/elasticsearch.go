package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v6/esapi"
	"github.com/linclaus/elasticprom/pkg/model"

	es6 "github.com/elastic/go-elasticsearch/v6"
)

var (
	client            *es6.Client
	ch                chan model.ElasticMetric
	dateTemplate      = "2006-01-02T15:04:05"
	indexDateTemplate = "2006.01.02"
	indexPrefix       = "filebeat-6.8.3-"
)

// Init elastic
func Init(metricChan chan model.ElasticMetric, addresses []string) {
	ch = metricChan
	cfg := es6.Config{
		Addresses: addresses,
	}
	client, _ = es6.NewClient(cfg)

	res, err := client.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()
	log.Println(res)
	AddMetric(make(chan struct{}), 5*time.Second, 1*time.Hour, "gotest", "hello", "123")
}

//AddMetric function
func AddMetric(quit <-chan struct{}, tickInterval time.Duration, esDuration time.Duration, container string, keyword string, strategyId string) {
	tick := time.NewTicker(tickInterval)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			count := countByKeyword(esDuration, container, keyword)
			em := model.ElasticMetric{
				Keyword:    keyword,
				StrategyId: strategyId,
				Count:      count,
			}
			select {
			case ch <- em:
			default:
				log.Println("send message timeout")
			}
		}
	}
}

func countByKeyword(d time.Duration, container string, keyword string) float64 {
	now := time.Now().UTC()
	from := now.Add(-1 * d).UTC()
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					map[string]interface{}{
						"term": map[string]interface{}{
							"kubernetes.container.name": container,
						}},
					map[string]interface{}{"match_phrase": map[string]interface{}{
						"message": keyword,
					}},
					map[string]interface{}{"range": map[string]interface{}{
						"@timestamp": map[string]interface{}{
							"gt": from.Format(dateTemplate),
							"lt": now.Format(dateTemplate),
						}},
					},
				},
			},
		},
	}
	jsonBody, _ := json.Marshal(query)

	log.Printf("jsonBody: %s", jsonBody)

	req := esapi.CountRequest{
		Index:        []string{strings.Join([]string{indexPrefix, from.Format(indexDateTemplate)}, ""), strings.Join([]string{indexPrefix, now.Format(indexDateTemplate)}, "")},
		DocumentType: []string{"doc"},
		Body:         bytes.NewReader(jsonBody),
	}
	res, err := req.Do(context.Background(), client)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	var jsonData map[string]interface{}
	if res.StatusCode == 200 {
		jsonResp, _ := ioutil.ReadAll(res.Body)
		json.Unmarshal([]byte(jsonResp), &jsonData)
		count, _ := jsonData["count"].(float64)
		log.Printf("count : %f", count)
		return count
	}
	log.Println(res.String())
	return 0
}
