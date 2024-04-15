package mastercontroller

import (
	exceptionsss_test "after-sales/api/expectionsss"
	jsonchecker "after-sales/api/helper/json/json-checker"
	helper_test "after-sales/api/helper_testt"
	"after-sales/api/payloads"
	masterpayloads "after-sales/api/payloads/master"
	masterservice "after-sales/api/services/master"
	"after-sales/api/validation"
	"net/http"
)

type CampaignMasterController interface {
}

type CampaignMasterControllerImpl struct {
	CampaignMasterService masterservice.CampaignMasterService
}

func NewCampaignMasterController(CampaignMasterService masterservice.CampaignMasterService) CampaignMasterController{
	return &CampaignMasterControllerImpl{
		CampaignMasterService: CampaignMasterService,
	}
}

func (r *CampaignMasterControllerImpl) SaveCampaignMaster(writer http.ResponseWriter, request *http.Request){
	var requestForm masterpayloads.CampaignMasterPost
	var message string

	err := jsonchecker.ReadFromRequestBody(request, &requestForm)
	if err != nil {
		exceptionsss_test.NewEntityException(writer, request, err)
		return
	}
	err = validation.ValidationForm(writer, request, requestForm)
	if err != nil {
		exceptionsss_test.NewBadRequestException(writer, request, err)
		return
	}

	create, err := r.CampaignMasterService.PostCampaignMaster(requestForm)
	if err != nil {
		helper_test.ReturnError(writer, request, err)
		return
	}

	if requestForm.CampaignId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)

}