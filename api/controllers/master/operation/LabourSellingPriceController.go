package masteroperationcontroller

import (
	exceptionsss_test "after-sales/api/expectionsss"
	helper_test "after-sales/api/helper_testt"
	jsonchecker "after-sales/api/helper_testt/json/json-checker"
	"after-sales/api/payloads"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	masteroperationservice "after-sales/api/services/master/operation"
	"after-sales/api/utils"
	"after-sales/api/validation"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type LabourSellingPriceController interface {
	GetLabourSellingPriceById(writer http.ResponseWriter, request *http.Request)
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

func (r *LabourSellingPriceControllerImpl) GetLabourSellingPriceById(writer http.ResponseWriter, request *http.Request) {

	labourSellingPriceId, _ := strconv.Atoi(chi.URLParam(request, "labour_selling_price_id"))

	result, err := r.LabourSellingPriceService.GetLabourSellingPriceById(labourSellingPriceId)
	if err != nil {
		helper_test.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, utils.ModifyKeysInResponse(result), "Get Data Successfully!", http.StatusOK)
}

func (r *LabourSellingPriceControllerImpl) SaveLabourSellingPrice(writer http.ResponseWriter, request *http.Request) {

	var formRequest masteroperationpayloads.LabourSellingPriceRequest
	err := jsonchecker.ReadFromRequestBody(request, &formRequest)
	var message string

	if err != nil {
		exceptionsss_test.NewEntityException(writer, request, err)
		return
	}
	err = validation.ValidationForm(writer, request, formRequest)
	if err != nil {
		exceptionsss_test.NewBadRequestException(writer, request, err)
		return
	}

	create, err := r.LabourSellingPriceService.SaveLabourSellingPrice(formRequest)

	if err != nil {
		helper_test.ReturnError(writer, request, err)
		return
	}

	message = "Create Data Successfully!"

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}
