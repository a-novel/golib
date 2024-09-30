package databasemocks

import "embed"

//go:embed 20200101120000_migration_1.down.sql 20200101120000_migration_1.up.sql 20200101130000_migration_2.down.sql 20200101130000_migration_2.up.sql 20200101140000_migration_3.down.sql 20200101140000_migration_3.up.sql
var MigrationsAll embed.FS

//go:embed 20200101120000_migration_1.down.sql 20200101120000_migration_1.up.sql 20200101130000_migration_2.down.sql 20200101130000_migration_2.up.sql
var MigrationsGroup1 embed.FS
