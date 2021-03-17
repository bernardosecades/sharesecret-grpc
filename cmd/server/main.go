package main

import (
	"context"
	"fmt"
	"log"
	"os"

	_ "github.com/bernardosecades/sharesecret/cmd"
	sharesecret "github.com/bernardosecades/sharesecret/internal"
	"github.com/bernardosecades/sharesecret/internal/server"
	"github.com/bernardosecades/sharesecret/internal/server/grpc"
	"github.com/bernardosecades/sharesecret/internal/server/http"
	"github.com/bernardosecades/sharesecret/internal/storage/mysql"
	"golang.org/x/sync/errgroup"
)

func main() {

	protocol := os.Getenv("SHARESECRET_SERVER_PROTOCOL")
	host := os.Getenv("SHARESECRET_SERVER_HOST")
	port := os.Getenv("SHARESECRET_SERVER_PORT")

	dbName := os.Getenv("DB_NAME")
	dbPass := os.Getenv("DB_PASS")
	dbUser := os.Getenv("DB_USER")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	secretKey := os.Getenv("SECRET_KEY")
	secretPassword := os.Getenv("SECRET_PASSWORD")

	secretRepository := mysql.NewMySQLSecretRepository(dbName, dbUser, dbPass, dbHost, dbPort)
	secretService := sharesecret.NewSecretService(secretRepository, secretKey, secretPassword)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// group context: https://bionic.fullstory.com/why-you-should-be-using-errgroup-withcontext-in-golang-server-handlers/
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		srvCfg := server.Config{Protocol: protocol, Host: host, Port: port}
		srv := grpc.NewServer(srvCfg, secretService)

		log.Printf("gRPC server running at %s://%s:%s ...\n", protocol, host, port)
		return srv.Serve()
	})

	g.Go(func() error {
		httpAddr := fmt.Sprintf(":%s", port)
		httpSrv := http.NewServer(httpAddr)

		log.Printf("HTTP server running at %s ...\n", httpAddr)
		return httpSrv.Serve(ctx)
	})

	log.Fatal(g.Wait())
}
