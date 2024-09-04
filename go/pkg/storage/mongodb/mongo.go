package mongodb

import (
	"context"
	"errors"
	"log"

	"GoNews/pkg/storage"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Storage представляет хранилище данных в MongoDB.
type Storage struct {
	client   *mongo.Client
	database *mongo.Database
	posts    *mongo.Collection
}

// Connect устанавливает соединение с базой данных и возвращает объект Storage.
// Connect устанавливает соединение с базой данных и возвращает объект Storage.
func Connect(uri, dbName string) (*Storage, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	// Проверка соединения
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	database := client.Database(dbName)
	log.Println("Подключение к базе данных mongo успешно установлено")
	return &Storage{
		client:   client,
		database: database,
		posts:    database.Collection("posts"),
	}, nil
}

// Posts возвращает список всех публикаций из БД.
func (s *Storage) Posts() ([]storage.Post, error) {
	var posts []storage.Post

	cursor, err := s.posts.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var post storage.Post
		err := cursor.Decode(&post)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

// AddPost добавляет новую публикацию в БД.
func (s *Storage) AddPost(post storage.Post) error {
	_, err := s.posts.InsertOne(context.Background(), post)
	return err
}

// UpdatePost обновляет данные публикации в БД по её ID.
func (s *Storage) UpdatePost(post storage.Post) error {
	filter := bson.D{{Key: "id", Value: post.ID}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "title", Value: post.Title},
			{Key: "content", Value: post.Content},
			{Key: "author_id", Value: post.AuthorID},
			{Key: "published_at", Value: post.PublishedAt},
		}},
	}

	result, err := s.posts.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("нет такого документа для удаления")
	}
	return nil
}

// DeletePost удаляет публикацию из БД по её ID.
func (s *Storage) DeletePost(post storage.Post) error {
	filter := bson.D{{Key: "id", Value: post.ID}}

	result, err := s.posts.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("нет такого документа для удаления")
	}
	return nil
}
