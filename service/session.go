package service

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/oswee/server/models"
	session "github.com/oswee/stubs/session/v1"
	"golang.org/x/net/context"
)

// CreateSession creates new User session
func (s *Server) CreateSession(ctx context.Context, req *session.CreateSessionRequest) (*session.Session, error) {
	qry := `INSERT sessions SET
		user_session_id=?,
		user_id=?,
		permission_id=?`
	db := models.DBLoc()
	defer db.Close()

	stmt, err := db.Prepare(qry)
	if err != nil {
		log.Println(err)
	}
	r := req.Session
	res, err := stmt.Exec(
		r.UserSessionId,
		r.UserId,
		r.PermissionId)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)

	return &session.Session{}, nil
}

// GetSession ...
func (s *Server) GetSession(ctx context.Context, req *session.GetSessionRequest) (*session.Session, error) {
	var a *session.Session
	qry := `SELECT
		ID,
		user_session_id,
		user_id,
		permission_id,
		create_time
	FROM sessions
	WHERE user_session_id=?`
	db := models.DBLoc()
	row := db.QueryRow(qry, req.Session.UserSessionId)
	switch err := row.Scan(&a.Id, &a.UserSessionId, &a.UserId, &a.PermissionId, &a.CreateTime); err {
	case sql.ErrNoRows:
		log.Fatal("No Session record were returned!")
	case nil:
		fmt.Println("Session record were returned")
	default:
		panic(err)
	}
	defer db.Close()
	return a, nil
}
