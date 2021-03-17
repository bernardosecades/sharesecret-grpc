# ShareSecret (gRPC + API REST)

ShareSecret is a service to share sensitive information that's both simple and secure.

If you share some text will be display it once and then delete it. After that it's gone forever.

We keep secrets for up to 5 days.

`Note`: project based on https://onetimesecret.com to learn go, grpc and grpc gateway

## Why should I trust you?

General we can't do anything with your information even if we wanted to (which we don't). If it's a password for example, we don't know the username or even the application that the credentials are for.

If you include a password, we use it to encrypt the secret. We don't store the password (only a crypted hash) so we can never know what the secret is because we can't decrypt it.

## Demo

### gRPC: Create secret without password

User A create a new secret and send identifier of secret to user B. The secret only can be seen one time, use try to see a second time but the secret was deleted.

![Share Secret Demo](demo/grpc_create_secret_without_pass.gif)

### gRPC: Create secret with password

User A create a new secret with password and send identifier of secret to user B. The secret only can be seen one time and User B need to know the password to see the secret.

![Share Secret Demo](demo/grpc_create_secret_with_pass.gif)

### API REST (gRPC gateway)

We have the same funcionality with API REST using grpc-gateway

![Share Secret Demo](demo/rest_create_secret.gif)

# gRPC and gRPC-Gateway (RESTful API + gRPC)

The gRPC-Gateway is a plugin of the Google protocol buffers compiler protoc. It reads protobuf service definitions and generates a reverse-proxy server which translates a RESTful HTTP API into gRPC. This server is generated according to the google.api.http annotations in your service definition

This helps you provide your APIs in both gRPC and RESTful style at the same time.

File proto for this project:

```
syntax = "proto3";

package sharesecret;

import "google/api/annotations.proto";

option go_package = "genproto;proto";

service SecretService {
  rpc CreateSecret (CreateSecretRequest) returns (CreateSecretResponse) {
    option (google.api.http) = {
      post: "/v1/secret"
      body: "*"
    };
  }
  rpc SeeSecret (SeeSecretRequest) returns (SeeSecretResponse) {
    option (google.api.http) = {
      get: "/v1/secret/{id}"
      // Note: if secret require password to see the content client will send "grpc-metadata-password" and we got from context "password"
    };
  }
}

message CreateSecretRequest {
  string content = 1;
  string password = 2; // Optional
}

message CreateSecretResponse {
  string id = 1;
}

message SeeSecretRequest {
  string id = 1;
  string password = 2;
}

message SeeSecretResponse {
  string content = 1;
}

```

Without gateway (only gRPC), Generate secret.pb.go and secret_grpc.pb.go:

```bash
protoc -I=proto --go_out=. --go-grpc_opt=require_unimplemented_servers=false --go-grpc_out=. proto/secret.proto  
```

With gateway (API REST), Generate secret.pb.go, secret_grpc.pb.go and secret.pb.gw.go:

```bash 
protoc -I=proto -I /Users/admin/go/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --go_out=. --go-grpc_opt=require_unimplemented_servers=false --go-grpc_out=. proto/secret.proto --grpc-gateway_out=logtostderr=true:./genproto 
```

# Run the project

Only you will need execute:

```bash
make up
```

# Run tests

You can execute all tests or by type:

```bash
make test-all
make test-unit
make test-integration
make e2e
```

# Makefile

Up the service:

```bash
make up
```

Down the service:

```bash
make down
```

See status containers:

```bash
make ps
```

Execute client grpc to see an example:

```bash
make client-grpc-connection-example
```

Remove secrets expired:

```bash
make purge-secrets
```


Execute all tests and see coverage:

```bash
make test-coverage
```

Execute all tests:

```bash
make test-all
```

Execute all unit tests:

```bash
make test-unit
```

Execute all integration tests:

```bash
make test-integration
```

Execute all e2e tests:

```bash
make test-e2e
```

# Example go client gRPC to consume the service

The client `NewSecretServiceClient` (see folder `genproto`) gRPC was autogenerated by `protoc-gen-go` from our file `./proto/secret.proto`

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"
	"os"

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
```

# Docker

We have an unique docker file for development and production environments using multi-stage builds:

- Stage `builder` with golang image to can compile and execute tests (docker-compose use this stage).
- Stage `production` with alpine image with binaries generated in `builder` stage.

Note: we are using `upx` (see https://github.com/upx/upx) in the docker file to compress binary files generated.