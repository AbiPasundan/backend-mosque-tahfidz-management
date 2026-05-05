package seeds

import (
	"fmt"
	"log"

	"backend-mosque-tahfidz-management/pkg/utils"

	"github.com/jmoiron/sqlx"
)

// UserSeeder seeds the users table with admin and mentor accounts
type UserSeeder struct{}

func (s *UserSeeder) TableName() string {
	return "users"
}

type userSeed struct {
	Name     string
	Email    string
	Password string
	Role     string
}

func (s *UserSeeder) Seed(db *sqlx.DB) error {
	users := []userSeed{
		// Admin accounts
		{Name: "Admin Utama", Email: "admin@tahfidz.com", Password: "admin123", Role: "admin"},
		{Name: "Admin Sekretariat", Email: "sekretariat@tahfidz.com", Password: "admin123", Role: "admin"},

		// Mentor accounts
		{Name: "Ustadz Ahmad", Email: "ahmad@tahfidz.com", Password: "mentor123", Role: "mentor"},
		{Name: "Ustadz Farid", Email: "farid@tahfidz.com", Password: "mentor123", Role: "mentor"},
		{Name: "Ustadzah Aisyah", Email: "aisyah@tahfidz.com", Password: "mentor123", Role: "mentor"},
		{Name: "Ustadz Ibrahim", Email: "ibrahim@tahfidz.com", Password: "mentor123", Role: "mentor"},
		{Name: "Ustadzah Fatimah", Email: "fatimah@tahfidz.com", Password: "mentor123", Role: "mentor"},
	}

	query := `INSERT INTO users (name, email, password, role) 
	          VALUES ($1, $2, $3, $4) 
	          ON CONFLICT (email) DO NOTHING`

	for _, u := range users {
		hashedPassword, err := utils.HashPassword(u.Password)
		if err != nil {
			return fmt.Errorf("failed to hash password for %s: %w", u.Email, err)
		}

		result, err := db.Exec(query, u.Name, u.Email, hashedPassword, u.Role)
		if err != nil {
			return fmt.Errorf("failed to insert user %s: %w", u.Email, err)
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected > 0 {
			log.Printf("   → Created user: %s (%s) [%s]", u.Name, u.Email, u.Role)
		} else {
			log.Printf("   → Skipped user (already exists): %s (%s)", u.Name, u.Email)
		}
	}

	return nil
}
