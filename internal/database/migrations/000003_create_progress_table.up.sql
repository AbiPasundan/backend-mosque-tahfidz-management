CREATE TABLE IF NOT EXISTS progress (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    mentor_id UUID NOT NULL REFERENCES users(id) ON DELETE SET NULL,
    surah VARCHAR(100) NOT NULL,
    status VARCHAR(50) NOT NULL,
    ayat_start INT NOT NULL,
    ayat_end INT NOT NULL,
    notes TEXT,
    progress_date DATE DEFAULT CURRENT_DATE
);

CREATE INDEX IF NOT EXISTS idx_progress_student_id ON progress(student_id);
CREATE INDEX IF NOT EXISTS idx_progress_date ON progress(progress_date);