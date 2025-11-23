CREATE TABLE pr_reviewers (
    pr_id VARCHAR(50),
    user_id VARCHAR(50) REFERENCES users(user_id),
    PRIMARY KEY (pr_id, user_id),

    CONSTRAINT pr_reviewers_pr_id_fkey 
    FOREIGN KEY (pr_id) REFERENCES pull_requests(pr_id)
    DEFERRABLE INITIALLY DEFERRED     
);