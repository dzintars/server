package service

import (
	"log"

	"github.com/oswee/proto/shipping/go"
	"github.com/oswee/server/models"
	"golang.org/x/net/context"
)

// ListDeliveryOrders returns a list of all known films.
func (s *Server) ListDeliveryOrders(ctx context.Context, req *shipping.ListDeliveryOrdersRequest) (*shipping.ListDeliveryOrdersResponse, error) {
	listDeliveryOrders := `SELECT id, stakeholder_id, reference, destination_address, destination_zip, destination_lat, destination_lng, total_weight, routing_sequence FROM delivery_orders WHERE stakeholder_id = ? LIMIT ?;`
	db := models.DBLoc()
	rows, err := db.Query(listDeliveryOrders, req.StakeholderId, req.ResultPerPage)
	if err != nil {
		log.Fatalf("Failed to query database: %v", err)
	}
	defer rows.Close()
	r := []*shipping.DeliveryOrder{}
	for rows.Next() {
		s := &shipping.DeliveryOrder{}

		err := rows.Scan(&s.Id, &s.StakeholderId, &s.Reference, &s.DestinationAddress, &s.DestinationZip, &s.DestinationLat, &s.DestinationLng, &s.TotalWeight, &s.RoutingSequence)
		if err != nil {
			log.Fatalf("Failed to read records: %v", err)
		}

		r = append(r, s)
	}
	defer db.Close()
	return &shipping.ListDeliveryOrdersResponse{DeliveryOrders: r}, nil
}

// CreateDeliveryOrder creates new delivery order.
func (s *Server) CreateDeliveryOrder(ctx context.Context, req *shipping.CreateDeliveryOrderRequest) (*shipping.DeliveryOrder, error) {
	do := `INSERT delivery_orders SET
		reference=?,
		destination_address=?,
		destination_zip=?,
		destination_lat=?,
		destination_lng=?,
		total_weight=?,
		routing_sequence=?`
	db := models.DBLoc()
	defer db.Close()

	stmt, err := db.Prepare(do)
	if err != nil {
		log.Println(err)
	}
	r := req.DeliveryOrder
	_, err = stmt.Exec(
		r.Reference,
		r.DestinationAddress,
		r.DestinationZip,
		r.DestinationLat,
		r.DestinationLng,
		r.TotalWeight,
		r.RoutingSequence)
	if err != nil {
		log.Fatal(err)
	}

	return &shipping.DeliveryOrder{}, nil
}

// DeleteDeliveryOrder deletes delivery order of given ID
func (s *Server) DeleteDeliveryOrder(ctx context.Context, req *shipping.DeleteDeliveryOrderRequest) (*shipping.EmptyDeliveryOrder, error) {
	do := `DELETE FROM delivery_orders WHERE id=?`
	db := models.DBLoc()
	defer db.Close()

	stmt, err := db.Prepare(do)
	if err != nil {
		log.Println(err)
	}
	_, err = stmt.Exec(req.Id)
	if err != nil {
		log.Fatal(err)
	}

	return &shipping.EmptyDeliveryOrder{}, nil
}
