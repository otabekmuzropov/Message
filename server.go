package main

import (
	pb "bitbucket.org/alien_soft/Message/genproto/message"
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"log"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"math/rand"
	"net"
)

type server struct {
	db *sqlx.DB
}

func (s *server) Create(ctx context.Context, req *pb.Message) (*pb.Message, error) {
	var message pb.Message
	log.Println("request sever", req.Name)

	id := rand.Uint32()

	create := `insert into message(id, name, time) values($1, $2, $3)`

	_, err := s.db.Exec(create, id, req.Name, req.Time)

	if err != nil {
		return nil, err
	}

	row := s.db.QueryRow("select id, name, time from message where id = $1", id)

	err = row.Scan(&message.Id, &message.Name, &message.Time)

	if err != nil {
		return nil, err
	}

	return &message, nil
}

func (s *server) Update(ctx context.Context, req *pb.Message) (*pb.Message, error) {
	var message pb.Message
	log.Println("Received request")

	update := `update message set name = $2, time = $3 where id = $1`

	_, err := s.db.Exec(update, req.Id, req.Name, req.Time)

	if err != nil {
		return nil, err
	}

	row := s.db.QueryRow("select id, name, time from message where id = $1", req.Id)

	err = row.Scan(&message.Id, &message.Name, &message.Time)

	if err != nil {
		return nil, err
	}

	return &message, nil
}

func (s *server) Delete(ctx context.Context, req *pb.DeleteRequest) (*empty.Empty, error) {

	log.Println("Received request")

	delete := `delete from message where id = $1`

	_, err := s.db.Exec(delete, req.Id)

	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func main() {
	listen, err := net.Listen("tcp", ":50054")

	if err != nil {
		log.Println("error while listening port", err)
		return
	}

	db, err := sqlx.Connect("postgres", "user=postgres dbname=test password=123 sslmode=disable")

	if err != nil {
		log.Println("error while connecting db", err)
		return
	}

	s := grpc.NewServer()

	pb.RegisterMessageServiceServer(s, &server{db:db})

	log.Println("listening %d port", 50054)

	if err := s.Serve(listen); err != nil {
		log.Println("failed serve", err)
		return
	}
}

