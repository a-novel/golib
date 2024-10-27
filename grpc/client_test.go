package grpc_test

import (
	"context"
	"testing"
	"time"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	testgrpc "google.golang.org/grpc/interop/grpc_testing"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	anovelgrpc "github.com/a-novel/golib/grpc"
	grpcmocks "github.com/a-novel/golib/grpc/mocks"
	x509mocks "github.com/a-novel/golib/grpc/mocks/x509"
	"github.com/a-novel/golib/testutils"
)

type stubServerParams struct {
	insecure bool

	authority string
	token     *oauth2.Token
}

func setupStubServer(t *testing.T, params stubServerParams) *grpcmocks.StubServer {
	t.Helper()

	ss := &grpcmocks.StubServer{
		EmptyCallF: func(ctx context.Context, _ *testgrpc.Empty) (*testgrpc.Empty, error) {
			pr, ok := peer.FromContext(ctx)
			if !ok {
				return nil, status.Error(codes.DataLoss, "Failed to get peer from ctx")
			}

			expectedSecLevel := lo.Ternary(params.insecure, credentials.NoSecurity, credentials.PrivacyAndIntegrity)
			if err := credentials.CheckSecurityLevel(pr.AuthInfo, expectedSecLevel); err != nil {
				return nil, status.Errorf(codes.Unauthenticated, "Wrong security level: %s", err)
			}

			md, ok := metadata.FromIncomingContext(ctx)
			if !ok {
				return nil, status.Error(codes.DataLoss, "Failed to get metadata from ctx")
			}

			contentType, ok := md["content-type"]
			if !ok || len(contentType) == 0 {
				return nil, status.Error(codes.DataLoss, "Failed to get content type from metadata")
			}
			if contentType[0] != "application/grpc" {
				return nil, status.Errorf(codes.InvalidArgument, "Wrong content type: got %s, want application/grpc", contentType[0])
			}

			if !params.insecure {
				authority, ok := md[":authority"]
				if !ok || len(authority) == 0 {
					return nil, status.Error(codes.DataLoss, "Failed to get authority from metadata")
				}
				if authority[0] != params.authority {
					return nil, status.Errorf(codes.Unauthenticated, "Wrong authority: %s, want %s", authority[0], params.authority)
				}

				wantToken := params.token.TokenType + " " + params.token.AccessToken

				token, ok := md["authorization"]
				if !ok || len(token) == 0 {
					return nil, status.Error(codes.DataLoss, "Failed to get token from metadata")
				}
				if token[0] != wantToken {
					return nil, status.Errorf(codes.Unauthenticated, "Wrong token: %s, want %s", token[0], wantToken)
				}
			}

			return new(testgrpc.Empty), nil
		},
	}

	return ss
}

func TestConnDevOK(t *testing.T) {
	// Run in isolation to avoid messing with flags.
	testutils.RunCMD(t, &testutils.CMDConfig{
		CmdFn: func(t *testing.T) {
			anovelgrpc.SystemCertPool = grpcmocks.FakeClientCerts()
			anovelgrpc.NewTokenSource = grpcmocks.FakeTokenSource(nil)

			ss := setupStubServer(t, stubServerParams{insecure: true})
			clean, err := grpcmocks.FakeGRPCServer(ss, nil, nil)
			require.NoError(t, err)
			defer clean()

			connPool := anovelgrpc.NewConnPool()
			defer connPool.Close()

			client, err := connPool.Open("127.0.0.1", 8080, anovelgrpc.ProtocolHTTPS)
			require.NoError(t, err)

			c := testgrpc.NewTestServiceClient(client)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			_, err = c.EmptyCall(ctx, new(testgrpc.Empty))
			testutils.RequireGRPCCodesEqual(t, err, codes.OK)
		},
		MainFn: func(t *testing.T, res *testutils.CMDResult) {
			require.True(t, res.Success, res.STDErr, res.STDOut)
		},
		Env: []string{"ENV=dev"},
	})
}

func TestConnDevErrorSecureServer(t *testing.T) {
	// Run in isolation to avoid messing with flags.
	testutils.RunCMD(t, &testutils.CMDConfig{
		CmdFn: func(t *testing.T) {
			anovelgrpc.SystemCertPool = grpcmocks.FakeClientCerts(x509mocks.ServerCACertPEM)
			anovelgrpc.NewTokenSource = grpcmocks.FakeTokenSource(new(grpcmocks.IDTokenStub))

			ss := setupStubServer(
				t,
				stubServerParams{insecure: false, authority: "127.0.0.1:8080", token: grpcmocks.DefaultToken},
			)
			clean, err := grpcmocks.FakeGRPCServer(ss, x509mocks.Server1KeyPEM, x509mocks.Server1CertPEM)
			require.NoError(t, err)
			defer clean()

			connPool := anovelgrpc.NewConnPool()
			defer connPool.Close()

			client, err := connPool.Open("127.0.0.1", 8080, anovelgrpc.ProtocolHTTPS)
			require.NoError(t, err)

			c := testgrpc.NewTestServiceClient(client)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			_, err = c.EmptyCall(ctx, new(testgrpc.Empty))
			require.Error(t, err)
		},
		MainFn: func(t *testing.T, res *testutils.CMDResult) {
			require.True(t, res.Success, res.STDErr, res.STDOut)
		},
		Env: []string{"ENV=dev"},
	})
}

func TestConnReleaseOK(t *testing.T) {
	// Run in isolation to avoid messing with flags.
	testutils.RunCMD(t, &testutils.CMDConfig{
		CmdFn: func(t *testing.T) {
			anovelgrpc.SystemCertPool = grpcmocks.FakeClientCerts(x509mocks.ServerCACertPEM)
			anovelgrpc.NewTokenSource = grpcmocks.FakeTokenSource(new(grpcmocks.IDTokenStub))

			ss := setupStubServer(
				t,
				stubServerParams{insecure: false, authority: "127.0.0.1:8080", token: grpcmocks.DefaultToken},
			)
			clean, err := grpcmocks.FakeGRPCServer(ss, x509mocks.Server1KeyPEM, x509mocks.Server1CertPEM)
			require.NoError(t, err)
			defer clean()

			connPool := anovelgrpc.NewConnPool()
			defer connPool.Close()

			client, err := connPool.Open("127.0.0.1", 8080, anovelgrpc.ProtocolHTTPS)
			require.NoError(t, err)

			c := testgrpc.NewTestServiceClient(client)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			_, err = c.EmptyCall(ctx, new(testgrpc.Empty))
			testutils.RequireGRPCCodesEqual(t, err, codes.OK)
		},
		MainFn: func(t *testing.T, res *testutils.CMDResult) {
			require.True(t, res.Success, res.STDErr, res.STDOut)
		},
		Env: []string{"ENV=staging"},
	})
}
