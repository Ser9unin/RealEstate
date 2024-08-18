package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type SrvCfg struct {
	Host              string
	Port              string
	ReadHeaderTimeout time.Duration // Настраиваем тайм-аут ожидания заголовков
	ReadTimeout       time.Duration // Настраиваем общий тайм-аут запроса
	WriteTimeout      time.Duration // Настраиваем тайм-аут записи ответа
	IdleTimeout       time.Duration // Настраиваем тайм-аут простоя соединения
}

func NewSeverCfg() SrvCfg {
	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8000"
	}

	readHeader, err := strconv.Atoi(os.Getenv("ReadHeaderTimeout"))
	if err != nil || readHeader == 0 {
		log.Println(err)
		readHeader = 1
	}

	readTimeout, err := strconv.Atoi(os.Getenv("ReadTimeout"))
	if err != nil || readTimeout == 0 {
		log.Println(err)
		readTimeout = 2
	}

	writeTimeout, err := strconv.Atoi(os.Getenv("WriteTimeout"))
	if err != nil || writeTimeout == 0 {
		log.Println(err)
		writeTimeout = 2
	}

	idleTimeout, err := strconv.Atoi(os.Getenv("IdleTimeout"))
	if err != nil || idleTimeout == 0 {
		log.Println(err)
		idleTimeout = 3
	}

	server := SrvCfg{
		Host:              host,
		Port:              port,
		ReadHeaderTimeout: time.Duration(readHeader) * time.Second,
		ReadTimeout:       time.Duration(readTimeout) * time.Second,
		WriteTimeout:      time.Duration(writeTimeout) * time.Second,
		IdleTimeout:       time.Duration(idleTimeout) * time.Second,
	}

	return server
}

func NewDBCfg() string {
	host := os.Getenv("POSTGRES_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("POSTGRES_PORT")
	if port == "" {
		port = "5433"
	}
	user := os.Getenv("POSTGRES_USER")
	if user == "" {
		user = "dev"
	}
	pass := os.Getenv("POSTGRES_PASSWORD")
	if pass == "" {
		pass = "somepass"
	}
	dbName := os.Getenv("POSTGRES_DB")
	if dbName == "" {
		dbName = "apartmentsdb"
	}
	sslMode := os.Getenv("POSTGRES_SSL")
	if sslMode == "" {
		sslMode = "disable"
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode= %s",
		host,
		port,
		user,
		pass,
		dbName,
		sslMode)

	return dsn
}
