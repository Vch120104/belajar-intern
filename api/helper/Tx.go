package helper

import (
	"after-sales/api/exceptions"

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
