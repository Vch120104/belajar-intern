package masteroperationcontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	jsonchecker "after-sales/api/helper/json/json-checker"
	"after-sales/api/payloads"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	masteroperationservice "after-sales/api/services/master/operation"
	"after-sales/api/utils"
	"after-sales/api/validation"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

type LabourSellingPriceDetailController interface {
	GetAllSellingPriceDetailByHeaderId(writer http.ResponseWriter, request *http.Request)
	GetSellingPriceDetailById(writer http.ResponseWriter, request *http.Request)
	SaveLabourSellingPriceDetail(writer http.ResponseWriter, request *http.Request)
	Duplicate(writer http.ResponseWriter, request *http.Request)
	SaveDuplicate(writer http.ResponseWriter, request *http.Request)
	DeleteLabourSellingPriceDetail(writer http.ResponseWriter, request *http.Request)
}
type LabourSellingPriceDetailControllerImpl struct {
	LabourSellingPriceService masteroperationservice.LabourSellingPriceService
}

func NewLabourSellingPriceDetailController(LabourSellingPriceService masteroperationservice.LabourSellingPriceService) LabourSellingPriceDetailController {
	return &LabourSellingPriceDetailControllerImpl{
		LabourSellingPriceService: LabourSellingPriceService,
	}
}

// @Summary Get Selling Price Detail By ID
// @Description Get Selling Price Detail By ID
// @Tags Master : Labour Selling Price
// @Accept json
// @Produce json
// @Param labour_selling_price_detail_id path int true "Labour Selling Price Detail ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/labour-selling-price-detail/{labour_selling_price_detail_id} [get]
func (r *LabourSellingPriceDetailControllerImpl) GetSellingPriceDetailById(writer http.ResponseWriter, request *http.Request) {
	detailId, errA := strconv.Atoi(chi.URLParam(request, "labour_selling_price_detail_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	result, err := r.LabourSellingPriceService.GetSellingPriceDetailById(detailId)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "success", 200)
}

// @Summary Save Duplicate
// @Description Save Duplicate
// @Tags Master : Labour Selling Price
// @Accept json
// @Produce json
// @Param req body masteroperationpayloads.SaveDuplicateLabourSellingPrice true "Save Duplicate Labour Selling Price"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/labour-selling-price/save-duplicate [post]
func (r *LabourSellingPriceDetailControllerImpl) SaveDuplicate(writer http.ResponseWriter, request *http.Request) {
	var formRequest masteroperationpayloads.SaveDuplicateLabourSellingPrice

	err := jsonchecker.ReadFromRequestBody(request, &formRequest)

	if err != nil {
		exceptions.NewEntityException(writer, request, err)
		return
	}

	err = validation.ValidationForm(writer, request, formRequest)
	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	create, err := r.LabourSellingPriceService.SaveDuplicate(formRequest)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, create, "Save Duplicate", http.StatusOK)
}

// @Summary Duplicate
// @Description Duplicate
// @Tags Master : Labour Selling Price
// @Accept json
// @Produce json
// @Param labour_selling_price_id path int true "Labour Selling Price ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/labour-selling-price/duplicate/{labour_selling_price_id} [get]
func (r *LabourSellingPriceDetailControllerImpl) Duplicate(writer http.ResponseWriter, request *http.Request) {
	sellingPriceId, errA := strconv.Atoi(chi.URLParam(request, "labour_selling_price_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	result, err := r.LabourSellingPriceService.Duplicate(sellingPriceId)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, utils.ModifyKeysInResponse(result), "success", 200)
}

// @Summary Get All Selling Price Detail By Header ID
// @Description Get All Selling Price Detail By Header ID
// @Tags Master : Labour Selling Price
// @Accept json
// @Produce json
// @Param labour_selling_price_id path int true "Labour Selling Price ID"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Param sort_of query string false "Sort Of"
// @Param sort_by query string false "Sort By"
// @Success 200 {object} payloads.ResponsePagination
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/labour-selling-price-detail/{labour_selling_price_id} [get]
func (r *LabourSellingPriceDetailControllerImpl) GetAllSellingPriceDetailByHeaderId(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	sellingPriceId, errA := strconv.Atoi(chi.URLParam(request, "labour_selling_price_id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	pagination := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	result, err := r.LabourSellingPriceService.GetAllSellingPriceDetailByHeaderId(sellingPriceId, pagination)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(
		writer,
		result.Rows,
		"Get Data Successfully!",
		http.StatusOK,
		result.Limit,
		result.Page,
		int64(result.TotalRows),
		result.TotalPages,
	)
}

// @Summary Save Labour Selling Price Detail
// @Description Save Labour Selling Price Detail
// @Tags Master : Labour Selling Price
// @Accept json
// @Produce json
// @Param req body masteroperationpayloads.LabourSellingPriceDetailRequest true "Save Labour Selling Price Detail"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/labour-selling-price-detail [post]
func (r *LabourSellingPriceDetailControllerImpl) SaveLabourSellingPriceDetail(writer http.ResponseWriter, request *http.Request) {

	var formRequest masteroperationpayloads.LabourSellingPriceDetailRequest
	err := jsonchecker.ReadFromRequestBody(request, &formRequest)
	var message string

	if err != nil {
		exceptions.NewEntityException(writer, request, err)
		return
	}
	err = validation.ValidationForm(writer, request, formRequest)
	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	create, err := r.LabourSellingPriceService.SaveLabourSellingPriceDetail(formRequest)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	message = "Create Data Successfully!"

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

// @Summary Delete Labour Selling Price Detail
// @Description Delete Labour Selling Price Detail
// @Tags Master : Labour Selling Price
// @Accept json
// @Produce json
// @Param multi_id path string true "Multi ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/labour-selling-price-detail/{multi_id} [delete]
func (r *LabourSellingPriceDetailControllerImpl) DeleteLabourSellingPriceDetail(writer http.ResponseWriter, request *http.Request) {

	multiId := chi.URLParam(request, "multi_id")
	if multiId == "[]" {
		payloads.NewHandleError(writer, "Invalid service request detail multi ID", http.StatusBadRequest)
		return
	}

	multiId = strings.Trim(multiId, "[]")
	elements := strings.Split(multiId, ",")

	var intIds []int
	for _, element := range elements {
		num, err := strconv.Atoi(strings.TrimSpace(element))
		if err != nil {
			payloads.NewHandleError(writer, "Error converting data to integer", http.StatusBadRequest)
			return
		}
		intIds = append(intIds, num)
	}

	success, baseErr := r.LabourSellingPriceService.DeleteLabourSellingPriceDetail(intIds)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Labour selling price detail not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	if success {
		payloads.NewHandleSuccess(writer, success, "Labour selling price deleted successfully", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Failed to delete Labour selling price", http.StatusInternalServerError)
	}

}
