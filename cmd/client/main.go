package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	sharesecretgrpc "github.com/bernardosecades/sharesecret/genproto"
	"google.golang.org/grpc"
)

func main() {

	host := os.Getenv("SHARESECRET_SERVER_HOST")
	port := os.Getenv("SHARESECRET_SERVER_PORT")

	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", host, port), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := sharesecretgrpc.NewSecretServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	mySecret := "this is a my secret"

	fmt.Println("####### CREATE SECRET #######")
	fmt.Println("Call gRPC server to create a new secret -> '" + mySecret + "'")
	r1, err1 := client.CreateSecret(ctx, &sharesecretgrpc.CreateSecretRequest{Content: mySecret})
	if err1 != nil {
		log.Fatalf("CreateSecret: %v", err1)
	}

	fmt.Println("Response from gRPC server:")
	fmt.Println(r1.String())

	fmt.Println("####### READ SECRET #######")
	fmt.Println("Call gRPC server to read the secret with ID -> '" + r1.GetId() + "'")
	r2, err2 := client.SeeSecret(ctx, &sharesecretgrpc.SeeSecretRequest{Id: r1.GetId()})

	if err2 != nil {
		log.Fatalf("SeeSecret: %v", err2)
	}

	fmt.Println("Response from gRPC server:")
	fmt.Println(r2.String())
}
