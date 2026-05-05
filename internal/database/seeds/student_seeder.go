package seeds

import (
	"fmt"
	"log"
	"time"

	"backend-mosque-tahfidz-management/pkg/utils"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// StudentSeeder seeds the students table with sample student data
type StudentSeeder struct{}

func (s *StudentSeeder) TableName() string {
	return "students"
}

type studentSeed struct {
	Name          string
	Password      string
	ProfileImg    string
	CoverImg      string
	Age           int
	LearningLevel string
	Fluency       string
	Status        string
	Contact       string
	JoinDate      time.Time
}

func (s *StudentSeeder) Seed(db *sqlx.DB) error {
	// Get mentor IDs to assign students to mentors
	var mentorIDs []uuid.UUID
	err := db.Select(&mentorIDs, `SELECT id FROM users WHERE role = 'mentor' AND deleted_at IS NULL ORDER BY name`)
	if err != nil {
		return fmt.Errorf("failed to fetch mentors: %w", err)
	}

	if len(mentorIDs) == 0 {
		return fmt.Errorf("no mentors found — please seed users first")
	}

	students := []studentSeed{
		{
			Name: "Muhammad Rizki", Password: "student123",
			Age: 12, LearningLevel: "Juz 1-5", Fluency: "Lancar",
			Status: "active", Contact: "081234567001",
			JoinDate: time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			Name: "Aisyah Putri", Password: "student123",
			Age: 10, LearningLevel: "Juz 1-3", Fluency: "Cukup Lancar",
			Status: "active", Contact: "081234567002",
			JoinDate: time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			Name: "Abdullah Hasan", Password: "student123",
			Age: 14, LearningLevel: "Juz 6-10", Fluency: "Lancar",
			Status: "active", Contact: "081234567003",
			JoinDate: time.Date(2024, 8, 20, 0, 0, 0, 0, time.UTC),
		},
		{
			Name: "Fatimah Zahra", Password: "student123",
			Age: 11, LearningLevel: "Juz 1-5", Fluency: "Perlu Bimbingan",
			Status: "active", Contact: "081234567004",
			JoinDate: time.Date(2025, 3, 10, 0, 0, 0, 0, time.UTC),
		},
		{
			Name: "Umar Farooq", Password: "student123",
			Age: 13, LearningLevel: "Juz 11-15", Fluency: "Lancar",
			Status: "active", Contact: "081234567005",
			JoinDate: time.Date(2024, 6, 5, 0, 0, 0, 0, time.UTC),
		},
		{
			Name: "Khadijah Aminah", Password: "student123",
			Age: 9, LearningLevel: "Juz 1-3", Fluency: "Perlu Bimbingan",
			Status: "active", Contact: "081234567006",
			JoinDate: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			Name: "Ali Rahman", Password: "student123",
			Age: 15, LearningLevel: "Juz 16-20", Fluency: "Sangat Lancar",
			Status: "active", Contact: "081234567007",
			JoinDate: time.Date(2023, 11, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			Name: "Zainab Sari", Password: "student123",
			Age: 12, LearningLevel: "Juz 1-5", Fluency: "Cukup Lancar",
			Status: "inactive", Contact: "081234567008",
			JoinDate: time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			Name: "Hamza Pratama", Password: "student123",
			Age: 10, LearningLevel: "Juz 1-3", Fluency: "Lancar",
			Status: "active", Contact: "081234567009",
			JoinDate: time.Date(2025, 1, 20, 0, 0, 0, 0, time.UTC),
		},
		{
			Name: "Maryam Dewi", Password: "student123",
			Age: 13, LearningLevel: "Juz 6-10", Fluency: "Cukup Lancar",
			Status: "active", Contact: "081234567010",
			JoinDate: time.Date(2024, 7, 10, 0, 0, 0, 0, time.UTC),
		},
		{
			Name: "Yusuf Hakim", Password: "student123",
			Age: 11, LearningLevel: "Juz 1-5", Fluency: "Perlu Bimbingan",
			Status: "inactive", Contact: "081234567011",
			JoinDate: time.Date(2025, 2, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			Name: "Salma Nurul", Password: "student123",
			Age: 14, LearningLevel: "Juz 11-15", Fluency: "Lancar",
			Status: "active", Contact: "081234567012",
			JoinDate: time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	query := `INSERT INTO students (mentor_id, name, password, profile_img, cover_img, age, learning_level, fluency, status, contact, join_date) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	          ON CONFLICT DO NOTHING`

	for i, st := range students {
		hashedPassword, err := utils.HashPassword(st.Password)
		if err != nil {
			return fmt.Errorf("failed to hash password for %s: %w", st.Name, err)
		}

		// Distribute students across mentors in round-robin fashion
		mentorID := mentorIDs[i%len(mentorIDs)]

		_, err = db.Exec(query,
			mentorID, st.Name, hashedPassword,
			st.ProfileImg, st.CoverImg,
			st.Age, st.LearningLevel, st.Fluency, st.Status,
			st.Contact, st.JoinDate,
		)
		if err != nil {
			return fmt.Errorf("failed to insert student %s: %w", st.Name, err)
		}

		log.Printf("   → Created student: %s (age %d, level: %s, mentor: %s)", st.Name, st.Age, st.LearningLevel, mentorID)
	}

	return nil
}
