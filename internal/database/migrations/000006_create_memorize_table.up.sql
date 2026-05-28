CREATE TABLE IF NOT EXISTS memorize (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    verified_by UUID REFERENCES users(id) ON DELETE SET NULL,
    surah VARCHAR(100) NOT NULL,
    surah_number INT NOT NULL,
    ayat_start INT NOT NULL,
    ayat_end INT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'in_progress',
    notes TEXT,
    memorized_at TIMESTAMP,
    last_reviewed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_memorize_student_id ON memorize(student_id);
CREATE INDEX idx_memorize_student_surah ON memorize(student_id, surah_number);
CREATE INDEX idx_memorize_status ON memorize(status);
