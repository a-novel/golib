package grpcf

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"

	"google.golang.org/api/idtoken"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/credentials/oauth"
)

var (
	_ CredentialsProvider = (*LocalCredentialsProvider)(nil)
	_ CredentialsProvider = (*GcloudCredentialsProvider)(nil)
)

type CredentialsProvider interface {
	Options(ctx context.Context) ([]grpc.DialOption, error)
}

type LocalCredentialsProvider struct{}

func (provider *LocalCredentialsProvider) Options(_ context.Context) ([]grpc.DialOption, error) {
	return []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}, nil
}

type GcloudCredentialsProvider struct {
	Host string
}

func (provider *GcloudCredentialsProvider) Options(ctx context.Context) ([]grpc.DialOption, error) {
	systemRoots, err := x509.SystemCertPool()
	if err != nil {
		return nil, fmt.Errorf("error getting system cert: %w", err)
	}

	tokenSource, err := idtoken.NewTokenSource(ctx, "https://"+provider.Host)
	if err != nil {
		return nil, fmt.Errorf("error getting token source: %w", err)
	}

	cred := credentials.NewTLS(&tls.Config{
		MinVersion: tls.VersionTLS12,
		RootCAs:    systemRoots,
	})

	return []grpc.DialOption{
		grpc.WithTransportCredentials(cred),
		grpc.WithAuthority(provider.Host + ":443"),
		grpc.WithPerRPCCredentials(oauth.TokenSource{TokenSource: tokenSource}),
	}, nil
}
