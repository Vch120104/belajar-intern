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

// GetAllDetail implements ContractServiceDetailController.
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

// GetById implements ContractServiceDetailController.
func (c *ContractServiceDetailControllerImpl) GetById(writer http.ResponseWriter, request *http.Request) {
	Id, _ := strconv.Atoi(chi.URLParam(request, "contract_service_package_detail_system_number"))

	result, err := c.ContractServiceDetailService.GetById(Id)

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, result, "Get Data Successfully", http.StatusOK)
}

// SaveDetail implements ContractServiceDetailController.
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

// UpdateDetail implements ContractServiceDetailController.
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

// DeleteDetail implements ContractServiceDetailController.
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
