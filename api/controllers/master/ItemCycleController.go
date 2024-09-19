package mastercontroller

import (
	"after-sales/api/helper"
	"after-sales/api/payloads"
	masterpayloads "after-sales/api/payloads/master"
	masterservice "after-sales/api/services/master"
	"net/http"
)

type ItemCycleController interface {
	ItemCycleInsert(writer http.ResponseWriter, request *http.Request)
}

type ItemCycleControllerImpl struct {
	ItemCycleService masterservice.ItemCycleService
}

func NewItemCycleController(ItemCycleService masterservice.ItemCycleService) ItemCycleController {
	return &ItemCycleControllerImpl{ItemCycleService: ItemCycleService}
}

// ItemCycleInsert
//
// @Summary			Item Cycle Insert
// @Description		Item Cycle Insert
// @Accept			json
// @Produce			json
// @Tags			Master : Item Cycle
// @Param			reqBody					body	masterpayloads.ItemCycleInsertPayloads	true	"Item Cycle Insert"
// @Success		201						{object}	payloads.Response
// @Failure		500,400,401,404,403,422	{object}	exceptions.BaseErrorResponse
// @Router			/v1/item-cycle [post]
func (i *ItemCycleControllerImpl) ItemCycleInsert(writer http.ResponseWriter, request *http.Request) {
	var itemCycleInsert masterpayloads.ItemCycleInsertPayloads
	helper.ReadFromRequestBody(request, &itemCycleInsert)
	success, err := i.ItemCycleService.ItemCycleInsert(itemCycleInsert)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, success, "Insert Item Cycle Success", http.StatusCreated)
}
