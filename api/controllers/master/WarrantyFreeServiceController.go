package mastercontroller

import (
	"after-sales/api/helper"
	"after-sales/api/payloads"
	masterpayloads "after-sales/api/payloads/master"
	masterservice "after-sales/api/services/master"
	"after-sales/api/utils"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type WarrantyFreeServiceController interface {
	GetWarrantyFreeServiceByID(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	SaveWarrantyFreeService(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}
type WarrantyFreeServiceControllerImpl struct {
	WarrantyFreeServiceService masterservice.WarrantyFreeServiceService
}

func NewWarrantyFreeServiceController(warrantyFreeServiceService masterservice.WarrantyFreeServiceService) WarrantyFreeServiceController {
	return &WarrantyFreeServiceControllerImpl{
		WarrantyFreeServiceService: warrantyFreeServiceService,
	}
}

func (r *WarrantyFreeServiceControllerImpl) GetWarrantyFreeServiceByID(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	warrantyFreeServiceId, _ := strconv.Atoi(params.ByName("warranty_free_services_id"))

	result := r.WarrantyFreeServiceService.GetWarrantyFreeServiceById(warrantyFreeServiceId)

	payloads.NewHandleSuccess(writer, utils.ModifyKeysInResponse(result), "Get Data Successfully!", http.StatusOK)
}

func (r *WarrantyFreeServiceControllerImpl) SaveWarrantyFreeService(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	var formRequest masterpayloads.WarrantyFreeServiceRequest
	helper.ReadFromRequestBody(request, &formRequest)
	var message = ""

	create := r.WarrantyFreeServiceService.SaveWarrantyFreeService(formRequest)

	if formRequest.WarrantyFreeServicesId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}
