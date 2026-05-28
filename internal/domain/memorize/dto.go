package memorize

import "github.com/google/uuid"

type CreateMemorizeRequest struct {
	StudentID   uuid.UUID `json:"student_id" validate:"required"`
	Surah       string    `json:"surah" validate:"required"`
	SurahNumber int       `json:"surah_number" validate:"required,min=1,max=114"`
	AyatStart   int       `json:"ayat_start" validate:"required,min=1"`
	AyatEnd     int       `json:"ayat_end" validate:"required,min=1,gtefield=AyatStart"`
	Status      string    `json:"status" validate:"required,oneof=memorized in_progress murojaah forgotten"`
	Notes       string    `json:"notes"`
}

type UpdateMemorizeStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=memorized in_progress murojaah forgotten"`
	Notes  string `json:"notes"`
}

type BulkUpdateStatusRequest struct {
	IDs    []uuid.UUID `json:"ids" validate:"required,min=1,dive,required"`
	Status string      `json:"status" validate:"required,oneof=memorized in_progress murojaah forgotten"`
}

type MemorizeResponse struct {
	ID             uuid.UUID  `json:"id"`
	StudentID      uuid.UUID  `json:"student_id"`
	VerifiedBy     *uuid.UUID `json:"verified_by,omitempty"`
	VerifierName   string     `json:"verifier_name,omitempty"`
	Surah          string     `json:"surah"`
	SurahNumber    int        `json:"surah_number"`
	AyatStart      int        `json:"ayat_start"`
	AyatEnd        int        `json:"ayat_end"`
	Status         string     `json:"status"`
	Notes          string     `json:"notes,omitempty"`
	MemorizedAt    *string    `json:"memorized_at,omitempty"`
	LastReviewedAt *string    `json:"last_reviewed_at,omitempty"`
}

type StudentMemorizeSummaryResponse struct {
	StudentID   uuid.UUID      `json:"student_id"`
	StudentName string         `json:"student_name"`
	Surahs      []SurahSummary `json:"surahs"`
}
