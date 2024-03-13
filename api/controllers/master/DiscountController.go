package mastercontroller

import (
	"after-sales/api/helper"
	"after-sales/api/payloads"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterservice "after-sales/api/services/master"
	"after-sales/api/utils"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type DiscountController interface {
	GetAllDiscount(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	GetAllDiscountIsActive(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	GetDiscountByCode(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	SaveDiscount(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	ChangeStatusDiscount(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}

type DiscountControllerImpl struct {
	discountservice masterservice.DiscountService
}

func NewDiscountController(discountService masterservice.DiscountService) DiscountController {
	return &DiscountControllerImpl{
		discountservice: discountService,
	}
}

// @Summary Get All Discount
// @Description REST API Discount
// @Accept json
// @Produce json
// @Tags Master : Discount
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Param is_active query string false "is_active" Enums(true, false)
// @Param discount_code_value query string false "discount_code_value"
// @Param discount_code_description query string false "discount_code_description"
// @Param sort_by query string false "sort_by"
// @Param sort_of query string false "sort_of"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/discount [get]
func (r *DiscountControllerImpl) GetAllDiscount(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	query := request.URL.Query()
	queryParams := map[string]string{
		"is_active":                 query.Get("is_active"),
		"discount_code_value":       query.Get("discount_code_value"),
		"discount_code_description": query.Get("discount_code_description"),
	}

	pagination := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(query, "limit"),
		Page:   utils.NewGetQueryInt(query, "page"),
		SortOf: query.Get("sort_of"),
		SortBy: query.Get("sort_by"),
	}

	filterCondition := utils.BuildFilterCondition(queryParams)

	result := r.discountservice.GetAllDiscount(filterCondition, pagination)

	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data Successfully!", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)
}

// @Summary Get All Discount drop down
// @Description REST API Discount
// @Accept json
// @Produce json
// @Tags Master : Discount
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/discount-drop-down/ [get]
func (r *DiscountControllerImpl) GetAllDiscountIsActive(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	result := r.discountservice.GetAllDiscountIsActive()

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Get Discount By Code
// @Description REST API Discount
// @Accept json
// @Produce json
// @Tags Master : Discount
// @Param discount_code path string true "discount_code"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/discount-by-code/{discount_code} [get]
func (r *DiscountControllerImpl) GetDiscountByCode(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	query := request.URL.Query()

	operationGroupCode := query.Get("discount_code")
	result := r.discountservice.GetDiscountByCode(operationGroupCode)

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Save Discount
// @Description REST API Discount
// @Accept json
// @Produce json
// @Tags Master : Discount
// @param reqBody body masterpayloads.DiscountResponse true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/discount [post]
func (r *DiscountControllerImpl) SaveDiscount(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	var requestForm masterpayloads.DiscountResponse
	var message = ""

	helper.ReadFromRequestBody(request, &requestForm)

	create := r.discountservice.SaveDiscount(requestForm)

	if requestForm.DiscountCodeId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

// @Summary Change Status Discount
// @Description REST API Discount
// @Accept json
// @Produce json
// @Tags Master : Discount
// @param discount_code_id path int true "discount_code_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/discount/{discount_code_id} [patch]
func (r *DiscountControllerImpl) ChangeStatusDiscount(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	uomId, _ := strconv.Atoi(params.ByName("discount_code_id"))

	response := r.discountservice.ChangeStatusDiscount(int(uomId))

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}
