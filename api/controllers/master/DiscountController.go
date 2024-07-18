package mastercontroller

import (
	exceptions "after-sales/api/exceptions"
	"fmt"

	// "after-sales/api/helper"
	helper "after-sales/api/helper"
	jsonchecker "after-sales/api/helper/json/json-checker"
	"after-sales/api/payloads"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterservice "after-sales/api/services/master"
	"after-sales/api/utils"
	"after-sales/api/validation"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	// "github.com/julienschmidt/httprouter"
)

type DiscountController interface {
	GetAllDiscount(writer http.ResponseWriter, request *http.Request)
	GetAllDiscountIsActive(writer http.ResponseWriter, request *http.Request)
	GetDiscountByCode(writer http.ResponseWriter, request *http.Request)
	GetDiscountById(writer http.ResponseWriter, request *http.Request)
	SaveDiscount(writer http.ResponseWriter, request *http.Request)
	ChangeStatusDiscount(writer http.ResponseWriter, request *http.Request)
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
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/discount [get]
func (r *DiscountControllerImpl) GetAllDiscount(writer http.ResponseWriter, request *http.Request) {
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

	result, err := r.discountservice.GetAllDiscount(filterCondition, pagination)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data Successfully!", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)
}

// @Summary Get All Discount drop down
// @Description REST API Discount
// @Accept json
// @Produce json
// @Tags Master : Discount
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/discount/drop-down/ [get]
func (r *DiscountControllerImpl) GetAllDiscountIsActive(writer http.ResponseWriter, request *http.Request) {

	result, err := r.discountservice.GetAllDiscountIsActive()
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Get Discount By ID
// @Description REST API Discount
// @Accept json
// @Produce json
// @Tags Master : Discount
// @Param id path string true "id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/discount/by-id/{id} [get]
func (r *DiscountControllerImpl) GetDiscountById(writer http.ResponseWriter, request *http.Request) {
	discountId, errr := strconv.Atoi(chi.URLParam(request, "id"))
	if errr != nil {
		fmt.Print(errr)
		return
	}
	discountResponse, errors := r.discountservice.GetDiscountById(int(discountId))

	if errors != nil {
		helper.ReturnError(writer, request, errors)
		return
	}
	payloads.NewHandleSuccess(writer, discountResponse, utils.GetDataSuccess, http.StatusOK)
}

// @Summary Get Discount By Code
// @Description REST API Discount
// @Accept json
// @Produce json
// @Tags Master : Discount
// @Param discount_code path string true "discount_code"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/discount/by-code/{discount_code} [get]
func (r *DiscountControllerImpl) GetDiscountByCode(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	// discountCode, _ := strconv.Atoi(chi.URLParam(request, "discount_code"))

	discountCode := query.Get("discount_code_value")
	result, err := r.discountservice.GetDiscountByCode(discountCode)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Save Discount
// @Description REST API Discount
// @Accept json
// @Produce json
// @Tags Master : Discount
// @param reqBody body masterpayloads.DiscountResponse true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/discount [post]
func (r *DiscountControllerImpl) SaveDiscount(writer http.ResponseWriter, request *http.Request) {

	var requestForm masterpayloads.DiscountResponse
	var message string

	err := jsonchecker.ReadFromRequestBody(request, &requestForm)
	if err != nil {
		exceptions.NewEntityException(writer, request, err)
		return
	}
	err = validation.ValidationForm(writer, request, requestForm)
	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	create, err := r.discountservice.SaveDiscount(requestForm)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

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
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/discount/{discount_code_id} [patch]
func (r *DiscountControllerImpl) ChangeStatusDiscount(writer http.ResponseWriter, request *http.Request) {
	discountId, _ := strconv.Atoi(chi.URLParam(request, "id"))

	response, err := r.discountservice.ChangeStatusDiscount(int(discountId))

	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}
