package db

import (
	"bytes"
	"context"
	"strings"
)

// LoadSchema is used to load a SQL schema.
func LoadSchema(schema []byte) {
	parts := bytes.Split(schema, []byte(";"))
	for _, v := range parts {
		stmt := strings.TrimSpace(string(v))
		if _, err := conn.Exec(context.TODO(), stmt); err != nil {
			panic(err)
		}
	}
}
