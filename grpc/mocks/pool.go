package grpcmocks

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"

	"golang.org/x/oauth2"
	"google.golang.org/api/idtoken"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	testgrpc "google.golang.org/grpc/interop/grpc_testing"
)

func FakeClientCerts(mocked ...[]byte) func() (*x509.CertPool, error) {
	return func() (*x509.CertPool, error) {
		pool := x509.NewCertPool()

		for _, cert := range mocked {
			if !pool.AppendCertsFromPEM(cert) {
				return nil, fmt.Errorf("append certificate to default pool")
			}
		}

		return pool, nil
	}
}

func FakeTokenSource(
	token oauth2.TokenSource,
) func(ctx context.Context, audience string, opts ...idtoken.ClientOption) (oauth2.TokenSource, error) {
	return func(_ context.Context, _ string, _ ...idtoken.ClientOption) (oauth2.TokenSource, error) {
		return token, nil
	}
}

func FakeGRPCServer(srv testgrpc.TestServiceServer, keyFile, certFile []byte) (func(), error) {
	var sOpts []grpc.ServerOption

	if keyFile == nil || certFile == nil {
		sOpts = append(sOpts, grpc.Creds(insecure.NewCredentials()))
	} else {
		cert, err := tls.X509KeyPair(certFile, keyFile)
		if err != nil {
			return nil, err
		}

		transport := credentials.NewTLS(&tls.Config{Certificates: []tls.Certificate{cert}})
		sOpts = append(sOpts, grpc.Creds(transport))
	}

	s := grpc.NewServer(sOpts...)

	testgrpc.RegisterTestServiceServer(s, srv)

	lis, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		return nil, err
	}

	go func() {
		_ = s.Serve(lis)
	}()

	stop := func() {
		s.Stop()
		_ = lis.Close()
	}

	return stop, nil
}
