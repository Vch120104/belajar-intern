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

// GetSellingPriceDetailById implements LabourSellingPriceDetailController.
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

// SaveDuplicate implements LabourSellingPriceDetailController.
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

// Duplicate implements LabourSellingPriceDetailController.
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
