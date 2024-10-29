package grpc

import (
	"context"
	"crypto/x509"
	_ "embed"
	"errors"
	"fmt"
	"sync"

	"google.golang.org/api/idtoken"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/credentials/oauth"

	"github.com/a-novel/golib/deploy"
)

var ErrConnectionPoolClosed = errors.New("connection pool is closed")

type Protocol string

const (
	ProtocolHTTP  Protocol = "http"
	ProtocolHTTPS Protocol = "https"
)

func (p *Protocol) WithAddr(host string) string {
	return fmt.Sprintf("%s://%s", *p, host)
}

//go:embed grpc-config.json
var grpcConfig string //nolint:unused

// SystemCertPool helps to mock the default behavior to retrieve system certificates. We do this because Go's
// implementation is OS specific.
var SystemCertPool = x509.SystemCertPool

// NewTokenSource helps to mock token source, otherwise tests would need to rely on actual Google Cloud
// credentials.
var NewTokenSource = idtoken.NewTokenSource

// ConnPool manages client connections to GRPC services. It is thread-safe, and should be shared for opening
// connections to different services.
type ConnPool interface {
	// Close terminates all connections in the pool. Use this method for graceful shutdowns.
	Close()
	// Open opens a new connection to an existing GRPC service.
	// This method automatically handles authentication under GCP environments.
	Open(host string, port int, protocol Protocol) (*grpc.ClientConn, error)
}

type connPoolImpl struct {
	conns []*grpc.ClientConn

	certs *x509.CertPool

	mu sync.Mutex

	closed bool
}

// Make sure the pool is properly initialized when used.
func (pool *connPoolImpl) ensureInit() error {
	// Load certificates from environment, if available.
	if pool.certs == nil {
		certs, err := SystemCertPool()
		if err != nil {
			return fmt.Errorf("load system root CA certificates: %w", err)
		}

		pool.certs = certs
	}

	return nil
}

func (pool *connPoolImpl) Close() {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	// No-op if already closed.
	if pool.closed {
		return
	}

	// Close all connections.
	for _, conn := range pool.conns {
		_ = conn.Close()
	}

	pool.closed = true
}

func (pool *connPoolImpl) getConnOptions(host string, port int, protocol Protocol) ([]grpc.DialOption, error) {
	pool.mu.Lock()
	defer pool.mu.Unlock()
	// Following configuration comes from official documentation.
	// https://cloud.google.com/run/docs/triggering/grpc?hl=fr

	// host should be of the form domain:port, e.g., example.com:443
	hostWithPort := fmt.Sprintf("%s:%d", host, port)
	// audience must be the auto-assigned URL of a Cloud Run service or HTTP Cloud Function without port number.
	audience := protocol.WithAddr(host)

	opts := []grpc.DialOption{grpc.WithAuthority(hostWithPort)}

	// Ignore authentication in non-release environments.
	if !deploy.IsReleaseEnv() {
		return append(opts, grpc.WithTransportCredentials(insecure.NewCredentials())), nil
	}

	// Instead of declaring a token source per request, use a global one. It has the benefit of
	// re-using and auto-refreshing tokens.
	tokenSource, err := NewTokenSource(context.Background(), audience)
	if err != nil {
		return nil, fmt.Errorf("create token source: %w", err)
	}

	// Basically the same thing as the docs, but allows us to override the trusted CA check when running
	// in tests.
	transport := credentials.NewClientTLSFromCert(pool.certs, hostWithPort)

	// Configure GRPC requests to be automatically authenticated, so Cloud credentials don't have to be
	// managed manually.
	return append(
		opts,
		grpc.WithTransportCredentials(transport),
		// Automate the step that adds the bearer token to the context of a request.
		grpc.WithPerRPCCredentials(oauth.TokenSource{TokenSource: tokenSource}),
	), nil
}

func (pool *connPoolImpl) Open(host string, port int, protocol Protocol) (*grpc.ClientConn, error) {
	// Ensure the pool has not been closed before trying anything.
	pool.mu.Lock()
	if pool.closed {
		pool.mu.Unlock()
		return nil, ErrConnectionPoolClosed
	}
	pool.mu.Unlock()

	// Make sure the pool is properly loaded. This settles the environment the first time it is called.
	if err := pool.ensureInit(); err != nil {
		return nil, fmt.Errorf("initialize connection pool: %w", err)
	}

	options, err := pool.getConnOptions(host, port, protocol)
	if err != nil {
		return nil, fmt.Errorf("get connection options: %w", err)
	}

	// TODO: uncomment when the app grows. This is temporarily deactivated because it drives GCP costs up.
	// options = append(options, grpc.WithDefaultServiceConfig(grpcConfig))

	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", host, port), options...)
	if err != nil {
		return nil, fmt.Errorf("open connection: %w", err)
	}

	// Append the new connection to the current pool.
	pool.mu.Lock()
	pool.conns = append(pool.conns, conn)
	pool.mu.Unlock()

	return conn, nil
}

// NewConnPool creates a new connection pool for GRPC services.
func NewConnPool() ConnPool {
	return &connPoolImpl{}
}
