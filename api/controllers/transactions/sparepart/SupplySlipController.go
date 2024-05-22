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

// GetSupplySlipByID retrieves a supply slip by ID
// @Summary Get Supply Slip By ID
// @Description Retrieve a supply slip by its ID
// @Accept json
// @Produce json
// @Tags Transaction : Spare Part Supply Slip
// @Param supply_slip_id path int true "Supply Slip ID"
// @Success 200 {object} payloads.Response
// @Failure 500,404 {object} exceptions.BaseErrorResponse
// @Router /v1/supply-slip/{supply_slip_id} [get]
func (r *SupplySlipControllerImpl) GetSupplySlipByID(writer http.ResponseWriter, request *http.Request) {
	// Get ID from URL
	// id := mux.Vars(request)["id"]

	// Get data from service
	// data, err := r.supplyslipservice.GetSupplySlipByID(id)
	// if err != nil {
	// 	// Return error
	// 	exceptions.NewNotFoundException(writer, request, err)
	// 	return
	// }

	// Return success
	// payloads.NewHandleSuccess(writer, data, "Get Data Successfully", http.StatusOK)
}
