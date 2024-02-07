package masteritemcontroller

// import (
// 	"after-sales/api/payloads"
// 	"after-sales/api/payloads/pagination"
// 	masteritemservice "after-sales/api/services/master/item"
// 	"after-sales/api/utils"
// 	"net/http"

// 	"github.com/julienschmidt/httprouter"
// )

// type DiscountPercentController interface {
// 	GetAllDiscountPercent(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
// }
// type DiscountPercentControllerImpl struct {
// 	DiscountPercentService masteritemservice.DiscountPercentService
// }

// func NewDiscountPercentController(discountPercentService masteritemservice.DiscountPercentService) DiscountPercentController {
// 	return &DiscountPercentControllerImpl{
// 		DiscountPercentService: discountPercentService,
// 	}
// }

// // @Summary Get All Discount Percent
// // @Description REST API Discount Percent
// // @Accept json
// // @Produce json
// // @Tags Master : Discount Percent
// // @Param page query string true "page"
// // @Param limit query string true "limit"
// // @Param is_active query string false "is_active" Enums(true, false)
// // @Param discount_code_value query string false "discount_code_value"
// // @Param discount_code_description query string false "discount_code_description"
// // @Param order_type_name query string false "order_type_name"
// // @Param discount query float64 false "discount"
// // @Param sort_by query string false "sort_by"
// // @Param sort_of query string false "sort_of"
// // @Success 200 {object} payloads.Response
// // @Failure 500,400,401,404,403,422 {object} exceptions.Error
// // @Router /aftersales-service/api/aftersales/discount-percent [get]
// func (r *DiscountPercentControllerImpl) GetAllDiscountPercent(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

// 	queryParams := map[string]string{
// 		"mtr_discount.discount_code_value":       params.ByName("discount_code_value"),
// 		"mtr_discount.discount_code_description": params.ByName("discount_code_description"),
// 		"order_type_name":                        params.ByName("order_type_name"),
// 		"mtr_discount_percent.discount":          params.ByName("discount"),
// 		"mtr_discount_percent.is_active":         params.ByName("is_active"),
// 	}

// 	paginate := pagination.Pagination{
// 		Limit:  utils.NewGetQueryInt(params, "limit"),
// 		Page:   utils.NewGetQueryInt(params, "page"),
// 		SortOf: params.ByName("sort_of"),
// 		SortBy: params.ByName("sort_by"),
// 	}

// 	criteria := utils.BuildFilterCondition(queryParams)

// 	paginatedData, totalPages, totalRows := r.DiscountPercentService.GetAllDiscountPercent(criteria, paginate)

// 	payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(paginatedData), "success", 200, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
// }

// // // @Summary Get Discount Percent By ID
// // // @Description REST API Discount Percent
// // // @Accept json
// // // @Produce json
// // // @Tags Master : Discount Percent
// // // @Param discount_percent_id path int true "discount_percent_id"
// // // @Success 200 {object} payloads.Response
// // // @Failure 500,400,401,404,403,422 {object} exceptions.Error
// // // @Router /aftersales-service/api/aftersales/discount-percent/{discount_percent_id} [get]
// // func (r *DiscountPercentController) GetDiscountPercentByID(c *gin.Context) {
// // 	trxHandle := c.MustGet("db_trx").(*gorm.DB)
// // 	discountPercentId, _ := strconv.Atoi(c.Param("discount_percent_id"))
// // 	result, err := r.discountpercentservice.WithTrx(trxHandle).GetDiscountPercentById(int(discountPercentId))
// // 	if err != nil {
// // 		exceptions.NotFoundException(c, err.Error())
// // 		return
// // 	}

// // 	payloads.HandleSuccess(c, result, "Get Data Successfully!", http.StatusOK)
// // }

// // // @Summary Save Discount Percent
// // // @Description REST API Discount Percent
// // // @Accept json
// // // @Produce json
// // // @Tags Master : Discount Percent
// // // @param reqBody body masteritempayloads.DiscountPercentResponse true "Form Request"
// // // @Success 200 {object} payloads.Response
// // // @Failure 500,400,401,404,403,422 {object} exceptions.Error
// // // @Router /aftersales-service/api/aftersales/discount-percent [post]
// // func (r *DiscountPercentController) SaveDiscountPercent(c *gin.Context) {
// // 	trxHandle := c.MustGet("db_trx").(*gorm.DB)
// // 	var request masteritempayloads.DiscountPercentResponse
// // 	var message = ""

// // 	if err := c.ShouldBindJSON(&request); err != nil {
// // 		exceptions.EntityException(c, err.Error())
// // 		return
// // 	}

// // 	if int(request.DiscountPercentId) != 0 {
// // 		result, err := r.discountpercentservice.WithTrx(trxHandle).GetDiscountPercentById(int(request.DiscountPercentId))

// // 		if err != nil {
// // 			exceptions.AppException(c, err.Error())
// // 			return
// // 		}

// // 		if result.DiscountPercentId == 0 {
// // 			exceptions.NotFoundException(c, err.Error())
// // 			return
// // 		}
// // 	}

// // 	create, err := r.discountpercentservice.WithTrx(trxHandle).SaveDiscountPercent(request)
// // 	if err != nil {
// // 		exceptions.AppException(c, err.Error())
// // 		return
// // 	}

// // 	if request.DiscountPercentId == 0 {
// // 		message = "Create Data Successfully!"
// // 	} else {
// // 		message = "Update Data Successfully!"
// // 	}

// // 	payloads.HandleSuccess(c, create, message, http.StatusOK)
// // }

// // // @Summary Change Status Discount Percent
// // // @Description REST API Discount Percent
// // // @Accept json
// // // @Produce json
// // // @Tags Master : Discount Percent
// // // @param discount_percent_id path int true "discount_percent_id"
// // // @Success 200 {object} payloads.Response
// // // @Failure 500,400,401,404,403,422 {object} exceptions.Error
// // // @Router /aftersales-service/api/aftersales/discount-percent/{discount_percent_id} [patch]
// // func (r *DiscountPercentController) ChangeStatusDiscountPercent(c *gin.Context) {
// // 	trxHandle := c.MustGet("db_trx").(*gorm.DB)
// // 	discountPercentId, err := strconv.Atoi(c.Param("discount_percent_id"))
// // 	if err != nil {
// // 		exceptions.EntityException(c, err.Error())
// // 		return
// // 	}
// // 	//id check
// // 	result, err := r.discountpercentservice.WithTrx(trxHandle).GetDiscountPercentById(int(discountPercentId))
// // 	if err != nil || result.DiscountPercentId == 0 {
// // 		exceptions.NotFoundException(c, err.Error())
// // 		return
// // 	}

// // 	response, err := r.discountpercentservice.WithTrx(trxHandle).ChangeStatusDiscountPercent(int(discountPercentId))
// // 	if err != nil {
// // 		exceptions.AppException(c, err.Error())
// // 		return
// // 	}

// // 	payloads.HandleSuccess(c, response, "Update Data Successfully!", http.StatusOK)
// // }
