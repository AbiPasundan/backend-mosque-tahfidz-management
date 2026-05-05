package seeds

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

// Seeder defines the interface for all seeders
type Seeder interface {
	Seed(db *sqlx.DB) error
	TableName() string
}

// RunAll executes all seeders in the correct order (respecting FK constraints)
func RunAll(db *sqlx.DB) error {
	seeders := []Seeder{
		&UserSeeder{},
		&StudentSeeder{},
		&ProgressSeeder{},
	}

	for _, s := range seeders {
		log.Printf("🌱 Seeding table: %s ...", s.TableName())
		if err := s.Seed(db); err != nil {
			return fmt.Errorf("failed to seed %s: %w", s.TableName(), err)
		}
		log.Printf("✅ Table %s seeded successfully", s.TableName())
	}

	log.Println("🎉 All seeds completed successfully!")
	return nil
}

// CleanAll truncates all tables in reverse order (respecting FK constraints)
func CleanAll(db *sqlx.DB) error {
	tables := []string{"progress", "students", "users"}

	for _, table := range tables {
		log.Printf("🧹 Cleaning table: %s ...", table)
		if _, err := db.Exec(fmt.Sprintf("DELETE FROM %s", table)); err != nil {
			return fmt.Errorf("failed to clean %s: %w", table, err)
		}
	}

	log.Println("🧹 All tables cleaned!")
	return nil
}

// FreshSeed cleans all tables then re-seeds everything
func FreshSeed(db *sqlx.DB) error {
	if err := CleanAll(db); err != nil {
		return err
	}
	return RunAll(db)
}
