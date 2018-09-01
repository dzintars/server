package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/oswee/server/service"
	app "github.com/oswee/stubs/app/v1"
	dms "github.com/oswee/stubs/dms/v1"
	metric "github.com/oswee/stubs/metric/v1"
	session "github.com/oswee/stubs/session/v1"
	//signin "github.com/oswee/stubs/signin/v1"
	signup "github.com/oswee/stubs/signup/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	port := flag.Int("port", 8080, "The port on which gRPC server will listen")
	flag.Parse()

	// We're not providing TLS options, so server will use plaintext.
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// create a server instance
	s := service.Server{}
	// Create the TLS credentials
	creds, err := credentials.NewServerTLSFromFile("cert/server.crt", "cert/server.key")
	if err != nil {
		log.Fatalf("could not load TLS keys: %s", err)
	}
	// Create an array of gRPC options with the credentials
	opts := []grpc.ServerOption{grpc.Creds(creds)}

	fmt.Printf("Server listening on %v\n", lis.Addr())

	grpcServer := grpc.NewServer(opts...)

	// Register our service implementation
	app.RegisterApplicationServiceServer(grpcServer, &s)
	dms.RegisterShippingServiceServer(grpcServer, &s)
	metric.RegisterMetricServer(grpcServer, &s)
	//signin.RegisterSigninServiceServer(grpcServer, &service.Server{})
	signup.RegisterSignupServiceServer(grpcServer, &s)
	session.RegisterSessionServiceServer(grpcServer, &s)

	// trap SIGINT / SIGTERM to exit cleanly
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("Shutting down the server...")
		grpcServer.GracefulStop()
	}()

	// finally, run the server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
}
func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
