package transactionworkshopcontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/payloads"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshopservice "after-sales/api/services/transaction/workshop"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"gorm.io/gorm"
)

type BookingEstimationController struct {
	bookingEstimationService transactionworkshopservice.BookingEstimationService
}

func OpenBookingEstimationRoutes(
	db *gorm.DB,
	r chi.Router,
	bookingEstimationService transactionworkshopservice.BookingEstimationService,
) {
	handler := BookingEstimationController{
		bookingEstimationService: bookingEstimationService,
	}

	r.Post("/save-booking-estimation", func(w http.ResponseWriter, req *http.Request) {
		handler.Save(w, req, db)
	})
}

// Save Booking Estimation
func (r *BookingEstimationController) Save(w http.ResponseWriter, req *http.Request, db *gorm.DB) {
	trxHandle := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			trxHandle.Rollback()
		}
	}()

	var requestBody transactionworkshoppayloads.SaveBookingEstimationRequest

	if err := req.ParseForm(); err != nil {
		exceptions.EntityException(w, err.Error())
		return
	}

	if err := req.ParseMultipartForm(10 << 20); err != nil {
		exceptions.EntityException(w, err.Error())
		return
	}

	if err := render.Decode(req, &requestBody); err != nil {
		exceptions.EntityException(w, err.Error())
		return
	}

	save, err := r.bookingEstimationService.WithTrx(trxHandle).Save(requestBody)
	if err != nil {
		exceptions.AppException(w, err.Error())
		return
	}

	if err := trxHandle.Commit().Error; err != nil {
		exceptions.AppException(w, err.Error())
		return
	}

	payloads.NewHandleSuccess(w, save, "Insert Successfully", http.StatusOK)
}
