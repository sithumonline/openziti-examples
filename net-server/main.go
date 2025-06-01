package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"

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

	_, ok := ztx.GetService(*service)
	if !ok {
		log.Fatalf("%s service not found", *service)
	}

	listener, err := ztx.Listen(*service)
	if err != nil {
		log.Fatalf("failed to listen on service %s: %v", *service, err)
	}
	defer listener.Close()

	log.Printf("Listening on service %s", *service)

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("failed to accept connection: %v", err)
				break
			}
			if conn == nil {
				log.Println("failed to accept connection, exiting")
				break
			}
			log.Printf("accepted connection from %s", conn.RemoteAddr())

			go func(c net.Conn) {
				defer c.Close()
				log.Printf("handling connection from %s", c.RemoteAddr())

				for {
					buf := make([]byte, 128)
					n, err := c.Read(buf)
					if err != nil {
						log.Printf("error reading from connection: %v", err)
						return
					}
					if n > 0 {
						log.Printf("read %d bytes from connection, data: %s", n, string(buf[:n]))
						if string(buf[:n]) == "ping" {
							log.Println("received ping, sending pong")
							_, err := c.Write([]byte("pong"))
							if err != nil {
								log.Printf("error writing to connection: %v", err)
								return
							}
						}
						if string(buf[:n]) == "close" {
							log.Println("received close command, closing connection")
							c.Close()
							return
						}
					}
				}
			}(conn)
		}
	}()

	<-c
	log.Println("interrupt received, shutting down")
}
