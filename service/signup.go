package service

import (
	"fmt"
	"log"

	"github.com/oswee/server/models"
	signup "github.com/oswee/stubs/signup/v1"
	"golang.org/x/net/context"
)

// CreateSignup creates new Signup record
func (s *Server) CreateSignup(ctx context.Context, req *signup.CreateSignupRequest) (*signup.Signup, error) {
	sql := `INSERT user_signups SET
		first_name=?,
		last_name=?,
		email=?,
		username=?,
		password=?,
		status=?`
	db := models.DBLoc()
	defer db.Close()

	stmt, err := db.Prepare(sql)
	if err != nil {
		log.Println(err)
	}
	r := req.Signup
	res, err := stmt.Exec(
		r.FistName,
		r.LastName,
		r.Email,
		r.Username,
		r.Password,
		r.Status)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)

	return &signup.Signup{}, nil
}
