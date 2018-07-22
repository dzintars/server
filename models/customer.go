package api

import (
	"context"
	"strings"

	pb "github.com/oswee/proto/customer/go"
)

// Server is used to implement customer.CustomerServer.
type Server struct {
	savedCustomers []*pb.CustomerRequest
}

// GetCustomers returns all customers by given filter
func (s *Server) GetCustomers(filter *pb.CustomerFilter, stream pb.Customer_GetCustomersServer) error {
	for _, customer := range s.savedCustomers {
		if filter.Keyword != "" {
			if !strings.Contains(customer.Name, filter.Keyword) {
				continue
			}
		}
		if err := stream.Send(customer); err != nil {
			return err
		}
	}
	return nil
}

// CreateCustomer creates a new Customer
func (s *Server) CreateCustomer(ctx context.Context, in *pb.CustomerRequest) (*pb.CustomerResponse, error) {
	pb.RegisterCustomerServer(s, Server{})
	customer := &pb.CustomerRequest{
		Id:    101,
		Name:  "Shiju Varghese",
		Email: "shiju@xyz.com",
		Phone: "732-757-2923",
	}
	s.savedCustomers = append(s.savedCustomers, in)
	return &pb.CustomerResponse{Id: in.Id, Success: true}, nil
}
