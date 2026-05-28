package memorize

import (
	"time"

	"github.com/google/uuid"
)

type Memorize struct {
	ID             uuid.UUID  `db:"id"`
	StudentID      uuid.UUID  `db:"student_id"`
	VerifiedBy     *uuid.UUID `db:"verified_by"`
	VerifierName   string     `db:"verifier_name"`
	Surah          string     `db:"surah"`
	SurahNumber    int        `db:"surah_number"`
	AyatStart      int        `db:"ayat_start"`
	AyatEnd        int        `db:"ayat_end"`
	Status         string     `db:"status"`
	Notes          string     `db:"notes"`
	MemorizedAt    *time.Time `db:"memorized_at"`
	LastReviewedAt *time.Time `db:"last_reviewed_at"`
	CreatedAt      time.Time  `db:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at"`
}

type SurahSummary struct {
	SurahNumber    int    `db:"surah_number" json:"surah_number"`
	Surah          string `db:"surah" json:"surah"`
	TotalAyat      int    `json:"total_ayat"`
	MemorizedCount int    `json:"memorized_count"`
	Status         string `json:"status"`
}
