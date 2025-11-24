CREATE TABLE pull_requests (
    pr_id VARCHAR(50) PRIMARY KEY,
    pr_name VARCHAR(255),
    author_id VARCHAR(50) REFERENCES users(user_id),
    status VARCHAR(50),
    created_at TIMESTAMP,
    merged_at TIMESTAMP NULL
);
