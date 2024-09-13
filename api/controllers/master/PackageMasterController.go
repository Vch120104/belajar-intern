package mastercontroller

import (
	"after-sales/api/exceptions"
	helper "after-sales/api/helper"
	"after-sales/api/payloads"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterservice "after-sales/api/services/master"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type PackageMasterController interface {
	GetAllPackageMaster(writer http.ResponseWriter, request *http.Request)
	GetAllPackageMasterDetail(writer http.ResponseWriter, request *http.Request)
	GetByIdPackageMaster(writer http.ResponseWriter, request *http.Request)
	GetByIdPackageMasterDetail(writer http.ResponseWriter, request *http.Request)
	SavepackageMaster(writer http.ResponseWriter, request *http.Request)
	SavePackageMasterDetail(writer http.ResponseWriter, request *http.Request)
	ChangeStatusPackageMaster(writer http.ResponseWriter, request *http.Request)
	ActivateMultiIdPackageMasterDetail(writer http.ResponseWriter, request *http.Request)
	DeactivateMultiIdPackageMasterDetail(writer http.ResponseWriter, request *http.Request)
	CopyToOtherModel(writer http.ResponseWriter, request *http.Request)
}

type PackageMasterControllerImpl struct {
	PackageMasterService masterservice.PackageMasterService
}

func NewPackageMasterController(packageMasterService masterservice.PackageMasterService) PackageMasterController {
	return &PackageMasterControllerImpl{
		PackageMasterService: packageMasterService,
	}
}

func (r *PackageMasterControllerImpl) GetAllPackageMaster(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	queryParams := map[string]string{
		"mtr_package.package_name":     queryValues.Get("package_name"),
		"mtr_package.package_code":     queryValues.Get("package_code"),
		"mtr_package.profit_center_id": queryValues.Get("profit_center_id"),
		"mtr_package.model_id":         queryValues.Get("model_id"),
		"mtr_package.variant_id":       queryValues.Get("variant_id"),
		"mtr_package.package_price":    queryValues.Get("package_price"),
		"mtr_package.is_active":        queryValues.Get("is_active"),
	}

	pagination := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	filterCondition := utils.BuildFilterCondition(queryParams)

	result, totalPages, totalRows, err := r.PackageMasterService.GetAllPackageMaster(filterCondition, pagination)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(result), "success", 200, pagination.Limit, pagination.Page, int64(totalRows), totalPages)
}

func (r *PackageMasterControllerImpl) GetAllPackageMasterDetail(writer http.ResponseWriter, request *http.Request) {
	PackageMasterId, errA := strconv.Atoi(chi.URLParam(request, "package_id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}
	queryValues := request.URL.Query()
	pagination := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	result, totalPages, totalRows, err := r.PackageMasterService.GetAllPackageMasterDetail(pagination, PackageMasterId)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(result), "success", 200, pagination.Limit, pagination.Page, int64(totalRows), totalPages)
}

func (r *PackageMasterControllerImpl) GetByIdPackageMaster(writer http.ResponseWriter, request *http.Request) {
	PackageMasterId, errA := strconv.Atoi(chi.URLParam(request, "package_id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}
	result, err := r.PackageMasterService.GetByIdPackageMaster(PackageMasterId)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

func (r *PackageMasterControllerImpl) GetByIdPackageMasterDetail(writer http.ResponseWriter, request *http.Request) {
	PackageMasterDetailId, errA := strconv.Atoi(chi.URLParam(request, "package_detail_id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}
	result, err := r.PackageMasterService.GetByIdPackageMasterDetail(PackageMasterDetailId)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

func (r *PackageMasterControllerImpl) SavepackageMaster(writer http.ResponseWriter, request *http.Request) {
	var formRequest masterpayloads.PackageMasterResponse
	helper.ReadFromRequestBody(request, &formRequest)
	var message string

	create, err := r.PackageMasterService.PostPackageMaster(formRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	if formRequest.PackageId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

func (r *PackageMasterControllerImpl) SavePackageMasterDetail(writer http.ResponseWriter, request *http.Request) {
	var formRequest masterpayloads.PackageMasterDetail
	helper.ReadFromRequestBody(request, &formRequest)
	packageId,_ := strconv.Atoi(chi.URLParam(request,"package_id"))

	create, err := r.PackageMasterService.PostPackageMasterDetail(formRequest,packageId)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, create, "create data successfully", http.StatusOK)
}

func (r *PackageMasterControllerImpl) ChangeStatusPackageMaster(writer http.ResponseWriter, request *http.Request) {
	PackageMasterId, errA := strconv.Atoi(chi.URLParam(request, "package_id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}
	result, err := r.PackageMasterService.ChangeStatusItemPackage(PackageMasterId)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

func (r *PackageMasterControllerImpl) ActivateMultiIdPackageMasterDetail(writer http.ResponseWriter, request *http.Request) {
	PackageDetailId := chi.URLParam(request, "package_detail_id")

	response, err := r.PackageMasterService.ActivateMultiIdPackageMasterDetail(PackageDetailId)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}

func (r *PackageMasterControllerImpl) DeactivateMultiIdPackageMasterDetail(writer http.ResponseWriter, request *http.Request) {
	PackageDetailId := chi.URLParam(request, "package_detail_id")
	
	response, err := r.PackageMasterService.DeactivateMultiIdPackageMasterDetail(PackageDetailId)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}

func (r *PackageMasterControllerImpl) CopyToOtherModel(writer http.ResponseWriter, request *http.Request) {
	PackageDetailId := chi.URLParam(request, "package_name")
	PackageMasterId, errA := strconv.Atoi(chi.URLParam(request, "package_id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}
	ModelId, errA := strconv.Atoi(chi.URLParam(request, "model_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	ressult, err := r.PackageMasterService.CopyToOtherModel(PackageMasterId, PackageDetailId, ModelId)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, ressult, "Update Data Successfully!", http.StatusOK)
}
