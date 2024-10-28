package transactionworkshopcontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads"
	"after-sales/api/payloads/pagination"
	transactionworkshopservice "after-sales/api/services/transaction/workshop"
	"after-sales/api/utils"
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
}

func NewContractServiceDetailController(contractServiceDetailService transactionworkshopservice.ContractServiceDetailService) ContractServiceDetailController {
	return &ContractServiceDetailControllerImpl{
		ContractServiceDetailService: contractServiceDetailService,
	}
}

// GetAllDetail implements ContractServiceDetailController.
func (c *ContractServiceDetailControllerImpl) GetAllDetail(writer http.ResponseWriter, request *http.Request) {
	// Mengambil query parameters
	queryValues := request.URL.Query()

	// Mengambil contract_service_system_number dari URL path parameter
	contractServiceSystemNumberStr := chi.URLParam(request, "contract_service_package_detail_system_number")
	contractServiceSystemNumber, err := strconv.Atoi(contractServiceSystemNumberStr)
	if err != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("failed to read request param, please check your param input"),
		})
		return
	}

	// Mempersiapkan pagination
	pagination := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	// Mendapatkan filter kondisi jika ada
	queryParams := map[string]string{
		"contract_service_system_number": contractServiceSystemNumberStr,
		// Tambahkan kondisi filter lain jika diperlukan
	}
	filterCondition := utils.BuildFilterCondition(queryParams)

	// Memanggil service untuk mendapatkan data
	result, totalPages, totalRows, errs := c.ContractServiceDetailService.GetAllDetail(contractServiceSystemNumber, filterCondition, pagination)
	if errs != nil {
		helper.ReturnError(writer, request, errs)
		return
	}

	// Mengembalikan response sukses dengan data paginated
	if len(result) > 0 {
		payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(result), "Get Data Successfully", http.StatusOK, pagination.Limit, pagination.Page, int64(totalRows), totalPages)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}
}
