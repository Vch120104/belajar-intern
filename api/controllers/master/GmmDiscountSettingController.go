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

// @Summary Get All Gmm Discount Setting
// @Description REST API Gmm Discount Setting
// @Accept json
// @Produce json
// @Tags Master : Gmm Discount Setting
// @Security AuthorizationKeyAuth
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/gmm-discount-setting [get]
func (c *GmmDiscountSettingControllerImpl) GetAllGmmDiscountSetting(writer http.ResponseWriter, request *http.Request) {
	result, err := c.GmmDiscountSettingService.GetAllGmmDiscountSetting()
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}
