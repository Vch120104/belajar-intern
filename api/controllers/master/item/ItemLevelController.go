package masteritemcontroller

import (
	"after-sales/api/exceptions"
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

	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	sortOf := c.Query("sort_of")
	sortBy := c.Query("sort_by")
	itemLevel := c.Query("item_level")
	itemClassCode := c.Query("item_class_code")
	itemLevelParent := c.Query("item_level_parent")
	itemLevelCode := c.Query("item_level_code")
	itemLevelName := c.Query("item_level_name")
	isActive := c.Query("is_active")

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

	payloads.HandleSuccessPagination(c, get.Rows, "Get Data Successfully!", 200, get.Limit, get.Page, get.TotalRows, get.TotalPages)
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

	itemLevelId, _ := strconv.Atoi(c.Param("item_level_id"))

	get := r.itemLevelService.GetById(itemLevelId)

	if get.ItemLevelId == 0 {
		exceptions.NotFoundException(c, "Item Level Data Not Found!")
		return
	}

	payloads.HandleSuccess(c, get, "Get Data Successfully!", http.StatusOK)
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

	var request masteritemlevelpayloads.SaveItemLevelRequest
	var message = ""

	if err := c.ShouldBindJSON(&request); err != nil {
		exceptions.EntityException(c, err.Error())
		return
	}

	create := r.itemLevelService.Save(request)

	if request.ItemLevelId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.HandleSuccess(c, create, message, http.StatusOK)
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

	itemLevelId, _ := strconv.Atoi(c.Param("item_level_id"))

	response := r.itemLevelService.ChangeStatus(int(itemLevelId))

	payloads.HandleSuccess(c, response, "Update Data Successfully!", http.StatusOK)
}
