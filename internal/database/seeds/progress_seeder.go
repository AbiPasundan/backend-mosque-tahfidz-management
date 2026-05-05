package seeds

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// ProgressSeeder seeds the progress table with sample hafalan progress records
type ProgressSeeder struct{}

func (s *ProgressSeeder) TableName() string {
	return "progress"
}

type studentMentorPair struct {
	StudentID uuid.UUID `db:"id"`
	MentorID  uuid.UUID `db:"mentor_id"`
}

type progressSeed struct {
	Surah        string
	Status       string
	AyatStart    int
	AyatEnd      int
	Notes        string
	DaysAgo      int // relative to today, for realistic dates
}

func (s *ProgressSeeder) Seed(db *sqlx.DB) error {
	// Fetch active students with their assigned mentors
	var pairs []studentMentorPair
	err := db.Select(&pairs, `SELECT id, mentor_id FROM students WHERE deleted_at IS NULL ORDER BY name`)
	if err != nil {
		return fmt.Errorf("failed to fetch students: %w", err)
	}

	if len(pairs) == 0 {
		return fmt.Errorf("no students found — please seed students first")
	}

	// Template progress records — each student will get a few of these
	progressTemplates := []progressSeed{
		{Surah: "Al-Fatihah", Status: "done", AyatStart: 1, AyatEnd: 7, Notes: "Bacaan sangat baik, tajwid sempurna", DaysAgo: 30},
		{Surah: "Al-Baqarah", Status: "done", AyatStart: 1, AyatEnd: 20, Notes: "Hafalan lancar, perlu perbaikan makhroj huruf", DaysAgo: 25},
		{Surah: "Al-Baqarah", Status: "done", AyatStart: 21, AyatEnd: 40, Notes: "Sangat baik, lanjut ke ayat berikutnya", DaysAgo: 20},
		{Surah: "Al-Baqarah", Status: "pending", AyatStart: 41, AyatEnd: 60, Notes: "Masih perlu pengulangan, belum lancar", DaysAgo: 15},
		{Surah: "Ali Imran", Status: "done", AyatStart: 1, AyatEnd: 15, Notes: "Progres baik, tajwid perlu perhatian", DaysAgo: 10},
		{Surah: "Ali Imran", Status: "pending", AyatStart: 16, AyatEnd: 30, Notes: "Sedang dalam proses menghafal", DaysAgo: 5},
		{Surah: "An-Nisa", Status: "done", AyatStart: 1, AyatEnd: 10, Notes: "Hafalan kuat, perbaikan pada idgham", DaysAgo: 3},
		{Surah: "An-Nisa", Status: "pending", AyatStart: 11, AyatEnd: 20, Notes: "Baru mulai hafalan surah ini", DaysAgo: 1},
		{Surah: "Al-Mulk", Status: "done", AyatStart: 1, AyatEnd: 30, Notes: "Khatam surah Al-Mulk, sangat lancar", DaysAgo: 12},
		{Surah: "Yasin", Status: "done", AyatStart: 1, AyatEnd: 20, Notes: "Progres cepat, tajwid baik", DaysAgo: 7},
		{Surah: "Yasin", Status: "pending", AyatStart: 21, AyatEnd: 40, Notes: "Masih menghafal, perlu muroja'ah", DaysAgo: 2},
		{Surah: "Ar-Rahman", Status: "done", AyatStart: 1, AyatEnd: 30, Notes: "Bacaan indah, sangat bagus", DaysAgo: 18},
	}

	query := `INSERT INTO progress (student_id, mentor_id, surah, status, ayat_start, ayat_end, notes, progress_date) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	now := time.Now()
	count := 0

	for i, pair := range pairs {
		// Each student gets 3-5 progress records from the template pool
		numRecords := 3 + (i % 3) // varies between 3, 4, 5
		startIdx := (i * 2) % len(progressTemplates)

		for j := 0; j < numRecords; j++ {
			tmpl := progressTemplates[(startIdx+j)%len(progressTemplates)]
			progressDate := now.AddDate(0, 0, -tmpl.DaysAgo)

			_, err := db.Exec(query,
				pair.StudentID, pair.MentorID,
				tmpl.Surah, tmpl.Status,
				tmpl.AyatStart, tmpl.AyatEnd,
				tmpl.Notes, progressDate,
			)
			if err != nil {
				return fmt.Errorf("failed to insert progress for student %s: %w", pair.StudentID, err)
			}
			count++
		}

		log.Printf("   → Created %d progress records for student: %s", numRecords, pair.StudentID)
	}

	log.Printf("   → Total progress records created: %d", count)
	return nil
}
