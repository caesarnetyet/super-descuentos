CREATE TABLE users (
                       id TEXT PRIMARY KEY,
                       name TEXT NOT NULL,
                       email TEXT NOT NULL UNIQUE
);

CREATE TABLE posts (
                       id TEXT PRIMARY KEY,
                       title TEXT NOT NULL,
                       description TEXT NOT NULL,
                       url TEXT NOT NULL,
                       author_id TEXT NOT NULL,
                       likes INTEGER NOT NULL DEFAULT 0,
                       expire_time DATETIME NOT NULL,
                       creation_time DATETIME NOT NULL,
                       FOREIGN KEY (author_id) REFERENCES users (id) ON DELETE CASCADE
);

-- Indexes for faster queries
CREATE INDEX idx_posts_author_id ON posts (author_id);
CREATE INDEX idx_posts_creation_time ON posts (creation_time);
CREATE INDEX idx_posts_expire_time ON posts (expire_time);