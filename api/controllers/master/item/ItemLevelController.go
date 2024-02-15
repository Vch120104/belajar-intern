package masteritemcontroller

import (
	"after-sales/api/helper"
	"after-sales/api/payloads"
	"strconv"

	// masteritemlevelentities "after-sales/api/entities/master/item_level"
	masteritemlevelpayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"

	// masteritemlevelrepo "after-sales/api/repositories/master/item_level"
	masteritemlevelservice "after-sales/api/services/master/item"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type ItemLevelController interface {
	GetAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	GetById(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Save(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	ChangeStatus(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}

type ItemLevelControllerImpl struct {
	itemLevelService masteritemlevelservice.ItemLevelService
}

func NewItemLevelController(ItemLevelService masteritemlevelservice.ItemLevelService) ItemLevelController {
	return &ItemLevelControllerImpl{
		itemLevelService: ItemLevelService,
	}
}

// @Summary Get All Item Level
// @Description Get All Item Level
// @Accept json
// @Produce json
// @Tags Master : Item Level
// @Security BearerAuth
// @Success 200 {object} payloads.Response
// @Param page query string true "Page"
// @Param limit query string true "Limit"
// @Param sort_by query string false "Sort Of: {column}"
// @Param sort_of query string false "Sort By: {asc}"
// @Param item_level query string false "Item Level"
// @Param item_class_code query string false "Item Class Code"
// @Param item_level_parent query string false "Item Level Parent"
// @Param item_level_code query string false "Item Level Code"
// @Param item_level_name query string false "Item Level Name"
// @Param is_active query bool false "Is Active"
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/item-level [get]
func (r *ItemLevelControllerImpl) GetAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	queryValues := request.URL.Query()
	page, _ := strconv.Atoi(queryValues.Get("page"))
	limit, _ := strconv.Atoi(queryValues.Get("limit"))
	sortOf := queryValues.Get("sort_of")
	sortBy := queryValues.Get("sort_by")
	itemLevel := queryValues.Get("item_level")
	itemClassCode := queryValues.Get("item_class_code")
	itemLevelParent := queryValues.Get("item_level_parent")
	itemLevelCode := queryValues.Get("item_level_code")
	itemLevelName := queryValues.Get("item_level_name")
	isActive := queryValues.Get("is_active")

	get := r.itemLevelService.GetAll(masteritemlevelpayloads.GetAllItemLevelResponse{
		ItemLevel:       itemLevel,
		ItemClassCode:   itemClassCode,
		ItemLevelParent: itemLevelParent,
		ItemLevelCode:   itemLevelCode,
		ItemLevelName:   itemLevelName,
		IsActive:        isActive,
	}, pagination.Pagination{
		Limit:  limit,
		SortOf: sortOf,
		SortBy: sortBy,
		Page:   page,
	})

	payloads.NewHandleSuccessPagination(writer, get.Rows, "Get Data Successfully!", 200, get.Limit, get.Page, get.TotalRows, get.TotalPages)
}

// @Summary Get Item Level By Id
// @Description Get Item Level By Id
// @Accept json
// @Produce json
// @Tags Master : Item Level
// @Security BearerAuth
// @Param item_level_id path string true "item_level_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/item-level-by-id [get]
func (r *ItemLevelControllerImpl) GetById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	itemLevelId, _ := strconv.Atoi(params.ByName("item_level_id"))

	get := r.itemLevelService.GetById(itemLevelId)

	payloads.NewHandleSuccess(writer, get, "Get Data Successfully!", http.StatusOK)
}

// @Summary Save Item Level
// @Description Save Item Level
// @Accept json
// @Produce json
// @Tags Master : Item Level
// @Security BearerAuth
// @param reqBody body masteritemlevelpayloads.SaveItemLevelRequest true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/item-level [post]
func (r *ItemLevelControllerImpl) Save(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	var formRequest masteritemlevelpayloads.SaveItemLevelRequest
	var message = ""

	helper.ReadFromRequestBody(request, &formRequest)

	create := r.itemLevelService.Save(formRequest)

	if formRequest.ItemLevelId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

// @Summary Change Item Level Status By Id
// @Description Change Item Level Status By Id
// @Accept json
// @Produce json
// @Tags Master : Item Level
// @Security BearerAuth
// @Param item_level_id path string true "item_level_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/item-level/{item_level_id} [patch]
func (r *ItemLevelControllerImpl) ChangeStatus(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	itemLevelId, _ := strconv.Atoi(params.ByName("item_level_id"))

	response := r.itemLevelService.ChangeStatus(int(itemLevelId))

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}
