// Package db provides a high-level, abstracted interface for all database operations.
package db

import (
	"csrspServer/db/sqlc"
	"database/sql"
)

// DB is the primary database interface for the application.
// It embeds sqlc.Queries to provide all generated query methods
// and holds the raw database connection.
type DB struct {
	*sqlc.Queries
	conn *sql.DB
}

var (
	// global is the unexported global instance of our database interface.
	// All functions in this package will operate on this instance.
	global *DB
)

// Init initializes the database connection and the global DB object.
// It should be called once at application startup from the main package.
func Init(user, password, dbName string, ips []string) error {
	dbConn, err := connect(user, password, dbName, ips)
	if err != nil {
		return err
	}

	global = &DB{
		Queries: sqlc.New(dbConn),
		conn:    dbConn,
	}
	return nil
}

// Close closes the database connection. It should be called on application shutdown.
func Close() {
	if global != nil && global.conn != nil {
		global.conn.Close()
	}
}
