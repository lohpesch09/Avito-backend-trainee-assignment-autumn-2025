CREATE TABLE users (
    user_id VARCHAR(50) PRIMARY KEY,
    user_name VARCHAR(255) UNIQUE,
    team_name VARCHAR(255),
    is_active BOOLEAN DEFAULT true,

    CONSTRAINT users_team_name_fkey 
    FOREIGN KEY (team_name) REFERENCES teams(team_name)
    DEFERRABLE INITIALLY DEFERRED    
); 
