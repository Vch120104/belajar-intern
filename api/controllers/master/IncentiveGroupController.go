package mastercontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/middlewares"
	"after-sales/api/payloads"
	masterservice "after-sales/api/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IncentiveGroupController struct {
	incentivegroupservice masterservice.IncentiveGroupService
}

func StartIncentiveGroupRoutes(
	db *gorm.DB,
	r *gin.RouterGroup,
	incentivegroupservice masterservice.IncentiveGroupService,
) {
	incentivegroupHandler := IncentiveGroupController{incentivegroupservice: incentivegroupservice}
	r.GET("/incentive-group/is-active", middlewares.DBTransactionMiddleware(db), incentivegroupHandler.GetAllIncentiveGroupIsActive)
}

// @Summary Get All Incentive Group Is Active
// @Description REST API Incentive Group
// @Accept json
// @Produce json
// @Tags Master : Incentive Group
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/incentive-group/is-active [get]
func (r *IncentiveGroupController) GetAllIncentiveGroupIsActive(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	result, err := r.incentivegroupservice.WithTrx(trxHandle).GetAllIncentiveGroupIsActive()
	if err != nil {
		exceptions.NotFoundException(c, err.Error())
		return
	}
	payloads.HandleSuccess(c, result, "Get Data Successfully!", http.StatusOK)
}
