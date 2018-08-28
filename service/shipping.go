package service

import (
	"fmt"
	"log"
	"time"

	"github.com/oswee/server/models"
	pb "github.com/oswee/stubs"
	dms "github.com/oswee/stubs/dms/v1"
	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
)

// MapsAPIkey is a Google Maps API key for Oswee project (limited access)
var MapsAPIkey string = "AIzaSyBslJsVcubCFlQvF36XuxXbrEOm588gSa4"

// ListDeliveryOrders returns a list of all known films.
func (s *Server) ListDeliveryOrders(ctx context.Context, req *dms.ListDeliveryOrdersRequest) (*dms.ListDeliveryOrdersResponse, error) {
	qReqTime := time.Now()
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

	rx := []*dms.DeliveryOrder{}

	for rows.Next() {
		r := &dms.DeliveryOrder{}
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
	qResTime := time.Now()
	fmt.Println("ListDeliveryOrders executed at: ", qReqTime.Format("2006-01-02 15:04:05"), "; Took: ", qResTime.Sub(qReqTime))
	return &dms.ListDeliveryOrdersResponse{DeliveryOrders: rx}, nil
}

// CreateDeliveryOrder creates new delivery order.
func (s *Server) CreateDeliveryOrder(ctx context.Context, req *dms.CreateDeliveryOrderRequest) (*dms.DeliveryOrder, error) {
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

	return &dms.DeliveryOrder{}, nil
}

// UpdateDeliveryOrder ...
func (s *Server) UpdateDeliveryOrder(ctx context.Context, req *dms.UpdateDeliveryOrderRequest) (*dms.DeliveryOrder, error) {
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

	return &dms.DeliveryOrder{}, nil
}

// DeleteDeliveryOrder deletes delivery order of given ID
func (s *Server) DeleteDeliveryOrder(ctx context.Context, req *dms.DeleteDeliveryOrderRequest) (*pb.Empty, error) {
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

	return &pb.Empty{}, nil
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
func (s *Server) GeoCodeDeliveryOrder(ctx context.Context, req *dms.GeoCodeDeliveryOrderRequest) (*dms.DeliveryOrder, error) {

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

	return &dms.DeliveryOrder{}, nil
}
