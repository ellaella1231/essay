CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    phone VARCHAR(20) UNIQUE NOT NULL,
    free_trial_count INT DEFAULT 3,
    subscription_expiry TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS grade_standards (
    grade VARCHAR(50) PRIMARY KEY,
    rubric_text TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS essays (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    image_url TEXT,
    prompt_text TEXT,
    student_content TEXT,
    perfect_version TEXT,
    score INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS errors (
    id SERIAL PRIMARY KEY,
    essay_id INT REFERENCES essays(id) ON DELETE CASCADE,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    original_segment TEXT NOT NULL,
    suggested_segment TEXT NOT NULL,
    explanation TEXT NOT NULL,
    is_learned BOOLEAN DEFAULT FALSE
);
