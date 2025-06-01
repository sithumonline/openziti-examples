package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/openziti/sdk-golang/ziti"
)

var (
	identity = flag.String("identity", "", "Ziti Identity file")
	service  = flag.String("service", "", "Ziti Service")
)

func main() {
	flag.Parse()

	if *identity == "" {
		log.Fatal("identity file must be specified with -identity flag")
	}

	if *service == "" {
		log.Fatal("service must be specified with -service flag")
	}

	cfg, err := ziti.NewConfigFromFile(*identity)
	if err != nil {
		log.Fatalf("failed to load ziti identity{%v}: %v", identity, err)
	}

	ztx, err := ziti.NewContext(cfg)
	if err != nil {
		log.Fatalf("failed to create ziti context: %v", err)
	}

	err = ztx.Authenticate()
	if err != nil {
		log.Fatalf("failed to authenticate: %v", err)
	}

	//--------

	listener, err := ztx.Listen(*service)
	if err != nil {
		log.Fatalf("could not bind service %s: %v", *service, err)
	}

	http.HandleFunc("/api/time", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "{\"time\": \"%s\"}", time.Now().Format(time.RFC3339))
	})

	// listener = tls.NewListener(listener, &tls.Config{ InsecureSkipVerify: true })
	server := &http.Server{
		Handler:   http.DefaultServeMux,
		TLSConfig: nil,
	}

	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
