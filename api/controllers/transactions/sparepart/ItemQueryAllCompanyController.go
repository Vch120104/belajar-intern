package transactionsparepartcontroller

import (
	"after-sales/api/helper"
	"after-sales/api/payloads"
	"after-sales/api/payloads/pagination"
	transactionsparepartservice "after-sales/api/services/transaction/sparepart"
	"after-sales/api/utils"
	"net/http"
) 

type ItemQueryAllCompanyControllerImpl struct {
	ItemQueryAllCompanyService transactionsparepartservice.ItemQueryAllCompanyService
}

type ItemQueryAllCompanyController interface {
	GetAllItemQueryAllCompany(writer http.ResponseWriter, request *http.Request)
}

func NewItemQueryAllCompanyController(
	itemQueryAllCompanyService transactionsparepartservice.ItemQueryAllCompanyService,
) ItemQueryAllCompanyController {
	return &ItemQueryAllCompanyControllerImpl{
		ItemQueryAllCompanyService: itemQueryAllCompanyService,
	}
}

func (c *ItemQueryAllCompanyControllerImpl) GetAllItemQueryAllCompany(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"company_id":    queryValues.Get("company_id"),
		"brand_id":      queryValues.Get("brand_id"),
		"item_code_1":   queryValues.Get("item_code_1"),
		"item_code_2":   queryValues.Get("item_code_2"),
		"item_code_3":   queryValues.Get("item_code_3"),
		"item_code_4":   queryValues.Get("item_code_4"),
		"moving_code_1": queryValues.Get("moving_code_1"),
		"moving_code_2": queryValues.Get("moving_code_2"),
		"moving_code_3": queryValues.Get("moving_code_3"),
		"moving_code_4": queryValues.Get("moving_code_4"),
		"moving_code_5": queryValues.Get("moving_code_5"),
		"moving_code_6": queryValues.Get("moving_code_6"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	response, err := c.ItemQueryAllCompanyService.GetAllItemQueryAllCompany(criteria, paginate)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(
		writer,
		response.Rows,
		"Get Data Successfully!",
		http.StatusOK,
		response.Limit,
		response.Page,
		int64(response.TotalRows),
		response.TotalPages,
	)
}
