package database_test

import (
	"context"
	"github.com/a-novel/golib/database"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestFreezeTime(t *testing.T) {
	db, closer, err := database.OpenTestDB(nil)
	require.NoError(t, err)
	defer closer()

	t.Log("GetTime/Legacy")
	now := time.Now()

	var dbTime time.Time
	require.NoError(t, db.NewSelect().ColumnExpr("NOW()").Scan(context.Background(), &dbTime))

	require.True(t, now.Before(dbTime))

	t.Log("FreezeTime")

	require.NoError(t, database.FreezeTime(db, time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)))
	require.NoError(t, db.NewSelect().ColumnExpr("NOW()").Scan(context.Background(), &dbTime))

	require.Equal(t, time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), dbTime)

	t.Log("UnfreezeTime")

	require.NoError(t, database.RestoreTime(db))
	require.NoError(t, db.NewSelect().ColumnExpr("NOW()").Scan(context.Background(), &dbTime))

	require.True(t, now.Before(dbTime))
}
