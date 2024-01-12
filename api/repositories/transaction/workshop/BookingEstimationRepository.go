package transactionworkshoprepository

import (
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"

	"gorm.io/gorm"
)

type BookingEstimationRepository interface {
	WithTrx(Trxhandle *gorm.DB) BookingEstimationRepository
	Save(transactionworkshoppayloads.SaveBookingEstimationRequest) (bool, error)
}
