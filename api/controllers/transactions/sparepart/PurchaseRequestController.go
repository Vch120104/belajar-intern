package transactionsparepartcontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionsparepartservice "after-sales/api/services/transaction/sparepart"
	"after-sales/api/utils"
	"after-sales/api/validation"
	"fmt"
	"net/http"
	"strconv"
	"time"

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
	VoidDetail(writer http.ResponseWriter, request *http.Request)

	UpdatePurchaseRequestHeader(writer http.ResponseWriter, request *http.Request)
	UpdatePurchaseRequestDetail(writer http.ResponseWriter, request *http.Request)
	SubmitPurchaseRequest(writer http.ResponseWriter, request *http.Request)
	SubmitPurchaseRequestDetail(writer http.ResponseWriter, request *http.Request)
	GetAllItemTypePr(writer http.ResponseWriter, request *http.Request)
	GetByIdItemTypePr(writer http.ResponseWriter, request *http.Request)
	GetByCodeItemTypePr(writer http.ResponseWriter, request *http.Request)
}

type PurchaseRequestControllerImpl struct {
	PurchaseRequestService transactionsparepartservice.PurchaseRequestService
}

func NewPurchaseRequestController(PurchaseRequestService transactionsparepartservice.PurchaseRequestService) PurchaseRequestController {
	return &PurchaseRequestControllerImpl{PurchaseRequestService: PurchaseRequestService}
}

// GetAllItemTypePr
//
//	@Summary		Get All Item For LookUp Purchase Request
//	@Description	REST API All Item For LookUp Purchase Request
//	@Accept			json
//	@Produce		json
//
// @Tags Transaction Sparepart: Purchase Request
//
//	@Param			page								query		string	true	"page"
//	@Param			limit								query		string	true	"limit"
//	@Param			item_code							query		string	false	"item_code"
//	@Param			company_id							query		string	true	"company_id"
//	@Param			item_name							query		string	false	"item_name"
//	@Param			item_class_name						query		string	false	"item_class_name"
//	@Param			item_type							query		string	false	"item_type"
//	@Param			item_level_1						query		string	false	"item_level_1"
//	@Param			item_level_2						query		string	false	"item_level_2"
//	@Param			item_level_3						query		string	false	"item_level_3"
//	@Param			item_level_4						query		string	false	"item_level_4"
//	@Param			quantity							query		int		false	"quantity"
//	@Param			sort_by								query		string	false	"sort_by"
//	@Param			sort_of								query		string	false	"sort_of"
//	@Success		200									{object}	[]transactionsparepartpayloads.PurchaseRequestItemGetAll
//	@Failure		500,400,401,404,403,422				{object}	exceptions.BaseErrorResponse
//	@Router			/v1/purchase-request/item [get]
func (controller *PurchaseRequestControllerImpl) GetAllItemTypePr(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	queryParams := map[string]string{
		// tambahin item group, stock keep, brand_id, stock_keeping = 1
		// item detail tidak boleh duplicate
		//
		"mtr_item.item_group_id": queryValues.Get("item_group_id"),
		"DT.brand_id":            queryValues.Get("brand_id"),
		"item_code":              queryValues.Get("item_code"),
		"item_name":              queryValues.Get("item_name"),
		"item_class_name":        queryValues.Get("item_class_name"),
		"IT.item_type_code":      queryValues.Get("item_type_code"),
		"L1.item_level_1_code":   queryValues.Get("item_level_1"),
		"L2.item_level_2_code":   queryValues.Get("item_level_2"),
		"L3.item_level_3_code":   queryValues.Get("item_level_3"),
		"L4.item_level_4_code":   queryValues.Get("item_level_4"),
		"quantity":               queryValues.Get("quantity"),
	}
	compid, _ := strconv.Atoi(queryValues.Get("company_id"))
	pagination := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}
	filter := utils.BuildFilterCondition(queryParams)
	result, err := controller.PurchaseRequestService.GetAllItemTypePurchaseRequest(filter, pagination, compid)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data Successfull", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)
}

// GetAllPurchaseRequest
//
//	@Summary		Get All Purchase Request
//	@Description	REST API Purchase Request
//	@Accept			json
//	@Produce		json
//
// @Tags Transaction Sparepart: Purchase Request
//
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
	if DateParams["purchase_request_date_from"] != "" {
		purchase_request_date_from, errsparseTime := time.Parse("2006-01-02T15:04:05.000Z", DateParams["purchase_request_date_from"])
		if errsparseTime != nil {
			payloads.NewHandleError(writer, "Invalid effective date", http.StatusBadRequest)
			return
		}
		DateParams["purchase_request_date_from"] = purchase_request_date_from.Format("2006-01-02")
	}
	if DateParams["purchase_request_date_to"] != "" {

		purchase_request_date_to, errsparse := time.Parse("2006-01-02T15:04:05.000Z", DateParams["purchase_request_date_to"])
		if errsparse != nil {
			payloads.NewHandleError(writer, "Invalid effective date", http.StatusBadRequest)
			return
		}
		DateParams["purchase_request_date_to"] = purchase_request_date_to.Format("2006-01-02")
		DateParams["purchase_request_date_to"] = DateParams["purchase_request_date_to"] + " 23:59:59.999"
	}

	fmt.Println(DateParams["purchase_request_date_to"])
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
//
// @Tags Transaction Sparepart: Purchase Request
//
//	@Param			purchase_request_system_number	path		int	true	"purchase_request_system_number"
//	@Success		200								{object}	transactionsparepartpayloads.PurchaseRequestGetAllListResponses
//	@Failure		500,400,401,404,403,422			{object}	exceptions.BaseErrorResponse
//	@Router			/v1/purchase-request/by-id/{purchase_request_system_number} [get]
func (controller *PurchaseRequestControllerImpl) GetByIdPurchaseRequest(writer http.ResponseWriter, request *http.Request) {
	PurchaseRequestSystemNumber, _ := strconv.Atoi(chi.URLParam(request, "purchase_request_system_number"))
	result, err := controller.PurchaseRequestService.GetByIdPurchaseRequest(PurchaseRequestSystemNumber)
	if err != nil {
		//err.Message = "Id Not Found"
		//err.Data = "Id Not Found"
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
//
// @Tags Transaction Sparepart: Purchase Request
//
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
		//helper.ReturnError(writer, request, err)
		exceptions.NewAppException(writer, request, err)
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
//
// @Tags Transaction Sparepart: Purchase Request
//
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

// NewPurchaseRequestHeader
//
//	@Summary		Create New Purchase Request
//	@Description	Create a new SaveHeader
//	@Accept			json
//	@Produce		json
//
// @Tags Transaction Sparepart: Purchase Request
//
//	@Param			reqBody					body		transactionsparepartpayloads.PurchaseRequestHeaderSaveRequest	true	"Purchase Request Header Data"
//	@Success		201						{object}	payloads.Response
//	@Failure		500,400,401,404,403,422	{object}	exceptions.BaseErrorResponse
//	@Router			/v1/purchase-request [post]
func (controller *PurchaseRequestControllerImpl) NewPurchaseRequestHeader(writer http.ResponseWriter, request *http.Request) {
	var purchaseRequest transactionsparepartpayloads.PurchaseRequestHeaderSaveRequest

	helper.ReadFromRequestBody(request, &purchaseRequest)
	if validationErr := validation.ValidationForm(writer, request, &purchaseRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}

	success, err := controller.PurchaseRequestService.NewPurchaseRequestHeader(purchaseRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, success, "save success", http.StatusOK)

}

// NewPurchaseRequestDetail
//
//	@Summary		Create New Purchase Request Detail
//	@Description	Create a new Save Detail
//	@Accept			json
//	@Produce		json
//
// @Tags Transaction Sparepart: Purchase Request
//
//	@Param			reqBody					body		transactionsparepartpayloads.PurchaseRequestSaveDetailRequestPayloads	true	"Purchase Request Header Data"
//	@Success		201						{object}	payloads.Response
//	@Failure		500,400,401,404,403,422	{object}	exceptions.BaseErrorResponse
//	@Router			/v1/purchase-request/detail [post]
func (controller *PurchaseRequestControllerImpl) NewPurchaseRequestDetail(writer http.ResponseWriter, request *http.Request) {
	var purchaseRequest transactionsparepartpayloads.PurchaseRequestSaveDetailRequestPayloads
	//var purchaseRequest transactionsparepartpayloads.PurchaseRequestHeaderSaveRequest

	helper.ReadFromRequestBody(request, &purchaseRequest)
	if validationErr := validation.ValidationForm(writer, request, &purchaseRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}

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
//
// @Tags Transaction Sparepart: Purchase Request
//
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
	if validationErr := validation.ValidationForm(writer, request, &puchaseRequestHeader); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}
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
// @Tags Transaction Sparepart: Purchase Request
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
	if validationErr := validation.ValidationForm(writer, request, &puchaseRequestDetail); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}
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
// @Tags Transaction Sparepart: Purchase Request
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

// SubmitPurchaseRequest
//
//	@Summary		Submit Purchase Request Header
//	@Description	Submit Purchase Request Header
//	@Accept			json
//	@Produce		json
//
// @Tags Transaction Sparepart: Purchase Request
//
//	@Param			purchase_request_system_number	path		int	true	"purchase_request_system_number"
//	@Success		201						{object}	payloads.Response
//	@Failure		500,400,401,404,403,422	{object}	exceptions.BaseErrorResponse
//	@Router			/v1/purchase-request/submit/{purchase_request_system_number} [post]
func (controller *PurchaseRequestControllerImpl) SubmitPurchaseRequest(writer http.ResponseWriter, request *http.Request) {
	var puchaseRequestHeader transactionsparepartpayloads.PurchaseRequestHeaderSaveRequest
	PurchaseRequestSystemNumber, _ := strconv.Atoi(chi.URLParam(request, "purchase_request_system_number"))
	success, err := controller.PurchaseRequestService.SubmitPurchaseRequest(puchaseRequestHeader, PurchaseRequestSystemNumber)
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
	if validationErr := validation.ValidationForm(writer, request, &puchaseRequestDetail); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}
	success, err := controller.PurchaseRequestService.InsertPurchaseRequestUpdateDetail(puchaseRequestDetail, PurchaseRequestSystemNumberDetail)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, success, "save success", http.StatusOK)
}

// GetByIdItemTypePr
//
//	@Summary		Get By Id Purchase Request Item Lookup
//	@Description	REST API Purchase Request Item Lookup
//	@Accept			json
//	@Produce		json
//
// @Tags Transaction Sparepart: Purchase Request
//
//	@Param			company_id		path		int	true	"company_id"
//	@Param			item_id			path		int	true	"item_id"
//	@Success		200								{object}	transactionsparepartpayloads.PurchaseRequestItemGetAll
//	@Failure		500,400,401,404,403,422			{object}	exceptions.BaseErrorResponse
//	@Router			/v1/purchase-request/item/by-id/{company_id}/{item_id} [get]
func (controller *PurchaseRequestControllerImpl) GetByIdItemTypePr(writer http.ResponseWriter, request *http.Request) {
	ItemId, _ := strconv.Atoi(chi.URLParam(request, "item_id"))
	CompId, _ := strconv.Atoi(chi.URLParam(request, "company_id"))

	result, err := controller.PurchaseRequestService.GetByIdItemTypePurchaseRequest(CompId, ItemId)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, result, "Get Data Successfully", http.StatusOK)
}

// GetByCodeItemTypePr
//
//	@Summary		Get By Id Purchase Request Item Lookup
//	@Description	REST API Purchase Request Item Lookup
//	@Accept			json
//	@Produce		json
//
// @Tags Transaction Sparepart: Purchase Request
//
//	@Param			company_id		path		int	true	"company_id"
//	@Param			item_code		path		string	true	"item_code"
//	@Success		200								{object}	transactionsparepartpayloads.PurchaseRequestItemGetAll
//	@Failure		500,400,401,404,403,422			{object}	exceptions.BaseErrorResponse
//	@Router			/v1/purchase-request/item/by-code/{company_id}/{item_code} [get]
func (controller *PurchaseRequestControllerImpl) GetByCodeItemTypePr(writer http.ResponseWriter, request *http.Request) {

	queryValues := request.URL.Query()

	//myfilter:=utils.FilterCondition{}
	ItemCode := queryValues.Get("item_code")
	CompId := utils.NewGetQueryInt(queryValues, "company_id")
	brandId := utils.NewGetQueryInt(queryValues, "brand_id")
	//filter := utils.BuildFilterCondition(queryParams)

	result, err := controller.PurchaseRequestService.GetByCodeItemTypePurchaseRequest(CompId, ItemCode, brandId)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, result, "Get Data Successfully", http.StatusOK)
}

// VoidDetail
//
// @Summary		Void Request Purchase Detail Multi Id
// @Description	Void Request Purchase Detail Multi Id
// @Accept			json
// @Produce		json
// @Tags Transaction Sparepart: Purchase Request
// @Param			purchase_request_detail_system_number	path		string true	"purchase_request_detail_system_number"
// @Success		201						{object}	payloads.Response
// @Failure		500,400,401,404,403,422	{object}	exceptions.BaseErrorResponse
// @Router			/v1/purchase-request/detail/{purchase_request_detail_system_number} [delete]
func (controller *PurchaseRequestControllerImpl) VoidDetail(writer http.ResponseWriter, request *http.Request) {
	// Void work order
	PurchaseRequesDetailId := chi.URLParam(request, "purchase_request_detail_system_number")

	success, baseErr := controller.PurchaseRequestService.VoidPurchaseRequestDetail(PurchaseRequesDetailId)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			helper.ReturnError(writer, request, baseErr)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}
	if success {
		payloads.NewHandleSuccess(writer, nil, "Purchase Request Detail voided successfully", http.StatusOK)
	} else {
		helper.ReturnError(writer, request, baseErr)
	}
}
