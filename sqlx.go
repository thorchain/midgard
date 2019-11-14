package main

import (
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"
)

type Event struct {
	Total int64 `db:"total"`
}

func main() {
	db, err := sqlx.Connect("postgres", "user=postgres password=password dbname=midgard sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	var e Event
	row := db.QueryRowx("SELECT SUM(coin) AS total FROM (SELECT (jsonb_array_elements(in_tx->'coins')->>'amount')::numeric AS coin FROM events) x;")
	err = row.StructScan(&e)

	fmt.Printf("%v\n", e.Total)
}
