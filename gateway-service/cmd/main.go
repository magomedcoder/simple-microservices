package main

import (
	"fmt"
	"github.com/magomedcoder/simple-microservice/gateway-service/internal"
	"log"
	"net/http"
	"os"
)

func main() {
	rabbitConn, err := internal.Connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()
	app := internal.Config{
		Rabbit: rabbitConn,
	}
	port := os.Getenv("PORT")
	log.Printf("Запуск шлюза на порту %s\n", port)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.Routes(),
	}
	if err = srv.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}
