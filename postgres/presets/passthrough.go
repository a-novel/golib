package postgrespresets

import (
	"context"

	"github.com/uptrace/bun"
)

type PassthroughConfig struct {
	db *bun.DB
}

func NewPassthroughConfig(db *bun.DB) PassthroughConfig {
	return PassthroughConfig{db: db}
}

func (config PassthroughConfig) DB() (*bun.DB, error) {
	return config.db, nil
}

func (config PassthroughConfig) RunMigrations(_ context.Context, _ *bun.DB) error {
	return nil
}
