package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	// file driver for reading the migration files
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// MigrateDb migrates the database to the latest version
// NOTE: this is creating its own connection because calling m.Close closes the db connection
func MigrateDb(dbURI string) error {
	db, err := sql.Open("postgres", dbURI)
	if err != nil {
		return err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	migrationDirectory, err := migrationDir()
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationDirectory), "postgres", driver)
	if err != nil {
		return err
	}
	defer m.Close()

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

// getProjectRoot traverses upwards to find the project root directory
func getProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir, nil
		}

		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			break
		}
		dir = parentDir
	}

	return "", fmt.Errorf("project root not found")
}

func migrationDir() (string, error) {
	root, err := getProjectRoot()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, "database", "migrations"), nil
}
