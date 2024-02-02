package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

func InitDB() *pgx.Conn {
	conn := connectToDB()
	if conn == nil {
		log.Panic("Can't connect to database")
	}
	return conn
}

func connectToDB() *pgx.Conn {
	counts := 0
	dns := os.Getenv("DSN")

	for {
		connection, err := openDB(dns)
		if err != nil {
			log.Println("Postgres not yet ready...")
			fmt.Println(err)
		} else {
			log.Println("Connected to database")
			return connection
		}

		if counts > 10 {
			return nil
		}
		log.Println("Backing off for 1 second")
		time.Sleep(time.Second)

		counts++
	}
}

func openDB(dns string) (*pgx.Conn, error) {
	fmt.Println("dns: ", dns)
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, dns)

	if err != nil {
		return nil, err
	}

	err = conn.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return conn, nil

}
