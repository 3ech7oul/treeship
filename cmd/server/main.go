package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	pb "treeship/api/gen/v1"
	"treeship/helper"
	"treeship/server"

	"github.com/gorilla/mux"
	"github.com/spf13/pflag"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
	lvl  = pflag.StringP("log-level", "l", "info", "Log level")
)

func main() {
	var dir string

	logger, err := helper.NewLogger(*lvl)
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}

	flag.StringVar(&dir, "dir", ".", "the directory to serve files from. Defaults to the current dir")
	flag.Parse()
	r := mux.NewRouter()

	registry := server.NewAgentManager(logger)
	routes := server.NewRoutes(registry)

	r.HandleFunc("/agents", routes.ListAgents)
	r.HandleFunc("/messages", routes.SendMessage)

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		log.Printf("Starting server on port %v", *port)
		log.Fatal(srv.ListenAndServe())
	}()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterAgentServiceServer(s, registry)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
