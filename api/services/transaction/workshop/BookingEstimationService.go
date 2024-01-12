package transactionworkshopservice

import (
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"

	"gorm.io/gorm"
)

type BookingEstimationService interface {
	WithTrx(Trxhandle *gorm.DB) BookingEstimationService
	Save(transactionworkshoppayloads.SaveBookingEstimationRequest) (bool, error)
}
