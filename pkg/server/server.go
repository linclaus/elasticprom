package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/linclaus/elasticprom/pkg/elastic"
	"github.com/linclaus/elasticprom/pkg/model"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	r                *mux.Router
	ElasticMetricMap map[string]*model.StrategyMetic
	debug            bool
}

func New(debug bool) Server {
	r := mux.NewRouter()
	s := Server{
		debug:            debug,
		r:                r,
		ElasticMetricMap: make(map[string]*model.StrategyMetic),
	}
	r.Handle("/metrics", promhttp.Handler())
	r.HandleFunc("/add_metric", s.AddStrategyMetric).Methods("POST")
	r.HandleFunc("/delete_metric", s.DeleteStrategyMetric).Methods("DELETE")

	return s
}

// Start starts a new server on the given address
func (s Server) Start(address string) {
	log.Println("Starting listener on", address)
	log.Fatal(http.ListenAndServe(address, s.r))
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
	sm := s.ElasticMetricMap["123"]
	if sm == nil {
		sm = &model.StrategyMetic{
			StrategyId:   "123",
			Container:    "gotest",
			Keyword:      "hello",
			TickInterval: 5 * time.Second,
			ESDuration:   2 * time.Hour,
			Quit:         make(chan struct{}),
		}
		s.ElasticMetricMap["123"] = sm
	}
	go elastic.AddMetric(sm)
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
	sm := s.ElasticMetricMap["123"]
	if sm != nil {
		close(sm.Quit)
		delete(s.ElasticMetricMap, "123")
	}
}
