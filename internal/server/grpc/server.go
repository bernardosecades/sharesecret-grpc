package grpc

import (
	"fmt"
	sharesecretgrpc "github.com/bernardosecades/sharesecret/genproto"
	sharesecret "github.com/bernardosecades/sharesecret/internal"
	"github.com/bernardosecades/sharesecret/internal/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"net"
	"os"
)

type grpcServer struct {
	config        server.Config
	secretService sharesecret.SecretService
}

func NewServer(config server.Config, ss sharesecret.SecretService) server.Server {
	return &grpcServer{config: config, secretService: ss}
}

func (s *grpcServer) Serve() error {
	addr := fmt.Sprintf("%s:%s", s.config.Host, s.config.Port)
	listener, err := net.Listen(s.config.Protocol, addr)
	if err != nil {
		return err
	}

	grpcLog := grpclog.NewLoggerV2(os.Stdout, os.Stderr, os.Stderr)
	grpclog.SetLoggerV2(grpcLog)

	srv := grpc.NewServer()

	serviceServer := NewShareSecretServer(s.secretService)
	sharesecretgrpc.RegisterSecretServiceServer(srv, serviceServer)

	if err := srv.Serve(listener); err != nil {
		return err
	}

	return nil
}
