package database_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/a-novel/golib/database"
)

func TestOpenDB(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping database test in short mode.")
	}

	t.Run("OK", func(t *testing.T) {
		dsn := "postgres://test:test@localhost:5432/test?sslmode=disable"

		db, closer, err := database.OpenDB(dsn)
		require.NoError(t, err)
		defer closer()

		// Try a simple query to check if connection is working.
		_, err = db.Exec("SELECT 1")
		require.NoError(t, err)

		// Close the connection.
		require.NoError(t, db.Close())

		// Try to query again.
		_, err = db.Exec("SELECT 1")
		require.Error(t, err)
	})

	t.Run("WrongDSN", func(t *testing.T) {
		dsn := "postgres://test:test@localhost:1234/test?sslmode=disable"

		_, _, err := database.OpenDB(dsn)
		require.Error(t, err)
	})
}
