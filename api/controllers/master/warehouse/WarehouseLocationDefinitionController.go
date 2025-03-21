package masterwarehousecontroller

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads"
	"after-sales/api/utils"
	"after-sales/api/validation"
	"encoding/json"
	"errors"
	"fmt"

	"strconv"

	// masteritemlevelentities "after-sales/api/entities/master/item_level"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	pagination "after-sales/api/payloads/pagination"

	// masteritemlevelrepo "after-sales/api/repositories/master/item_level"
	masterwarehouseservice "after-sales/api/services/master/warehouse"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type WarehouseLocationDefinitionControllerImpl struct {
	WarehouseLocationDefinitionService masterwarehouseservice.WarehouseLocationDefinitionService
}

type WarehouseLocationDefinitionController interface {
	GetByLevel(writer http.ResponseWriter, request *http.Request)
	GetAll(writer http.ResponseWriter, request *http.Request)
	GetById(writer http.ResponseWriter, request *http.Request)
	Save(writer http.ResponseWriter, request *http.Request)
	SaveData(writer http.ResponseWriter, request *http.Request)
	ChangeStatus(writer http.ResponseWriter, request *http.Request)
	PopupWarehouseLocationLevel(writer http.ResponseWriter, request *http.Request)
}

func NewWarehouseLocationDefinitionController(WarehouseLocationDefinitionService masterwarehouseservice.WarehouseLocationDefinitionService) WarehouseLocationDefinitionController {
	return &WarehouseLocationDefinitionControllerImpl{
		WarehouseLocationDefinitionService: WarehouseLocationDefinitionService,
	}
}

// @Summary Get All Warehouse Location
// @Description Get All Warehouse Location
// @Accept json
// @Produce json
// @Tags Master : Warehouse Location Definition
// @Security AuthorizationKeyAuth
// @Success 200 {object} payloads.Response
// @Param page query string true "Page"
// @Param limit query string true "Limit"
// @Param is_active query bool false "Is Active"
// @Param warehouse_location_code query string false "Warehouse Location Code"
// @Param warehouse_location_name query string false "Warehouse Location Name"
// @Param company_id query string false "Company Id"
// @Param warehouse_location_detail_name query string false "Warehouse Location Detail Name"
// @Param sort_by query string false "Sort Of: {column}"
// @Param sort_of query string false "Sort By: {asc}"
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/warehouse-location-definition/ [get]
func (r *WarehouseLocationDefinitionControllerImpl) GetAll(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"mtr_warehouse_location_definition.warehouse_location_definition_level":       queryValues.Get("warehouse_location_definition_level"),
		"mtr_warehouse_location_definition.warehouse_location_definition_level_code":  queryValues.Get("warehouse_location_definition_level_code"),
		"mtr_warehouse_location_definition.warehouse_location_definition_description": queryValues.Get("warehouse_location_definition_description"),
		"mtr_warehouse_location_definition.is_active":                                 queryValues.Get("is_active"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	paginatedData, totalPages, totalRows, err := r.WarehouseLocationDefinitionService.GetAll(criteria, paginate)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(paginatedData), "Get Data Successfully", http.StatusOK, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
}

// @Summary Get Warehouse Location By Level Id
// @Description Get Warehouse Location By Level Id
// @Accept json
// @Produce json
// @Tags Master : Warehouse Location Definition
// @Security AuthorizationKeyAuth
// @Param warehouse_location_definition_level_id path int true "Warehouse Location Definition Level ID"
// @Param warehouse_location_definition_level_code path string true "Warehouse Location Definition ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/warehouse-location-definition/by-level/{warehouse_location_definition_level_id}/{warehouse_location_definition_level_code} [get]
func (r *WarehouseLocationDefinitionControllerImpl) GetByLevel(writer http.ResponseWriter, request *http.Request) {
	warehouseLocationDefinitionLevelID, errA := strconv.Atoi(chi.URLParam(request, "warehouse_location_definition_level_id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}
	warehouseLocationDefinitionID := chi.URLParam(request, "warehouse_location_definition_level_code") // Menggunakan nilai string langsung

	get, err := r.WarehouseLocationDefinitionService.GetByLevel(warehouseLocationDefinitionLevelID, warehouseLocationDefinitionID)

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, get, "Get Data Successfully!", http.StatusOK)
}

// @Summary Get Warehouse Location By Id
// @Description Get Warehouse Location By Id
// @Accept json
// @Produce json
// @Tags Master : Warehouse Location Definition
// @Security AuthorizationKeyAuth
// @Param warehouse_location_definition_id path int true "warehouse_location_definition_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/warehouse-location-definition/by-id/{warehouse_location_definition_id} [get]
func (r *WarehouseLocationDefinitionControllerImpl) GetById(writer http.ResponseWriter, request *http.Request) {

	WarehouseLocationDefinitionId, errA := strconv.Atoi(chi.URLParam(request, "warehouse_location_definition_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	get, err := r.WarehouseLocationDefinitionService.GetById(WarehouseLocationDefinitionId)

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, get, "Get Data Successfully!", http.StatusOK)

}

// @Summary Save Warehouse Location
// @Description Save Warehouse Location
// @Accept json
// @Produce json
// @Tags Master : Warehouse Location Definition
// @Security AuthorizationKeyAuth
// @param reqBody body masterwarehousepayloads.WarehouseLocationDefinitionResponse true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/warehouse-location-definition/ [post]
func (r *WarehouseLocationDefinitionControllerImpl) Save(writer http.ResponseWriter, request *http.Request) {
	var message string
	var formRequest masterwarehousepayloads.WarehouseLocationDefinitionResponse
	helper.ReadFromRequestBody(request, &formRequest)
	if validationErr := validation.ValidationForm(writer, request, &formRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}
	save, err := r.WarehouseLocationDefinitionService.Save(formRequest)

	if formRequest.WarehouseLocationDefinitionId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, save, message, http.StatusOK)

}

// @Summary Save Data Warehouse Location
// @Description Save Data Warehouse Location
// @Accept json
// @Produce json
// @Tags Master : Warehouse Location Definition
// @Security AuthorizationKeyAuth
// @Param warehouse_location_definition_id path int true "Warehouse Location ID"
// @param reqBody body masterwarehousepayloads.WarehouseLocationDefinitionResponse true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/warehouse-location-definition/{warehouse_location_id} [put]
func (r *WarehouseLocationDefinitionControllerImpl) SaveData(writer http.ResponseWriter, request *http.Request) {
	warehouseLocationID := chi.URLParam(request, "warehouse_location_definition_id")
	id, err := strconv.Atoi(warehouseLocationID)
	if err != nil {
		errResponse := &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        fmt.Errorf("invalid warehouse_location_id"),
		}
		exceptions.NewBadRequestException(writer, request, errResponse)
		return
	}

	var formRequest masterwarehousepayloads.WarehouseLocationDefinitionResponse
	if err := json.NewDecoder(request.Body).Decode(&formRequest); err != nil {
		errResponse := &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        fmt.Errorf("invalid request body"),
		}
		exceptions.NewBadRequestException(writer, request, errResponse)
		return
	}
	formRequest.WarehouseLocationDefinitionId = id

	save, saveErr := r.WarehouseLocationDefinitionService.SaveData(formRequest)
	if saveErr != nil {
		exceptions.NewNotFoundException(writer, request, saveErr)
		return
	}

	var message string
	if id == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, save, message, http.StatusOK)
}

// @Summary Change Warehouse Location Status By Id
// @Description Change Warehouse Location Status By Id
// @Accept json
// @Produce json
// @Tags Master : Warehouse Location Definition
// @Security AuthorizationKeyAuth
// @Param warehouse_location_definition_id path int true "Warehouse Location Id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/warehouse-location-definition/{warehouse_location_definition_id} [patch]
func (r *WarehouseLocationDefinitionControllerImpl) ChangeStatus(writer http.ResponseWriter, request *http.Request) {

	WarehouseLocationDefinitionId, errA := strconv.Atoi(chi.URLParam(request, "warehouse_location_definition_id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	entity, err := r.WarehouseLocationDefinitionService.ChangeStatus(WarehouseLocationDefinitionId)
	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	responseData := map[string]interface{}{
		"is_active":                        entity.IsActive,
		"warehouse_location_definition_id": entity.WarehouseLocationDefinitionId,
	}

	payloads.NewHandleSuccess(writer, responseData, "Updated successfully", http.StatusOK)

}

// @Summary Get All Warehouse Location Popup
// @Description REST API Warehouse Location Popup
// @Accept json
// @Produce json
// @Tags Master : Warehouse Location Definition
// @Security AuthorizationKeyAuth
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Param warehouse_location_definition_level_id query string false "warehouse_location_definition_level_id"
// @Param sort_by query string false "sort_by"
// @Param sort_of query string false "sort_of"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/warehouse-location-definition/popup-level [get]
func (r *WarehouseLocationDefinitionControllerImpl) PopupWarehouseLocationLevel(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"mtr_warehouse_location_definition_level.warehouse_location_definition_level_id":          queryValues.Get("warehouse_location_definition_level_id"),
		"mtr_warehouse_location_definition_level.warehouse_location_definition_level_description": queryValues.Get("warehouse_location_definition_level_description"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	paginatedData, totalPages, totalRows, err := r.WarehouseLocationDefinitionService.PopupWarehouseLocationLevel(criteria, paginate)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(paginatedData), "Get Data Successfully", http.StatusOK, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
}
