package main

import (
	"database/sql"
	"fmt"
	"github.com/pradeepitm12/compaaa/internal-http/internal/pkg/logger"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	internalhttp "github.com/pradeepitm12/compaaa/internal-http/internal/api/http"
	"github.com/pradeepitm12/compaaa/internal-http/internal/db"
	"github.com/pradeepitm12/compaaa/internal-http/internal/transfer"
)

func main() {
	logger.Init()
	defer logger.L().Sync()

	log := logger.L()
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	pgsql, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to connect to DB: %v", err))
	}
	defer pgsql.Close()
	pgsql.SetMaxOpenConns(10)
	pgsql.SetMaxIdleConns(5)
	pgsql.SetConnMaxLifetime(time.Minute * 5)
	log.Info("Connected to DB")

	accountRepo := db.NewAccountRepository(pgsql)
	transactionRepo := db.NewTransactionRepository(pgsql)
	txManager := db.NewTxManager(pgsql)

	transferService := transfer.NewService(accountRepo, transactionRepo, txManager, log)

	handler := internalhttp.NewHandler(transferService, accountRepo)
	router := internalhttp.NewRouter(handler)

	addr := ":8080"
	log.Info(fmt.Sprintf("starting server at %s", addr))
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal(fmt.Sprintf("server failed: %v", err))
	}
}
