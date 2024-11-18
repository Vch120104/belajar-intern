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

func (c *GmmPriceCodeControllerImpl) GetAllGmmPriceCode(writer http.ResponseWriter, request *http.Request) {
	results, err := c.GmmPriceCodeService.GetAllGmmPriceCode()
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, results, "Get Data Successfully!", http.StatusOK)
}

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

func (c *GmmPriceCodeControllerImpl) GetGmmPriceCodeDropdown(writer http.ResponseWriter, request *http.Request) {
	results, err := c.GmmPriceCodeService.GetGmmPriceCodeDropdown()
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, results, "Get Data Successfully!", http.StatusOK)
}

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
