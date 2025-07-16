package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/uptrace/bun"
)

const PingTimeout = 10 * time.Second

// Ping a database connection until it succeeds or the timeout is reached.
func Ping(ctx context.Context, client *bun.DB) error {
	start := time.Now()

	for err := client.PingContext(ctx); err != nil; err = client.PingContext(ctx) {
		if time.Since(start) > PingTimeout {
			return fmt.Errorf("ping database: %w", err)
		}
	}

	return nil
}
