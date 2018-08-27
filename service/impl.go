package service

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/oswee/server/models"
	app "github.com/oswee/stubs/app/v1"
	"golang.org/x/net/context"
)

// Server ...
type Server struct {
	//films []*proto.Film
}

// GetApplication returns requested application data
func (s *Server) GetApplication(ctx context.Context, req *app.GetApplicationRequest) (*app.GetApplicationResponse, error) {
	var a app.Application
	getApp := `SELECT
			id,
			parent_id,
			name,
			full_name,
			permalink,
			type
		FROM os_applications
		WHERE os_applications.id =?`
	db := models.DBLoc()
	row := db.QueryRow(getApp, req.Id)
	switch err := row.Scan(&a.Id, &a.ParentId, &a.Name, &a.FullName, &a.Permalink, &a.Type); err {
	case sql.ErrNoRows:
		log.Fatal("No Application record were returned!")
	case nil:
		fmt.Println("Application record were returned")
	default:
		panic(err)
	}
	defer db.Close()
	return &app.GetApplicationResponse{Application: &a}, nil
}

// ListApplications returns a list of all known films.
func (s *Server) ListApplications(ctx context.Context, req *app.ListApplicationsRequest) (*app.ListApplicationsResponse, error) {
	listApps := `SELECT id, parent_id, name, full_name, permalink, type FROM os_applications LIMIT ?;`
	db := models.DBLoc()
	rows, err := db.Query(listApps, req.ResultPerPage)
	if err != nil {
		log.Fatalf("Failed to query database: %v", err)
	}
	defer rows.Close()
	r := []*app.Application{}
	for rows.Next() {
		a := &app.Application{}

		err := rows.Scan(&a.Id, &a.ParentId, &a.Name, &a.FullName, &a.Permalink, &a.Type)
		if err != nil {
			log.Fatalf("Failed to read records: %v", err)
		}

		r = append(r, a)
	}
	defer db.Close()
	return &app.ListApplicationsResponse{Applications: r}, nil
}

// compile-type check that our new type provides the correct server interface
//var _ proto.StarfriendsServer = (*Server)(nil)
