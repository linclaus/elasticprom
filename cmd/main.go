package main

import (
	"flag"
	"os"
	"strings"

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

	elasticUrls := strings.Split(args.ElasticsearchUrl, ",")
	s := server.New(args.Debug, elasticUrls)
	s.Start(args.Addr)
}
