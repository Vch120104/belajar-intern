package transactionsparepartcontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/payloads"
	transactionsparepartservice "after-sales/api/services/transaction/sparepart"

	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type SupplySlipController struct {
	supplyslipservice transactionsparepartservice.SupplySlipService
}

func StartSupplySlipRoutes(
	db *gorm.DB,
	r chi.Router,
	supplyslipservice transactionsparepartservice.SupplySlipService,
) {
	supplySlipHandler := SupplySlipController{supplyslipservice: supplyslipservice}
	r.Get("/supply-slip/{supply_system_number}", supplySlipHandler.GetSupplySlipByID)
}

// Get Supply Slip By ID
func (r *SupplySlipController) GetSupplySlipByID(w http.ResponseWriter, req *http.Request) {
	supplySystemNumber, _ := strconv.Atoi(chi.URLParam(req, "supply_system_number"))
	result, err := r.supplyslipservice.GetSupplySlipById(int32(supplySystemNumber))
	if err != nil {
		exceptions.NotFoundException(w, err.Error())
		return
	}
	payloads.NewHandleSuccess(w, result, "Get Data Successfully!", http.StatusOK)
}
