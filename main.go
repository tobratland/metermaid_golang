package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"

	_ "github.com/lib/pq"
	deliveryGrpc "github.com/tobra/metermaid/delivery/grpc"
	"github.com/tobra/metermaid/repository"
	"github.com/tobra/metermaid/usecase"
	"google.golang.org/grpc"
)

func main() {
	serverAddress := "127.0.0.1:8080"
	dbUser := "postgres"
	dbPass := "test123"
	dbName := "metermaid"
	dbSslMode := "disable"

	connStr := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=%s", dbUser, dbName, dbPass, dbSslMode)

	dbConn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

	tr := repository.NewPostgresTimeSeriesRepository(dbConn)
	tu := usecase.NewTimeseriesUsecase(tr)
	list, err := net.Listen("tcp", serverAddress)
	if err != nil {
		fmt.Println("Something happened when opening server")
	}

	server := grpc.NewServer()
	deliveryGrpc.NewTimeSeriesServerGrpc(server, tu)
	fmt.Println("Server running at ", serverAddress)

	err = server.Serve(list)
	if err != nil {
		fmt.Println("Unexpected error", err)
	}
}
