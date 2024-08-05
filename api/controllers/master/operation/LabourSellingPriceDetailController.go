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
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type LabourSellingPriceDetailController interface {
	GetAllSellingPriceDetailByHeaderId(writer http.ResponseWriter, request *http.Request)
	SaveLabourSellingPriceDetail(writer http.ResponseWriter, request *http.Request)
}
type LabourSellingPriceDetailControllerImpl struct {
	LabourSellingPriceService masteroperationservice.LabourSellingPriceService
}

func NewLabourSellingPriceDetailController(LabourSellingPriceService masteroperationservice.LabourSellingPriceService) LabourSellingPriceDetailController {
	return &LabourSellingPriceDetailControllerImpl{
		LabourSellingPriceService: LabourSellingPriceService,
	}
}

func (r *LabourSellingPriceDetailControllerImpl) GetAllSellingPriceDetailByHeaderId(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	sellingPriceId, _ := strconv.Atoi(chi.URLParam(request, "labour_selling_price_id"))

	pagination := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	paginatedData, totalPages, totalRows, err := r.LabourSellingPriceService.GetAllSellingPriceDetailByHeaderId(sellingPriceId, pagination)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(paginatedData), "success", 200, pagination.Limit, pagination.Page, int64(totalRows), totalPages)
}

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
