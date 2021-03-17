package grpc

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sharesecretgrpc "github.com/bernardosecades/sharesecret/genproto"
	sharesecret "github.com/bernardosecades/sharesecret/internal"
	_ "google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	_ "google.golang.org/grpc/status"
)

type shareSecretHandler struct {
	secretService sharesecret.SecretService
}

func NewShareSecretServer(s sharesecret.SecretService) sharesecretgrpc.SecretServiceServer {
	return &shareSecretHandler{secretService: s}
}

func (s shareSecretHandler) CreateSecret(ctx context.Context, req *sharesecretgrpc.CreateSecretRequest) (*sharesecretgrpc.CreateSecretResponse, error) {

	secret, err := s.secretService.CreateSecret(req.Content, req.Password)
	if err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}

	r := &sharesecretgrpc.CreateSecretResponse{}
	r.Id = secret.ID

	return r, nil
}

func (s shareSecretHandler) SeeSecret(ctx context.Context, req *sharesecretgrpc.SeeSecretRequest) (*sharesecretgrpc.SeeSecretResponse, error) {

	reqHeaders, ok := metadata.FromIncomingContext(ctx) // In postman su need put prefix: grpc-metadata-{yourHeaderName}. Example: grpc-metadata-password

	if !ok {
		return nil, errors.New("Error context")
	}

	var password = ""
	if pass, ok := reqHeaders["password"]; ok {
		password = pass[0]
	} else {
		password = req.Password
	}

	c, err := s.secretService.GetContentSecret(req.Id, password)

	if err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}

	r := &sharesecretgrpc.SeeSecretResponse{}
	r.Content = c

	return r, nil
}
