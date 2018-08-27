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
	"google.golang.org/grpc"
)

func main() {
	port := flag.Int("port", 8080, "The port on which gRPC server will listen")
	flag.Parse()

	// We're not providing TLS options, so server will use plaintext.
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Printf("Server listening on %v\n", lis.Addr())
	grpcServer := grpc.NewServer()

	// Register our service implementation
	app.RegisterApplicationServiceServer(grpcServer, &service.Server{})
	dms.RegisterShippingServiceServer(grpcServer, &service.Server{})
	metric.RegisterMetricServer(grpcServer, &service.Server{})

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
