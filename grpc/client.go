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
// connections to the same service.
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

	if pool.closed {
		return
	}

	for _, conn := range pool.conns {
		_ = conn.Close()
	}

	pool.closed = true
}

func (pool *connPoolImpl) getConnOptions(host string, port int, protocol Protocol) ([]grpc.DialOption, error) {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	if !deploy.IsReleaseEnv() {
		return []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}, nil
	}

	tokenSource, err := NewTokenSource(context.Background(), protocol.WithAddr(host))
	if err != nil {
		return nil, fmt.Errorf("create token source: %w", err)
	}

	transport := credentials.NewClientTLSFromCert(pool.certs, fmt.Sprintf("%s:%d", host, port))

	// Configure GRPC requests to be automatically authenticated, so Cloud credentials don't have to be
	// managed manually.
	return []grpc.DialOption{
		grpc.WithTransportCredentials(transport),
		grpc.WithAuthority(fmt.Sprintf("%s:%d", host, port)),
		grpc.WithPerRPCCredentials(oauth.TokenSource{TokenSource: tokenSource}),
	}, nil
}

func (pool *connPoolImpl) Open(host string, port int, protocol Protocol) (*grpc.ClientConn, error) {
	pool.mu.Lock()
	if pool.closed {
		pool.mu.Unlock()
		return nil, ErrConnectionPoolClosed
	}
	pool.mu.Unlock()

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

	pool.mu.Lock()
	pool.conns = append(pool.conns, conn)
	pool.mu.Unlock()

	return conn, nil
}

// NewConnPool creates a new connection pool for GRPC services.
func NewConnPool() ConnPool {
	return &connPoolImpl{}
}
