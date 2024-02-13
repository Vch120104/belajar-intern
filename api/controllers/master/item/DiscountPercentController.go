package masteritemcontroller

import (
	"after-sales/api/helper"
	"after-sales/api/payloads"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemservice "after-sales/api/services/master/item"
	"after-sales/api/utils"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type DiscountPercentController interface {
	GetAllDiscountPercent(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	GetDiscountPercentByID(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	SaveDiscountPercent(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	ChangeStatusDiscountPercent(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
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
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/discount-percent [get]
func (r *DiscountPercentControllerImpl) GetAllDiscountPercent(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"mtr_discount.discount_code_value":       queryValues.Get("discount_code_value"),
		"mtr_discount.discount_code_description": queryValues.Get("discount_code_description"),
		"order_type_name":                        queryValues.Get("order_type_name"),
		"mtr_discount_percent.discount":          queryValues.Get("discount"),
		"mtr_discount_percent.is_active":         queryValues.Get("is_active"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	paginatedData, totalPages, totalRows := r.DiscountPercentService.GetAllDiscountPercent(criteria, paginate)

	payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(paginatedData), "success", 200, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
}

// @Summary Get Discount Percent By ID
// @Description REST API Discount Percent
// @Accept json
// @Produce json
// @Tags Master : Discount Percent
// @Param discount_percent_id path int true "discount_percent_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/discount-percent/{discount_percent_id} [get]
func (r *DiscountPercentControllerImpl) GetDiscountPercentByID(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	discountPercentId, _ := strconv.Atoi(params.ByName("discount_percent_id"))

	result := r.DiscountPercentService.GetDiscountPercentById(discountPercentId)

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Save Discount Percent
// @Description REST API Discount Percent
// @Accept json
// @Produce json
// @Tags Master : Discount Percent
// @param reqBody body masteritempayloads.DiscountPercentResponse true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/discount-percent [post]
func (r *DiscountPercentControllerImpl) SaveDiscountPercent(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	var formRequest masteritempayloads.DiscountPercentResponse
	helper.ReadFromRequestBody(request, &formRequest)
	var message = ""

	create := r.DiscountPercentService.SaveDiscountPercent(formRequest)

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
// @param discount_percent_id path int true "discount_percent_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/discount-percent/{discount_percent_id} [patch]
func (r *DiscountPercentControllerImpl) ChangeStatusDiscountPercent(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	discountPercentId, _ := strconv.Atoi(params.ByName("discount_percent_id"))

	response := r.DiscountPercentService.ChangeStatusDiscountPercent(int(discountPercentId))

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}
