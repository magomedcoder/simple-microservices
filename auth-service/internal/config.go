package internal

import (
	"database/sql"
	"log"
	"os"
	"time"
)

type Config struct {
	DB     *sql.DB
	Models Models
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

var count int64

func ConnectToDB() *sql.DB {
	dsn := os.Getenv("DSN_POSTGRES")
	for {
		conn, err := openDB(dsn)
		if err != nil {
			log.Println("PostgreSQL еще не готов. Повторная попытка... ")
			count++
		} else {
			log.Println("Подключено к PostgreSQL")
			return conn
		}
		
		if count > 20 {
			log.Println(err)
			return nil
		}

		log.Println("Откладываем на 2 секунды...")
		time.Sleep(2 * time.Second)
		continue
	}
}
