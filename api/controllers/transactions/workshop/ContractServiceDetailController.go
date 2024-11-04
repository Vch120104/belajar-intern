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

	result, totalPages, totalRows, errs := c.ContractServiceDetailService.GetAllDetail(contractServiceSystemNumber, filterCondition, pagination)
	if errs != nil {
		helper.ReturnError(writer, request, errs)
		return
	}

	if len(result) > 0 {
		payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(result), "Get Data Successfully", http.StatusOK, pagination.Limit, pagination.Page, int64(totalRows), totalPages)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}
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
