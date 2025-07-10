package postgrespresets

import (
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

func (config PassthroughConfig) RunMigrations(_ bun.IDB) error {
	return nil
}
