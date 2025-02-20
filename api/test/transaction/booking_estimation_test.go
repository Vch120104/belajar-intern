package test

import (
	"after-sales/api/config"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshoprepositoryimpl "after-sales/api/repositories/transaction/workshop/repositories-workshop-impl"
	transactionworkshopserviceimpl "after-sales/api/services/transaction/workshop/services-workshop-impl"
	"fmt"
	"testing"
)

func TestSaveBookingEstimation(t *testing.T) {
	// Initialize environment configurations
	config.InitEnvConfigs(true, "")

	// Initialize the repository and service
	bookingEstimationRepo := transactionworkshoprepositoryimpl.OpenBookingEstimationRepositoryImpl()
	bookingEstimationService := transactionworkshopserviceimpl.OpenBookingEstimationServiceImpl(bookingEstimationRepo, nil, nil)

	// Create a request object for booking estimation
	request := transactionworkshoppayloads.BookingEstimationSaveRequest{

		CampaignId:             1,
		CompanyId:              1,
		DealerRepresentativeId: 1,
	}

	// Call the Save method and capture the return values
	result, err := bookingEstimationService.Save(request, 1)

	// Check if there was an error
	if err != nil {
		t.Errorf("Error saving booking estimation: %v", err)
		return
	}

	// Print the result for debugging purposes
	fmt.Println("Booking estimation saved successfully:", result)
}
