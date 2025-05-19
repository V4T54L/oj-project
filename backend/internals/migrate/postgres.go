package migrate

import (
	"context"
	"database/sql"
	"time"
)

func MigratePostgres(ctx context.Context, db *sql.DB) error {
	query := schemaQuery + seedingQuery
	// _ = seedingQuery // skipping seeeding the database for prod
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	_, err := db.ExecContext(ctx, query)

	return err

}

const (
	schemaQuery = `
		-- users
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username TEXT UNIQUE NOT NULL,
			email TEXT UNIQUE NOT NULL,
			avatar_url TEXT,
			hashed_password TEXT NOT NULL
		);

		-- difficulties
		CREATE TABLE IF NOT EXISTS difficulties (
			id SERIAL PRIMARY KEY,
			name TEXT UNIQUE NOT NULL
		);

		-- tags
		CREATE TABLE IF NOT EXISTS tags (
			id SERIAL PRIMARY KEY,
			name TEXT UNIQUE NOT NULL
		);

		-- status_ids
		CREATE TABLE IF NOT EXISTS status_ids (
			id SERIAL PRIMARY KEY,
			name TEXT UNIQUE NOT NULL
		);

		-- programming_languages
		CREATE TABLE IF NOT EXISTS programming_languages (
			id SERIAL PRIMARY KEY,
			name TEXT UNIQUE NOT NULL
		);

		-- problems
		CREATE TABLE IF NOT EXISTS problems (
			id SERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			description TEXT NOT NULL,
			difficulty_id INT REFERENCES difficulties(id),
			acceptance_rate NUMERIC,
			constraints TEXT,
			time_limit_ms INT,
			memory_limit_kb INT
		);

		-- problem_examples
		CREATE TABLE IF NOT EXISTS problem_examples (
			id SERIAL PRIMARY KEY,
			problem_id INT REFERENCES problems(id) ON DELETE CASCADE,
			input TEXT,
			expected_output TEXT,
			explanation TEXT
		);

		-- problem_tags
		CREATE TABLE IF NOT EXISTS problem_tags (
			tag_id INT REFERENCES tags(id) ON DELETE CASCADE,
			problem_id INT REFERENCES problems(id) ON DELETE CASCADE,
			PRIMARY KEY (tag_id, problem_id)
		);

		-- code_submissions
		CREATE TABLE IF NOT EXISTS code_submissions (
			id SERIAL PRIMARY KEY,
			user_id INT REFERENCES users(id),
			problem_id INT REFERENCES problems(id),
			language_id INT REFERENCES programming_languages(id),
			code TEXT NOT NULL,
			status_id INT REFERENCES status_ids(id),
			runtime_ms INT,
			memory_kb INT,
			message TEXT,
			created_at TIMESTAMPTZ DEFAULT now()
		);

		-- code_solutions
		CREATE TABLE IF NOT EXISTS code_solutions (
			id SERIAL PRIMARY KEY,
			problem_id INT REFERENCES problems(id),
			language_id INT REFERENCES programming_languages(id),
			code TEXT NOT NULL,
			explanation TEXT
		);

		-- hidden_test_cases
		CREATE TABLE IF NOT EXISTS hidden_test_cases (
			id SERIAL PRIMARY KEY,
			problem_id INT REFERENCES problems(id) ON DELETE CASCADE,
			input TEXT NOT NULL,
			expected_output TEXT NOT NULL
		);

		-- Indexes
		CREATE INDEX IF NOT EXISTS idx_code_submissions_user_id ON code_submissions(user_id);
		CREATE INDEX IF NOT EXISTS idx_code_submissions_problem_id ON code_submissions(problem_id);
		CREATE INDEX IF NOT EXISTS idx_problems_difficulty_id ON problems(difficulty_id);
		CREATE INDEX IF NOT EXISTS idx_problem_tags_tag_id ON problem_tags(tag_id);
	`

	seedingQuery = `
		-- Seed difficulties
		INSERT INTO difficulties (name) VALUES
			('Easy'),
			('Medium'),
			('Hard')
		ON CONFLICT DO NOTHING;

		-- Seed status_ids
		INSERT INTO status_ids (name) VALUES
			('Pending'),
			('Accepted'),
			('Wrong Answer'),
			('Runtime Error'),
			('Compilation Error'),
			('Time Limit Exceeded')
		ON CONFLICT DO NOTHING;

		-- Seed programming_languages
		INSERT INTO programming_languages (name) VALUES
			('Python'),
			('Java'),
			('C++'),
			('JavaScript'),
			('Go'),
			('C#'),
			('Rust'),
			('Kotlin')
		ON CONFLICT DO NOTHING;

		-- Example tags (can be expanded)
		INSERT INTO tags (name) VALUES
			('Array'),
			('String'),
			('Binary Search'),
			('Dynamic Programming'),
			('Graph'),
			('Heap')
		ON CONFLICT DO NOTHING;
	`
)
