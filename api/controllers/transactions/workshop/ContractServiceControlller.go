package transactionworkshopcontroller

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads"
	"after-sales/api/payloads/pagination"
	transactionworkshopservice "after-sales/api/services/transaction/workshop"
	"after-sales/api/utils"
	"fmt"
	"net/http"
)

type ContractServiceControllerImpl struct {
	ContractServiceService transactionworkshopservice.ContractServiceService
}

type ContractServiceController interface {
	GetAll(writer http.ResponseWriter, request *http.Request)
	// GetById(writer http.ResponseWriter, request *http.Request)
	// Qcpass(writer http.ResponseWriter, request *http.Request)
	// Reorder(writer http.ResponseWriter, request *http.Request)
}

func NewContractServiceController(ContractServiceService transactionworkshopservice.ContractServiceService) ContractServiceController {
	return &ContractServiceControllerImpl{
		ContractServiceService: ContractServiceService,
	}
}

func (r *ContractServiceControllerImpl) GetAll(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	// Mengambil query params sesuai dengan SQL VB yang diberikan
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

	// Pagination params
	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	// Menggunakan utils untuk membangun filter condition dari queryParams
	criteria := utils.BuildFilterCondition(queryParams)
	fmt.Println("Filter Conditions:", criteria)

	// Panggil service untuk mendapatkan data sesuai filter
	paginatedData, totalPages, totalRows, err := r.ContractServiceService.GetAll(criteria, paginate)
	if err != nil {
		// Jika ada error
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	// Jika data ditemukan
	if len(paginatedData) > 0 {
		payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(paginatedData), "Get Data Successfully", http.StatusOK, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
	} else {
		// Jika data tidak ditemukan
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}
}
