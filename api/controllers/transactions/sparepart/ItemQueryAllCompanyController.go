package transactionsparepartcontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads"
	"after-sales/api/payloads/pagination"
	transactionsparepartservice "after-sales/api/services/transaction/sparepart"
	"after-sales/api/utils"
	"bytes"
	"net/http"
	"time"
)

type ItemQueryAllCompanyControllerImpl struct {
	ItemQueryAllCompanyService transactionsparepartservice.ItemQueryAllCompanyService
}

type ItemQueryAllCompanyController interface {
	GetAllItemQueryAllCompany(writer http.ResponseWriter, request *http.Request)
	GetItemQueryAllCompanyDownload(writer http.ResponseWriter, request *http.Request)
}

func NewItemQueryAllCompanyController(
	itemQueryAllCompanyService transactionsparepartservice.ItemQueryAllCompanyService,
) ItemQueryAllCompanyController {
	return &ItemQueryAllCompanyControllerImpl{
		ItemQueryAllCompanyService: itemQueryAllCompanyService,
	}
}

// @Summary Get All Item Query All Company
// @Description Get All Item Query All Company
// @Tags Transaction : Sparepart Item Query All Company
// @Security AuthorizationKeyAuth
// @Accept json
// @Produce json
// @Param company_id query string false "Company ID"
// @Param brand_id query string false "Brand ID"
// @Param item_code_1 query string false "Item Code 1"
// @Param item_code_2 query string false "Item Code 2"
// @Param item_code_3 query string false "Item Code 3"
// @Param item_code_4 query string false "Item Code 4"
// @Param moving_code_1 query string false "Moving Code 1"
// @Param moving_code_2 query string false "Moving Code 2"
// @Param moving_code_3 query string false "Moving Code 3"
// @Param moving_code_4 query string false "Moving Code 4"
// @Param moving_code_5 query string false "Moving Code 5"
// @Param moving_code_6 query string false "Moving Code 6"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Param sort_of query string false "Sort Of"
// @Param sort_by query string false "Sort By"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-query-all-company [get]
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

// DownloadTemplate godoc
// @Summary Download Item Query All Company
// @Description REST API Download Template
// @Security AuthorizationKeyAuth
// @Accept json
// @Produce json
// @Tags Transaction : Sparepart Item Query All Company
// @Param company_id query int true "Company ID"
// @Param item_code_1 query string true "Item Code 1"
// @Param item_code_2 query string false "Item Code 2"
// @Param item_code_3 query string false "Item Code 3"
// @Param item_code_4 query string false "Item Code 4"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-query-all-company/download [get]
func (c *ItemQueryAllCompanyControllerImpl) GetItemQueryAllCompanyDownload(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"company_id":  queryValues.Get("company_id"),
		"item_code_1": queryValues.Get("item_code_1"),
		"item_code_2": queryValues.Get("item_code_2"),
		"item_code_3": queryValues.Get("item_code_3"),
		"item_code_4": queryValues.Get("item_code_4"),
	}

	if queryParams["company_id"] == "" {
		payloads.NewHandleError(writer, "company_id cannot be empty!", http.StatusBadRequest)
		return
	}
	if queryParams["item_code_1"] == "" {
		payloads.NewHandleError(writer, "item_code_1 cannot be empty!", http.StatusBadRequest)
		return
	}

	criteria := utils.BuildFilterCondition(queryParams)
	f, err := c.ItemQueryAllCompanyService.GetItemQueryAllCompanyDownload(criteria)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	var b bytes.Buffer
	if err := f.Write(&b); err != nil {
		// Create BaseErrorResponse for file write error
		baseErr := &exceptions.BaseErrorResponse{
			Err:        err,
			StatusCode: http.StatusInternalServerError,
		}
		exceptions.NewAppException(writer, request, baseErr)
		return
	}

	downloadName := time.Now().Local().Format("2006-01-02_15-04-05") + "_Item_Query_All_Company.xlsx"
	writer.Header().Set("Content-Description", "File Transfer")
	writer.Header().Set("Content-Disposition", "attachment; filename="+downloadName)
	writer.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	writer.Header().Set("Content-Transfer-Encoding", "binary")
	writer.Header().Set("Expires", "0")
	writer.Header().Set("Cache-Control", "must-revalidate")
	writer.Header().Set("Pragma", "public")

	// Write the buffer to the HTTP response
	_, writeErr := writer.Write(b.Bytes())
	if writeErr != nil {
		// Create BaseErrorResponse for writer.Write error
		baseErr := &exceptions.BaseErrorResponse{
			Err:        writeErr,
			StatusCode: http.StatusInternalServerError,
		}
		// Use a generic error handling function to respond with the error
		exceptions.NewAppException(writer, request, baseErr)
		return
	}
}
