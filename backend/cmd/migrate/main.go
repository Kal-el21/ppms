// Command migrate applies or rolls back PPMS-Kal's SQL migrations (backend/migrations) against the
// configured Postgres database.
//
// Usage:
//
//	go run ./cmd/migrate -direction up            # apply all pending migrations
//	go run ./cmd/migrate -direction down           # roll back all migrations
//	go run ./cmd/migrate -direction up -steps 1    # apply the next single migration
//	go run ./cmd/migrate -direction down -steps 1  # roll back the last single migration
package main

import (
	"errors"
	"flag"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/Kal-el21/backend/configs"
)

func main() {
	direction := flag.String("direction", "up", "migration direction: up or down")
	steps := flag.Int("steps", 0, "number of steps to move (0 = all the way)")
	flag.Parse()

	cfg := configs.Load()
	if err := cfg.Validate(); err != nil {
		log.Fatalf("load config: %v", err)
	}

	m, err := migrate.New("file://migrations", cfg.MigrateURL())
	if err != nil {
		log.Fatalf("init migrate: %v", err)
	}

	switch *direction {
	case "up":
		if *steps > 0 {
			err = m.Steps(*steps)
		} else {
			err = m.Up()
		}
	case "down":
		if *steps > 0 {
			err = m.Steps(-*steps)
		} else {
			err = m.Down()
		}
	default:
		log.Fatalf("unknown -direction %q (expected up or down)", *direction)
	}

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("migrate %s: %v", *direction, err)
	}

	log.Printf("migrate %s: done", *direction)
}
