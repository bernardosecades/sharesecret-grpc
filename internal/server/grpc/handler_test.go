// +build e2e

package grpc

import (
	"context"
	"log"
	"net"
	"os"
	"testing"

	sharesecretgrpc "github.com/bernardosecades/sharesecret/genproto"
	sharesecret "github.com/bernardosecades/sharesecret/internal"
	"github.com/bernardosecades/sharesecret/internal/storage/mysql"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {

	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()

	dbName := os.Getenv("DB_NAME")
	dbPass := os.Getenv("DB_PASS")
	dbUser := os.Getenv("DB_USER")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	secretKey := os.Getenv("SECRET_KEY")
	secretPassword := os.Getenv("SECRET_PASSWORD")

	secretRepository := mysql.NewMySQLSecretRepository(dbName, dbUser, dbPass, dbHost, dbPort)
	secretService := sharesecret.NewSecretService(secretRepository, secretKey, secretPassword)

	sharesecretgrpc.RegisterSecretServiceServer(s, NewShareSecretServer(secretService))
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestCreateAndSeeSecretWithoutPassword(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := sharesecretgrpc.NewSecretServiceClient(conn)
	resp1, err1 := client.CreateSecret(ctx, &sharesecretgrpc.CreateSecretRequest{Content: "This is my secret"})
	if err1 != nil {
		t.Fatalf("CreateSecret failed: %v", err1)
	}

	assert.Len(t, resp1.GetId(), 36)

	resp2, err2 := client.SeeSecret(ctx, &sharesecretgrpc.SeeSecretRequest{Id: resp1.GetId()})

	if err2 != nil {
		t.Fatalf("CreateSecret failed: %v", err1)
	}

	assert.Equal(t, "This is my secret", resp2.GetContent())

	resp3, err3 := client.SeeSecret(ctx, &sharesecretgrpc.SeeSecretRequest{Id: resp1.GetId()})

	assert.Nil(t, resp3)
	assert.NotNil(t, err3)
}

func TestCreateAndSeeSecretWithPassword(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := sharesecretgrpc.NewSecretServiceClient(conn)
	resp1, err1 := client.CreateSecret(ctx, &sharesecretgrpc.CreateSecretRequest{Content: "This is my secret", Password: "1234"})
	if err1 != nil {
		t.Fatalf("CreateSecret failed: %v", err1)
	}

	assert.Len(t, resp1.GetId(), 36)

	resp2, err2 := client.SeeSecret(ctx, &sharesecretgrpc.SeeSecretRequest{Id: resp1.GetId(), Password: "1234"})

	if err2 != nil {
		t.Fatalf("CreateSecret failed: %v", err1)
	}

	assert.Equal(t, "This is my secret", resp2.GetContent())

	resp3, err3 := client.SeeSecret(ctx, &sharesecretgrpc.SeeSecretRequest{Id: resp1.GetId(), Password: "1234"})

	assert.Nil(t, resp3)
	assert.NotNil(t, err3)
}
