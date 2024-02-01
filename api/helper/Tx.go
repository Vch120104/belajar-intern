package helper

import (
	"gorm.io/gorm"
)

func CommitOrRollback(tx *gorm.DB) {
	err := recover()
	if err != nil {
		errorRollback := tx.Rollback()
		PanicIfError(errorRollback.Error)
		panic(err)
	} else {
		errorCommit := tx.Commit()
		PanicIfError(errorCommit.Error)
	}
}
