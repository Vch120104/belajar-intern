package helper

import (
	"after-sales/api/exceptions"
	"log"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// CommitOrRollback commits the transaction if no error or panic occurs, otherwise rolls back.
// It also recovers from panics and logs the error.
func CommitOrRollback(tx *gorm.DB, err *exceptions.BaseErrorResponse) {
	if r := recover(); r != nil {
		tx.Rollback()
		logrus.Info("Recovered from panic:", r)
		return
	}

	if err != nil {
		tx.Rollback()
		logrus.Info(err)
	} else {
		tx.Commit()
	}
}

func CommitOrRollbackTrx(tx *gorm.DB) *exceptions.BaseErrorResponse {
	const maxRetries = 3
	var rollbackErr error

	for i := 0; i < maxRetries; i++ {
		if i > 0 {
			time.Sleep(time.Millisecond * 100) // set delay 100ms between retries
		}

		if err := tx.Commit().Error; err != nil {
			log.Printf("Failed to commit transaction: %v", err)
			rollbackErr = tx.Rollback().Error
			if rollbackErr != nil {
				log.Printf("Failed to rollback after commit failure: %v", rollbackErr)
			} else {
				log.Println("Transaction rolled back after commit failure.")
			}
		} else {
			log.Println("Transaction committed successfully.")
			return nil
		}
	}

	// Handle context cancellation
	select {
	case <-tx.Statement.Context.Done():
		if rollbackErr != nil {
			return &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to rollback transaction after context cancellation.",
				Err:        rollbackErr,
			}
		} else {
			return &exceptions.BaseErrorResponse{
				StatusCode: http.StatusRequestTimeout,
				Message:    "Transaction rolled back due to context cancellation.",
			}
		}
	default:
		return &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to process transaction after multiple attempts.",
		}
	}
}
