package main

import (
	"context"
	"flag"
	"log"

	"github.com/openziti/sdk-golang/ziti"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

var (
	identity = flag.String("identity", "", "Ziti Identity file")
	service  = flag.String("service", "", "Ziti Service")
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(_ context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

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

	lis, err := ztx.Listen(*service)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
