package masteritemcontroller

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	jsonchecker "after-sales/api/helper/json/json-checker"
	"after-sales/api/payloads"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemservice "after-sales/api/services/master/item"
	"after-sales/api/utils"
	"after-sales/api/validation"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type DiscountPercentController interface {
	GetAllDiscountPercent(writer http.ResponseWriter, request *http.Request)
	GetDiscountPercentByID(writer http.ResponseWriter, request *http.Request)
	SaveDiscountPercent(writer http.ResponseWriter, request *http.Request)
	ChangeStatusDiscountPercent(writer http.ResponseWriter, request *http.Request)
}
type DiscountPercentControllerImpl struct {
	DiscountPercentService masteritemservice.DiscountPercentService
}

func NewDiscountPercentController(discountPercentService masteritemservice.DiscountPercentService) DiscountPercentController {
	return &DiscountPercentControllerImpl{
		DiscountPercentService: discountPercentService,
	}
}

// @Summary Get All Discount Percent
// @Description REST API Discount Percent
// @Accept json
// @Produce json
// @Tags Master : Discount Percent
// @Security AuthorizationKeyAuth
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Param is_active query string false "is_active" Enums(true, false)
// @Param discount_code_value query string false "discount_code_value"
// @Param discount_code_description query string false "discount_code_description"
// @Param order_type_name query string false "order_type_name"
// @Param discount query float64 false "discount"
// @Param sort_by query string false "sort_by"
// @Param sort_of query string false "sort_of"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/discount-percent/ [get]
func (r *DiscountPercentControllerImpl) GetAllDiscountPercent(writer http.ResponseWriter, request *http.Request) {

	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"mtr_discount.discount_code_id":      queryValues.Get("discount_code_id"),
		"mtr_discount.discount_code":         queryValues.Get("discount_code"),
		"mtr_discount.discount_description":  queryValues.Get("discount_description"),
		"mtr_discount_percent.order_type_id": queryValues.Get("order_type_id"),
		"order_type_name":                    queryValues.Get("order_type_name"),
		"mtr_discount_percent.discount":      queryValues.Get("discount"),
		"mtr_discount_percent.is_active":     queryValues.Get("is_active"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	paginatedData, err := r.DiscountPercentService.GetAllDiscountPercent(criteria, paginate)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(writer, paginatedData.Rows, "Get Data Successfully!", http.StatusOK, paginate.Limit, paginate.Page, paginatedData.TotalRows, paginatedData.TotalPages)
}

// @Summary Get Discount Percent By ID
// @Description REST API Discount Percent
// @Accept json
// @Produce json
// @Tags Master : Discount Percent
// @Security AuthorizationKeyAuth
// @Param discount_percent_id path int true "discount_percent_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/discount-percent/{discount_percent_id} [get]
func (r *DiscountPercentControllerImpl) GetDiscountPercentByID(writer http.ResponseWriter, request *http.Request) {

	discountPercentId, errA := strconv.Atoi(chi.URLParam(request, "discount_percent_id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	result, err := r.DiscountPercentService.GetDiscountPercentById(discountPercentId)

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Save Discount Percent
// @Description REST API Discount Percent
// @Accept json
// @Produce json
// @Tags Master : Discount Percent
// @Security AuthorizationKeyAuth
// @param reqBody body masteritempayloads.DiscountPercentResponse true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/discount-percent [post]
func (r *DiscountPercentControllerImpl) SaveDiscountPercent(writer http.ResponseWriter, request *http.Request) {

	var formRequest masteritempayloads.DiscountPercentResponse
	err := jsonchecker.ReadFromRequestBody(request, &formRequest)
	var message = ""

	if err != nil {
		exceptions.NewEntityException(writer, request, err)
		return
	}
	err = validation.ValidationForm(writer, request, formRequest)
	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	create, err := r.DiscountPercentService.SaveDiscountPercent(formRequest)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	if formRequest.DiscountPercentId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

// @Summary Change Status Discount Percent
// @Description REST API Discount Percent
// @Accept json
// @Produce json
// @Tags Master : Discount Percent
// @Security AuthorizationKeyAuth
// @param discount_percent_id path int true "discount_percent_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/discount-percent/{discount_percent_id} [patch]
func (r *DiscountPercentControllerImpl) ChangeStatusDiscountPercent(writer http.ResponseWriter, request *http.Request) {

	discountPercentId, errA := strconv.Atoi(chi.URLParam(request, "discount_percent_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	response, err := r.DiscountPercentService.ChangeStatusDiscountPercent(int(discountPercentId))

	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}
