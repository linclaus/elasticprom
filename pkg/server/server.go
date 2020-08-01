package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/linclaus/elasticprom/pkg/elastic"
	"github.com/linclaus/elasticprom/pkg/metrics"
	"github.com/linclaus/elasticprom/pkg/model"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	r                *mux.Router
	elasticMetricMap map[string]*model.StrategyMetic
	elasticDB        *elastic.ElasticDB
	debug            bool
	metricChan       chan model.ElasticMetric
}

func New(debug bool, elasticAddress []string) Server {
	r := mux.NewRouter()
	elasticDB, _ := elastic.ConnectES(elasticAddress)
	s := Server{
		debug:            debug,
		r:                r,
		elasticDB:        elasticDB,
		elasticMetricMap: make(map[string]*model.StrategyMetic),
		metricChan:       make(chan model.ElasticMetric),
	}
	r.Handle("/metrics", promhttp.Handler())
	r.HandleFunc("/add_metric", s.AddStrategyMetric).Methods("POST")
	r.HandleFunc("/delete_metric", s.DeleteStrategyMetric).Methods("DELETE")
	return s
}

// Start starts a new server on the given address
func (s Server) Start(address string) {
	go InitMetric()
	log.Println("Starting listener on", address)
	log.Fatal(http.ListenAndServe(address, s.r))
}

func InitMetric() {
	metricChan := make(chan model.ElasticMetric, 1024)
	for em := range metricChan {
		metrics.MyMetricGauge.Inc()
		metrics.MyMetricGaugeVec.WithLabelValues("l1", "l2").Inc()
		metrics.MyMetricGaugeVec.WithLabelValues("l2", "l3").Inc()
		metrics.ElasticMetricCountVec.WithLabelValues(em.Keyword, em.StrategyId).Set(em.Count)
	}
}

func (s Server) AddStrategyMetric(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to read payload: %s\n", err)
		http.Error(w, fmt.Sprintf("Failed to read payload: %s", err), http.StatusBadRequest)
		return
	}

	if s.debug {
		log.Println("Received webhook payload", string(body))
	}
	sm := s.elasticMetricMap["123"]
	if sm == nil {
		sm = &model.StrategyMetic{
			StrategyId:   "123",
			Container:    "gotest",
			Keyword:      "hello",
			TickInterval: 5 * time.Second,
			ESDuration:   2 * time.Hour,
			Quit:         make(chan struct{}),
		}
		s.elasticMetricMap["123"] = sm
	}
	go s.elasticDB.GetMetric(s.metricChan, sm)
}

func (s Server) DeleteStrategyMetric(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to read payload: %s\n", err)
		http.Error(w, fmt.Sprintf("Failed to read payload: %s", err), http.StatusBadRequest)
		return
	}

	if s.debug {
		log.Println("Received webhook payload", string(body))
	}
	sm := s.elasticMetricMap["123"]
	if sm != nil {
		close(sm.Quit)
		delete(s.elasticMetricMap, "123")
	}
}
