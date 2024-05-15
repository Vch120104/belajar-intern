package transactionsparepartcontroller

import (
	transactionsparepartservice "after-sales/api/services/transaction/sparepart"
	"net/http"
)

type SupplySlipControllerImpl struct {
	supplyslipservice transactionsparepartservice.SupplySlipService
}

type SupplySlipController interface {
	GetSupplySlipByID(writer http.ResponseWriter, request *http.Request)
}

func NewSupplySlipController(supplyslipservice transactionsparepartservice.SupplySlipService) SupplySlipController {
	return &SupplySlipControllerImpl{
		supplyslipservice: supplyslipservice,
	}
}

// Get Supply Slip By ID
func (r *SupplySlipControllerImpl) GetSupplySlipByID(writer http.ResponseWriter, request *http.Request) {
	// Get ID from URL
	// id := mux.Vars(request)["id"]

	// Get data from service
	// data, err := r.supplyslipservice.GetSupplySlipByID(id)
	// if err != nil {
	// 	// Return error
	// 	exceptionsss_test.NewNotFoundException(writer, request, err)
	// 	return
	// }

	// Return success
	// payloads.NewHandleSuccess(writer, data, "Get Data Successfully", http.StatusOK)
}
