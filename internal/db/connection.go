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

// test
//var url = "postgresql://postgres.yzgljcgjdvtzycdnsxnu:aizm.amin1234@aws-0-us-east-1.pooler.supabase.com:6543/postgres"

// prod
var url = "postgresql://postgres.lbummzbrrxzywwizdtej:aizm.admin1234@aws-0-ap-southeast-1.pooler.supabase.com:6543/postgres"

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
