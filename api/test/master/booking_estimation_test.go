package test

// import (
// 	"after-sales/api/config"
// 	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
// 	transactionworkshoprepositoryimpl "after-sales/api/repositories/transaction/workshop/repositories-workshop-impl"
// 	transactionworkshopserviceimpl "after-sales/api/services/transaction/workshop/services-workshop-impl"
// 	"fmt"
// 	"testing"
// 	"time"
// )

// func TestSaveBookingEstimation(t *testing.T) {
// 	config.InitEnvConfigs(true, "")
// 	bookingEstimationRepo := transactionworkshoprepositoryimpl.OpenBookingEstimationRepositoryImpl()
// 	bookingEstimationService := transactionworkshopserviceimpl.OpenBookingEstimationServiceImpl(bookingEstimationRepo, nil, nil)

// 	// objek *gorm.DB dari manajer koneksi database
// 	db := config.InitDB()

// 	save, err := bookingEstimationService.Save(db, transactionworkshoppayloads.BookingEstimationRequest{
// 		BatchSystemNumber:              1,
// 		BookingSystemNumber:            1,
// 		BrandId:                        1,
// 		ModelId:                        1,
// 		VariantId:                      1,
// 		VehicleId:                      1,
// 		EstimationSystemNumber:         1,
// 		PdiSystemNumber:                1,
// 		ServiceRequestSystemNumber:     1,
// 		ContractSystemNumber:           1,
// 		AgreementId:                    1,
// 		CampaignId:                     1,
// 		CompanyId:                      1,
// 		ProfitCenterId:                 1,
// 		DealerRepresentativeId:         1,
// 		CustomerId:                     1,
// 		DocumentStatusId:               1,
// 		BookingEstimationBatchDate:     time.Now(),
// 		BookingEstimationVehicleNumber: "TEST1",
// 		AgreementNumberBr:              "TEST1",
// 		IsUnregistered:                 "T",
// 		ContactPersonName:              "TEST1",
// 		ContactPersonPhone:             "TEST1",
// 		ContactPersonMobile:            "TEST1",
// 		ContactPersonVia:               "TEST1",
// 		InsurancePolicyNo:              "TEST1",
// 		InsuranceExpiredDate:           time.Now(),
// 		InsuranceClaimNo:               "TEST1",
// 		InsurancePic:                   "TEST",
// 	})

// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println(save)
// }
