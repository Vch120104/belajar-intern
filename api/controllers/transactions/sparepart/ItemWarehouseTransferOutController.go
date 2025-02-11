package transactionsparepartcontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionsparepartservice "after-sales/api/services/transaction/sparepart"
	"after-sales/api/validation"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type ItemWarehouseTransferOutController interface {
	InsertHeader(writer http.ResponseWriter, request *http.Request)
	InsertDetail(writer http.ResponseWriter, request *http.Request)
	InsertDetailFromReceipt(writer http.ResponseWriter, request *http.Request)
	GetTransferOutById(writer http.ResponseWriter, request *http.Request)
	GetAllTransferOut(writer http.ResponseWriter, request *http.Request)
	GetAllTransferOutDetail(writer http.ResponseWriter, request *http.Request)
	SubmitTransferOut(writer http.ResponseWriter, request *http.Request)
	UpdateTransferOutDetail(writer http.ResponseWriter, request *http.Request)
	DeleteTransferOutDetail(writer http.ResponseWriter, request *http.Request)
	DeleteTransferOut(writer http.ResponseWriter, request *http.Request)
}

func NewItemWarehouseTransferOutControllerImpl(itemWarehouseTransferOutService transactionsparepartservice.ItemWarehouseTransferOutService) ItemWarehouseTransferOutController {
	return &ItemWarehouseTransferOutControllerImpl{
		ItemWarehouseTransferOutService: itemWarehouseTransferOutService,
	}
}

type ItemWarehouseTransferOutControllerImpl struct {
	ItemWarehouseTransferOutService transactionsparepartservice.ItemWarehouseTransferOutService
}

// DeleteTransferOut implements ItemWarehouseTransferOutController.
func (r *ItemWarehouseTransferOutControllerImpl) DeleteTransferOut(writer http.ResponseWriter, request *http.Request) {
	panic("unimplemented")
}

// DeleteTransferOutDetail implements ItemWarehouseTransferOutController.
func (r *ItemWarehouseTransferOutControllerImpl) DeleteTransferOutDetail(writer http.ResponseWriter, request *http.Request) {
	panic("unimplemented")
}

// GetAllTransferOut implements ItemWarehouseTransferOutController.
func (r *ItemWarehouseTransferOutControllerImpl) GetAllTransferOut(writer http.ResponseWriter, request *http.Request) {
	panic("unimplemented")
}

// GetAllTransferOutDetail implements ItemWarehouseTransferOutController.
func (r *ItemWarehouseTransferOutControllerImpl) GetAllTransferOutDetail(writer http.ResponseWriter, request *http.Request) {
	panic("unimplemented")
}

// GetTransferOutById implements ItemWarehouseTransferOutController.
func (r *ItemWarehouseTransferOutControllerImpl) GetTransferOutById(writer http.ResponseWriter, request *http.Request) {
	panic("unimplemented")
}

// InsertDetail implements ItemWarehouseTransferOutController.
func (r *ItemWarehouseTransferOutControllerImpl) InsertDetail(writer http.ResponseWriter, request *http.Request) {
	panic("unimplemented")
}

// InsertDetailFromReceipt implements ItemWarehouseTransferOutController.
func (r *ItemWarehouseTransferOutControllerImpl) InsertDetailFromReceipt(writer http.ResponseWriter, request *http.Request) {
	panic("unimplemented")
}

// InsertHeader implements ItemWarehouseTransferOutController.
func (r *ItemWarehouseTransferOutControllerImpl) InsertHeader(writer http.ResponseWriter, request *http.Request) {
	var transferRequest transactionsparepartpayloads.InsertItemWarehouseHeaderTransferOutRequest

	helper.ReadFromRequestBody(request, &transferRequest)
	if validationErr := validation.ValidationForm(writer, request, &transferRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}

	success, err := r.ItemWarehouseTransferOutService.InsertHeader(transferRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, success, "save success", http.StatusCreated)
}

// SubmitTransferOut implements ItemWarehouseTransferOutController.
func (r *ItemWarehouseTransferOutControllerImpl) SubmitTransferOut(writer http.ResponseWriter, request *http.Request) {
	transferRequestSystemNumber, errA := strconv.Atoi(chi.URLParam(request, "id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	var transferRequest transactionsparepartpayloads.SubmitItemWarehouseTransferRequest

	helper.ReadFromRequestBody(request, &transferRequest)
	if validationErr := validation.ValidationForm(writer, request, &transferRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}

	success, err := r.ItemWarehouseTransferOutService.SubmitTransferOut(transferRequestSystemNumber)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, success, "Get Data Success", http.StatusCreated)
}

// UpdateTransferOutDetail implements ItemWarehouseTransferOutController.
func (r *ItemWarehouseTransferOutControllerImpl) UpdateTransferOutDetail(writer http.ResponseWriter, request *http.Request) {
	panic("unimplemented")
}
