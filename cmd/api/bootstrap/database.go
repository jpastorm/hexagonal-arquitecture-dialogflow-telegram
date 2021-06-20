package bootstrap

import (
	"database/sql"
	"log"
)

func newPSQLDatabase(conf Configuration) *sql.DB {
	psqlDB, err := NewPsql(
		NewDBConfig(
			conf.Database.Engine,
			conf.Database.User,
			conf.Database.Password,
			conf.Database.Server,
			conf.Database.Port,
			conf.Database.Name,
			conf.Database.SSLMode,
		),
	)
	if err != nil {
		log.Fatalf("cannot connect to Postgres database: %v", err)
	}

	return psqlDB.GetConnection()
}
