package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var (
	testQueries *Queries
	testDB      *sql.DB
)

func TestMain(m *testing.M) {
	var err error
	connStr := "postgresql://root:1337@127.0.0.1:5433/simple_bank?sslmode=disable"
	testDB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Panicf("failed to make connection to db")
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
