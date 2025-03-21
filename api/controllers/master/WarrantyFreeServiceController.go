package mastercontroller

import (
	masterentities "after-sales/api/entities/master"
	"after-sales/api/exceptions"
	helper "after-sales/api/helper"
	"after-sales/api/payloads"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterservice "after-sales/api/services/master"
	"after-sales/api/utils"
	"after-sales/api/validation"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type WarrantyFreeServiceController interface {
	GetAllWarrantyFreeService(writer http.ResponseWriter, request *http.Request)
	GetWarrantyFreeServiceByID(writer http.ResponseWriter, request *http.Request)
	SaveWarrantyFreeService(writer http.ResponseWriter, request *http.Request)
	ChangeStatusWarrantyFreeService(writer http.ResponseWriter, request *http.Request)
	UpdateWarrantyFreeService(writer http.ResponseWriter, request *http.Request)
}
type WarrantyFreeServiceControllerImpl struct {
	WarrantyFreeServiceService masterservice.WarrantyFreeServiceService
}

func NewWarrantyFreeServiceController(warrantyFreeServiceService masterservice.WarrantyFreeServiceService) WarrantyFreeServiceController {
	return &WarrantyFreeServiceControllerImpl{
		WarrantyFreeServiceService: warrantyFreeServiceService,
	}
}

// @Summary Get All Warranty Free Services
// @Description Retrieve all warranty free services with optional filtering and pagination
// @Accept json
// @Produce json
// @Tags Master : Warranty Free Service
// @Security AuthorizationKeyAuth
// @Param page query string true "Page number"
// @Param limit query string true "Items per page"
// @Param is_active query string false "Is active" Enums(true,false)
// @Param effective_date query string false "Effective date"
// @Param brand_code query string false "Brand code"
// @Param model_code query string false "Model code"
// @Param warranty_free_service_type_code query string false "Warranty free service type code"
// @Param sort_by query string false "Field to sort by"
// @Param sort_of query string false "Sort order (asc/desc)"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/warranty-free-service/ [get]
func (r *WarrantyFreeServiceControllerImpl) GetAllWarrantyFreeService(writer http.ResponseWriter, request *http.Request) {

	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"mtr_warranty_free_service.is_active":      queryValues.Get("is_active"),
		"mtr_warranty_free_service.effective_date": queryValues.Get("effective_date"),
		"brand_code":                      queryValues.Get("brand_code"),
		"model_code":                      queryValues.Get("model_code"),
		"warranty_free_service_type_code": queryValues.Get("warranty_free_service_type_code"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	result, err := r.WarrantyFreeServiceService.GetAllWarrantyFreeService(criteria, paginate)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(
		writer,
		result.Rows,
		"Get Data Successfully!",
		http.StatusOK,
		result.Limit,
		result.Page,
		int64(result.TotalRows),
		result.TotalPages,
	)
}

// @Summary Get Warranty Free Service By ID
// @Description Retrieve a warranty free service by its ID
// @Accept json
// @Produce json
// @Tags Master : Warranty Free Service
// @Security AuthorizationKeyAuth
// @Param warranty_free_services_id path int true "Warranty Free Service ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/warranty-free-service/{warranty_free_services_id} [get]
func (r *WarrantyFreeServiceControllerImpl) GetWarrantyFreeServiceByID(writer http.ResponseWriter, request *http.Request) {

	warrantyFreeServiceId, errA := strconv.Atoi(chi.URLParam(request, "warranty_free_services_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	result, err := r.WarrantyFreeServiceService.GetWarrantyFreeServiceById(warrantyFreeServiceId)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, utils.ModifyKeysInResponse(result), "Get Data Successfully!", http.StatusOK)
}

// @Summary Save Warranty Free Service
// @Description Create or update a warranty free service
// @Accept json
// @Produce json
// @Tags Master : Warranty Free Service
// @Security AuthorizationKeyAuth
// @Param reqBody body masterpayloads.WarrantyFreeServiceRequest true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/warranty-free-service/ [post]
func (r *WarrantyFreeServiceControllerImpl) SaveWarrantyFreeService(writer http.ResponseWriter, request *http.Request) {

	var formRequest masterpayloads.WarrantyFreeServiceRequest
	helper.ReadFromRequestBody(request, &formRequest)
	if validationErr := validation.ValidationForm(writer, request, &formRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}
	var message string

	create, err := r.WarrantyFreeServiceService.SaveWarrantyFreeService(formRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	if formRequest.WarrantyFreeServicesId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

// @Summary Change Status Warranty Free Service
// @Description Change the status of a warranty free service
// @Accept json
// @Produce json
// @Tags Master : Warranty Free Service
// @Security AuthorizationKeyAuth
// @Param warranty_free_services_id path int true "Warranty Free Service ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/warranty-free-service/{warranty_free_services_id} [patch]
func (r *WarrantyFreeServiceControllerImpl) ChangeStatusWarrantyFreeService(writer http.ResponseWriter, request *http.Request) {

	warrantyFreeServiceId, errA := strconv.Atoi(chi.URLParam(request, "warranty_free_services_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	response, err := r.WarrantyFreeServiceService.ChangeStatusWarrantyFreeService(warrantyFreeServiceId)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}

// @Summary Update Warranty Free Service
// @Description Update a warranty free service
// @Accept json
// @Produce json
// @Tags Master : Warranty Free Service
// @Security AuthorizationKeyAuth
// @Param warranty_free_services_id path int true "Warranty Free Service ID"
// @Param reqBody body masterentities.WarrantyFreeService true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/warranty-free-service/{warranty_free_services_id} [put]
func (r *WarrantyFreeServiceControllerImpl) UpdateWarrantyFreeService(writer http.ResponseWriter, request *http.Request) {
	warranty_free_services_id, errA := strconv.Atoi(chi.URLParam(request, "warranty_free_services_id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}
	var formRequest masterentities.WarrantyFreeService
	helper.ReadFromRequestBody(request, &formRequest)
	if validationErr := validation.ValidationForm(writer, request, &formRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}
	result, err := r.WarrantyFreeServiceService.UpdateWarrantyFreeService(formRequest, warranty_free_services_id)
	if err != nil {
		exceptions.NewConflictException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Update Data Successfully!", http.StatusOK)
}
