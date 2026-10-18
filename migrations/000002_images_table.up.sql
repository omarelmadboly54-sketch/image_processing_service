CREATE TABLE IF NOT EXISTS images(
    id BIGSERIAL PRIMARY KEY,
    file_url  VARCHAR(255)NOT NULL,
    original_name VARCHAR(255) NOT NULL,
    file_size BIGINT not NULL,
    format VARCHAR(20)NOT NULL,
    width INT NOT NULL,
    height INT NOT NULL,
    user_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    Foreign Key (user_id) REFERENCES users(id) ON DELETE CASCADE
);
