package masteroperationcontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	jsonchecker "after-sales/api/helper/json/json-checker"
	"after-sales/api/payloads"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	"errors"

	// "after-sales/api/payloads/pagination"
	masteroperationservice "after-sales/api/services/master/operation"
	"after-sales/api/utils"
	"after-sales/api/validation"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type LabourSellingPriceController interface {
	GetLabourSellingPriceById(writer http.ResponseWriter, request *http.Request)
	GetAllSellingPrice(writer http.ResponseWriter, request *http.Request)
	SaveLabourSellingPrice(writer http.ResponseWriter, request *http.Request)
}
type LabourSellingPriceControllerImpl struct {
	LabourSellingPriceService masteroperationservice.LabourSellingPriceService
}

func NewLabourSellingPriceController(LabourSellingPriceService masteroperationservice.LabourSellingPriceService) LabourSellingPriceController {
	return &LabourSellingPriceControllerImpl{
		LabourSellingPriceService: LabourSellingPriceService,
	}
}

func (r *LabourSellingPriceControllerImpl) GetAllSellingPrice(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"mtr_labour_selling_price.company_id":     queryValues.Get("company_id"),
		"mtr_labour_selling_price.effective_date": queryValues.Get("effective_date"),
		"mtr_labour_selling_price.bill_to_id":     queryValues.Get("bill_to_id"),
		"mtr_labour_selling_price.job_type_id":    queryValues.Get("job_type_id"),
		"mtr_labour_selling_price.description":    queryValues.Get("description"),
		"mtr_labour_selling_price.brand_id":       queryValues.Get("brand_id"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	result, err := r.LabourSellingPriceService.GetAllSellingPrice(criteria, paginate)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data Successfully!", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)

}

func (r *LabourSellingPriceControllerImpl) GetLabourSellingPriceById(writer http.ResponseWriter, request *http.Request) {

	labourSellingPriceId, errA := strconv.Atoi(chi.URLParam(request, "labour_selling_price_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	result, err := r.LabourSellingPriceService.GetLabourSellingPriceById(labourSellingPriceId)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, utils.ModifyKeysInResponse(result), "Get Data Successfully!", http.StatusOK)
}

func (r *LabourSellingPriceControllerImpl) SaveLabourSellingPrice(writer http.ResponseWriter, request *http.Request) {

	var formRequest masteroperationpayloads.LabourSellingPriceRequest
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

	create, err := r.LabourSellingPriceService.SaveLabourSellingPrice(formRequest)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	message = "Create Data Successfully!"

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}
