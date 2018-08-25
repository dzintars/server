package service

import (
	"log"

	"github.com/oswee/proto/shipping/go"
	"github.com/oswee/server/models"
	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
)

// MapsAPIkey is a Google Maps API key for Oswee project (limited access)
var MapsAPIkey string = "AIzaSyBslJsVcubCFlQvF36XuxXbrEOm588gSa4"

// ListDeliveryOrders returns a list of all known films.
func (s *Server) ListDeliveryOrders(ctx context.Context, req *shipping.ListDeliveryOrdersRequest) (*shipping.ListDeliveryOrdersResponse, error) {
	listDeliveryOrdersSQL := `SELECT
		id,
		stakeholder_id,
		reference,
		destination_address,
		destination_zip,
		destination_lat,
		destination_lng,
		total_weight,
		routing_sequence
	FROM delivery_orders
	WHERE stakeholder_id = ?
	ORDER BY routing_sequence ASC
	LIMIT ?;`
	db := models.DBLoc()
	rows, err := db.Query(listDeliveryOrdersSQL, req.StakeholderId, req.ResultPerPage)
	if err != nil {
		log.Fatalf("Failed to query database: %v", err)
	}
	defer rows.Close()

	rx := []*shipping.DeliveryOrder{}

	for rows.Next() {
		r := &shipping.DeliveryOrder{}
		err := rows.Scan(
			&r.Id,
			&r.StakeholderId,
			&r.Reference,
			&r.DestinationAddress,
			&r.DestinationZip,
			&r.DestinationLat,
			&r.DestinationLng,
			&r.TotalWeight,
			&r.RoutingSequence,
		)
		if err != nil {
			log.Fatalf("Failed to read records: %v", err)
		}

		rx = append(rx, r)
	}
	defer db.Close()
	return &shipping.ListDeliveryOrdersResponse{DeliveryOrders: rx}, nil
}

// CreateDeliveryOrder creates new delivery order.
func (s *Server) CreateDeliveryOrder(ctx context.Context, req *shipping.CreateDeliveryOrderRequest) (*shipping.DeliveryOrder, error) {
	do := `INSERT delivery_orders
	SET
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

// UpdateDeliveryOrder ...
func (s *Server) UpdateDeliveryOrder(ctx context.Context, req *shipping.UpdateDeliveryOrderRequest) (*shipping.DeliveryOrder, error) {
	do := `UPDATE delivery_orders
	SET
		reference=?,
		destination_address=?,
		destination_zip=?,
		destination_lat=?,
		destination_lng=?,
		total_weight=?,
		routing_sequence=?
	WHERE id=?`
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
		r.RoutingSequence,
		r.Id)
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

// GeoCode ...
func geoCode(a string) ([]maps.GeocodingResult, error) {
	c, err := maps.NewClient(maps.WithAPIKey(MapsAPIkey))
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}
	r := &maps.GeocodingRequest{
		Address: a,
	}
	loc, err := c.Geocode(context.Background(), r)
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}
	return loc, nil
}

// GeoCodeDeliveryOrder ...
func (s *Server) GeoCodeDeliveryOrder(ctx context.Context, req *shipping.GeoCodeDeliveryOrderRequest) (*shipping.DeliveryOrder, error) {

	sql := `UPDATE delivery_orders 
	SET
		destination_lat=?,
		destination_lng=?
		WHERE id=?`

	db := models.DBLoc()
	defer db.Close()

	g, err := geoCode(req.Address)
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	stmt, err := db.Prepare(sql)
	if err != nil {
		log.Println(err)
	}
	_, err = stmt.Exec(
		g[0].Geometry.Location.Lat,
		g[0].Geometry.Location.Lng,
		req.Id)
	if err != nil {
		log.Fatal(err)
	}

	return &shipping.DeliveryOrder{}, nil
}
