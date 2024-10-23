package transactionsparepartcontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/payloads"
	"after-sales/api/payloads/pagination"
	transactionsparepartservice "after-sales/api/services/transaction/sparepart"
	"after-sales/api/utils"
	"net/http"
)

type ItemInquiryController interface {
	GetAllItemInquiry(writer http.ResponseWriter, request *http.Request)
}

type ItemInquiryControllerImpl struct {
	ItemInquiryService transactionsparepartservice.ItemInquiryService
}

func NewItemInquiryController(itemInquiryService transactionsparepartservice.ItemInquiryService) ItemInquiryController {
	return &ItemInquiryControllerImpl{
		ItemInquiryService: itemInquiryService,
	}
}

func (i *ItemInquiryControllerImpl) GetAllItemInquiry(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	queryParams := map[string]string{
		"company_id":         queryValues.Get("company_id"),
		"company_session_id": queryValues.Get("company_session_id"),
		"mi.item_id":         queryValues.Get("item_id"),
	}

	if queryParams["company_id"] == "" {
		payloads.NewHandleError(writer, "company_id is required", http.StatusBadRequest)
		return
	}

	if queryParams["company_session_id"] == "" {
		payloads.NewHandleError(writer, "company_session_id is required", http.StatusBadRequest)
		return
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	result, err := i.ItemInquiryService.GetAllItemInquiry(criteria, paginate)

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(result.Rows), "Get Data Successfully", http.StatusOK, result.Limit, result.Page, result.TotalRows, result.TotalPages)
}
