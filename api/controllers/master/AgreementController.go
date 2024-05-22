package mastercontroller

import (

	// "after-sales/api/middlewares"

	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterservice "after-sales/api/services/master"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type AgreementController interface {
	GetAgreementById(writer http.ResponseWriter, request *http.Request)
	SaveAgreement(writer http.ResponseWriter, request *http.Request)
	ChangeStatusAgreement(writer http.ResponseWriter, request *http.Request)
	GetAllAgreement(writer http.ResponseWriter, request *http.Request)

	GetAllDiscountGroup(writer http.ResponseWriter, request *http.Request)
	GetDiscountGroupAgreementById(writer http.ResponseWriter, request *http.Request)
	AddDiscountGroup(writer http.ResponseWriter, request *http.Request)
	DeleteDiscountGroup(writer http.ResponseWriter, request *http.Request)

	GetAllItemDiscount(writer http.ResponseWriter, request *http.Request)
	GetDiscountItemAgreementById(writer http.ResponseWriter, request *http.Request)
	AddItemDiscount(writer http.ResponseWriter, request *http.Request)
	DeleteItemDiscount(writer http.ResponseWriter, request *http.Request)

	GetAllDiscountValue(writer http.ResponseWriter, request *http.Request)
	GetDiscountValueAgreementById(writer http.ResponseWriter, request *http.Request)
	AddDiscountValue(writer http.ResponseWriter, request *http.Request)
	DeleteDiscountValue(writer http.ResponseWriter, request *http.Request)
}

type AgreementControllerImpl struct {
	AgreementService masterservice.AgreementService
}

func NewAgreementController(AgreementService masterservice.AgreementService) AgreementController {
	return &AgreementControllerImpl{
		AgreementService: AgreementService,
	}
}

// @Summary Get Agreement By Id
// @Description Retrieve an agreement by its ID
// @Accept json
// @Produce json
// @Tags Master : Agreement
// @Param agreement_id path int true "Agreement ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/agreement/{agreement_id} [get]
func (r *AgreementControllerImpl) GetAgreementById(writer http.ResponseWriter, request *http.Request) {

	AgreementId, _ := strconv.Atoi(chi.URLParam(request, "agreement_id"))

	result, err := r.AgreementService.GetAgreementById(int(AgreementId))
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Save Agreement
// @Description Create or update an agreement
// @Accept json
// @Produce json
// @Tags Master : Agreement
// @Param reqBody body masterpayloads.AgreementRequest true "Agreement Data"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/agreement/ [post]
func (r *AgreementControllerImpl) SaveAgreement(writer http.ResponseWriter, request *http.Request) {

	var formRequest masterpayloads.AgreementRequest
	helper.ReadFromRequestBody(request, &formRequest)
	var message = ""

	create, err := r.AgreementService.SaveAgreement(formRequest)
	if err != nil {
		exceptions.NewConflictException(writer, request, err)
		return
	}

	if formRequest.AgreementId == 0 {
		message = "Create Data Successfully!"
		payloads.NewHandleSuccess(writer, create, message, http.StatusCreated)
	} else {
		message = "Update Data Successfully!"
		payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
	}
}

// @Summary Change Status Agreement
// @Description Change the status of an agreement
// @Accept json
// @Produce json
// @Tags Master : Agreement
// @Param agreement_id path int true "Agreement ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/agreement/{agreement_id} [patch]
func (r *AgreementControllerImpl) ChangeStatusAgreement(writer http.ResponseWriter, request *http.Request) {

	agreement_id, _ := strconv.Atoi(chi.URLParam(request, "agreement_id"))

	response, err := r.AgreementService.ChangeStatusAgreement(int(agreement_id))
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}

// @Summary Get All Agreements
// @Description Retrieve all agreements with optional filtering and pagination
// @Accept json
// @Produce json
// @Tags Master : Agreement
// @Param page query string true "Page number"
// @Param limit query string true "Items per page"
// @Param is_active query string false "Agreement status"
// @Param supplier_name query string false "Supplier name"
// @Param agreement_code query string false "Agreement code"
// @Param brand_id query string false "Brand ID"
// @Param customer_id query string false "Customer ID"
// @Param profit_center_id query string false "Profit center ID"
// @Param dealer_id query string false "Dealer ID"
// @Param top_id query string false "Top ID"
// @Param sort_by query string false "Field to sort by"
// @Param sort_of query string false "Sort order (asc/desc)"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/agreement [get]
func (r *AgreementControllerImpl) GetAllAgreement(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query() // Retrieve query parameters

	queryParams := map[string]string{
		"mtr_agreement.agreement_id": queryValues.Get("agreement_id"),
		"mtr_agreement.brand_id":     queryValues.Get("brand_id"),
		"mtr_agreement.customer_id":  queryValues.Get("customer_id"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: chi.URLParam(request, "sort_of"),
		SortBy: chi.URLParam(request, "sort_by"),
	}
	print(queryParams)

	criteria := utils.BuildFilterCondition(queryParams)
	paginatedData, totalPages, totalRows, err := r.AgreementService.GetAllAgreement(criteria, paginate)

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	if len(paginatedData) > 0 {
		payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(paginatedData), "Get Data Successfully", http.StatusOK, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}
}

// @Summary Add Discount Group
// @Description Add a new discount group to an agreement by its ID
// @Accept json
// @Produce json
// @Tags Master : Agreement
// @Param agreement_id path int true "Agreement ID"
// @Param reqBody body masterpayloads.DiscountGroupRequest true "Discount Group Data"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/agreement/{agreement_id}/discount/group [post]
func (r *AgreementControllerImpl) AddDiscountGroup(writer http.ResponseWriter, request *http.Request) {
	agreementID, _ := strconv.Atoi(chi.URLParam(request, "agreement_id"))

	var groupRequest masterpayloads.DiscountGroupRequest
	helper.ReadFromRequestBody(request, &groupRequest)

	if err := r.AgreementService.AddDiscountGroup(int(agreementID), groupRequest); err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, nil, "Discount group added successfully", http.StatusCreated)
}

// @Summary Delete Discount Group
// @Description Delete a discount group from an agreement by its ID
// @Accept json
// @Produce json
// @Tags Master : Agreement
// @Param agreement_id path int true "Agreement ID"
// @Param agreement_discount_group_id path int true "Group ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/agreement/{agreement_id}/discount/group/{agreement_discount_group_id} [delete]
func (r *AgreementControllerImpl) DeleteDiscountGroup(writer http.ResponseWriter, request *http.Request) {
	agreementID, _ := strconv.Atoi(chi.URLParam(request, "agreement_id"))
	groupID, _ := strconv.Atoi(chi.URLParam(request, "agreement_discount_group_id"))

	if err := r.AgreementService.DeleteDiscountGroup(int(agreementID), int(groupID)); err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, nil, "Discount group deleted successfully", http.StatusOK)
}

// @Summary Add Item Discount
// @Description Add a new item discount to an agreement by its ID
// @Accept json
// @Produce json
// @Tags Master : Agreement
// @Param agreement_id path int true "Agreement ID"
// @Param reqBody body masterpayloads.ItemDiscountRequest true "Item Discount Data"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/agreement/{agreement_id}/discount/item [post]
func (r *AgreementControllerImpl) AddItemDiscount(writer http.ResponseWriter, request *http.Request) {
	agreementID, _ := strconv.Atoi(chi.URLParam(request, "agreement_id"))

	var itemRequest masterpayloads.ItemDiscountRequest
	helper.ReadFromRequestBody(request, &itemRequest)

	if err := r.AgreementService.AddItemDiscount(int(agreementID), itemRequest); err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, nil, "Item discount added successfully", http.StatusCreated)
}

// @Summary Delete Item Discount
// @Description Delete an item discount from an agreement by its ID
// @Accept json
// @Produce json
// @Tags Master : Agreement
// @Param agreement_id path int true "Agreement ID"
// @Param agreement_item_id path int true "Item ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/agreement/{agreement_id}/discount/item/{agreement_item_id} [delete]
func (r *AgreementControllerImpl) DeleteItemDiscount(writer http.ResponseWriter, request *http.Request) {
	agreementID, _ := strconv.Atoi(chi.URLParam(request, "agreement_id"))
	itemID, _ := strconv.Atoi(chi.URLParam(request, "agreement_item_id"))

	if err := r.AgreementService.DeleteItemDiscount(int(agreementID), int(itemID)); err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, nil, "Item discount deleted successfully", http.StatusOK)
}

// @Summary Add Discount Value
// @Description Add a new discount value to an agreement by its ID
// @Accept json
// @Produce json
// @Tags Master : Agreement
// @Param agreement_id path int true "Agreement ID"
// @Param reqBody body masterpayloads.DiscountValueRequest true "Discount Value Data"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/agreement/{agreement_id}/discount/value [post]
func (r *AgreementControllerImpl) AddDiscountValue(writer http.ResponseWriter, request *http.Request) {
	agreementID, _ := strconv.Atoi(chi.URLParam(request, "agreement_id"))

	var valueRequest masterpayloads.DiscountValueRequest
	helper.ReadFromRequestBody(request, &valueRequest)

	if err := r.AgreementService.AddDiscountValue(int(agreementID), valueRequest); err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, nil, "Discount value added successfully", http.StatusCreated)
}

// @Summary Delete Discount Value
// @Description Delete a discount value from an agreement by its ID
// @Accept json
// @Produce json
// @Tags Master : Agreement
// @Param agreement_id path int true "Agreement ID"
// @Param agreement_discount_id path int true "Value ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/agreement/{agreement_id}/discount/value/{agreement_discount_id} [delete]
func (r *AgreementControllerImpl) DeleteDiscountValue(writer http.ResponseWriter, request *http.Request) {
	agreementID, _ := strconv.Atoi(chi.URLParam(request, "agreement_id"))
	valueID, _ := strconv.Atoi(chi.URLParam(request, "agreement_discount_id"))

	if err := r.AgreementService.DeleteDiscountValue(int(agreementID), int(valueID)); err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, nil, "Discount value deleted successfully", http.StatusOK)
}

// @Summary Get All Discount Group
// @Description Retrieve all discount groups from an agreement by its ID
// @Accept json
// @Produce json
// @Tags Master : Agreement
// @Param agreement_id path int true "Agreement ID"
// @Param page query string true "Page number"
// @Param limit query string true "Items per page"
// @Param sort_by query string false "Field to sort by"
// @Param sort_of query string false "Sort order (asc/desc)"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/agreement/{agreement_id}/discount/group [get]
func (r *AgreementControllerImpl) GetAllDiscountGroup(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query() // Retrieve query parameters

	queryParams := map[string]string{
		"agreement_id": queryValues.Get("agreement_id"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: chi.URLParam(request, "sort_of"),
		SortBy: chi.URLParam(request, "sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)
	paginatedData, totalPages, totalRows, err := r.AgreementService.GetAllDiscountGroup(criteria, paginate)

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	if len(paginatedData) > 0 {
		payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(paginatedData), "Get Data Successfully", http.StatusOK, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}
}

// @Summary Get Discount Group By Id
// @Description Retrieve a discount group from an agreement by its ID
// @Accept json
// @Produce json
// @Tags Master : Agreement
// @Param agreement_id path int true "Agreement ID"
// @Param agreement_discount_group_id path int true "Group ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/agreement/{agreement_id}/discount/group/{agreement_discount_group_id} [get]
func (r *AgreementControllerImpl) GetDiscountGroupAgreementById(writer http.ResponseWriter, request *http.Request) {
	agreementID, _ := strconv.Atoi(chi.URLParam(request, "agreement_id"))
	groupID, _ := strconv.Atoi(chi.URLParam(request, "agreement_discount_group_id"))

	result, err := r.AgreementService.GetDiscountGroupAgreementById(int(agreementID), int(groupID))
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Get All Discount Item
// @Description Retrieve all discount items from an agreement by its ID
// @Accept json
// @Produce json
// @Tags Master : Agreement
// @Param agreement_id path int true "Agreement ID"
// @Param page query string true "Page number"
// @Param limit query string true "Items per page"
// @Param sort_by query string false "Field to sort by"
// @Param sort_of query string false "Sort order (asc/desc)"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/agreement/{agreement_id}/discount/item [get]
func (r *AgreementControllerImpl) GetAllItemDiscount(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query() // Retrieve query parameters

	queryParams := map[string]string{
		"agreement_id": queryValues.Get("agreement_id"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: chi.URLParam(request, "sort_of"),
		SortBy: chi.URLParam(request, "sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)
	paginatedData, totalPages, totalRows, err := r.AgreementService.GetAllItemDiscount(criteria, paginate)

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	if len(paginatedData) > 0 {
		payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(paginatedData), "Get Data Successfully", http.StatusOK, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}
}

// @Summary Get Discount Item By Id
// @Description Retrieve a discount item from an agreement by its ID
// @Accept json
// @Produce json
// @Tags Master : Agreement
// @Param agreement_id path int true "Agreement ID"
// @Param agreement_item_id path int true "Item ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/agreement/{agreement_id}/discount/item/{agreement_item_id} [get]
func (r *AgreementControllerImpl) GetDiscountItemAgreementById(writer http.ResponseWriter, request *http.Request) {
	agreementID, _ := strconv.Atoi(chi.URLParam(request, "agreement_id"))
	itemID, _ := strconv.Atoi(chi.URLParam(request, "agreement_item_id"))

	result, err := r.AgreementService.GetDiscountItemAgreementById(int(agreementID), int(itemID))
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Get All Discount Value
// @Description Retrieve all discount values from an agreement by its ID
// @Accept json
// @Produce json
// @Tags Master : Agreement
// @Param agreement_id path int true "Agreement ID"
// @Param page query string true "Page number"
// @Param limit query string true "Items per page"
// @Param sort_by query string false "Field to sort by"
// @Param sort_of query string false "Sort order (asc/desc)"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/agreement/{agreement_id}/discount/value [get]
func (r *AgreementControllerImpl) GetAllDiscountValue(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query() // Retrieve query parameters

	queryParams := map[string]string{
		"agreement_id": queryValues.Get("agreement_id"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: chi.URLParam(request, "sort_of"),
		SortBy: chi.URLParam(request, "sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)
	paginatedData, totalPages, totalRows, err := r.AgreementService.GetAllDiscountValue(criteria, paginate)

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	if len(paginatedData) > 0 {
		payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(paginatedData), "Get Data Successfully", http.StatusOK, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}
}

// @Summary Get Discount Value By Id
// @Description Retrieve a discount value from an agreement by its ID
// @Accept json
// @Produce json
// @Tags Master : Agreement
// @Param agreement_id path int true "Agreement ID"
// @Param agreement_discount_id path int true "Value ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/agreement/{agreement_id}/discount/value/{agreement_discount_id} [get]
func (r *AgreementControllerImpl) GetDiscountValueAgreementById(writer http.ResponseWriter, request *http.Request) {
	agreementID, _ := strconv.Atoi(chi.URLParam(request, "agreement_id"))
	valueID, _ := strconv.Atoi(chi.URLParam(request, "agreement_discount_id"))

	result, err := r.AgreementService.GetDiscountValueAgreementById(int(agreementID), int(valueID))
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}
