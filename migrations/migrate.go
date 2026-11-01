// Package migrations applies the versioned SQL files embedded in this directory.
package migrations

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"sort"
	"strings"
	"time"
)

//go:embed *.sql
var files embed.FS

// Run applies all pending migrations for up, or reverts the latest one for down.
// Each SQL file and its history update run in one transaction. A session advisory
// lock serializes migrations across API instances and the standalone command.
func Run(ctx context.Context, db *sql.DB, direction string) error {
	if direction != "up" && direction != "down" {
		return fmt.Errorf("migration direction must be up or down")
	}
	conn, err := db.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	if _, err := conn.ExecContext(ctx, "SELECT pg_advisory_lock(823491)"); err != nil {
		return err
	}
	defer func() {
		unlockCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, _ = conn.ExecContext(unlockCtx, "SELECT pg_advisory_unlock(823491)")
	}()
	if _, err := conn.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations (
        name TEXT PRIMARY KEY, applied_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
    )`); err != nil {
		return err
	}
	rows, err := conn.QueryContext(ctx, "SELECT name FROM schema_migrations ORDER BY name")
	if err != nil {
		return err
	}
	applied := make(map[string]bool)
	var latest string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			rows.Close()
			return err
		}
		applied[name] = true
		latest = name
	}
	err = rows.Err()
	rows.Close()
	if err != nil {
		return err
	}
	entries, err := files.ReadDir(".")
	if err != nil {
		return err
	}
	var pending []string
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".up.sql") {
			name := strings.TrimSuffix(entry.Name(), ".up.sql")
			if !applied[name] {
				pending = append(pending, name)
			}
		}
	}
	sort.Strings(pending)
	if direction == "down" {
		if latest == "" {
			return nil
		}
		pending = []string{latest}
	}
	for _, name := range pending {
		if err := apply(ctx, conn, name, direction); err != nil {
			return fmt.Errorf("migration %s %s: %w", name, direction, err)
		}
	}
	return nil
}

func apply(ctx context.Context, conn *sql.Conn, name, direction string) error {
	contents, err := files.ReadFile(name + "." + direction + ".sql")
	if err != nil {
		return err
	}
	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(ctx, string(contents)); err != nil {
		return err
	}
	if direction == "up" {
		_, err = tx.ExecContext(ctx, "INSERT INTO schema_migrations (name) VALUES ($1)", name)
	} else {
		_, err = tx.ExecContext(ctx, "DELETE FROM schema_migrations WHERE name = $1", name)
	}
	if err != nil {
		return err
	}
	return tx.Commit()
}
