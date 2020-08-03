package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/linclaus/elasticprom/pkg/db"
	"github.com/linclaus/elasticprom/pkg/metrics"
	"github.com/linclaus/elasticprom/pkg/model"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	r                *mux.Router
	elasticMetricMap *model.ElasticMetricMap
	db               db.Storer
	debug            bool
	metricChan       chan model.ElasticMetric
}

func New(debug bool, db db.Storer) Server {
	r := mux.NewRouter()
	s := Server{
		debug:            debug,
		r:                r,
		db:               db,
		elasticMetricMap: &model.ElasticMetricMap{},
		metricChan:       make(chan model.ElasticMetric, 1024),
	}
	r.Handle("/metrics", s.metricHandler(promhttp.Handler()))
	r.HandleFunc("/add_metric", s.handleFuncInterceptor(s.AddStrategyMetric)).Methods("POST")
	r.HandleFunc("/delete_metric/{id}", s.DeleteStrategyMetric).Methods("DELETE")
	return s
}

// Start starts a new server on the given address
func (s Server) Start(address string) {
	go s.initMetric()
	log.Println("Starting listener on", address)
	log.Fatal(http.ListenAndServe(address, s.r))
}

func (s Server) initMetric() {
	for em := range s.metricChan {
		metrics.MyMetricGauge.Inc()
		metrics.MyMetricGaugeVec.WithLabelValues("l1", "l2").Inc()
		metrics.MyMetricGaugeVec.WithLabelValues("l2", "l3").Inc()
		metrics.ElasticMetricCountVec.WithLabelValues(em.Keyword, em.StrategyId).Set(em.Count)
	}
}

func (s Server) metricHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("before metric handler")
		next.ServeHTTP(w, r)
		log.Println("after metric handler")
	})

}

func (s Server) handleFuncInterceptor(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("before handlerFunc")
		h(w, r)
		log.Println("after handlerFunc")
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
	smr := &model.StrategyMetricRequest{}
	json.Unmarshal([]byte(body), smr)
	sm := s.elasticMetricMap.Get(smr.StrategyId)
	if sm == nil {
		sm = &model.StrategyMetic{
			StrategyId:   smr.StrategyId,
			Container:    smr.Container,
			Keyword:      smr.Keyword,
			TickInterval: time.Duration(smr.TickInterval) * time.Second,
			ESDuration:   time.Duration(smr.ESDuration) * time.Hour,
			Quit:         make(chan struct{}),
		}
		s.elasticMetricMap.Set(smr.StrategyId, sm)
	}
	go s.db.GetMetric(s.metricChan, sm)
}

func (s Server) DeleteStrategyMetric(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sm := s.elasticMetricMap.Get(vars["id"])
	if sm != nil {
		close(sm.Quit)
		s.elasticMetricMap.Delete(sm.StrategyId)
	}
}
