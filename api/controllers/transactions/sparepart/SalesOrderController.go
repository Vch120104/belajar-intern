package transactionsparepartcontroller

import (
	transactionsparepartservice "after-sales/api/services/transaction/sparepart"
	"net/http"
)

// SalesOrderController
type SalesOrderController struct {
	salesOrderService transactionsparepartservice.SalesOrderService
}

// StartSalesOrderController
func NewSalesOrderController(salesOrderService transactionsparepartservice.SalesOrderService) SalesOrderController {
	return SalesOrderController{
		salesOrderService: salesOrderService,
	}
}

// GetSalesOrderByID retrieves a sales order by ID
// @Summary Get Sales Order By ID
// @Description Retrieve a sales order by its ID
// @Accept json
// @Produce json
// @Tags Transaction : Spare Part Sales Order
// @Param sales_order_id path int true "Sales Order ID"
// @Success 200 {object} payloads.Response
// @Failure 500,404 {object} exceptions.BaseErrorResponse
// @Router /v1/sales-order/{sales_order_id} [get]
func (c *SalesOrderController) GetSalesOrderByID(w http.ResponseWriter, r *http.Request) {

}
