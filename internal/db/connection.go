package db

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var conn *pgxpool.Pool

var url = os.Getenv("DATABASE_URL")

func NewCon() *pgxpool.Pool {
	var err error

	conf, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Fatalf("Unable to parse connection string: %v", err)
	}
	conf.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeCacheDescribe

	conn, err = pgxpool.NewWithConfig(context.Background(), conf)
	if err != nil {
		log.Printf("Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	log.Println("Created connection to the database...")

	return conn
}

type Dao struct {
	Mu sync.Mutex
}
