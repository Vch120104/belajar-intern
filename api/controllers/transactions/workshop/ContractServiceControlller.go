package transactionworkshopcontroller

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshopservice "after-sales/api/services/transaction/workshop"
	"after-sales/api/validation"

	"after-sales/api/utils"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type ContractServiceControllerImpl struct {
	ContractServiceService transactionworkshopservice.ContractServiceService
}

type ContractServiceController interface {
	GetAll(writer http.ResponseWriter, request *http.Request)
	GetById(writer http.ResponseWriter, request *http.Request)
	Save(writer http.ResponseWriter, request *http.Request)
	Void(writer http.ResponseWriter, request *http.Request)
	Submit(writer http.ResponseWriter, request *http.Request)
}

func NewContractServiceController(ContractServiceService transactionworkshopservice.ContractServiceService) ContractServiceController {
	return &ContractServiceControllerImpl{
		ContractServiceService: ContractServiceService,
	}
}

func (r *ContractServiceControllerImpl) GetAll(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"trx_contract_service.company_id":                       queryValues.Get("company_code"),
		"trx_contract_service.contract_service_document_number": queryValues.Get("contract_serv_doc_no"),
		"trx_contract_service.contract_service_from":            queryValues.Get("date_from"),
		"trx_contract_service.contract_service_to":              queryValues.Get("date_to"),
		"trx_contract_service.brand_id":                         queryValues.Get("vehicle_brand"),
		"trx_contract_service.model_id":                         queryValues.Get("model_code"),
		"mtr_vehicle_registration_certificate.vehicle_tnkb":     queryValues.Get("tnkb"),
		"trx_contract_service.contract_service_status_id":       queryValues.Get("contract_serv_status"),
	}

	fmt.Println("Query Params:", queryParams)

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)
	fmt.Println("Filter Conditions:", criteria)

	result, err := r.ContractServiceService.GetAll(criteria, paginate)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
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

// GetById implements ContractServiceController.
func (r *ContractServiceControllerImpl) GetById(writer http.ResponseWriter, request *http.Request) {
	idstr := chi.URLParam(request, "contract_service_system_number")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid request ID", http.StatusBadRequest)
		return
	}

	queryValues := request.URL.Query()

	var filterConditions []utils.FilterCondition
	for field, value := range map[string]string{
		"trx_contract_service.contract_service_system_number": queryValues.Get("contract_service_system_number"),
	} {
		if value != "" {
			filterConditions = append(filterConditions, utils.FilterCondition{
				ColumnField: field,
				ColumnValue: value,
			})
		}
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	data, baseErr := r.ContractServiceService.GetById(id, filterConditions, paginate)
	if baseErr != nil {
		exceptions.NewAppException(writer, request, baseErr)
		return
	}

	payloads.NewHandleSuccess(writer, data, "Get Data Successfully", http.StatusOK)
}

// Save implements ContractServiceController.
func (r *ContractServiceControllerImpl) Save(writer http.ResponseWriter, request *http.Request) {
	var contractServiceInsert transactionworkshoppayloads.ContractServiceInsert
	helper.ReadFromRequestBody(request, &contractServiceInsert)
	if validationErr := validation.ValidationForm(writer, request, &contractServiceInsert); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}

	result, saveErr := r.ContractServiceService.Save(contractServiceInsert)
	if saveErr != nil {
		helper.ReturnError(writer, request, saveErr)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Contract Service Saved Successfully", http.StatusOK)
}

// Void implements ContractServiceController.
func (r *ContractServiceControllerImpl) Void(writer http.ResponseWriter, request *http.Request) {
	workOrderIdStr := chi.URLParam(request, "contract_service_system_number")
	workOrderId, err := strconv.Atoi(workOrderIdStr)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid contratct service ID", http.StatusBadRequest)
		return
	}

	success, baseErr := r.ContractServiceService.Void(workOrderId)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, baseErr.Message, http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	if success {
		payloads.NewHandleSuccess(writer, success, "contract service voided successfully", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Failed to void contract service", http.StatusInternalServerError)
	}
}

// Submit implements ContractServiceController.
func (r *ContractServiceControllerImpl) Submit(writer http.ResponseWriter, request *http.Request) {
	Id, _ := strconv.Atoi(chi.URLParam(request, "contract_service_system_number"))
	res, err := r.ContractServiceService.Submit(Id)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, res, "Successfully Submit Contract Service", http.StatusOK)
}
