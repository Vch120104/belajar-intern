package transactionsparepartcontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/payloads"
	transactionsparepartservice "after-sales/api/services/transaction/sparepart"

	// "after-sales/api/middlewares"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SupplySlipController struct {
	supplyslipservice transactionsparepartservice.SupplySlipService
}

func StartSupplySlipRoutes(
	db *gorm.DB,
	r *gin.RouterGroup,
	supplyslipservice transactionsparepartservice.SupplySlipService,
) {
	supplySlipHandler := SupplySlipController{supplyslipservice: supplyslipservice}
	r.GET("/supply-slip/:supply_system_number", supplySlipHandler.GetSupplySlipByID)
}

// @Summary Get Supply Slip By ID
// @Description REST API Supply Slip
// @Accept json
// @Produce json
// @Tags Transaction : Supply Slip
// @Param supply_system_number path int true "supply_system_number"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/supply-slip/{supply_system_number} [get]
func (r *SupplySlipController) GetSupplySlipByID(c *gin.Context) {
	SupplySystemNumber, _ := strconv.Atoi(c.Param("supply_system_number"))
	result, err := r.supplyslipservice.GetSupplySlipById(int32(SupplySystemNumber))
	if err != nil {
		exceptions.NotFoundException(c, err.Error())
		return
	}
	payloads.HandleSuccess(c, result, "Get Data Successfully!", http.StatusOK)
}
