package transactionjpcbrepositoryimpl

import (
	"after-sales/api/config"
	transactionjpcbentities "after-sales/api/entities/transaction/JPCB"
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionjpcbpayloads "after-sales/api/payloads/transaction/JPCB"
	transactionjpcbrepository "after-sales/api/repositories/transaction/JPCB"
	"after-sales/api/utils"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type CarWashImpl struct{}

func NewCarWashRepositoryImpl() transactionjpcbrepository.CarWashRepository {
	return &CarWashImpl{}
}

func (*CarWashImpl) GetAll(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	// select field that is missing rn : model_color, model, color, tnkb, creation_time
	joinQuery := tx.Table("trx_car_wash").
		Select(`trx_work_order.work_order_system_number, trx_work_order.work_order_document_number, trx_work_order.model_id, trx_work_order.vehicle_id,
				trx_work_order.promise_time, trx_work_order.promise_date, trx_car_wash.car_wash_bay_id, mtr_car_wash_bay.car_wash_bay_description,trx_car_wash.car_wash_status_id, 
				mtr_car_wash_status.car_wash_status_description, trx_car_wash.start_time, trx_car_wash.end_time, trx_car_wash.car_wash_priority_id, 
				mtr_car_wash_priority.car_wash_priority_description`).
		Joins("LEFT JOIN trx_work_order ON trx_car_wash.work_order_system_number = trx_work_order.work_order_system_number AND trx_car_wash.company_id = trx_work_order.company_id").
		Joins("LEFT JOIN mtr_car_wash_priority ON trx_car_wash.car_wash_priority_id = mtr_car_wash_priority.car_wash_priority_id").
		Joins("LEFT JOIN mtr_car_wash_status ON trx_car_wash.car_wash_status_id = mtr_car_wash_status.car_wash_status_id").
		Joins("LEFT JOIN mtr_car_wash_bay ON trx_car_wash.car_wash_bay_id = mtr_car_wash_bay.car_wash_bay_id")

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
			&carWashPayload.CarWashBayId, &carWashPayload.CarWashBayDescription, &carWashPayload.CarWashStatusId, &carWashPayload.CarWashStatusDescription, &carWashPayload.StartTime,
			&carWashPayload.EndTime, &carWashPayload.CarWashPriorityId, &carWashPayload.CarWashPriorityDescription,
		)
		if err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}

		//Fetch data Model from external services
		ModelURL := config.EnvConfigs.SalesServiceUrl + "unit-model/" + strconv.Itoa(modelId)
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
			CarWashBayDescription:      carWashPayload.CarWashBayDescription,
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
			"CarWashBayDescription":      response.CarWashBayDescription,
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
		Where("work_order_system_number = ? AND car_wash_status_id = 2", workOrderSystemNumber).Find(&carWashEntities)
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

func (r *CarWashImpl) GetAllCarWashPriorityDropDown(tx *gorm.DB) ([]transactionjpcbpayloads.CarWashPriorityDropDownResponse, *exceptions.BaseErrorResponse) {
	var entities transactionjpcbentities.CarWashPriority
	var responses []transactionjpcbpayloads.CarWashPriorityDropDownResponse
	rows, err := tx.Model(&entities).Rows()
	if err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	defer rows.Close()

	for rows.Next() {
		var carWashPriorityId int
		var carWashPriorityCode, carWashPriorityDescription string
		var isActive bool

		err := rows.Scan(&isActive, &carWashPriorityId, &carWashPriorityCode, &carWashPriorityDescription)
		if err != nil {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}

		response := transactionjpcbpayloads.CarWashPriorityDropDownResponse{
			CarWashPriorityId:          carWashPriorityId,
			CarWashPriorityDescription: carWashPriorityDescription,
			CarWashPriorityCode:        carWashPriorityCode,
			IsActive:                   isActive,
		}
		responses = append(responses, response)
	}

	return responses, nil
}

func (r *CarWashImpl) DeleteCarWash(tx *gorm.DB, workOrderSystemNumber int) (bool, *exceptions.BaseErrorResponse) {
	mainTable := "trx_car_wash"
	mainAlias := "carwash"

	joinTables := []utils.JoinTable{
		{Table: "mtr_car_wash_bay", Alias: "bay", ForeignKey: mainAlias + ".car_wash_bay_id", ReferenceKey: "bay.car_wash_bay_id"},
		{Table: "mtr_car_wash_status", Alias: "status", ForeignKey: mainAlias + ".car_wash_status_id", ReferenceKey: "status.car_wash_status_id"},
		{Table: "trx_work_order", Alias: "wo", ForeignKey: mainAlias + ".work_order_system_number", ReferenceKey: "wo.work_order_system_number"},
	}

	joinQuery := utils.CreateJoin(tx, mainTable, mainAlias, joinTables...)

	keyAttributes := []string{
		"wo.work_order_system_number",
		"wo.work_order_document_number",
		"bay.car_wash_bay_description",
		"carwash.car_wash_status_id",
		"status.car_wash_status_description",
	}

	var result transactionjpcbpayloads.CarWashErrorDetail
	joinQuery = joinQuery.Select(keyAttributes).Where("carwash.work_order_system_number = ?", workOrderSystemNumber).
		Scan(&result)
	if joinQuery.Error != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        joinQuery.Error,
		}
	}

	if result.WorkOrderSystemNumber == 0 {
		return false, nil
	}

	DRAFT := 1
	if result.CarWashStatusId != DRAFT {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusOK,
			Message: "Can't delete Car Wash Allocation, Car Wash status for " + result.WorkOrderDocumentNumber + " : " +
				result.CarWashBayDescription + " is already " + result.CarWashStatusDescription,
			Err: fmt.Errorf(
				"failed to delete car wash",
			),
		}
	}

	var carWashEntity transactionjpcbentities.CarWash

	err := tx.Model(&carWashEntity).Where("work_order_system_number = ?", workOrderSystemNumber).First(&carWashEntity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}
	err = tx.Delete(&carWashEntity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			Message: "Failed to delete car wash",
			Err:     err,
		}
	}

	var workOrderEntity transactionworkshopentities.WorkOrder
	whereQuery := tx.Model(&workOrderEntity).Where("work_order_system_number = ?", workOrderSystemNumber).First(&workOrderEntity)
	if whereQuery.Error != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}
	whereQuery = whereQuery.Update("car_wash", 0)
	if whereQuery.Error != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}
	return true, nil
}

func (r *CarWashImpl) PostCarWash(tx *gorm.DB, workOrderSystemNumber int) (transactionjpcbpayloads.CarWashPostResponse, *exceptions.BaseErrorResponse) {
	const (
		qcPass         = 6
		statusDraft    = 1
		priorityNormal = 2
	)

	var workOrderEntity transactionworkshopentities.WorkOrder
	var workOrderResponse transactionjpcbpayloads.CarWashWorkOrder

	err := tx.Model(&workOrderEntity).Select("car_wash, company_id, work_order_status_id").Where("work_order_system_number = ?", workOrderSystemNumber).Scan(&workOrderResponse).Error
	if err != nil {
		errorHelperInternalServerError(err)
	}

	//Fetch data Model from external services
	CompanyURL := config.EnvConfigs.GeneralServiceUrl + "company-detail/" + strconv.Itoa(workOrderResponse.CompanyId)
	var getCompanyResponse transactionjpcbpayloads.CarWashCompanyResponse
	errFetchCompany := utils.Get(CompanyURL, &getCompanyResponse, nil)
	getCompanyResponse.IsUseJPCB = true //TODO remove later
	if errFetchCompany != nil {
		return transactionjpcbpayloads.CarWashPostResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Company Data from external service",
			Err:        err,
		}
	}

	if getCompanyResponse.IsUseJPCB { //TODO check if company use jpcb
		if workOrderResponse.WorkOrderStatusId == qcPass {
			if workOrderResponse.CarWash {
				var workOrder int
				result := tx.Model(&transactionjpcbentities.CarWash{}).Select("work_order_system_number").
					Where("work_order_system_number = ?", workOrderSystemNumber).Scan(&workOrder)
				if result.Error != nil {
					errorHelperInternalServerError(result.Error)
				}
				if result.RowsAffected == 0 {
					newCarWash := transactionjpcbentities.CarWash{
						CompanyId:             workOrderResponse.CompanyId,
						WorkOrderSystemNumber: workOrderSystemNumber,
						StatusId:              statusDraft,
						PriorityId:            priorityNormal,
						CarWashDate:           time.Now(),
						BayId:                 nil,
					}

					err := tx.Create(&newCarWash).Error
					if err != nil {
						return transactionjpcbpayloads.CarWashPostResponse{}, &exceptions.BaseErrorResponse{
							StatusCode: http.StatusOK,
							Err:        err,
						}
					}

					return transactionjpcbpayloads.CarWashPostResponse{
						CarWashId:             newCarWash.CarWashId,
						CompanyId:             newCarWash.CompanyId,
						WorkOrderSystemNumber: newCarWash.WorkOrderSystemNumber,
						BayId:                 newCarWash.BayId,
						StatusId:              newCarWash.StatusId,
						PriorityId:            newCarWash.PriorityId,
						CarWashDate:           newCarWash.CarWashDate,
						StartTime:             newCarWash.StartTime,
						EndTime:               newCarWash.EndTime,
						ActualTime:            newCarWash.ActualTime,
					}, nil
				}

				return errorHelperDataAlreadyExist()
			}
		} else {
			if workOrderResponse.CarWash {
				lineTypeOperationId := 1
				workOrderDetail := transactionworkshopentities.WorkOrderDetail{}

				resultLineTypeOperation := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
					Where("work_order_system_number = ? AND line_type_id = ?", workOrderSystemNumber, lineTypeOperationId).First(&workOrderDetail)
				if resultLineTypeOperation.Error != nil {
					errorHelperInternalServerError(resultLineTypeOperation.Error)
				}

				if workOrderDetail.WorkOrderSystemNumber == 0 {
					result := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
						Where("work_order_system_number = ? AND frt_quantity <> supply_quantity", workOrderSystemNumber).First(&workOrderDetail)
					if result.Error != nil {
						errorHelperInternalServerError(result.Error)
					}

					if workOrderDetail.WorkOrderSystemNumber == 0 {
						carWash := transactionjpcbentities.CarWash{}
						result := tx.Model(&transactionjpcbentities.CarWash{}).
							Where("work_order_system_number = ?", workOrderSystemNumber).First(&carWash)
						if result.Error != nil {
							errorHelperInternalServerError(result.Error)
						}

						if carWash.CarWashId == 0 {
							newCarWash := transactionjpcbentities.CarWash{
								CompanyId:             workOrderResponse.CompanyId,
								WorkOrderSystemNumber: workOrderSystemNumber,
								StatusId:              statusDraft,
								PriorityId:            priorityNormal,
								CarWashDate:           time.Now(),
								BayId:                 nil,
							}

							err := tx.Create(&newCarWash).Error
							if err != nil {
								return transactionjpcbpayloads.CarWashPostResponse{}, &exceptions.BaseErrorResponse{
									StatusCode: http.StatusOK,
									Err:        err,
								}
							}

							return transactionjpcbpayloads.CarWashPostResponse{
								CarWashId:             newCarWash.CarWashId,
								CompanyId:             newCarWash.CompanyId,
								WorkOrderSystemNumber: newCarWash.WorkOrderSystemNumber,
								BayId:                 newCarWash.BayId,
								StatusId:              newCarWash.StatusId,
								PriorityId:            newCarWash.PriorityId,
								CarWashDate:           newCarWash.CarWashDate,
								StartTime:             newCarWash.StartTime,
								EndTime:               newCarWash.EndTime,
								ActualTime:            newCarWash.ActualTime,
							}, nil
						}

						return errorHelperDataAlreadyExist()
					}
				} else {
					var deleteCarWash transactionjpcbentities.CarWash

					result := tx.Model(&deleteCarWash).Select("work_order_system_number").
						Where("work_order_system_number = ?", workOrderSystemNumber).First(&deleteCarWash)
					if result.Error != nil {
						errorHelperInternalServerError(result.Error)
					}

					if deleteCarWash.CarWashId != 0 {
						err = tx.Delete(&deleteCarWash).Error
						if err != nil {
							return transactionjpcbpayloads.CarWashPostResponse{}, &exceptions.BaseErrorResponse{
								Message: "Failed to delete car wash",
								Err:     err,
							}
						}
					}

					return transactionjpcbpayloads.CarWashPostResponse{}, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusBadRequest,
						Message:    "Can't add work order with line type operation",
						Err:        fmt.Errorf("cant add work order with line type operation"),
					}
				}
			}
		}
	}

	return errorHelperBadRequest()
}

func errorHelperDataAlreadyExist() (transactionjpcbpayloads.CarWashPostResponse, *exceptions.BaseErrorResponse) {
	return transactionjpcbpayloads.CarWashPostResponse{}, &exceptions.BaseErrorResponse{
		StatusCode: http.StatusBadRequest,
		Message:    "data already exist",
		Err:        fmt.Errorf("data already exist"),
	}
}

func errorHelperBadRequest() (transactionjpcbpayloads.CarWashPostResponse, *exceptions.BaseErrorResponse) {
	return transactionjpcbpayloads.CarWashPostResponse{}, &exceptions.BaseErrorResponse{
		StatusCode: http.StatusBadRequest,
		Message:    "Failed to create car wash",
		Err:        fmt.Errorf("fail to create car wash"),
	}
}

func errorHelperInternalServerError(err error) (transactionjpcbpayloads.CarWashPostResponse, *exceptions.BaseErrorResponse) {
	return transactionjpcbpayloads.CarWashPostResponse{}, &exceptions.BaseErrorResponse{
		StatusCode: http.StatusInternalServerError,
		Err:        err,
	}
}

func errorHelperNotFound(err error) (transactionjpcbpayloads.CarWashPostResponse, *exceptions.BaseErrorResponse) {
	return transactionjpcbpayloads.CarWashPostResponse{}, &exceptions.BaseErrorResponse{
		StatusCode: http.StatusNotFound,
		Err:        err,
	}
}
