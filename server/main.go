package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net"

	_ "github.com/mattn/go-sqlite3"

	pb "github.com/maki3cat/toy-temporal/api-go/workflow"
	"google.golang.org/grpc"
)

var db *sql.DB
var storage *workflowStorage

func init() {
	temp, err := sql.Open("sqlite3", "./toy-temporal.db")
	if err != nil {
		log.Fatal(err)
	}
	db = temp
	InitDB(db)
}

func main() {
	defer db.Close()

	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterWorkflowServer(s, &workflowServer{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
