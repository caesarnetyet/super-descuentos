-- testdata/seed.sql

-- Insert test users
INSERT INTO users (id, name, email) VALUES
                                        ('123e4567-e89b-12d3-a456-426614174000', 'Test User 1', 'user1@example.com'),
                                        ('223e4567-e89b-12d3-a456-426614174000', 'Test User 2', 'user2@example.com'),
                                        ('323e4567-e89b-12d3-a456-426614174000', 'Test User 3', 'user3@example.com');

-- Insert test posts
INSERT INTO posts (id, title, description, url, author_id, likes, expire_time, creation_time) VALUES
                                                                                                  (
                                                                                                      '423e4567-e89b-12d3-a456-426614174000',
                                                                                                      'First Test Post',
                                                                                                      'This is the first test post description',
                                                                                                      'https://example.com/post1',
                                                                                                      '123e4567-e89b-12d3-a456-426614174000',
                                                                                                      10,
                                                                                                      datetime('now', '+1 day'),
                                                                                                      datetime('now')
                                                                                                  ),
                                                                                                  (
                                                                                                      '523e4567-e89b-12d3-a456-426614174000',
                                                                                                      'Second Test Post',
                                                                                                      'This is the second test post description',
                                                                                                      'https://example.com/post2',
                                                                                                      '223e4567-e89b-12d3-a456-426614174000',
                                                                                                      5,
                                                                                                      datetime('now', '+2 days'),
                                                                                                      datetime('now')
                                                                                                  ),
                                                                                                  (
                                                                                                      '623e4567-e89b-12d3-a456-426614174000',
                                                                                                      'Third Test Post',
                                                                                                      'This is the third test post description',
                                                                                                      'https://example.com/post3',
                                                                                                      '323e4567-e89b-12d3-a456-426614174000',
                                                                                                      15,
                                                                                                      datetime('now', '+3 days'),
                                                                                                      datetime('now')
                                                                                                  );