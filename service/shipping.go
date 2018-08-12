package service

import (
	"log"

	"github.com/oswee/proto/shipping/go"
	"github.com/oswee/server/models"
	"golang.org/x/net/context"
)

// ListDeliveryOrders returns a list of all known films.
func (s *Server) ListDeliveryOrders(ctx context.Context, req *shipping.ListDeliveryOrdersRequest) (*shipping.ListDeliveryOrdersResponse, error) {
	listDeliveryOrders := `SELECT id, stakeholder_id, destination_address, total_weight, routing_sequence FROM delivery_orders WHERE stakeholder_id = ? LIMIT ?;`
	db := models.DBLoc()
	rows, err := db.Query(listDeliveryOrders, req.StakeholderId, req.ResultPerPage)
	if err != nil {
		log.Fatalf("Failed to query database: %v", err)
	}
	defer rows.Close()
	r := []*shipping.DeliveryOrder{}
	for rows.Next() {
		s := &shipping.DeliveryOrder{}

		err := rows.Scan(&s.Id, &s.StakeholderId, &s.DestinationAddress, &s.TotalWeight, &s.RoutingSequence)
		if err != nil {
			log.Fatalf("Failed to read records: %v", err)
		}

		r = append(r, s)
	}
	defer db.Close()
	return &shipping.ListDeliveryOrdersResponse{DeliveryOrders: r}, nil
}
