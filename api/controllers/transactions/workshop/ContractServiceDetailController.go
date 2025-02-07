package transactionworkshopcontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	jsonchecker "after-sales/api/helper/json/json-checker"
	"after-sales/api/payloads"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshopservice "after-sales/api/services/transaction/workshop"
	"after-sales/api/utils"
	"after-sales/api/validation"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type ContractServiceDetailControllerImpl struct {
	ContractServiceDetailService transactionworkshopservice.ContractServiceDetailService
}

type ContractServiceDetailController interface {
	GetAllDetail(writer http.ResponseWriter, request *http.Request)
	GetById(writer http.ResponseWriter, request *http.Request)
	SaveDetail(writer http.ResponseWriter, request *http.Request)
	UpdateDetail(writer http.ResponseWriter, request *http.Request)
	DeleteDetail(writer http.ResponseWriter, request *http.Request)
}

func NewContractServiceDetailController(contractServiceDetailService transactionworkshopservice.ContractServiceDetailService) ContractServiceDetailController {
	return &ContractServiceDetailControllerImpl{
		ContractServiceDetailService: contractServiceDetailService,
	}
}

// @Summary Get All Contract Service Detail
// @Description Retrieve all contract service detail with optional filtering and pagination
// @Accept json
// @Produce json
// @Tags Transaction Workshop : Contract Service Detail
// @Param contract_service_system_number path string true "Contract Service System Number"
// @Param limit query string true "Items per page"
// @Param page query string true "Page number"
// @Param sort_of query string false "Sort order (asc/desc)"
// @Param sort_by query string false "Field to sort by"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/contract-service-detail/{contract_service_system_number} [get]
func (c *ContractServiceDetailControllerImpl) GetAllDetail(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	contractServiceSystemNumberStr := chi.URLParam(request, "contract_service_system_number")
	contractServiceSystemNumber, err := strconv.Atoi(contractServiceSystemNumberStr)
	if err != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("failed to read request param, please check your param input"),
		})
		return
	}

	pagination := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	queryParams := map[string]string{
		"contract_service_system_number": contractServiceSystemNumberStr,
	}
	filterCondition := utils.BuildFilterCondition(queryParams)

	result, errs := c.ContractServiceDetailService.GetAllDetail(contractServiceSystemNumber, filterCondition, pagination)
	if errs != nil {
		helper.ReturnError(writer, request, errs)
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

// @Summary Get Contract Service Detail By ID
// @Description Retrieve contract service detail by ID
// @Accept json
// @Produce json
// @Tags Transaction Workshop : Contract Service Detail
// @Param contract_service_package_detail_system_number path string true "Contract Service Package Detail System Number"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/contract-service-detail/{contract_service_package_detail_system_number} [get]
func (c *ContractServiceDetailControllerImpl) GetById(writer http.ResponseWriter, request *http.Request) {
	Id, _ := strconv.Atoi(chi.URLParam(request, "contract_service_package_detail_system_number"))

	result, err := c.ContractServiceDetailService.GetById(Id)

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, result, "Get Data Successfully", http.StatusOK)
}

// @Summary Save Contract Service Detail
// @Description Save contract service detail
// @Accept json
// @Produce json
// @Tags Transaction Workshop : Contract Service Detail
// @Param contract_service_system_number path string true "Contract Service System Number"
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/contract-service-detail/{contract_service_system_number} [post]
func (c *ContractServiceDetailControllerImpl) SaveDetail(writer http.ResponseWriter, request *http.Request) {
	formRequest := transactionworkshoppayloads.ContractServiceIdResponse{}

	err := jsonchecker.ReadFromRequestBody(request, &formRequest)
	if err != nil {
		exceptions.NewEntityException(writer, request, err)
		return
	}

	err = validation.ValidationForm(writer, request, formRequest)
	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	create, err := c.ContractServiceDetailService.SaveDetail(formRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, create, "Create Data Successfully", http.StatusCreated) // Menggunakan StatusCreated (201)
}

// @Summary Update Contract Service Detail
// @Description Update contract service detail
// @Accept json
// @Produce json
// @Tags Transaction Workshop : Contract Service Detail
// @Param contract_service_system_number path string true "Contract Service System Number"
// @Param contract_service_line path string true "Contract Service Line"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/contract-service-detail/{contract_service_system_number}/{contract_service_line} [put]
func (c *ContractServiceDetailControllerImpl) UpdateDetail(writer http.ResponseWriter, request *http.Request) {
	contractServiceSystemNumber, err := strconv.Atoi(chi.URLParam(request, "contract_service_system_number"))
	if err != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err,
			Message:    "Invalid contract service system number",
		})
		return
	}

	contractServiceLine := chi.URLParam(request, "contract_service_line")
	if contractServiceLine == "" {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid contract service line",
		})
		return
	}

	var detailRequest transactionworkshoppayloads.ContractServiceDetailRequest
	helper.ReadFromRequestBody(request, &detailRequest)
	if validationErr := validation.ValidationForm(writer, request, &detailRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err,
			Message:    "Validation error",
		})
		return
	}

	updatedDetail, iferr := c.ContractServiceDetailService.UpdateDetail(contractServiceSystemNumber, contractServiceLine, detailRequest)
	if iferr != nil {
		exceptions.NewAppException(writer, request, iferr)
		return
	}

	if updatedDetail.ContractServiceSystemNumber > 0 {
		payloads.NewHandleSuccess(writer, updatedDetail, "Detail updated successfully", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}
}

// @Summary Delete Contract Service Detail
// @Description Delete contract service detail
// @Accept json
// @Produce json
// @Tags Transaction Workshop : Contract Service Detail
// @Param contract_service_system_number path string true "Contract Service System Number"
// @Param package_code path string true "Package Code"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/contract-service-detail/{contract_service_system_number}/{package_code} [delete]
func (c *ContractServiceDetailControllerImpl) DeleteDetail(writer http.ResponseWriter, request *http.Request) {
	contractServiceSystemNumberStr := chi.URLParam(request, "contract_service_system_number")
	packageCode := chi.URLParam(request, "package_code")

	contractServiceSystemNumber, err := strconv.Atoi(contractServiceSystemNumberStr)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid contract service system number", http.StatusBadRequest)
		return
	}

	success, baseErr := c.ContractServiceDetailService.DeleteDetail(contractServiceSystemNumber, packageCode)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, baseErr.Message, http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}
	if success {
		payloads.NewHandleSuccess(writer, success, "Contract service detail deleted successfully", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Failed to delete contract service detail", http.StatusInternalServerError)
	}
}
