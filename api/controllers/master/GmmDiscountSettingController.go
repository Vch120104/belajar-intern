package mastercontroller

import (
	"after-sales/api/helper"
	"after-sales/api/payloads"
	masterservice "after-sales/api/services/master"
	"net/http"
)

type GmmDiscountSettingController interface {
	GetAllGmmDiscountSetting(writer http.ResponseWriter, request *http.Request)
}

type GmmDiscountSettingControllerImpl struct {
	GmmDiscountSettingService masterservice.GmmDiscountSettingService
}

func NewGmmDiscountSettingControllerImpl(gmmDiscountSettingService masterservice.GmmDiscountSettingService) GmmDiscountSettingController {
	return &GmmDiscountSettingControllerImpl{
		GmmDiscountSettingService: gmmDiscountSettingService,
	}
}

func (c *GmmDiscountSettingControllerImpl) GetAllGmmDiscountSetting(writer http.ResponseWriter, request *http.Request) {
	result, err := c.GmmDiscountSettingService.GetAllGmmDiscountSetting()
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}
