package test

import (
	"after-sales/api/config"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshoprepositoryimpl "after-sales/api/repositories/transaction/workshop/repositories-workshop-impl"
	transactionworkshopserviceimpl "after-sales/api/services/transaction/workshop/services-workshop-impl"
	"fmt"
	"testing"
	"time"
)

func TestSaveBookingEstimation(t *testing.T) {
	// Initialize environment configurations
	config.InitEnvConfigs(true, "")

	// Initialize the repository and service
	bookingEstimationRepo := transactionworkshoprepositoryimpl.OpenBookingEstimationRepositoryImpl()
	bookingEstimationService := transactionworkshopserviceimpl.OpenBookingEstimationServiceImpl(bookingEstimationRepo, nil, nil)

	// Initialize the database connection
	db := config.InitDB()

	// Create a request object for booking estimation
	request := transactionworkshoppayloads.BookingEstimationRequest{
		BatchSystemNumber:              1,
		BookingSystemNumber:            1,
		BrandId:                        1,
		ModelId:                        1,
		VariantId:                      1,
		VehicleId:                      1,
		EstimationSystemNumber:         1,
		PdiSystemNumber:                1,
		ServiceRequestSystemNumber:     1,
		ContractSystemNumber:           1,
		AgreementId:                    1,
		CampaignId:                     1,
		CompanyId:                      1,
		ProfitCenterId:                 1,
		DealerRepresentativeId:         1,
		CustomerId:                     1,
		DocumentStatusId:               1,
		BookingEstimationBatchDate:     time.Now(),
		BookingEstimationVehicleNumber: "TEST1",
		AgreementNumberBr:              "TEST1",
		IsUnregistered:                 true,
		ContactPersonName:              "TEST1",
		ContactPersonPhone:             "TEST1",
		ContactPersonMobile:            "TEST1",
		ContactPersonViaId:             1,
		InsurancePolicyNo:              "TEST1",
		InsuranceExpiredDate:           time.Now(),
		InsuranceClaimNo:               "TEST1",
		InsurancePic:                   "TEST",
	}

	// Call the Save method and capture the return values
	result, err := bookingEstimationService.Save(db, request)

	// Check if there was an error
	if err != nil {
		t.Errorf("Error saving booking estimation: %v", err)
		return
	}

	// Print the result for debugging purposes
	fmt.Println("Booking estimation saved successfully:", result)
}
