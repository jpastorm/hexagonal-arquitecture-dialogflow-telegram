package bootstrap

import (
	"database/sql"
	"fmt"
)

type Psql struct {
	db *sql.DB
}

func NewPsql(config DBConfig) (*Psql, error) {
	if config.SSLMode == "" {
		config.SSLMode = "disable"
	}

	dns := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		config.User,
		config.Password,
		config.Server,
		config.Port,
		config.DBName,
		config.SSLMode,
	)

	dbCon, err := sql.Open("postgres", dns)
	return &Psql{db: dbCon}, err
}

func (p *Psql) GetConnection() *sql.DB {
	return p.db
}

type DBConfig struct {
	Driver   string
	User     string
	Password string
	Server   string
	Port     uint
	DBName   string
	SSLMode  string
}

func NewDBConfig(driver string, user string, password string, server string, port uint, DBName string, SSLMode string) DBConfig {
	return DBConfig{Driver: driver, User: user, Password: password, Server: server, Port: port, DBName: DBName, SSLMode: SSLMode}
}