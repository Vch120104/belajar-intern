package mastercontroller

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
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
)

type MovingCodeController interface {
	GetAllMovingCode(writer http.ResponseWriter, request *http.Request)
	PushMovingCodePriority(writer http.ResponseWriter, request *http.Request)
	CreateMovingCode(writer http.ResponseWriter, request *http.Request)
	UpdateMovingCode(writer http.ResponseWriter, request *http.Request)
	GetMovingCodebyId(writer http.ResponseWriter, request *http.Request)
	ChangeStatusMovingCode(writer http.ResponseWriter, request *http.Request)
	GetDropdownMovingCode(writer http.ResponseWriter, request *http.Request)
	ActivateMovingCode(writer http.ResponseWriter, request *http.Request)
	DeactiveMovingCode(writer http.ResponseWriter, request *http.Request)
}

type MovingCodeControllerImpl struct {
	MovingCodeService masterservice.MovingCodeService
}

// ActivateMovingCode implements MovingCodeController.
func (r *MovingCodeControllerImpl) ActivateMovingCode(writer http.ResponseWriter, request *http.Request) {
	id := chi.URLParam(request, "moving_code_id")

	response, err := r.MovingCodeService.ActivateMovingCode(id)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Activate Status Successfully!", http.StatusOK)
}

// DeactiveMovingCode implements MovingCodeController.
func (r *MovingCodeControllerImpl) DeactiveMovingCode(writer http.ResponseWriter, request *http.Request) {
	id := chi.URLParam(request, "moving_code_id")

	response, err := r.MovingCodeService.DeactiveMovingCode(id)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Deactive Status Successfully!", http.StatusOK)
}

// GetDropdownMovingCode implements MovingCodeController.
func (r *MovingCodeControllerImpl) GetDropdownMovingCode(writer http.ResponseWriter, request *http.Request) {
	companyId, _ := strconv.Atoi(chi.URLParam(request, "company_id"))

	result, err := r.MovingCodeService.GetDropdownMovingCode(companyId)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// ChangeStatusMovingCode implements MovingCodeController.
func (r *MovingCodeControllerImpl) ChangeStatusMovingCode(writer http.ResponseWriter, request *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(request, "moving_code_id"))

	response, err := r.MovingCodeService.ChangeStatusMovingCode(id)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Change Status Successfully!", http.StatusOK)
}

// CreateMovingCode implements MovingCodeController.
func (r *MovingCodeControllerImpl) CreateMovingCode(writer http.ResponseWriter, request *http.Request) {
	var formRequest masterpayloads.MovingCodeListRequest
	err := jsonchecker.ReadFromRequestBody(request, &formRequest)

	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	err = validation.ValidationForm(writer, request, formRequest)

	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	create, err := r.MovingCodeService.CreateMovingCode(formRequest)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, create, "Create Data Successfully!", http.StatusOK)
}

// GetAllMovingCode implements MovingCodeController.
func (r *MovingCodeControllerImpl) GetAllMovingCode(writer http.ResponseWriter, request *http.Request) {

	companyId, _ := strconv.Atoi(chi.URLParam(request, "company_id"))

	queryValues := request.URL.Query()

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: "asc",
		SortBy: "Priority",
	}

	result, err := r.MovingCodeService.GetAllMovingCode(companyId, paginate)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data Successfully!", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)

}

// GetMovingCodebyId implements MovingCodeController.
func (r *MovingCodeControllerImpl) GetMovingCodebyId(writer http.ResponseWriter, request *http.Request) {
	movingCodeId, _ := strconv.Atoi(chi.URLParam(request, "moving_code_id"))

	result, err := r.MovingCodeService.GetMovingCodebyId(movingCodeId)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// PushMovingCodePriority implements MovingCodeController.
func (r *MovingCodeControllerImpl) PushMovingCodePriority(writer http.ResponseWriter, request *http.Request) {
	itemPackageId, _ := strconv.Atoi(chi.URLParam(request, "moving_code_id"))
	companyId, _ := strconv.Atoi(chi.URLParam(request, "company_id"))

	result, err := r.MovingCodeService.PushMovingCodePriority(companyId, itemPackageId)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Push Priority Successfull!", http.StatusOK)
}

// UpdateMovingCode implements MovingCodeController.
func (r *MovingCodeControllerImpl) UpdateMovingCode(writer http.ResponseWriter, request *http.Request) {
	var formRequest masterpayloads.MovingCodeListUpdate
	err := jsonchecker.ReadFromRequestBody(request, &formRequest)

	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	err = validation.ValidationForm(writer, request, formRequest)

	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	create, err := r.MovingCodeService.UpdateMovingCode(formRequest)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, create, "Update Data Successfully!", http.StatusOK)
}

func NewMovingCodeController(MovingCodeService masterservice.MovingCodeService) MovingCodeController {
	return &MovingCodeControllerImpl{
		MovingCodeService: MovingCodeService,
	}
}
