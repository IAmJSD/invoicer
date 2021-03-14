package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"reflect"
	"time"
)

// InsertFailure is used to insert a failure into the database.
func InsertFailure(ctx context.Context, insertingErr error) error {
	message := reflect.TypeOf(&insertingErr).Elem().Name() + ": " + insertingErr.Error()
	_ = fmt.Errorf("⚠️ %s", message)
	_, err := conn.Exec(ctx, "INSERT INTO failures (failure_time, message) VALUES (NOW(), $1)", message)
	return err
}

// Failure is used to define a database failure.
type Failure struct {
	FailureTime time.Time `json:"failure_time"`
	Message     string    `json:"message"`
}

// FailureRows is used to get the rows from the failures table.
type FailureRows struct {
	rows pgx.Rows
}

// Close is used to close the pointer.
func (e FailureRows) Close() {
	e.rows.Close()
}

// Next passes through the Next function from pgx.Rows.
func (e FailureRows) Next() bool {
	return e.rows.Next()
}

// Err passes through the Err function from pgx.Rows.
func (e FailureRows) Err() error {
	return e.rows.Err()
}

// Scan is used to scan the row into a Failure object.
func (e FailureRows) Scan() (*Failure, error) {
	var x Failure
	if err := e.rows.Scan(&x.FailureTime, &x.Message); err != nil {
		return nil, err
	}
	return &x, nil
}

// All is used to get all of the failures from the database.
func (e FailureRows) All() ([]*Failure, error) {
	a := make([]*Failure, 0)
	for e.Next() {
		email, err := e.Scan()
		if err != nil {
			return nil, err
		}
		a = append(a, email)
	}
	return a, nil
}

// SelectAllFailures is used to select all of the failures from the database.
func SelectAllFailures(ctx context.Context) (*FailureRows, error) {
	rows, err := conn.Query(ctx, "SELECT failure_time, message FROM failures ORDER BY failure_time DESC")
	if err != nil {
		return nil, err
	}
	return &FailureRows{rows: rows}, nil
}
