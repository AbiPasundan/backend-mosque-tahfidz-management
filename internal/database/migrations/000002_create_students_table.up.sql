CREATE TABLE IF NOT EXISTS students (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    mentor_id UUID REFERENCES users(id) ON DELETE SET NULL,
    name VARCHAR(255) NOT NULL,
    password TEXT NOT NULL,
    profile_img VARCHAR(255),
    cover_img VARCHAR(255),
    age INT,
    learning_level VARCHAR(100),
    fluency VARCHAR(100),
    status VARCHAR(50),
    contact VARCHAR(50),
    join_date DATE DEFAULT CURRENT_DATE,
    last_progress TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);