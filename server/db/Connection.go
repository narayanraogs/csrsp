// Package db contains the database schema, queries, and connection logic.
package db

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// connect loops through a list of IPs and attempts to establish a database connection.
// It returns the first successful connection.
func connect(user, password, dbName string, ips []string) (*sql.DB, error) {
	// The DSN should have parseTime=true to properly handle DATETIME columns.
	const dsnFormat = "%s:%s@tcp(%s:3306)/%s?interpolateParams=true&parseTime=true"

	for _, ip := range ips {
		dsn := fmt.Sprintf(dsnFormat, user, password, ip, dbName)

		db, err := sql.Open("mysql", dsn)
		if err != nil {
			slog.Warn("Failed to create DB handle", "ip", ip, "error", err)
			continue // Try next IP
		}

		// Ping the database to verify the connection is alive.
		if err := db.Ping(); err != nil {
			slog.Warn("Failed to ping database", "ip", ip, "error", err)
			db.Close()
			continue // Try next IP
		}

		// Recommended settings for connection pooling.
		db.SetConnMaxLifetime(time.Minute * 3)
		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(10)

		slog.Info("Database connection successful", "ip", ip)
		return db, nil
	}
	return nil, fmt.Errorf("could not connect to any of the provided database IPs")
}
