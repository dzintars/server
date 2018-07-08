package models

import (
	"database/sql"
	"fmt"

	pb "github.com/oswee/proto/customer/go"
)

// Server represents the gRPC server
type Server struct {
}

// GetCustomer returns full ist of all customers
func (s *Server) GetCustomer(r *pb.GetCustomerRequest) (*pb.GetCustomerResponse, error) {
	customer, err := getSome(string(r.Id))
	if err != nil {
		fmt.Println(err)
	}
	return &pb.GetCustomerResponse{*pb.Customer}, nil
}

func getSome(customerID string) (*pb.Customer, error) {
	var customer *pb.Customer

	getCustomer := `SELECT
			id,
			name,
			created_at
		FROM os_customers
		WHERE os_customers.id = ?`

	db := dbLoc()
	row := db.QueryRow(getCustomer, customerID)
	switch err := row.Scan(&customer.Id, &customer.Name, &customer.CreateTime); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
	case nil:
		fmt.Println("Customer:", customer, "were returned")
	default:
		panic(err)
	}
	defer db.Close()
	return customer, nil
}
