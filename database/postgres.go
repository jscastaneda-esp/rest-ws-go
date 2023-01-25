package database

import (
	"context"
	"database/sql"
	"log"

	"github.com/jscastaneda-esp/rest-ws-go/models"
	_ "github.com/lib/pq"
)

type PostgresRepository struct {
	db *sql.DB
}

func (repo *PostgresRepository) Close() error {
	return repo.db.Close()
}

func (repo *PostgresRepository) InsertUser(ctx context.Context, user *models.User) error {
	_, err := repo.db.ExecContext(ctx, "INSERT INTO users (id, email, password) VALUES ($1, $2, $3)", user.Id, user.Email, user.Password)
	return err
}

func (repo *PostgresRepository) GetUserById(ctx context.Context, id string) (*models.User, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT id, email FROM users WHERE id = $1 LIMIT 1", id)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			log.Println(err)
		}
	}()

	if err := rows.Err(); err != nil {
		return nil, err
	}

	user := new(models.User)
	if rows.Next() {
		if err := rows.Scan(&user.Id, &user.Email); err != nil {
			return nil, err
		}
	}

	return user, nil
}

func (repo *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT id, email, password FROM users WHERE email = $1 LIMIT 1", email)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			log.Println(err)
		}
	}()

	if err := rows.Err(); err != nil {
		return nil, err
	}

	user := new(models.User)
	if rows.Next() {
		if err := rows.Scan(&user.Id, &user.Email, &user.Password); err != nil {
			return nil, err
		}
	}

	return user, nil
}

func (repo *PostgresRepository) InsertPost(ctx context.Context, post *models.Post) error {
	_, err := repo.db.ExecContext(ctx, "INSERT INTO posts (id, post_content, user_id) VALUES ($1, $2, $3)", post.Id, post.PostContent, post.UserId)
	return err
}

func (repo *PostgresRepository) GetPostById(ctx context.Context, id string) (*models.Post, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT id, post_content, created_at, user_id FROM posts WHERE id = $1 LIMIT 1", id)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			log.Println(err)
		}
	}()

	if err := rows.Err(); err != nil {
		return nil, err
	}

	post := new(models.Post)
	if rows.Next() {
		if err := rows.Scan(&post.Id, &post.PostContent, &post.CreatedAt, &post.UserId); err != nil {
			return nil, err
		}
	}

	return post, nil
}

func (repo *PostgresRepository) ListPosts(ctx context.Context, page uint64, rowsFetch uint64) ([]*models.Post, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT id, post_content, created_at, user_id FROM posts LIMIT $1 OFFSET $2", rowsFetch, (page-1)*rowsFetch)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			log.Println(err)
		}
	}()

	if err := rows.Err(); err != nil {
		return nil, err
	}

	posts := []*models.Post{}
	for rows.Next() {
		post := new(models.Post)
		if err := rows.Scan(&post.Id, &post.PostContent, &post.CreatedAt, &post.UserId); err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (repo *PostgresRepository) UpdatePost(ctx context.Context, post *models.Post) error {
	_, err := repo.db.ExecContext(ctx, "UPDATE posts SET post_content = $1 WHERE id = $2 and user_id = $3", post.PostContent, post.Id, post.UserId)
	return err
}

func (repo *PostgresRepository) DeletePost(ctx context.Context, id string, userId string) error {
	_, err := repo.db.ExecContext(ctx, "DELETE FROM posts WHERE id = $1 and user_id = $2", id, userId)
	return err
}

func NewPostgresRepository(url string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	return &PostgresRepository{db}, nil
}
