package http

import (
	"context"
	"net/http"

	sharesecretgrpc "github.com/bernardosecades/sharesecret/genproto"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	httpAddr string
}

func NewServer(httpAddr string) *Server {
	return &Server{httpAddr: httpAddr}
}

func (s *Server) Serve(ctx context.Context) error {

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure(), grpc.WithReturnConnectionError()}

	err := sharesecretgrpc.RegisterSecretServiceHandlerFromEndpoint(ctx, mux, s.httpAddr, opts)
	if err != nil {
		return err
	}

	return http.ListenAndServe(s.httpAddr, mux)
}
