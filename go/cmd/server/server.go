package main

import (
	"GoNews/pkg/api"
	"GoNews/pkg/storage"
	"GoNews/pkg/storage/memdb"
	"GoNews/pkg/storage/mongodb"
	"GoNews/pkg/storage/postgresql"
	"log"
	"net/http"
)

// Сервер GoNews.
type server struct {
	db  storage.Interface
	api *api.API
}

func main() {
	// Создаём объект сервера.
	var srv server

	// Создаём объекты баз данных.
	//
	// БД в памяти.
	db := memdb.New()

	//Реляционная БД PostgreSQL.
	dbInfoPSQL := "host=localhost user=postgres password=admin dbname=postgres port=5432 sslmode=disable"
	db2, err := postgresql.Connect(dbInfoPSQL)
	if err != nil {
		log.Fatal(err)
	}
	// Документная БД MongoDB.
	uri := "mongodb://admin:secret@localhost:27017/"
	dbName := "posts" // Замените на имя вашей базы данных

	db3, err := mongodb.Connect(uri, dbName)
	if err != nil {
		log.Fatal(err)
	}
	_, _ = db2, db3

	// Инициализируем хранилище сервера конкретной БД.
	srv.db = db

	// Создаём объект API и регистрируем обработчики.
	srv.api = api.New(srv.db)

	// Запускаем веб-сервер на порту 8080 на всех интерфейсах.
	// Предаём серверу маршрутизатор запросов,
	// поэтому сервер будет все запросы отправлять на маршрутизатор.
	// Маршрутизатор будет выбирать нужный обработчик.
	http.ListenAndServe(":8080", srv.api.Router())
}
