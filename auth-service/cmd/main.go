package main

import (
	"fmt"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/magomedcoder/simple-microservice/auth-service/internal"
	"log"
	"net/http"
	"os"
)

func main() {
	conn := internal.ConnectToDB()
	if conn == nil {
		log.Panic("Невозможно подключиться к PostgreSQL")
	}
	app := internal.Config{
		DB:     conn,
		Models: internal.New(conn),
	}
	port := os.Getenv("PORT")
	log.Printf("Запуск службы аутентификации на порту %s\n", port)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.Routes(),
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}
