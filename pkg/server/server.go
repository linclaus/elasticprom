package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/linclaus/elasticprom/pkg/elastic"
	"github.com/linclaus/elasticprom/pkg/model"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	r                *mux.Router
	elasticMetricMap *model.StrategyMetricMap
	debug            bool
}

func New(debug bool) Server {
	r := mux.NewRouter()
	s := Server{
		debug: debug,
		r:     r,
		elasticMetricMap:make(map[string]model.StrategyMetricMap)
	}
	r.Handle("/metrics", promhttp.Handler())

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
	sm:=&StrategyMetric{
		StrategyId   string
		container    string
		Keyword      string
		tickInterval time.Duration
		esDuration   time.Duration
		quit         chan struct{}
	}

	elastic.AddMetric(sm)
}
