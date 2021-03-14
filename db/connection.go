package db

import (
	"context"
	"github.com/jackc/pgx/v4"
	"os"
)

// Defines the connection.
var conn *pgx.Conn

// Creates the DB connection.
func init() {
	var err error
	if conn, err = pgx.Connect(context.TODO(), os.Getenv("CONNECTION_STRING")); err != nil {
		panic(err)
	}
	if err = conn.Ping(context.TODO()); err != nil {
		panic(err)
	}
}
