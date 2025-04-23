package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DataBase struct {
	conn *sqlx.DB
	cfg  config
}

type config struct {
	Dsn string `env:"PQ_DSN,required,notEmpty" envDefault:"host=127.0.0.1 port=5432 user=auth_user password=auth_password dbname=auth_db sslmode=disable"`
}
