package transactionworkshopserviceimpl

import (
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshoprepository "after-sales/api/repositories/transaction/workshop"
	transactionworkshopservice "after-sales/api/services/transaction/workshop"

	"gorm.io/gorm"
)

type BookingEstimationServiceImpl struct {
	structBookingEstimationRepo transactionworkshoprepository.BookingEstimationRepository
}

func OpenBookingEstimationServiceImpl(bookingEstimationRepo transactionworkshoprepository.BookingEstimationRepository) transactionworkshopservice.BookingEstimationService {
	return &BookingEstimationServiceImpl{
		structBookingEstimationRepo: bookingEstimationRepo,
	}
}

func (s *BookingEstimationServiceImpl) WithTrx(Trxhandle *gorm.DB) transactionworkshopservice.BookingEstimationService {
	s.structBookingEstimationRepo = s.structBookingEstimationRepo.WithTrx(Trxhandle)
	return s
}

func (s *BookingEstimationServiceImpl) Save(request transactionworkshoppayloads.SaveBookingEstimationRequest) (bool, error) {
	save, err := s.structBookingEstimationRepo.Save(request)

	if err != nil {
		return false, err
	}

	return save, nil

}
