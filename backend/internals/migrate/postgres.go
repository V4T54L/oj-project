package migrate

import (
	"context"
	"database/sql"
	"log"
	"time"
)

func MigratePostgres(ctx context.Context, db *sql.DB) error {
	query := schemaQuery + seedingQuery + devQuery
	_ = seedingQuery + devQuery // skipping seeding the database for prod
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	_, err := db.ExecContext(ctx, query)
	if err != nil {
		log.Println("Error migrating db")
	}

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

	devQuery = `
		-- Seed Users for Development
		INSERT INTO users (username, email, avatar_url, hashed_password) VALUES
			('devuser1', 'devuser1@example.com', 'https://avatar.com/devuser1.png', 'hashed_password_1'),
			('devuser2', 'devuser2@example.com', 'https://avatar.com/devuser2.png', 'hashed_password_2'),
			('devuser3', 'devuser3@example.com', 'https://avatar.com/devuser3.png', 'hashed_password_3')
		ON CONFLICT DO NOTHING;

		-- Seed Problems for Development
		INSERT INTO problems (title, description, difficulty_id, acceptance_rate, constraints, time_limit_ms, memory_limit_kb) VALUES
			('Two Sum', 'Given an array of integers, return indices of the two numbers such that they add up to a specific target. Input is taken via stdin.', 1, 0.9, 'N/A', 2000, 256),
			('Longest Substring Without Repeating Characters', 'Find the length of the longest substring without repeating characters. Input is taken via stdin.', 2, 0.75, '1 ≤ s.length ≤ 5 * 10^4', 3000, 512),
			('Merge Intervals', 'Given a collection of intervals, merge all overlapping intervals. Input is taken via stdin.', 3, 0.65, '1 ≤ intervals.length ≤ 10^4', 5000, 1024)
		ON CONFLICT DO NOTHING;

		-- Seed Problem Examples for Development (with stdin-like input and expected output)
		INSERT INTO problem_examples (problem_id, input, expected_output, explanation) VALUES
			(1, '4\n2 7 11 15\n9', '0 1', 'The sum of nums[0] and nums[1] equals the target 9.'),
			(2, '4\nabcabcbb', '3', 'The longest substring without repeating characters is "abc", with length 3.'),
			(3, '4\n1 3 5 7\n2 6 8 10\n3 5 7 9\n', '1 6\n8 10\n', 'Merging overlapping intervals: [1, 6] and [8, 10].')
		ON CONFLICT DO NOTHING;

		-- Seed Code Solutions for Development (standard code handling stdin)
		INSERT INTO code_solutions (problem_id, language_id, code, explanation) VALUES
			(1, 1, 'def twoSum(nums, target):\n    for i in range(len(nums)):\n        for j in range(i+1, len(nums)):\n            if nums[i] + nums[j] == target:\n                return [i, j]\n\ntarget = int(input())\nnums = list(map(int, input().split()))\nprint(twoSum(nums, target))', 'Solution for finding two indices in a list where the sum equals the target.'),
			(2, 2, 'def lengthOfLongestSubstring(s):\n    start, max_len = 0, 0\n    seen = {}\n    for end, char in enumerate(s):\n        if char in seen and seen[char] >= start:\n            start = seen[char] + 1\n        seen[char] = end\n        max_len = max(max_len, end - start + 1)\n    return max_len\n\ns = input().strip()\nprint(lengthOfLongestSubstring(s))', 'Solution to find the longest substring without repeating characters.'),
			(3, 3, 'def merge(intervals):\n    intervals.sort(key=lambda x: x[0])\n    merged = []\n    for interval in intervals:\n        if not merged or merged[-1][1] < interval[0]:\n            merged.append(interval)\n        else:\n            merged[-1][1] = max(merged[-1][1], interval[1])\n    return merged\n\nintervals = [list(map(int, input().split())) for _ in range(int(input()))]\nmerged = merge(intervals)\nfor interval in merged:\n    print(" ".join(map(str, interval)))', 'Solution to merge overlapping intervals, reads input from stdin.')
		ON CONFLICT DO NOTHING;

		-- Seed Tags for Development
		INSERT INTO problem_tags (tag_id, problem_id) VALUES
			(1, 1), -- Array
			(2, 1), -- String
			(3, 2), -- Dynamic Programming
			(4, 3)  -- Binary Search
		ON CONFLICT DO NOTHING;
	`
)
