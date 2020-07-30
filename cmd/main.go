package main

import (
	"flag"
	"os"
	"strings"

	"github.com/linclaus/elasticprom/pkg/elastic"
	"github.com/linclaus/elasticprom/pkg/metrics"
	"github.com/linclaus/elasticprom/pkg/model"
	"github.com/linclaus/elasticprom/pkg/server"
)

type Args struct {
	Addr             string
	Debug            bool
	ElasticsearchUrl string
}

func main() {
	args := Args{
		ElasticsearchUrl: os.Getenv("ELASTICSEARCH-URL"),
	}
	flag.StringVar(&args.Addr, "listen-address", ":8080", "The address to listen on for HTTP requests.")
	flag.BoolVar(&args.Debug, "debug", true, "debug or not.")

	flag.Parse()

	metricChan := make(chan model.ElasticMetric, 1024)

	elasticUrls := strings.Split(args.ElasticsearchUrl, ",")
	go elastic.Init(metricChan, elasticUrls)
	go metrics.Init(metricChan)

	s := server.New(args.Debug)
	s.Start(args.Addr)
}
