package transactionjpcbrepositoryimpl

import (
	"after-sales/api/config"
	transactionjpcbentities "after-sales/api/entities/transaction/JPCB"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionjpcbpayloads "after-sales/api/payloads/transaction/JPCB"
	transactionjpcbrepository "after-sales/api/repositories/transaction/JPCB"
	"after-sales/api/utils"
	"fmt"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type CarWashImpl struct{}

// GetById implements transactionjpcbrepository.CarWashRepository.
func (*CarWashImpl) GetById(tx *gorm.DB, id int) (transactionjpcbentities.CarWash, *exceptions.BaseErrorResponse) {
	panic("unimplemented")
}

func NewCarWashRepositoryImpl() transactionjpcbrepository.CarWashRepository {
	return &CarWashImpl{}
}

func (*CarWashImpl) GetAll(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	// select field that is missing rn : model_color, model, color, tnkb, creation_time
	joinQuery := tx.Table("trx_car_wash").
		Select(`trx_work_order.work_order_system_number, trx_work_order.work_order_document_number, trx_work_order.model_id, trx_work_order.vehicle_id,
				trx_work_order.promise_time, trx_work_order.promise_date, trx_car_wash.car_wash_bay_id, trx_car_wash.car_wash_status_id, mtr_car_wash_status.car_wash_status_description,
				trx_car_wash.start_time, trx_car_wash.end_time, trx_car_wash.car_wash_priority_id, mtr_car_wash_priority.car_wash_priority_description`).
		Joins("LEFT JOIN trx_work_order ON trx_car_wash.work_order_system_number = trx_work_order.work_order_system_number AND trx_car_wash.company_id = trx_work_order.company_id").
		Joins("LEFT JOIN mtr_car_wash_priority ON trx_car_wash.car_wash_priority_id = mtr_car_wash_priority.car_wash_priority_id").
		Joins("LEFT JOIN mtr_car_wash_status ON trx_car_wash.car_wash_status_id = mtr_car_wash_status.car_wash_status_id")

	joinQuery = utils.ApplyFilter(joinQuery, filterCondition)
	whereQuery := joinQuery.Where("trx_work_order.car_wash = 1 AND trx_car_wash.car_wash_status_id != 4")
	rows, err := whereQuery.Rows()
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	defer rows.Close()

	var convertedCarWashResponse []transactionjpcbpayloads.CarWashGetAllResponse
	for rows.Next() {
		var carWashPayload transactionjpcbpayloads.CarWashGetAllResponse
		var modelId int
		var vehicleId int

		err := rows.Scan(
			&carWashPayload.WorkOrderSystemNumber, &carWashPayload.WorkOrderDocumentNumber, &modelId, &vehicleId, &carWashPayload.PromiseTime, &carWashPayload.PromiseDate,
			&carWashPayload.CarWashBayId, &carWashPayload.CarWashStatusId, &carWashPayload.CarWashStatusDescription, &carWashPayload.StartTime, &carWashPayload.EndTime,
			&carWashPayload.CarWashPriorityId, &carWashPayload.CarWashPriorityDescription,
		)
		if err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}

		//Fetch data Model from external services
		ModelURL := config.EnvConfigs.SalesServiceUrl + "unit-model/" + strconv.Itoa(modelId)
		fmt.Println("Fetching Model data from:", ModelURL)

		var getModelResponse transactionjpcbpayloads.CarWashModelResponse
		errFetchModel := utils.Get(ModelURL, &getModelResponse, nil)
		if errFetchModel != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch brand data from external service",
				Err:        err,
			}
		}

		//Fetch data Color from vehicle master then unit color
		VehicleURL := config.EnvConfigs.SalesServiceUrl + "vehicle-master-by-chassis-number/" + strconv.Itoa(vehicleId)

		var getVehicleResponse transactionjpcbpayloads.CarWashVehicleResponse
		errFetchVehicle := utils.Get(VehicleURL, &getVehicleResponse, nil)
		if errFetchVehicle != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch vehicle data from external service",
				Err:        err,
			}
		}

		ColourUrl := config.EnvConfigs.SalesServiceUrl + "unit-colour/" + strconv.Itoa(getVehicleResponse.VehicleColourId)
		var getColourResponse transactionjpcbpayloads.CarWashColourResponse
		errFetchColour := utils.Get(ColourUrl, &getColourResponse, nil)
		if errFetchColour != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch colour data from external service",
			}
		}

		//Fetch tnkb, from mtr_vehicle_registration_certificate get by vehicle_id
		//TODO

		carWashResponse := transactionjpcbpayloads.CarWashGetAllResponse{
			WorkOrderSystemNumber:      carWashPayload.WorkOrderSystemNumber,
			WorkOrderDocumentNumber:    carWashPayload.WorkOrderDocumentNumber,
			Model:                      getModelResponse.ModelName,
			Color:                      getColourResponse.VariantColourName,
			Tnkb:                       "",
			PromiseTime:                carWashPayload.PromiseTime,
			PromiseDate:                carWashPayload.PromiseDate,
			CarWashBayId:               carWashPayload.CarWashBayId,
			CarWashStatusId:            carWashPayload.CarWashStatusId,
			CarWashStatusDescription:   carWashPayload.CarWashStatusDescription,
			StartTime:                  carWashPayload.StartTime,
			EndTime:                    carWashPayload.EndTime,
			CarWashPriorityId:          carWashPayload.CarWashPriorityId,
			CarWashPriorityDescription: carWashPayload.CarWashPriorityDescription,
		}
		convertedCarWashResponse = append(convertedCarWashResponse, carWashResponse)
	}

	var mapResponses []map[string]interface{}
	for _, response := range convertedCarWashResponse {
		responseMap := map[string]interface{}{
			"WorkOrderSystemNumber":      response.WorkOrderSystemNumber,
			"WorkOrderDocumentNumber":    response.WorkOrderDocumentNumber,
			"Model":                      response.Model,
			"Color":                      response.Color,
			"Tnkb":                       response.Tnkb,
			"PromiseTime":                response.PromiseTime,
			"PromiseDate":                response.PromiseDate,
			"CarWashBayId":               response.CarWashBayId,
			"CarWashStatusId":            response.CarWashStatusId,
			"CarWashStatusDescription":   response.CarWashStatusDescription,
			"StartTime":                  response.StartTime,
			"EndTime":                    response.EndTime,
			"CarWashPriorityId":          response.CarWashPriorityId,
			"CarWashPriorityDescription": response.CarWashPriorityDescription,
		}

		mapResponses = append(mapResponses, responseMap)
	}

	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(mapResponses, &pages)

	return paginatedData, totalPages, totalRows, nil
}

func (*CarWashImpl) UpdatePriority(tx *gorm.DB, workOrderSystemNumber, carWashPriorityId int) (transactionjpcbentities.CarWash, *exceptions.BaseErrorResponse) {
	var carWashEntities []transactionjpcbentities.CarWash

	checkBayStatusQuery := tx.Model(&carWashEntities).Select("car_wash_bay_id").
		Where("work_order_system_number = ? AND car_wash_status_id = 3", workOrderSystemNumber).Find(&carWashEntities)
	if checkBayStatusQuery.Error != nil {
		return transactionjpcbentities.CarWash{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        checkBayStatusQuery.Error,
		}
	}
	if len(carWashEntities) == 0 {
		checkBayAllocationQuery := tx.Model(&carWashEntities).Select("car_wash_bay_id").Where("work_order_system_number = ?", workOrderSystemNumber).Find(&carWashEntities)
		if checkBayAllocationQuery.Error != nil {
			return transactionjpcbentities.CarWash{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        checkBayAllocationQuery.Error,
			}
		}
		if len(carWashEntities) == 0 {
			return transactionjpcbentities.CarWash{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        fmt.Errorf("bay already allocated"),
			}
		}

		var resultEntities transactionjpcbentities.CarWash
		updateQuery := tx.Model(&resultEntities).Where("work_order_system_number = ?", workOrderSystemNumber).
			Updates(map[string]interface{}{"car_wash_priority_id": carWashPriorityId}).Preload("WorkOrder").Preload("CarWashBay").Preload("CarWashStatus").Preload("CarWashPriority").
			First(&resultEntities)
		if updateQuery.Error != nil {
			return transactionjpcbentities.CarWash{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        updateQuery.Error,
			}
		}

		return resultEntities, nil
	}

	return transactionjpcbentities.CarWash{}, &exceptions.BaseErrorResponse{
		StatusCode: http.StatusInternalServerError,
		Err:        fmt.Errorf("bay already started"),
	}
}
