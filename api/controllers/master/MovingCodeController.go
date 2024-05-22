package mastercontroller

import (
	exceptions "after-sales/api/exceptions"
	helper_test "after-sales/api/helper_testt"
	jsonchecker "after-sales/api/helper_testt/json/json-checker"
	"after-sales/api/payloads"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterservice "after-sales/api/services/master"
	"after-sales/api/utils"
	"after-sales/api/validation"
	"errors"
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
}

type MovingCodeControllerImpl struct {
	MovingCodeService masterservice.MovingCodeService
}

// ChangeStatusMovingCode implements MovingCodeController.
func (r *MovingCodeControllerImpl) ChangeStatusMovingCode(writer http.ResponseWriter, request *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(request, "item_package_detail_id"))

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
	queryValues := request.URL.Query()

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: "Priority",
		SortBy: "asc",
	}

	paginatedData, totalPages, totalRows, err := r.MovingCodeService.GetAllMovingCode(paginate)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(paginatedData), "Get Data Successfully!", 200, paginate.Limit, paginate.Page, int64(totalRows), totalPages)

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
	result, err := r.MovingCodeService.PushMovingCodePriority(itemPackageId)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Push Priority Successfull!", http.StatusOK)
}

// UpdateMovingCode implements MovingCodeController.
func (r *MovingCodeControllerImpl) UpdateMovingCode(writer http.ResponseWriter, request *http.Request) {
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
