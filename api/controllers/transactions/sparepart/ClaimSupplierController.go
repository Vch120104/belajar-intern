package transactionsparepartcontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionsparepartservice "after-sales/api/services/transaction/sparepart"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type ClaimSupplierController interface {
	InsertItemClaim(writer http.ResponseWriter, request *http.Request)
	InsertItemClaimDetail(writer http.ResponseWriter, request *http.Request)
	GetItemClaimById(writer http.ResponseWriter, request *http.Request)
}
type ClaimSupplierControllerImpl struct {
	service transactionsparepartservice.ClaimSupplierService
}

func NewClaimSupplierControllerImpl(service transactionsparepartservice.ClaimSupplierService) ClaimSupplierController {
	return &ClaimSupplierControllerImpl{service: service}
}
func (controller *ClaimSupplierControllerImpl) InsertItemClaim(writer http.ResponseWriter, request *http.Request) {
	var claimSupplierInsertPayload transactionsparepartpayloads.ClaimSupplierInsertPayload
	helper.ReadFromRequestBody(request, &claimSupplierInsertPayload)
	res, err := controller.service.InsertItemClaim(claimSupplierInsertPayload)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, res, "Successfully Inserted item claim Header", http.StatusCreated)
}
func (controller *ClaimSupplierControllerImpl) InsertItemClaimDetail(writer http.ResponseWriter, request *http.Request) {
	var claimSupplierDetailInsertPayload transactionsparepartpayloads.ClaimSupplierInsertDetailPayload
	helper.ReadFromRequestBody(request, &claimSupplierDetailInsertPayload)
	res, err := controller.service.InsertItemClaimDetail(claimSupplierDetailInsertPayload)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	//payloads.NewHandleSuccess(writer, res, "Successfully Inserted item claim Detail", http.StatusInternalServerError)
	payloads.NewHandleSuccess(writer, res, "Successfully Inserted item claim Header", http.StatusCreated)

}
func (controller *ClaimSupplierControllerImpl) GetItemClaimById(writer http.ResponseWriter, request *http.Request) {
	claimIdStr := chi.URLParam(request, "claim_system_number")
	claimId, errs := strconv.Atoi(claimIdStr)
	if errs != nil {
		helper.ReturnError(writer, request, &exceptions.BaseErrorResponse{Err: errs, StatusCode: http.StatusInternalServerError})
		return
	}
	res, err := controller.service.GetItemClaimById(claimId)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, res, "Successfully Inserted item claim Header", http.StatusCreated)
}
