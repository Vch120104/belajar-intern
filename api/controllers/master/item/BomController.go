package masteritemcontroller

import (
	"after-sales/api/helper"
	"after-sales/api/payloads"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemservice "after-sales/api/services/master/item"
	"after-sales/api/utils"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type BomController interface {
	GetBomMasterById(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	GetBomMasterList(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	SaveBomMaster(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	ChangeStatusBomMaster(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}

type BomControllerImpl struct {
	BomService masteritemservice.BomService
}

func NewBomController(bomService masteritemservice.BomService) BomController {
	return &BomControllerImpl{
		BomService: bomService,
	}
}

// @Summary Get All Bom Master
// @Description REST API Bom Master
// @Accept json
// @Produce json
// @Tags Master : Bom Master
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Param is_active query string false "is_active" Enums(true, false)
// @Param bom_master_code query string false "bom_master_code"
// @Param bom_master_qty query string false "bom_master_qty"
// @Param bom_master_uom query string false "bom_master_uom"
// @Param bom_master_effective_date query time false "bom_master_effective_date"
// @Param sort_by query string false "sort_by"
// @Param sort_of query string false "sort_of"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /bom [get]
func (r *BomControllerImpl) GetBomMasterList(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"is_active":                         queryValues.Get("is_active"),
		"mtr_bom.bom_master_id":             queryValues.Get("bom_master_id"),
		"mtr_bom.bom_master_effective_date": queryValues.Get("bom_master_effective_date"),
		"mtr_bom.item_id":                   queryValues.Get("mtr_bom.item_id"), // Ambil nilai item_id tanpa mtr_bom.
	}

	// // Get the value of bom_master_code from the URL query
	// bomMasterCode := queryValues.Get("mtr_bom.item_id") // Ambil nilai item_id tanpa mtr_bom.
	// log.Printf("Value of bom_master_code from URL query: %s", bomMasterCode)

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	//log.Printf("Criteria created: %+v", criteria) // Add logging to display the created criteria

	paginatedData, totalPages, totalRows := r.BomService.GetBomMasterList(criteria, paginate)

	payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(paginatedData), "Get Data Successfully", 200, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
}

// @Summary Get Bom Master By ID
// @Description REST API Bom Master
// @Accept json
// @Produce json
// @Tags Master : Bom Master
// @Param bom_master_id path int true "bom_master_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /bom/{bom_master_id} [get]
func (r *BomControllerImpl) GetBomMasterById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	bomMasterId, _ := strconv.Atoi(params.ByName("bom_master_id"))

	result := r.BomService.GetBomMasterById(bomMasterId)

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Save Bom Master
// @Description REST API Bom Master
// @Accept json
// @Produce json
// @Tags Master : Bom Master
// @param reqBody body masteritempayloads.BomMasterResponse true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /bom [post]
func (r *BomControllerImpl) SaveBomMaster(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	var formRequest masteritempayloads.BomMasterRequest
	var message = ""
	helper.ReadFromRequestBody(request, &formRequest)

	create := r.BomService.SaveBomMaster(formRequest)

	if formRequest.BomMasterId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

// @Summary Change Status Bom Master
// @Description REST API Bom Master
// @Accept json
// @Produce json
// @Tags Master : Bom Master
// @param bom_master_id path int true "bom_master_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /bom/{bom_master_id} [patch]
func (r *BomControllerImpl) ChangeStatusBomMaster(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	bomMasterId, _ := strconv.Atoi(params.ByName("bom_master_id"))

	response := r.BomService.ChangeStatusBomMaster(int(bomMasterId))

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}
