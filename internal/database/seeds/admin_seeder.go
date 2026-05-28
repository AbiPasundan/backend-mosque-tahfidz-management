package seeds

import (
	"fmt"
	"log"

	"backend-mosque-tahfidz-management/pkg/utils"

	"github.com/jmoiron/sqlx"
)

// AdminSeeder seeds the users table with a single admin account
type AdminSeeder struct{}

func (s *AdminSeeder) TableName() string {
	return "users"
}

func (s *AdminSeeder) Seed(db *sqlx.DB) error {
	admin := userSeed{
		Name:     "Admin",
		Email:    "admin@admin.com",
		Password: "admin123",
		Role:     "admin",
	}

	hashedPassword, err := utils.HashPassword(admin.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password for %s: %w", admin.Email, err)
	}

	query := `INSERT INTO users (name, email, password, role) 
	          VALUES ($1, $2, $3, $4) 
	          ON CONFLICT (email) DO NOTHING`

	result, err := db.Exec(query, admin.Name, admin.Email, hashedPassword, admin.Role)
	if err != nil {
		return fmt.Errorf("failed to insert admin %s: %w", admin.Email, err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected > 0 {
		log.Printf("   → Created admin: %s (%s)", admin.Name, admin.Email)
	} else {
		log.Printf("   → Skipped admin (already exists): %s (%s)", admin.Name, admin.Email)
	}

	return nil
}
