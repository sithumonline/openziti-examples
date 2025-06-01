package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
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

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	conn, err := ztx.Dial(*service)
	if err != nil {
		log.Fatalf("failed to dial %s: %v", *service, err)
	}
	defer conn.Close()

	log.Printf("Dialed service %s", *service)

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	// Goroutine for reading from the connection
	go func() {
		for {
			buf := make([]byte, 128)
			n, err := conn.Read(buf)
			if err != nil {
				log.Printf("error reading from connection: %v", err)
				return
			}
			if n > 0 {
				log.Printf("read %d bytes from connection, data: %s", n, string(buf[:n]))
				if string(buf[:n]) == "ping" {
					log.Println("received ping, sending pong")
					_, err := conn.Write([]byte("pong"))
					if err != nil {
						log.Printf("error writing to connection: %v", err)
						return
					}
				}
			}
		}
	}()

	// Main goroutine handles ping and signal
	for {
		select {
		case <-c:
			log.Println("received interrupt signal, closing connection")
			_, err := conn.Write([]byte("close"))
			if err != nil {
				log.Printf("error writing close message: %v", err)
			}
			conn.Close()
			return
		case <-ticker.C:
			_, err := conn.Write([]byte("ping"))
			if err != nil {
				log.Printf("error writing to connection: %v", err)
				return
			}
			log.Println("sent ping to service")
		}
	}
}
