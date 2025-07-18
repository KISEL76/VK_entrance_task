package main

import (
	"log"
	"marketplace/app/db"
	"marketplace/app/handlers"
	"marketplace/app/utils"
	"net/http"
)

func main() {
	database, err := db.NewPostgres()
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: : %v", err)
	}
	defer database.Close()

	Handler := &handlers.Handler{DB: database}

	http.HandleFunc("/register", utils.WithTokenIfPresent(Handler.Register))
	http.HandleFunc("/login", utils.WithTokenIfPresent(Handler.Login))
	http.HandleFunc("/create/add", utils.WithAuth(Handler.CreateAd))
	http.HandleFunc("/ads", utils.WithTokenIfPresent(Handler.GoodsList))

	log.Println("Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
