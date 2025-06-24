package postgres

import (
	"database/sql"
	"sync"

	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Postgres struct {
	dbConn *sql.DB
}

var once = sync.Once{}

func NewPostgres() *Postgres {
	var db *sql.DB
	once.Do(func() { db = createConnection() })
	return &Postgres{
		dbConn: db,
	}
}

func createConnection() *sql.DB {

	db, err := sql.Open("postgres", viper.GetString("postgresURL"))
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	logrus.Info("Connected to postgres", "URL", viper.GetString("postgresURL"))
	return db
}
