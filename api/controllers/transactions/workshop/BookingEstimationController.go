package transactionworkshopcontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/middlewares"
	"after-sales/api/payloads"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshopservice "after-sales/api/services/transaction/workshop"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BookingEstimationController struct {
	bookingEstimationService transactionworkshopservice.BookingEstimationService
}

func OpenBookingEstimationRoutes(
	db *gorm.DB,
	r *gin.RouterGroup,
	bookingEstimationService transactionworkshopservice.BookingEstimationService,
) {
	handler := BookingEstimationController{
		bookingEstimationService: bookingEstimationService,
	}

	r.POST("/save-booking-estimation", middlewares.DBTransactionMiddleware(db), handler.Save)
}

// @Summary Save Booking Estimation
// @Description Save Booking Estimation
// @Accept json
// @Produce json
// @Tags Master : Booking Estimation
// @Security BearerAuth
// @param reqBody body transactionworkshoppayloads.SaveBookingEstimationRequest true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/save-booking-estimation [post]
func (r *BookingEstimationController) Save(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	requestBody := transactionworkshoppayloads.SaveBookingEstimationRequest{}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		exceptions.EntityException(c, err.Error())
		return
	}

	save, err := r.bookingEstimationService.WithTrx(trxHandle).Save(requestBody)

	if err != nil {
		exceptions.AppException(c, err.Error())
		return
	}

	payloads.HandleSuccess(c, save, "Insert Successfully", http.StatusOK)

}
