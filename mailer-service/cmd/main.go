package main

import (
	"fmt"
	"github.com/magomedcoder/simple-microservice/mailer-service/internal"
	"log"
	"net/http"
	"os"
)

func main() {
	app := internal.Config{
		Mailer: internal.CreateMail(),
	}
	port := os.Getenv("PORT")
	log.Printf("Запуск службы почтовый сервиса на порту %s\n", port)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.Routes(),
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}
