package masterwarehousecontroller

import (
	"after-sales/api/helper"
	"after-sales/api/payloads"
	masterwarehouseservice "after-sales/api/services/master/warehouse"
	"net/http"

	"github.com/go-chi/chi/v5"
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

// @Summary Get Warehouse Costing Type By Code
// @Description Get Warehouse Costing Type By Code
// @Tags Master : Costing Type
// @Security AuthorizationKeyAuth
// @Accept json
// @Produce json
// @Param warehouse-costing-type-code path string true "Warehouse Costing Type Code"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/warehouse-costing-type/by-code/{warehouse-costing-type-code} [get]
func (w *WarehouseCostingTypeControllerImpl) GetWarehouseCostingTypeByCode(writer http.ResponseWriter, request *http.Request) {
	CostingTypeCode := chi.URLParam(request, "warehouse-costing-type-code")

	get, err := w.CostingService.GetByCodeWarehouseCostingType(CostingTypeCode)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, get, "Get Data Successfully!", http.StatusOK)

}
