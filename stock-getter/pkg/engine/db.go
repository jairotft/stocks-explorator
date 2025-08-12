package engine

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

func connectToDB() (*pgx.Conn, error) {

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	database := os.Getenv("DB_DATABASE")

	url := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=require", user, password, host, port, database)
	dsn := fmt.Sprintf(url)

	db, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}

	fmt.Println("Conectado a CockroachDB")

	return db, nil

}
