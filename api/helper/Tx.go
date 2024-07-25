package helper

import (
	"after-sales/api/exceptions"
	"net/http"

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
	defer func() {
		if r := recover(); r != nil {
			//log.Printf("Recovered in CommitOrRollbackTrx: %v", r)
			_ = tx.Rollback()
		}
	}()

	if err := tx.Commit().Error; err != nil {
		//log.Printf("Failed to commit transaction: %v", err)
		if rollbackErr := tx.Rollback().Error; rollbackErr != nil {
			//log.Printf("Failed to rollback after commit failure: %v", rollbackErr)
			return &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to rollback transaction after commit failure.",
				Err:        rollbackErr,
			}
		}
		return &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to commit transaction and rolled back successfully.",
			Err:        err,
		}
	}

	//log.Println("Transaction committed successfully.")
	return nil
}
