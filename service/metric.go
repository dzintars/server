package service

import (
	"fmt"
	"log"

	"github.com/oswee/server/models"
	protobuf "github.com/oswee/stubs"
	metric "github.com/oswee/stubs/metric/v1"
	"golang.org/x/net/context"
)

// CreatePageView returns a list of all known films.
func (s *Server) CreatePageView(ctx context.Context, req *metric.CreatePageViewRequest) (*protobuf.Empty, error) {
	pv := `INSERT page_views SET
		x_forwarded_host=?,
		x_forwarded_server=?,
		url=?,
		user_agent=?,
		x_forwarded_for=?,
		request_time=?,
		request_headers=?`
	db := models.DBLoc()
	defer db.Close()

	stmt, err := db.Prepare(pv)
	if err != nil {
		log.Println(err)
	}
	r := req.PageView
	res, err := stmt.Exec(
		r.XForwardedHost,
		r.XForwardedServer,
		r.Url,
		r.UserAgent,
		r.XForwardedFor,
		r.RequestTime,
		r.RequestHeaders)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)

	return &protobuf.Empty{}, nil
}
