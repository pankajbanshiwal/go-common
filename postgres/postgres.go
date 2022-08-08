package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/okcredit/go-common/config"
)

func New(conf config.Config) (*sql.DB, error) {
	host := conf.Get("database.host")
	port := conf.Get("database.port")
	username := conf.Get("database.username")
	password := conf.Get("database.password")
	dbname := conf.Get("database.dbname")
	dataSrc := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, username, password, dbname)
	db, err := sql.Open("postgres", dataSrc)
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(1)
	return db, nil
}
