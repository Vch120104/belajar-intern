package masterwarehousecontroller

import (
	"after-sales/api/helper"
	"after-sales/api/payloads"
	masterwarehouseservice "after-sales/api/services/master/warehouse"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type WarehouseCostingTypeController interface {
	GetWarehouseCostingTypeByCode(writer http.ResponseWriter, request *http.Request)
}
type WarehouseCostingTypeControllerImpl struct {
	CostingService masterwarehouseservice.WarehouseCostingTypeService
}

func NewWarehouseCostingTypeController(CostingService masterwarehouseservice.WarehouseCostingTypeService) WarehouseCostingTypeController {
	return &WarehouseCostingTypeControllerImpl{CostingService: CostingService}
}
func (w *WarehouseCostingTypeControllerImpl) GetWarehouseCostingTypeByCode(writer http.ResponseWriter, request *http.Request) {
	CostingTypeCode := chi.URLParam(request, "warehouse-costing-type-code")

	get, err := w.CostingService.GetByCodeWarehouseCostingType(CostingTypeCode)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, get, "Get Data Successfully!", http.StatusOK)

}
