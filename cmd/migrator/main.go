package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
)

func main() {
	var dsn, path string

	flag.StringVar(&dsn, "dsn", "", "dsn string")
	flag.StringVar(&path, "path", "", "path to migrations")
	flag.Parse()

	if dsn == "" {
		panic("dsn is required")
	}
	if path == "" {
		panic("path is required")
	}

	m, err := migrate.New(
		"file://"+path,
		dsn,
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")

			return
		}
		panic(err)
	}

	fmt.Println("migrations applied")
}
