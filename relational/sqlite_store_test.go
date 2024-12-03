package relational_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"super-descuentos/relational"
	"testing"
	"time"

	_ "modernc.org/sqlite"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"super-descuentos/model"
)

// testDB wraps the common test database structures
type testDB struct {
	db    *sql.DB
	store *relational.SQLStore
	ctx   context.Context
}

// testData contains the seed data references
type testData struct {
	users []model.User
	posts []model.Post
}

// loadSQLFile reads and executes a SQL file
func loadSQLFile(t *testing.T, db *sql.DB, filename string) {
	path := filepath.Join(filename)
	content, err := os.ReadFile(path)
	require.NoError(t, err, "Failed to read %s", filename)

	_, err = db.Exec(string(content))
	require.NoError(t, err, "Failed to execute %s", filename)
}

// setupTestDB initializes a test database with schema and seed data
func setupTestDB(t *testing.T) *testDB {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)

	// Load schema and seed data
	loadSQLFile(t, db, "./schema/schema.sql")
	loadSQLFile(t, db, "./seed/seed.sql")

	return &testDB{
		db:    db,
		store: relational.NewSQLStore(db),
		ctx:   context.Background(),
	}
}

// loadTestData reads the seeded data into Go structs
func loadTestData(t *testing.T, tdb *testDB) testData {
	var data testData

	// Load seeded users
	rows, err := tdb.db.Query("SELECT id, name, email FROM users")
	require.NoError(t, err)
	defer rows.Close()

	for rows.Next() {
		var user model.User
		var idStr string
		err := rows.Scan(&idStr, &user.Name, &user.Email)
		require.NoError(t, err)
		user.ID, err = uuid.Parse(idStr)
		require.NoError(t, err)
		data.users = append(data.users, user)
	}

	// Load seeded posts
	rows, err = tdb.db.Query(`
		SELECT p.id, p.title, p.description, p.url, p.author_id, p.likes, p.expire_time, p.creation_time,
		       u.name, u.email
		FROM posts p
		JOIN users u ON p.author_id = u.id
	`)
	require.NoError(t, err)
	defer rows.Close()

	for rows.Next() {
		var post model.Post
		var idStr, authorIDStr string
		var authorName, authorEmail string
		err := rows.Scan(
			&idStr,
			&post.Title,
			&post.Description,
			&post.Url,
			&authorIDStr,
			&post.Likes,
			&post.ExpireTime,
			&post.CreationTime,
			&authorName,
			&authorEmail,
		)
		require.NoError(t, err)
		post.ID, err = uuid.Parse(idStr)
		require.NoError(t, err)
		post.Author.ID, err = uuid.Parse(authorIDStr)
		require.NoError(t, err)
		post.Author.Name = authorName
		post.Author.Email = authorEmail
		data.posts = append(data.posts, post)
	}

	return data
}

// createTestPost is a helper to create a test post
func createTestPost(user model.User) model.Post {
	now := time.Now().UTC()
	return model.Post{
		ID:           uuid.New(),
		Title:        "Test Post",
		Description:  "Test Description",
		Url:          "https://example.com",
		Author:       user,
		Likes:        0,
		ExpireTime:   now.Add(24 * time.Hour),
		CreationTime: now,
	}
}

func TestSQLStoreIntegration(t *testing.T) {
	tdb := setupTestDB(t)
	defer tdb.db.Close()

	data := loadTestData(t, tdb)

	testCases := []struct {
		name     string
		testFunc func(t *testing.T, tdb *testDB, data testData)
	}{
		{
			name: "Create and Get Post",
			testFunc: func(t *testing.T, tdb *testDB, data testData) {
				post := createTestPost(data.users[0])

				err := tdb.store.CreatePost(tdb.ctx, post)
				require.NoError(t, err)

				fetchedPost, err := tdb.store.GetPost(tdb.ctx, post.ID)
				require.NoError(t, err)
				require.Equal(t, post.ID, fetchedPost.ID)
				require.Equal(t, post.Title, fetchedPost.Title)
				require.Equal(t, post.Author.ID, fetchedPost.Author.ID)
			},
		},
		{
			name: "Update Post",
			testFunc: func(t *testing.T, tdb *testDB, data testData) {
				post := data.posts[0]
				updatedPost := post
				updatedPost.Title = "Updated Title"
				updatedPost.Description = "Updated Description"
				updatedPost.Likes = 20

				err := tdb.store.UpdatePost(tdb.ctx, post.ID, updatedPost)
				require.NoError(t, err)

				fetchedPost, err := tdb.store.GetPost(tdb.ctx, post.ID)
				require.NoError(t, err)
				require.Equal(t, updatedPost.Title, fetchedPost.Title)
				require.Equal(t, updatedPost.Description, fetchedPost.Description)
				require.Equal(t, updatedPost.Likes, fetchedPost.Likes)
			},
		},
		{
			name: "Get Multiple Posts",
			testFunc: func(t *testing.T, tdb *testDB, data testData) {
				// Create additional posts
				for i := 0; i < 3; i++ {
					post := createTestPost(data.users[i%2])
					post.Title = fmt.Sprintf("Additional Post %d", i)
					err := tdb.store.CreatePost(tdb.ctx, post)
					require.NoError(t, err)
				}

				// Test pagination
				posts, err := tdb.store.GetPosts(tdb.ctx, 0, 3)
				require.NoError(t, err)
				require.Len(t, posts, 3)

				posts, err = tdb.store.GetPosts(tdb.ctx, 3, 2)
				require.NoError(t, err)
				require.NotEmpty(t, posts)
			},
		},
		{
			name: "Delete Post",
			testFunc: func(t *testing.T, tdb *testDB, data testData) {
				post := data.posts[1]

				err := tdb.store.DeletePost(tdb.ctx, post.ID)
				require.NoError(t, err)

				_, err = tdb.store.GetPost(tdb.ctx, post.ID)
				require.Error(t, err)
			},
		},
		{
			name: "Error Cases",
			testFunc: func(t *testing.T, tdb *testDB, data testData) {
				// Test getting non-existent post
				_, err := tdb.store.GetPost(tdb.ctx, uuid.New())
				require.Error(t, err)

				// Test creating post with non-existent author
				invalidPost := createTestPost(model.User{ID: uuid.New()})
				err = tdb.store.CreatePost(tdb.ctx, invalidPost)
				require.Error(t, err)

				// Test updating non-existent post
				err = tdb.store.UpdatePost(tdb.ctx, uuid.New(), data.posts[0])
				require.Error(t, err)

				// Test deleting non-existent post
				err = tdb.store.DeletePost(tdb.ctx, uuid.New())
				require.Error(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.testFunc(t, tdb, data)
		})
	}
}
