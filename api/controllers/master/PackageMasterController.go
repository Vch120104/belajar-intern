package mastercontroller

import (
	"after-sales/api/exceptions"
	helper "after-sales/api/helper"
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

type PackageMasterController interface {
	GetAllPackageMaster(writer http.ResponseWriter, request *http.Request)
	GetAllPackageMasterDetail(writer http.ResponseWriter, request *http.Request)
	GetByIdPackageMaster(writer http.ResponseWriter, request *http.Request)
	GetByIdPackageMasterDetail(writer http.ResponseWriter, request *http.Request)
	GetByCodePackageMaster(writer http.ResponseWriter, request *http.Request)
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

// @Summary Get All Package Master
// @Description Get All Package Master
// @Tags Master : Package Master
// @Accept json
// @Produce json
// @Param package_name query string false "package_name"
// @Param package_code query string false "package_code"
// @Param profit_center_id query string false "profit_center_id"
// @Param profit_center_name query string false "profit_center_name"
// @Param model_id query string false "model_id"
// @Param model_code query string false "model_code"
// @Param model_description query string false "model_description"
// @Param variant_id query string false "variant_id"
// @Param variant_description query string false "variant_description"
// @Param package_price query string false "package_price"
// @Param is_active query string false "is_active"
// @Param limit query int false "limit"
// @Param page query int false "page"
// @Param sort_of query string false "sort_of"
// @Param sort_by query string false "sort_by"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/package-master [get]
func (r *PackageMasterControllerImpl) GetAllPackageMaster(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	queryParams := map[string]string{
		"mtr_package.package_name":     queryValues.Get("package_name"),
		"mtr_package.package_code":     queryValues.Get("package_code"),
		"mtr_package.profit_center_id": queryValues.Get("profit_center_id"),
		"profit_center_name":           queryValues.Get("profit_center_name"),
		"mtr_package.model_id":         queryValues.Get("model_id"),
		"model_code":                   queryValues.Get("model_code"),
		"model_description":            queryValues.Get("model_description"),
		"mtr_package.variant_id":       queryValues.Get("variant_id"),
		"variant_description":          queryValues.Get("variant_description"),
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

	result, err := r.PackageMasterService.GetAllPackageMaster(filterCondition, pagination)
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

// @Summary Get All Package Master Detail
// @Description Get All Package Master Detail
// @Tags Master : Package Master
// @Accept json
// @Produce json
// @Param package_id path int true "package_id"
// @Param limit query int false "limit"
// @Param page query int false "page"
// @Param sort_of query string false "sort_of"
// @Param sort_by query string false "sort_by"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/package-master/detail/{package_id} [get]
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

	result, err := r.PackageMasterService.GetAllPackageMasterDetail(pagination, PackageMasterId)
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

// @Summary Get By Id Package Master
// @Description Get By Id Package Master
// @Tags Master : Package Master
// @Accept json
// @Produce json
// @Param package_id path int true "package_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/package-master/{package_id} [get]
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

// @Summary Get By Id Package Master Detail
// @Description Get By Id Package Master Detail
// @Tags Master : Package Master
// @Accept json
// @Produce json
// @Param package_detail_id path int true "package_detail_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/package-master/detail/by-id/{package_detail_id} [get]
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

// @Summary Get By Code Package Master
// @Description Get By Code Package Master
// @Tags Master : Package Master
// @Accept json
// @Produce json
// @Param package_code path string true "package_code"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/package-master/by-code/{package_code} [get]
func (r *PackageMasterControllerImpl) GetByCodePackageMaster(writer http.ResponseWriter, request *http.Request) {
	PackageMasterCode := chi.URLParam(request, "package_code")

	result, err := r.PackageMasterService.GetByCodePackageMaster(PackageMasterCode)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, result, "Get Data Successfully", http.StatusOK)
}

// @Summary Save Package Master
// @Description Save Package Master
// @Tags Master : Package Master
// @Accept json
// @Produce json
// @Param request body masterpayloads.PackageMasterResponse true "PackageMasterResponse"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/package-master [post]
func (r *PackageMasterControllerImpl) SavepackageMaster(writer http.ResponseWriter, request *http.Request) {
	var formRequest masterpayloads.PackageMasterResponse
	helper.ReadFromRequestBody(request, &formRequest)
	if validationErr := validation.ValidationForm(writer, request, &formRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}
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

// @Summary Save Package Master Detail
// @Description Save Package Master Detail
// @Tags Master : Package Master
// @Accept json
// @Produce json
// @Param package_id path int true "package_id"
// @Param request body masterpayloads.PackageMasterDetail true "PackageMasterDetail"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/package-master/detail/{package_id} [post]
func (r *PackageMasterControllerImpl) SavePackageMasterDetail(writer http.ResponseWriter, request *http.Request) {
	var formRequest masterpayloads.PackageMasterDetail
	helper.ReadFromRequestBody(request, &formRequest)
	if validationErr := validation.ValidationForm(writer, request, &formRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}
	packageId, _ := strconv.Atoi(chi.URLParam(request, "package_id"))

	create, err := r.PackageMasterService.PostPackageMasterDetail(formRequest, packageId)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, create, "create data successfully", http.StatusOK)
}

// @Summary Change Status Package Master
// @Description Change Status Package Master
// @Tags Master : Package Master
// @Accept json
// @Produce json
// @Param package_id path int true "package_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/package-master/{package_id} [patch]
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

// @Summary Activate Multi Id Package Master Detail
// @Description Activate Multi Id Package Master Detail
// @Tags Master : Package Master
// @Accept json
// @Produce json
// @Param package_detail_id path int true "package_detail_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/package-master/detail/activate/{package_id}/{package_detail_id} [patch]
func (r *PackageMasterControllerImpl) ActivateMultiIdPackageMasterDetail(writer http.ResponseWriter, request *http.Request) {
	PackageDetailId := chi.URLParam(request, "package_detail_id")

	response, err := r.PackageMasterService.ActivateMultiIdPackageMasterDetail(PackageDetailId)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}

// @Summary Deactivate Multi Id Package Master Detail
// @Description Deactivate Multi Id Package Master Detail
// @Tags Master : Package Master
// @Accept json
// @Produce json
// @Param package_detail_id path int true "package_detail_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/package-master/detail/deactivate/{package_id}/{package_detail_id} [patch]
func (r *PackageMasterControllerImpl) DeactivateMultiIdPackageMasterDetail(writer http.ResponseWriter, request *http.Request) {
	PackageDetailId := chi.URLParam(request, "package_detail_id")

	response, err := r.PackageMasterService.DeactivateMultiIdPackageMasterDetail(PackageDetailId)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}

// @Summary Copy To Other Model
// @Description Copy To Other Model
// @Tags Master : Package Master
// @Accept json
// @Produce json
// @Param package_id path int true "package_id"
// @Param package_name path string true "package_name"
// @Param model_id path int true "model_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/package-master/copy/{package_id}/{package_name}/{model_id} [get]
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
