package transactionsparepartcontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionsparepartservice "after-sales/api/services/transaction/sparepart"
	"after-sales/api/utils"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type PurchaseRequestController interface {
	GetAllPurchaseRequest(writer http.ResponseWriter, request *http.Request)
	GetByIdPurchaseRequest(writer http.ResponseWriter, request *http.Request)
	GetAllPurchaseRequestDetail(writer http.ResponseWriter, request *http.Request)
	GetByIdPurchaseRequestDetail(writer http.ResponseWriter, request *http.Request)
	NewPurchaseRequestHeader(writer http.ResponseWriter, request *http.Request)
	NewPurchaseRequestDetail(writer http.ResponseWriter, request *http.Request)
	Void(writer http.ResponseWriter, request *http.Request)
	UpdatePurchaseRequestHeader(writer http.ResponseWriter, request *http.Request)
	UpdatePurchaseRequestDetail(writer http.ResponseWriter, request *http.Request)
	SubmitPurchaseRequestHeader(writer http.ResponseWriter, request *http.Request)
	SubmitPurchaseRequestDetail(writer http.ResponseWriter, request *http.Request)
}

type PurchaseRequestControllerImpl struct {
	PurchaseRequestService transactionsparepartservice.PurchaseRequestService
}

func NewPurchaseRequestController(PurchaseRequestService transactionsparepartservice.PurchaseRequestService) PurchaseRequestController {
	return &PurchaseRequestControllerImpl{PurchaseRequestService: PurchaseRequestService}
}

// GetAllPurchaseRequest
//
//	@Summary		Get All Purchase Request
//	@Description	REST API Purchase Request
//	@Accept			json
//	@Produce		json
//	@Tags			Transaction : Purchase Request
//	@Param			page								query		string	true	"page"
//	@Param			limit								query		string	true	"limit"
//	@Param			purchase_request_document_number	query		string	false	"purchase_request_document_number"
//	@Param			purchase_order_no					query		string	false	"purchase_order_no"
//	@Param			purchase_request_date_from			query		string	false	"purchase_request_date_from"
//	@Param			purchase_request_date_to			query		string	false	"purchase_request_date_to"
//	@Param			item_group_id						query		int		false	"item_group_id"
//	@Param			reference_document_number			query		string	false	"reference_document_number"
//	@Param			purchase_request_document_status_id	query		int		false	"purchase_request_document_status_id"
//	@Param			created_by_user_id					query		int		false	"created_by_user_id"
//	@Param			cost_center_id						query		int		false	"cost_center_id"
//	@Param			sort_by								query		string	false	"sort_by"
//	@Param			sort_of								query		string	false	"sort_of"
//	@Success		200									{object}	[]transactionsparepartpayloads.PurchaseRequestGetAllListResponses
//	@Failure		500,400,401,404,403,422				{object}	exceptions.BaseErrorResponse
//	@Router			/v1/purchase-request/ [get]
func (controller *PurchaseRequestControllerImpl) GetAllPurchaseRequest(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	queryParams := map[string]string{
		"purchase_request_document_number":    queryValues.Get("purchase_request_document_number"),
		"item_group_id":                       queryValues.Get("item_group_id"),
		"reference_document_number":           queryValues.Get("reference_document_number"),
		"created_by_user_id":                  queryValues.Get("created_by_user_id"),
		"cost_center_id":                      queryValues.Get("cost_center_id"),
		"purchase_request_document_status_id": queryValues.Get("purchase_request_document_status_id"),
	}
	DateParams := map[string]string{
		"purchase_request_date_from": queryValues.Get("purchase_request_date_from"),
		"purchase_request_date_to":   queryValues.Get("purchase_request_date_to"),
	}
	fmt.Println(DateParams["purchase_request_date_from"])
	pagination := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}
	filter := utils.BuildFilterCondition(queryParams)
	result, err := controller.PurchaseRequestService.GetAllPurchaseRequest(filter, pagination, DateParams)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data Successfull", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)
}

// GetByIdPurchaseRequest
//
//	@Summary		Get By Id Purchase Request
//	@Description	REST API Purchase Request
//	@Accept			json
//	@Produce		json
//	@Tags			Transaction : Purchase Request
//	@Param			purchase_request_system_number	path		int	true	"purchase_request_system_number"
//	@Success		200								{object}	transactionsparepartpayloads.PurchaseRequestGetAllListResponses
//	@Failure		500,400,401,404,403,422			{object}	exceptions.BaseErrorResponse
//	@Router			/v1/purchase-request/by-id/{purchase_request_system_number} [get]
func (controller *PurchaseRequestControllerImpl) GetByIdPurchaseRequest(writer http.ResponseWriter, request *http.Request) {
	PurchaseRequestSystemNumber, _ := strconv.Atoi(chi.URLParam(request, "purchase_request_system_number"))
	result, err := controller.PurchaseRequestService.GetByIdPurchaseRequest(PurchaseRequestSystemNumber)
	if err != nil {
		err.Message = "Id Not Found"
		err.Data = "Id Not Found"
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// GetAllPurchaseRequestDetail
//
//	@Summary		Get All Purchase Request Detail
//	@Description	REST API Purchase Request Detail
//	@Accept			json
//	@Produce		json
//	@Tags			Transaction : Purchase Request
//	@Param			page							query		string	true	"page"
//	@Param			limit							query		string	true	"limit"
//	@Param			purchase_request_system_number	query		string	false	"purchase_request_system_number"
//	@Param			sort_by							query		string	false	"sort_by"
//	@Param			sort_of							query		string	false	"sort_of"
//	@Success		200								{object}	transactionsparepartpayloads.PurchaseRequestDetailResponsesPayloads
//	@Failure		500,400,401,404,403,422			{object}	exceptions.BaseErrorResponse
//	@Router			/v1/purchase-request/detail [get]
func (controller *PurchaseRequestControllerImpl) GetAllPurchaseRequestDetail(writer http.ResponseWriter, request *http.Request) {
	//TODO implement me
	//PurchaseRequestSystemNumber := chi.URLParam(request, "purchase_request_system_number")

	queryValues := request.URL.Query()
	queryParams := map[string]string{
		"purchase_request_system_number": queryValues.Get("purchase_request_system_number"),
	}
	pagination := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}
	filterCondition := utils.BuildFilterCondition(queryParams)
	result, err := controller.PurchaseRequestService.GetAllPurchaseRequestDetail(filterCondition, pagination)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data Successfully!", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)
}

// GetByIdPurchaseRequestDetail
//
//	@Summary		Get By Id Purchase Request Detail
//	@Description	REST API Purchase Request Detail
//	@Accept			json
//	@Produce		json
//	@Tags			Transaction : Purchase Request
//	@Param			purchase_request_system_number_detail	path		int	true	"purchase_request_system_number_detail"
//	@Success		200										{object}	transactionsparepartpayloads.PurchaseRequestDetailResponsesPayloads
//	@Failure		500,400,401,404,403,422					{object}	exceptions.BaseErrorResponse
//	@Router			/v1/purchase-request/by-id/{purchase_request_system_number_detail}/detail [get]
func (controller *PurchaseRequestControllerImpl) GetByIdPurchaseRequestDetail(writer http.ResponseWriter, request *http.Request) {
	//TODO implement mee
	PurchaseRequestSystemNumberDetail, _ := strconv.Atoi(chi.URLParam(request, "purchase_request_system_number_detail"))
	result, err := controller.PurchaseRequestService.GetByIdPurchaseRequestDetail(PurchaseRequestSystemNumberDetail)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, result, "Get Data Successfully", http.StatusOK)
	//panic("implement me"
	//}
}

// SavePurchaseRequestHeader
//
//	@Summary		Create New Purchase Request
//	@Description	Create a new SaveHeader
//	@Accept			json
//	@Produce		json
//	@Tags			Transaction : Purchase Request
//	@Param			reqBody					body		transactionsparepartpayloads.PurchaseRequestHeaderSaveRequest	true	"Purchase Request Header Data"
//	@Success		201						{object}	payloads.Response
//	@Failure		500,400,401,404,403,422	{object}	exceptions.BaseErrorResponse
//	@Router			/v1/purchase-request [post]
func (controller *PurchaseRequestControllerImpl) NewPurchaseRequestHeader(writer http.ResponseWriter, request *http.Request) {
	var purchaseRequest transactionsparepartpayloads.PurchaseRequestHeaderSaveRequest

	helper.ReadFromRequestBody(request, &purchaseRequest)

	success, err := controller.PurchaseRequestService.NewPurchaseRequestHeader(purchaseRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, success, "save success", http.StatusOK)

}

// SavePurchaseRequestDetail
//
//	@Summary		Create New Purchase Request Detail
//	@Description	Create a new Save Detail
//	@Accept			json
//	@Produce		json
//	@Tags			Transaction : Purchase Request
//	@Param			reqBody					body		transactionsparepartpayloads.PurchaseRequestSaveDetailRequestPayloads	true	"Purchase Request Header Data"
//	@Success		201						{object}	payloads.Response
//	@Failure		500,400,401,404,403,422	{object}	exceptions.BaseErrorResponse
//	@Router			/v1/purchase-request/detail [post]
func (controller *PurchaseRequestControllerImpl) NewPurchaseRequestDetail(writer http.ResponseWriter, request *http.Request) {
	var purchaseRequest transactionsparepartpayloads.PurchaseRequestSaveDetailRequestPayloads
	//var purchaseRequest transactionsparepartpayloads.PurchaseRequestHeaderSaveRequest

	helper.ReadFromRequestBody(request, &purchaseRequest)

	success, err := controller.PurchaseRequestService.NewPurchaseRequestDetail(purchaseRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, success, "save success", http.StatusOK)

}

// UpdatePurchaseRequestHeader
//
//	@Summary		Update Purchase Request Header
//	@Description	Update Purchase Request Header
//	@Accept			json
//	@Produce		json
//	@Tags			Transaction : Purchase Request
//	@Param			purchase_request_system_number	path		int	true	"purchase_request_system_number"
//	@Param			reqBody					body		transactionsparepartpayloads.PurchaseRequestHeaderSaveRequest	true	"Purchase Request Header Data"
//	@Success		201						{object}	payloads.Response
//	@Failure		500,400,401,404,403,422	{object}	exceptions.BaseErrorResponse
//	@Router			/v1/purchase-request/{purchase_request_system_number} [put]
func (controller *PurchaseRequestControllerImpl) UpdatePurchaseRequestHeader(writer http.ResponseWriter, request *http.Request) {
	//TODO implement me
	var puchaseRequestHeader transactionsparepartpayloads.PurchaseRequestHeaderSaveRequest
	PurchaseRequestSystemNumber, _ := strconv.Atoi(chi.URLParam(request, "purchase_request_system_number"))

	helper.ReadFromRequestBody(request, &puchaseRequestHeader)
	success, err := controller.PurchaseRequestService.SavePurchaseRequestUpdateHeader(puchaseRequestHeader, PurchaseRequestSystemNumber)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, success, "save success", http.StatusOK)

}

// UpdatePurchaseRequestDetail
//
// @Summary		Update Purchase Request Detail
// @Description	Update Purchase Request Detail
// @Accept			json
// @Produce		json
// @Tags			Transaction : Purchase Request
// @Param			purchase_request_detail_system_number	path		int	true	"purchase_request_detail_system_number"
// @Param			reqBody					body		transactionsparepartpayloads.PurchaseRequestSaveDetailRequestPayloads	true	"Purchase Request Header Data"
// @Success		201						{object}	transactionsparepartpayloads.PurchaseRequestSaveDetailRequestPayloads
// @Failure		500,400,401,404,403,422	{object}	exceptions.BaseErrorResponse
// @Router			/v1/purchase-request/detail/{purchase_request_detail_system_number} [put]
func (controller *PurchaseRequestControllerImpl) UpdatePurchaseRequestDetail(writer http.ResponseWriter, request *http.Request) {
	//TODO implement me
	PurchaseRequestSystemNumberDetail, _ := strconv.Atoi(chi.URLParam(request, "purchase_request_detail_system_number"))

	var puchaseRequestDetail transactionsparepartpayloads.PurchaseRequestSaveDetailRequestPayloads
	helper.ReadFromRequestBody(request, &puchaseRequestDetail)
	success, err := controller.PurchaseRequestService.SavePurchaseRequestUpdateDetail(puchaseRequestDetail, PurchaseRequestSystemNumberDetail)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, success, "save success", http.StatusOK)
}

// Void
//
// @Summary		Void Request Detail
// @Description	Void Request Detail
// @Accept			json
// @Produce		json
// @Tags			Transaction : Purchase Request
// @Param			purchase_request_system_number	path		int	true	"purchase_request_system_number"
// @Success		201						{object}	payloads.Response
// @Failure		500,400,401,404,403,422	{object}	exceptions.BaseErrorResponse
// @Router			/v1/purchase-request/{purchase_request_system_number} [delete]
func (controller *PurchaseRequestControllerImpl) Void(writer http.ResponseWriter, request *http.Request) {
	// Void work order
	PurchaseRequestId := chi.URLParam(request, "purchase_request_system_number")
	PurchaseRequestsysno, err := strconv.Atoi(PurchaseRequestId)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid work order ID", http.StatusBadRequest)
		return
	}
	success, baseErr := controller.PurchaseRequestService.VoidPurchaseRequest(PurchaseRequestsysno)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			helper.ReturnError(writer, request, baseErr)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	if success {
		payloads.NewHandleSuccess(writer, nil, "Purchase Order voided successfully", http.StatusOK)
	} else {
		helper.ReturnError(writer, request, baseErr)
	}
}

// SubmitPurchaseRequestHeader
//
//	@Summary		Submit Purchase Request Header
//	@Description	Submit Purchase Request Header
//	@Accept			json
//	@Produce		json
//	@Tags			Transaction : Purchase Request
//	@Param			purchase_request_system_number	path		int	true	"purchase_request_system_number"
//	@Param			reqBody					body		transactionsparepartpayloads.PurchaseRequestHeaderSaveRequest	true	"Purchase Request Header Data"
//	@Success		201						{object}	payloads.Response
//	@Failure		500,400,401,404,403,422	{object}	exceptions.BaseErrorResponse
//	@Router			/v1/purchase-request/submit/{purchase_request_system_number} [post]
func (controller *PurchaseRequestControllerImpl) SubmitPurchaseRequestHeader(writer http.ResponseWriter, request *http.Request) {
	var puchaseRequestHeader transactionsparepartpayloads.PurchaseRequestHeaderSaveRequest
	PurchaseRequestSystemNumber, _ := strconv.Atoi(chi.URLParam(request, "purchase_request_system_number"))

	helper.ReadFromRequestBody(request, &puchaseRequestHeader)
	success, err := controller.PurchaseRequestService.InsertPurchaseRequestUpdateHeader(puchaseRequestHeader, PurchaseRequestSystemNumber)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, success, "save success", http.StatusOK)

}

// SubmitPurchaseRequestDetail
//
// @Summary		Submit Purchase Request Detail
// @Description	Submit Purchase Request Detail
// @Accept			json
// @Produce		json
// @Tags			Transaction : Purchase Request
// @Param			purchase_request_detail_system_number	path		int	true	"purchase_request_detail_system_number"
// @Param			reqBody					body		transactionsparepartpayloads.PurchaseRequestSaveDetailRequestPayloads	true	"Purchase Request Header Data"
// @Success		201						{object}	payloads.Response
// @Failure		500,400,401,404,403,422	{object}	exceptions.BaseErrorResponse
// @Router			/v1/purchase-request/submit/detail/{purchase_request_detail_system_number} [post]
func (controller *PurchaseRequestControllerImpl) SubmitPurchaseRequestDetail(writer http.ResponseWriter, request *http.Request) {
	PurchaseRequestSystemNumberDetail, _ := strconv.Atoi(chi.URLParam(request, "purchase_request_detail_system_number"))

	var puchaseRequestDetail transactionsparepartpayloads.PurchaseRequestSaveDetailRequestPayloads
	helper.ReadFromRequestBody(request, &puchaseRequestDetail)
	success, err := controller.PurchaseRequestService.InsertPurchaseRequestUpdateDetail(puchaseRequestDetail, PurchaseRequestSystemNumberDetail)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, success, "save success", http.StatusOK)
}
