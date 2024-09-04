package postgresql

import (
	"GoNews/pkg/storage"
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
	// Импортируйте пакет с вашим интерфейсом и структурой Post
)

// Хранилище данных.
type Storage struct {
	db *pgxpool.Pool
}

// Connect устанавливает соединение с базой данных и возвращает объект Storage.
func Connect(constr string) (*Storage, error) {
	db, err := pgxpool.Connect(context.Background(), constr)
	if err != nil {
		return nil, err
	}
	log.Println("Подключение к базе данных postgres успешно установлено")
	return &Storage{db: db}, nil
}

// Posts возвращает список всех публикаций из БД.
func (s *Storage) Posts() ([]storage.Post, error) {
	rows, err := s.db.Query(context.Background(), `
		SELECT p.id, p.title, p.content, p.author_id, a.name, p.created_at, p.published_at
		FROM posts p
		JOIN authors a ON p.author_id = a.id
		ORDER BY p.created_at DESC;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []storage.Post
	for rows.Next() {
		var post storage.Post
		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.AuthorID,
			&post.AuthorName,
			&post.CreatedAt,
			&post.PublishedAt,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return posts, nil
}

// AddPost добавляет новую публикацию в БД и возвращает её ID.
func (s *Storage) AddPost(post storage.Post) error {
	_, err := s.db.Query(context.Background(), `
		INSERT INTO posts (title, content, author_id, created_at, published_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id;
	`,
		post.Title,
		post.Content,
		post.AuthorID,
		post.CreatedAt,
		post.PublishedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

// UpdatePost обновляет данные публикации в БД по её ID.
func (s *Storage) UpdatePost(post storage.Post) error {
	commandTag, err := s.db.Exec(context.Background(), `
		UPDATE posts
		SET title = $1, content = $2, author_id = $3, published_at = $4
		WHERE id = $5;
	`,
		post.Title,
		post.Content,
		post.AuthorID,
		post.PublishedAt,
		post.ID,
	)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return errors.New("no rows were updated")
	}
	return nil
}

// DeletePost удаляет публикацию из БД по её ID.
func (s *Storage) DeletePost(post storage.Post) error {
	commandTag, err := s.db.Exec(context.Background(), `
		DELETE FROM posts WHERE id = $1;
	`, post.ID)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return errors.New("no rows were deleted")
	}
	return nil
}
