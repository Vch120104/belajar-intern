package transactionworkshoprepositoryimpl

import (
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshoprepository "after-sales/api/repositories/transaction/workshop"
	"log"

	"gorm.io/gorm"
)

type BookingEstimationImpl struct {
	DB *gorm.DB
}

func OpenBookingEstimationImpl(db *gorm.DB) transactionworkshoprepository.BookingEstimationRepository {
	return &BookingEstimationImpl{DB: db}
}

func (r *BookingEstimationImpl) WithTrx(Trxhandle *gorm.DB) transactionworkshoprepository.BookingEstimationRepository {
	if Trxhandle != nil {
		log.Println("Transaction Database Not Found")
		return r
	}
	r.DB = Trxhandle
	return r
}

func (r *BookingEstimationImpl) Save(request transactionworkshoppayloads.SaveBookingEstimationRequest) (bool, error) {

	var bookingEstimationEntities = transactionworkshopentities.BookingEstimation{
		BatchSystemNumber:              request.BatchSystemNumber,
		BookingSystemNumber:            request.BookingSystemNumber,
		BrandId:                        request.BrandId,
		ModelId:                        request.ModelId,
		VariantId:                      request.VariantId,
		VehicleId:                      request.VehicleId,
		EstimationSystemNumber:         request.EstimationSystemNumber,
		PdiSystemNumber:                request.PdiSystemNumber,
		ServiceRequestSystemNumber:     request.ServiceRequestSystemNumber,
		ContractSystemNumber:           request.ContractSystemNumber,
		AgreementId:                    request.AgreementId,
		CampaignId:                     request.CampaignId,
		CompanyId:                      request.CompanyId,
		ProfitCenterId:                 request.ProfitCenterId,
		DealerRepresentativeId:         request.DealerRepresentativeId,
		CustomerId:                     request.CustomerId,
		DocumentStatusId:               request.DocumentStatusId,
		BookingEstimationBatchDate:     request.BookingEstimationBatchDate,
		BookingEstimationVehicleNumber: request.BookingEstimationVehicleNumber,
		AgreementNumberBr:              request.AgreementNumberBr,
		ContactPersonName:              request.ContactPersonName,
		ContactPersonPhone:             request.ContactPersonPhone,
		ContactPersonMobile:            request.ContactPersonMobile,
		ContactPersonVia:               request.ContactPersonVia,
		InsurancePolicyNo:              request.InsurancePolicyNo,
		InsuranceExpiredDate:           request.InsuranceExpiredDate,
		InsuranceClaimNo:               request.InsuranceClaimNo,
		InsurancePic:                   request.InsurancePic,
	}

	rows, err := r.DB.Model(&bookingEstimationEntities).
		Save(&bookingEstimationEntities).
		Rows()

	if err != nil {
		return false, err
	}

	defer rows.Close()

	return true, nil
}
