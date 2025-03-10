package mastercontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterservice "after-sales/api/services/master"
	"after-sales/api/utils"
	"after-sales/api/validation"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type PointProspectingController interface {
	CreatePointProspecting(writer http.ResponseWriter, request *http.Request)
	UpdatePointProspectingStatus(writer http.ResponseWriter, request *http.Request)
	GetAllPointProspecting(writer http.ResponseWriter, request *http.Request)
	UpdatePointProspectingData(writer http.ResponseWriter, request *http.Request)
	GetOnePointProspecting(writer http.ResponseWriter, request *http.Request)
}

type PointProspectingControllerImpl struct {
	PointProspectingService masterservice.PointProspectingService
}

func NewPointProspectingControllerImpl(service masterservice.PointProspectingService) PointProspectingController {
	return &PointProspectingControllerImpl{
		PointProspectingService: service,
	}
}

func (c *PointProspectingControllerImpl) CreatePointProspecting(writer http.ResponseWriter, request *http.Request) {
	var pointProspecting masterpayloads.PointProspectingRequest
	helper.ReadFromRequestBody(request, &pointProspecting)

	if validationErr := validation.ValidationForm(writer, request, &pointProspecting); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}

	isTrue, err := c.PointProspectingService.CreatePointProspecting(pointProspecting)
	if err != nil {
		exceptions.NewConflictException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, isTrue, "Point Prospecting created successfully", http.StatusCreated)

}

func (c *PointProspectingControllerImpl) UpdatePointProspectingStatus(writer http.ResponseWriter, request *http.Request) {
	pointVariable := chi.URLParam(request, "pointVariable")
	pointValue, errA := strconv.Atoi(chi.URLParam(request, "pointValue"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("check your params"),
		})
		return
	}

	var pointProspecting masterpayloads.PointProspectingUpdateStatus
	helper.ReadFromRequestBody(request, &pointProspecting)

	res, err := c.PointProspectingService.UpdatePointProspectingStatus(pointVariable, pointValue, pointProspecting)
	if err != nil {
		exceptions.NewConflictException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, res, "Point Prospecting updated successfully", http.StatusOK)
}

func (c *PointProspectingControllerImpl) GetAllPointProspecting(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	filterCondition := map[string]string{
		"RECORD_STATUS":  queryValues.Get("RECORD_STATUS"),
		"POINT_VARIABLE": queryValues.Get("POINT_VARIABLE"),
		"POINT_VALUE":    queryValues.Get("POINT_VALUE"),
	}

	paginate := pagination.Pagination{
		Limit: utils.NewGetQueryInt(queryValues, "limit"),
		Page:  utils.NewGetQueryInt(queryValues, "pages"),
	}

	filterConds := utils.BuildFilterCondition(filterCondition)

	res, err := c.PointProspectingService.GetAllPointProspecting(filterConds, paginate)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(
		writer,
		res.Rows,
		"Point Prospecting data retrieved successfully",
		http.StatusOK,
		res.Limit,
		res.Page,
		int64(res.TotalRows),
		res.TotalPages,
	)
	// fmt.Println(res.TotalRows)

}

func (c *PointProspectingControllerImpl) UpdatePointProspectingData(writer http.ResponseWriter, request *http.Request) {
	pointVariable := chi.URLParam(request, "pointVariable")
	pointValue, errA := strconv.Atoi(chi.URLParam(request, "pointValue"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("check your params"),
		})
		return
	}

	var pointProspecting masterpayloads.PointProspectingUpdateRequest
	helper.ReadFromRequestBody(request, &pointProspecting)

	res, err := c.PointProspectingService.UpdatePointProspectingData(pointVariable, pointValue, pointProspecting)
	if err != nil {
		exceptions.NewConflictException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, res, "Point Prospecting updated successfully", http.StatusOK)
}

func (c *PointProspectingControllerImpl) GetOnePointProspecting(writer http.ResponseWriter, request *http.Request) {
	pointVariable := chi.URLParam(request, "pointVariable")
	pointValue, errA := strconv.Atoi(chi.URLParam(request, "pointValue"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("check your params"),
		})
		return
	}

	res, err := c.PointProspectingService.GetOnePointProspecting(pointVariable, pointValue)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, res, "Point Prospecting data retrieved successfully", http.StatusOK)
}
