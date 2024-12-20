-- name: CreatePost :exec
INSERT INTO posts (id, title, description, url, author_id, likes, expire_time, creation_time)
VALUES (?, ?, ?, ?, ?, ?, ?, ?);

-- name: DeletePost :execresult
DELETE FROM posts WHERE id = ?;

-- name: UpdatePost :execresult
UPDATE posts
SET title = ?, description = ?, url = ?, likes = ?, expire_time = ?
WHERE id = ?;

-- name: GetPost :one
SELECT
    *
FROM
    posts
where id = ?;

-- name: GetPostsWithAuthor :many
SELECT
    sqlc.embed(posts), sqlc.embed(users)
FROM
    posts
join users
    on posts.author_id = users.id
ORDER BY
    creation_time DESC
LIMIT ? OFFSET ?;

-- name: GetUser :one
SELECT id, name, email
FROM users
WHERE id = ?;

-- name: CreateUser :exec
INSERT INTO users (id, name, email)
VALUES (?, ?, ?);

-- name: GetAuthors :many
SELECT id, name, email
FROM users
LIMIT ? OFFSET ?;

-- name: GetAuthorByEmail :one
SELECT id, name, email
FROM users
WHERE email = ?;

