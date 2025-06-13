package transfer

import (
	"context"
	"database/sql"
	"log"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/pradeepitm12/compaaa/internal-http/internal/pkg/logger"

	"github.com/pradeepitm12/compaaa/internal-http/internal/db"
	"github.com/pradeepitm12/compaaa/internal-http/internal/domain/model"
	"github.com/shopspring/decimal"

	_ "github.com/lib/pq"
)

func setupTestDB(t *testing.T) *sql.DB {
	//dsn := "postgresql://myuser:admin123@localhost:5432/mydatabase?sslmode=disable"
	dsn := "postgresql://postgres:secret@localhost:5432/transferdb?sslmode=disable"
	dbConn, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("failed to connect to DB: %v", err)
	}
	dbConn.SetMaxOpenConns(100)
	dbConn.SetMaxIdleConns(50)
	return dbConn
}

func TestConcurrentTransfers(t *testing.T) {
	dbConn := setupTestDB(t)
	defer dbConn.Close()

	_, _ = dbConn.Exec(`DELETE FROM transactions`)
	_, _ = dbConn.Exec(`DELETE FROM accounts`)

	ctx := context.Background()

	// Repositories and service setup
	accountRepo := db.NewAccountRepository(dbConn)
	txRepo := db.NewTransactionRepository(dbConn)
	txManager := db.NewTxManager(dbConn)
	logger.Init()

	svc := NewService(accountRepo, txRepo, txManager, logger.L())

	source := model.NewAccount(1, decimal.NewFromInt(1000))
	dest := model.NewAccount(2, decimal.NewFromInt(0))
	txManager.Do(ctx, func(ctx context.Context, tx *sql.Tx) error {
		if err := accountRepo.Create(ctx, tx, source); err != nil {
			t.Fatal(err)
		}
		if err := accountRepo.Create(ctx, tx, dest); err != nil {
			t.Fatal(err)
		}
		return nil
	})
	var wg sync.WaitGroup
	var successCount int64
	n := 50
	amount := decimal.NewFromInt(10)

	semaphore := make(chan struct{}, 10) // limit to 10 concurrent DB connections

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // increased timeout
			defer cancel()

			_, err := svc.Transfer(ctx, 1, 2, amount)
			if err == nil {
				atomic.AddInt64(&successCount, 1)
			} else {
				log.Printf("transfer failed: %v", err)
			}
		}()
	}
	wg.Wait()

	err := txManager.Do(ctx, func(ctx context.Context, tx *sql.Tx) error {
		srcAcc, err := accountRepo.GetByID(ctx, tx, 1)
		if err != nil {
			return err
		}
		destAcc, err := accountRepo.GetByID(ctx, tx, 2)
		if err != nil {
			return err
		}

		expectedSrc := decimal.NewFromInt(1000).Sub(amount.Mul(decimal.NewFromInt(successCount)))
		expectedDest := amount.Mul(decimal.NewFromInt(successCount))

		if !srcAcc.Balance().Equal(expectedSrc) {
			t.Errorf("unexpected source balance: got %s, want %s", srcAcc.Balance(), expectedSrc)
		}
		if !destAcc.Balance().Equal(expectedDest) {
			t.Errorf("unexpected destination balance: got %s, want %s", destAcc.Balance(), expectedDest)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("final balance check failed: %v", err)
	}
}
