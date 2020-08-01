package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/linclaus/elasticprom/pkg/db"
	"github.com/linclaus/elasticprom/pkg/server"
)

type Args struct {
	Addr             string
	Debug            bool
	ElasticsearchUrl string
	DryRun           bool
}

func main() {
	args := Args{
		ElasticsearchUrl: os.Getenv("ELASTICSEARCH-URL"),
	}
	flag.StringVar(&args.Addr, "listen-address", ":8080", "The address to listen on for HTTP requests.")
	flag.BoolVar(&args.Debug, "debug", true, "debug or not.")
	flag.BoolVar(&args.DryRun, "dryrun", false, "uses a null db driver that writes received webhooks to stdout")

	flag.Parse()

	var driver db.Storer
	if args.DryRun {
		log.Println("dry-run")
		driver = db.NullDB{}
	} else {
		elasticUrls := strings.Split(args.ElasticsearchUrl, ",")
		driver, _ = db.ConnectES(elasticUrls)
	}
	driver.GetVersion()

	s := server.New(args.Debug, driver)
	s.Start(args.Addr)
}
