package transactionworkshopcontroller

import (
	"after-sales/api/payloads"
	"after-sales/api/payloads/pagination"
	transactionworkshopservice "after-sales/api/services/transaction/workshop"
	"after-sales/api/utils"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type VehicleHistoryController interface {
	GetAllFieldVehicleHistory(writer http.ResponseWriter, request *http.Request)
	GetVehicleHistoryById(writer http.ResponseWriter, request *http.Request)
}
type VehicleHistoryControllerImpl struct {
	VehicleHIstoryService transactionworkshopservice.VehicleHistoryService
}

func NewVehicleHistoryController(VehicleHIstoryService transactionworkshopservice.VehicleHistoryService) VehicleHistoryController {
	return &VehicleHistoryControllerImpl{VehicleHIstoryService: VehicleHIstoryService}
}

// GetAllFieldVehicleHistory gets all work orders
//
//	@Summary		Get All Vehicle History
//	@Description	REST API Vehicle History
//	@Accept			json
//	@Produce		json
//	@Tags			Transaction : Vehicle History
//	@Param			page					query		string	true	"page"
//	@Param			limit					query		string	true	"limit"
//	@Param			vehicle_model			query		string	false	"vehicle_model"
//	@Param			chassis_no				query		string	false	"chassis_no"
//	@Param			brand_id				query		string	false	"brand_id"
//	@Param			sort_by					query		string	false	"sort_by"
//	@Param			sort_of					query		string	false	"sort_of"
//	@Success		200						{object}	payloads.Response
//	@Failure		500,400,401,404,403,422	{object}	exceptions.BaseErrorResponse
//	@Router			/v1/vehicle-history/ [get]
func (r *VehicleHistoryControllerImpl) GetAllFieldVehicleHistory(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	queryParams := map[string]string{
		"vehicle_model": queryValues.Get("vehicle_model"),
		"chassis_no":    queryValues.Get("chassis_no"),
		"brand_id":      queryValues.Get("brand_id"),
		//"approval_status_id":           queryValues.Get("approval_status_id"),
	}

	paginations := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}
	filterCondition := utils.BuildFilterCondition(queryParams)
	result, err := r.VehicleHIstoryService.GetAllVehicleHistory(filterCondition, paginations)
	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data Successfully!", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)
	fmt.Println(result, err)
}

// GetVehicleHistoryById gets all work orders
//
//	@Summary		Get Vehicle History Id
//	@Description	REST API Vehicle History
//	@Accept			json
//	@Produce		json
//	@Tags			Transaction : Vehicle History
//	@Param			work_order_system_number_id	path		int true	"work_order_system_number_id"
//	@Success		200							{object}	transactionworkshoppayloads.VehicleHistoryByIdResponses
//	@Failure		500,400,401,404,403,422		{object}	exceptions.BaseErrorResponse
//	@Router			/v1/vehicle-history/by-id/{work_order_system_number_id} [get]
func (r *VehicleHistoryControllerImpl) GetVehicleHistoryById(writer http.ResponseWriter, request *http.Request) {
	WorkOrderSystemNumber, _ := strconv.Atoi(chi.URLParam(request, "work_order_system_number_id"))
	result, err := r.VehicleHIstoryService.GetVehicleHistoryById(WorkOrderSystemNumber)
	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
	fmt.Println(result, err)
}
