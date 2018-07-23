package service

import (
	"log"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/oswee/proto"
	"github.com/oswee/server/models"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// To start with, we'll hardcode the database of films.
var films = []*proto.Film{
	&proto.Film{
		Id:          "4",
		Title:       "A New Hope",
		Director:    "George Lucas",
		Producer:    "Gary Kurtz, Rick McCallum",
		ReleaseDate: toProto(1977, 5, 25),
	},
	&proto.Film{
		Id:          "5",
		Title:       "The Empire Strikes Back",
		Director:    "Irvin Kershner",
		Producer:    "Gary Kurtz, Rick McCallum",
		ReleaseDate: toProto(1980, 5, 17),
	},
	&proto.Film{
		Id:          "6",
		Title:       "Return of the Jedi",
		Director:    "Richard Marquand",
		Producer:    "Howard G. Kazanjian, George Lucas, Rick McCallum",
		ReleaseDate: toProto(1983, 5, 25),
	},
}

func toProto(year, month, day int) *timestamp.Timestamp {
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	ts, err := ptypes.TimestampProto(t)
	if err != nil {
		panic(err)
	}
	return ts
}

type Server struct {
	films []*proto.Film
}

// GetFilm queries a film by ID or returns an error if not found.
func (s *Server) GetFilm(ctx context.Context,
	req *proto.GetFilmRequest) (*proto.GetFilmResponse, error) {
	var film *proto.Film
	for _, f := range films {
		if f.Id == req.Id {
			film = f
			break
		}
	}
	if film == nil {
		return nil, status.Errorf(codes.NotFound, "no film with id %q", req.Id)
	}
	return &proto.GetFilmResponse{Film: film}, nil
}

// ListFilms returns a list of all known films.
func (s *Server) ListFilms(ctx context.Context, req *proto.ListFilmsRequest) (*proto.ListFilmsResponse, error) {
	getFilms := `SELECT id, title, director, producer FROM films;`
	db := models.DBLoc()
	rows, err := db.Query(getFilms)
	if err != nil {
		log.Fatalf("Failed to query database: %v", err)
	}
	defer rows.Close()
	results := []*proto.Film{}
	for rows.Next() {
		film := proto.Film{}
		var (
			id       string
			title    string
			director string
			producer string
		)
		err := rows.Scan(&id, &title, &director, &producer)
		if err != nil {
			log.Fatalf("Failed to read records: %v", err)
		}
		film.Id = id
		film.Title = title
		film.Director = director
		film.Producer = producer

		results = append(results, &film)
	}
	defer db.Close()
	return &proto.ListFilmsResponse{Films: results}, nil
}

// compile-type check that our new type provides the correct server interface
var _ proto.StarfriendsServer = (*Server)(nil)
