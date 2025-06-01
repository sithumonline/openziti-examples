package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"

	sdk_golang "github.com/openziti/sdk-golang"
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

	// client := sdk_golang.NewHttpClient(ztx, &tls.Config{InsecureSkipVerify: true})
	client := sdk_golang.NewHttpClient(ztx, nil)
	resp, err := client.Get(fmt.Sprintf("http://%s/api/time", *service))
	if err != nil {
		log.Fatalf("failed to get %s: %v", *service, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("unexpected status code: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("failed to read response body: %v", err)
	}

	// Print the response body
	log.Printf("Response from %s: %s", *service, body)
}
