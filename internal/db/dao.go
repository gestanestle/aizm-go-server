package db

import (
	"context"
	"fmt"
	"gestanestle/aizm-server/internal/models"

	"github.com/jackc/pgx/v5"
)

func (d *Dao) PersistEvent(c models.Conditions) error {
	d.Mu.Lock()
	defer d.Mu.Unlock()

	db, err := conn.Acquire(context.Background())
	if err != nil {
		panic("Error acquiring connection to database.")
	}

	defer db.Release()

	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		fmt.Printf("conn.BeginTx \n%v", err)
	}
	defer tx.Rollback(context.Background())

	q := `INSERT INTO conditions (id, time, temp, humidity) VALUES ($1, $2, $3, $4);`
	_, err = tx.Exec(context.Background(), q, c.ID, c.Time, c.Temp, c.Humidity)

	if err != nil {
		fmt.Printf("tx.Exec \n%v", err)
	}

	if err := tx.Commit(context.Background()); err != nil {
		fmt.Printf("tx.Commit \n%v", err)
	}

	return nil
}
