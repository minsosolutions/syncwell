package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/minsosolutions/syncwell/internal/mockapi"
)

func main() {
	addr := flag.String("addr", ":8099", "listen address")
	data := flag.String("data", "data", "path to the data directory")
	flag.String("config", "config/syncwell.json", "unused by the mock API; accepted so the container command is uniform")
	flag.Parse()

	log.Printf("syncwell mock api listening on %s, serving %s", *addr, *data)
	if err := http.ListenAndServe(*addr, mockapi.New(*data)); err != nil {
		log.Fatal(err)
	}
}
