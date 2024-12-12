package transactionsparepartcontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/payloads"
	"after-sales/api/payloads/pagination"
	transactionsparepartservice "after-sales/api/services/transaction/sparepart"
	"after-sales/api/utils"
	"net/http"
)

type ItemLocationTransferControllerImpl struct {
	ItemLocationTransferService transactionsparepartservice.ItemLocationTransferService
}

type ItemLocationTransferController interface {
	GetAllItemLocationTransfer(writer http.ResponseWriter, request *http.Request)
}

func NewItemLocationTransferController(
	itemLocationTransferServiceService transactionsparepartservice.ItemLocationTransferService,
) ItemLocationTransferController {
	return &ItemLocationTransferControllerImpl{
		ItemLocationTransferService: itemLocationTransferServiceService,
	}
}
func (r *ItemLocationTransferControllerImpl) GetAllItemLocationTransfer(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"trx_item_warehouse_transfer_request.company_id":                       queryValues.Get("company_id"),
		"trx_item_warehouse_transfer_request.transfer_request_document_number": queryValues.Get("transfer_request_document_number"),
		"trx_item_warehouse_transfer_request.transfer_request_date_from":       queryValues.Get("transfer_request_date_from"),
		"trx_item_warehouse_transfer_request.transfer_request_date_to":         queryValues.Get("transfer_request_date_to"),
		"RequestFromWarehouse.warehouse_group_id":                              queryValues.Get("request_from_warehouse_group_id"),
		"RequestToWarehouse.warehouse_group_id":                                queryValues.Get("request_to_warehouse_group_id"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	response, err := r.ItemLocationTransferService.GetAllItemLocationTransfer(criteria, paginate)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(
		writer,
		response.Rows,
		"Get Data Successfully!",
		http.StatusOK,
		response.Limit,
		response.Page,
		int64(response.TotalRows),
		response.TotalPages,
	)
}
