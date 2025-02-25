package mastercontroller

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	jsonchecker "after-sales/api/helper/json/json-checker"
	"after-sales/api/payloads"
	masterpayloads "after-sales/api/payloads/master"
	masterservice "after-sales/api/services/master"
	"after-sales/api/validation"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type GmmPriceCodeController interface {
	GetAllGmmPriceCode(writer http.ResponseWriter, request *http.Request)
	GetGmmPriceCodeById(writer http.ResponseWriter, request *http.Request)
	GetGmmPriceCodeDropdown(writer http.ResponseWriter, request *http.Request)
	SaveGmmPriceCode(writer http.ResponseWriter, request *http.Request)
	UpdateGmmPriceCode(writer http.ResponseWriter, request *http.Request)
	ChangeStatusGmmPriceCode(writer http.ResponseWriter, request *http.Request)
	DeleteGmmPriceCode(writer http.ResponseWriter, request *http.Request)
}

type GmmPriceCodeControllerImpl struct {
	GmmPriceCodeService masterservice.GmmPriceCodeService
}

func NewGmmPriceCodeControllerImpl(gmmPriceCodeService masterservice.GmmPriceCodeService) GmmPriceCodeController {
	return &GmmPriceCodeControllerImpl{
		GmmPriceCodeService: gmmPriceCodeService,
	}
}

// @Summary Get All Gmm Price Code
// @Description REST API Gmm Price Code
// @Accept json
// @Produce json
// @Tags Master : Gmm Price Code
// @Security AuthorizationKeyAuth
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/gmm-price-code [get]
func (c *GmmPriceCodeControllerImpl) GetAllGmmPriceCode(writer http.ResponseWriter, request *http.Request) {
	results, err := c.GmmPriceCodeService.GetAllGmmPriceCode()
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, results, "Get Data Successfully!", http.StatusOK)
}

// @Summary Get Gmm Price Code By Id
// @Description REST API Gmm Price Code
// @Accept json
// @Produce json
// @Tags Master : Gmm Price Code
// @Security AuthorizationKeyAuth
// @Param gmm_price_code_id path int true "gmm_price_code_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/gmm-price-code/{gmm_price_code_id} [get]
func (c *GmmPriceCodeControllerImpl) GetGmmPriceCodeById(writer http.ResponseWriter, request *http.Request) {
	gmmPriceCodeId, errA := strconv.Atoi(chi.URLParam(request, "gmm_price_code_id"))
	if errA != nil {
		helper.ReturnError(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("failed to read url param, please check your param input"),
		})
		return
	}
	if gmmPriceCodeId == 0 {
		helper.ReturnError(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("ID cannot be 0"),
		})
		return
	}

	result, err := c.GmmPriceCodeService.GetGmmPriceCodeById(gmmPriceCodeId)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Get Gmm Price Code Dropdown
// @Description REST API Gmm Price Code
// @Accept json
// @Produce json
// @Tags Master : Gmm Price Code
// @Security AuthorizationKeyAuth
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/gmm-price-code/dropdown [get]
func (c *GmmPriceCodeControllerImpl) GetGmmPriceCodeDropdown(writer http.ResponseWriter, request *http.Request) {
	results, err := c.GmmPriceCodeService.GetGmmPriceCodeDropdown()
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, results, "Get Data Successfully!", http.StatusOK)
}

// @Summary Save Gmm Price Code
// @Description REST API Gmm Price Code
// @Accept json
// @Produce json
// @Tags Master : Gmm Price Code
// @Security AuthorizationKeyAuth
// @Param body body masterpayloads.GmmPriceCodeSaveRequest true "body"
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/gmm-price-code [post]
func (c *GmmPriceCodeControllerImpl) SaveGmmPriceCode(writer http.ResponseWriter, request *http.Request) {
	formRequest := masterpayloads.GmmPriceCodeSaveRequest{}
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

	create, err := c.GmmPriceCodeService.SaveGmmPriceCode(formRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, create, "Create Data Successfully!", http.StatusCreated)
}

// @Summary Update Gmm Price Code
// @Description REST API Gmm Price Code
// @Accept json
// @Produce json
// @Tags Master : Gmm Price Code
// @Security AuthorizationKeyAuth
// @Param gmm_price_code_id path int true "gmm_price_code_id"
// @Param body body masterpayloads.GmmPriceCodeUpdateRequest true "body"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/gmm-price-code/{gmm_price_code_id} [put]
func (c *GmmPriceCodeControllerImpl) UpdateGmmPriceCode(writer http.ResponseWriter, request *http.Request) {
	gmmPriceCodeId, errA := strconv.Atoi(chi.URLParam(request, "gmm_price_code_id"))
	if errA != nil {
		helper.ReturnError(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("failed to read url param, please check your param input"),
		})
		return
	}
	if gmmPriceCodeId == 0 {
		helper.ReturnError(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("ID cannot be 0"),
		})
		return
	}

	formRequest := masterpayloads.GmmPriceCodeUpdateRequest{}
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

	update, err := c.GmmPriceCodeService.UpdateGmmPriceCode(gmmPriceCodeId, formRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, update, "Update Data Successfully!", http.StatusOK)
}

// @Summary Change Status Gmm Price Code
// @Description REST API Gmm Price Code
// @Accept json
// @Produce json
// @Tags Master : Gmm Price Code
// @Security AuthorizationKeyAuth
// @Param gmm_price_code_id path int true "gmm_price_code_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/gmm-price-code/{gmm_price_code_id} [patch]
func (c *GmmPriceCodeControllerImpl) ChangeStatusGmmPriceCode(writer http.ResponseWriter, request *http.Request) {
	gmmPriceCodeId, errA := strconv.Atoi(chi.URLParam(request, "gmm_price_code_id"))
	if errA != nil {
		helper.ReturnError(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("failed to read url param, please check your param input"),
		})
		return
	}
	if gmmPriceCodeId == 0 {
		helper.ReturnError(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("ID cannot be 0"),
		})
		return
	}

	status, err := c.GmmPriceCodeService.ChangeStatusGmmPriceCode(gmmPriceCodeId)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, status, "Update Data Successfully!", http.StatusOK)
}

// @Summary Delete Gmm Price Code
// @Description REST API Gmm Price Code
// @Accept json
// @Produce json
// @Tags Master : Gmm Price Code
// @Security AuthorizationKeyAuth
// @Param gmm_price_code_id path int true "gmm_price_code_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/gmm-price-code/{gmm_price_code_id} [delete]
func (c *GmmPriceCodeControllerImpl) DeleteGmmPriceCode(writer http.ResponseWriter, request *http.Request) {
	gmmPriceCodeId, errA := strconv.Atoi(chi.URLParam(request, "gmm_price_code_id"))
	if errA != nil {
		helper.ReturnError(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("failed to read url param, please check your param input"),
		})
		return
	}
	if gmmPriceCodeId == 0 {
		helper.ReturnError(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("ID cannot be 0"),
		})
		return
	}

	delete, err := c.GmmPriceCodeService.DeleteGmmPriceCode(gmmPriceCodeId)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, delete, "Delete Data Successfully!", http.StatusOK)
}
