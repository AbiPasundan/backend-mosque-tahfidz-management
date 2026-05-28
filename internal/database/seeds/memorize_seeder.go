package seeds

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// MemorizeSeeder seeds the memorize table with sample hafalan snapshot records
type MemorizeSeeder struct{}

func (s *MemorizeSeeder) TableName() string {
	return "memorize"
}

type memorizeSeed struct {
	Surah       string
	SurahNumber int
	AyatStart   int
	AyatEnd     int
	Status      string
	Notes       string
	DaysAgo     int // relative to today
}

func (s *MemorizeSeeder) Seed(db *sqlx.DB) error {
	// Fetch active students with their assigned mentors
	var pairs []studentMentorPair
	err := db.Select(&pairs, `SELECT id, mentor_id FROM students WHERE deleted_at IS NULL ORDER BY name`)
	if err != nil {
		return fmt.Errorf("failed to fetch students: %w", err)
	}

	if len(pairs) == 0 {
		return fmt.Errorf("no students found — please seed students first")
	}

	// Templates for memorize records
	memorizeTemplates := []memorizeSeed{
		{Surah: "Al-Fatihah", SurahNumber: 1, AyatStart: 1, AyatEnd: 7, Status: "memorized", Notes: "Sangat lancar, makhraj sempurna", DaysAgo: 30},
		{Surah: "Al-Baqarah", SurahNumber: 2, AyatStart: 1, AyatEnd: 20, Status: "memorized", Notes: "Lancar dengan sedikit catatan tajwid", DaysAgo: 25},
		{Surah: "Al-Baqarah", SurahNumber: 2, AyatStart: 21, AyatEnd: 40, Status: "in_progress", Notes: "Sedang proses menyempurnakan hafalan", DaysAgo: 20},
		{Surah: "Al-Mulk", SurahNumber: 67, AyatStart: 1, AyatEnd: 30, Status: "memorized", Notes: "Khatam surah Al-Mulk dengan lancar", DaysAgo: 15},
		{Surah: "Yasin", SurahNumber: 36, AyatStart: 1, AyatEnd: 20, Status: "murojaah", Notes: "Perlu diulang-ulang agar lebih lancar", DaysAgo: 10},
		{Surah: "Yasin", SurahNumber: 36, AyatStart: 21, AyatEnd: 40, Status: "forgotten", Notes: "Lupa beberapa ayat tengah, perlu dihafal ulang", DaysAgo: 5},
	}

	query := `
		INSERT INTO memorize (
			id, student_id, verified_by, surah, surah_number, 
			ayat_start, ayat_end, status, notes, 
			memorized_at, last_reviewed_at, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	now := time.Now()
	count := 0

	for i, pair := range pairs {
		// Each student gets 3-5 memorize records
		numRecords := 3 + (i % 3) // varies between 3, 4, 5
		startIdx := (i * 2) % len(memorizeTemplates)

		for j := 0; j < numRecords; j++ {
			tmpl := memorizeTemplates[(startIdx+j)%len(memorizeTemplates)]
			createdAt := now.AddDate(0, 0, -tmpl.DaysAgo)

			var memorizedAt *time.Time
			if tmpl.Status == "memorized" {
				mTime := createdAt
				memorizedAt = &mTime
			}

			var lastReviewedAt *time.Time
			if tmpl.Status == "murojaah" {
				rTime := createdAt
				lastReviewedAt = &rTime
			}

			id := uuid.New()
			_, err := db.Exec(
				query,
				id,
				pair.StudentID,
				pair.MentorID,
				tmpl.Surah,
				tmpl.SurahNumber,
				tmpl.AyatStart,
				tmpl.AyatEnd,
				tmpl.Status,
				tmpl.Notes,
				memorizedAt,
				lastReviewedAt,
				createdAt,
				createdAt,
			)
			if err != nil {
				return fmt.Errorf("failed to insert memorize record for student %s: %w", pair.StudentID, err)
			}
			count++
		}

		log.Printf("   → Created %d memorize records for student: %s", numRecords, pair.StudentID)
	}

	log.Printf("   → Total memorize records created: %d", count)
	return nil
}
