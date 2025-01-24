package transactionworkshopcontroller

import (
	"after-sales/api/payloads/pagination"
	transactionworkshopservice "after-sales/api/services/transaction/workshop"
	"after-sales/api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PrintGatePassController struct {
	PrintGatePassService transactionworkshopservice.PrintGatePassService
}

func NewPrintGatePassController(service transactionworkshopservice.PrintGatePassService) *PrintGatePassController {
	return &PrintGatePassController{
		PrintGatePassService: service,
	}
}

// GetAll handles the request to get all Gate Pass data with filtering and pagination
func (ctrl *PrintGatePassController) GetAll(c *gin.Context) {
	// Parsing query params for filters and pagination
	workOrderNo := c.DefaultQuery("work_order_no", "")     // Filter by Work Order No.
	workOrderDate := c.DefaultQuery("work_order_date", "") // Filter by Work Order Date
	customer := c.DefaultQuery("customer", "")             // Filter by Customer
	noPolisi := c.DefaultQuery("no_polisi", "")            // Filter by No. Polisi
	printOption := c.DefaultQuery("print_option", "")      // Filter by Print Option

	// Parsing pagination params
	page := pagination.Pagination{}
	if err := c.ShouldBindQuery(&page); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pagination parameters"})
		return
	}

	// Create filter conditions based on provided query parameters
	filterConditions := []utils.FilterCondition{}

	if workOrderNo != "" {
		filterConditions = append(filterConditions, utils.FilterCondition{
			Field: "wo_sys_no", // Field name in the database
			Value: workOrderNo,
		})
	}

	if workOrderDate != "" {
		filterConditions = append(filterConditions, utils.FilterCondition{
			Field: "wo_date", // Field name for Work Order Date
			Value: workOrderDate,
		})
	}

	if customer != "" {
		filterConditions = append(filterConditions, utils.FilterCondition{
			Field: "customer_id", // Customer ID field
			Value: customer,
		})
	}

	if noPolisi != "" {
		filterConditions = append(filterConditions, utils.FilterCondition{
			Field: "tnkb", // No. Polisi field (TNKB)
			Value: noPolisi,
		})
	}

	if printOption != "" {
		filterConditions = append(filterConditions, utils.FilterCondition{
			Field: "approval_gatepass_status", // Approval Gatepass Status field
			Value: printOption,
		})
	}

	// Calling the service to fetch the data
	paginatedResult, errResponse := ctrl.PrintGatePassService.GetAll(filterConditions, page)
	if errResponse != nil {
		c.JSON(errResponse.StatusCode, gin.H{"error": errResponse.Message})
		return
	}

	// Return the paginated results
	c.JSON(http.StatusOK, gin.H{
		"data":  paginatedResult.Rows,
		"total": paginatedResult.TotalRows,
		"page":  paginatedResult.Page,
		"limit": paginatedResult.Limit,
	})
}
