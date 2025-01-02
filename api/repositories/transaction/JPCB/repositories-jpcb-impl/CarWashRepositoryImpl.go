package transactionjpcbrepositoryimpl

import (
	transactionjpcbentities "after-sales/api/entities/transaction/JPCB"
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionjpcbpayloads "after-sales/api/payloads/transaction/JPCB"
	transactionjpcbrepository "after-sales/api/repositories/transaction/JPCB"
	"after-sales/api/utils"
	generalserviceapiutils "after-sales/api/utils/general-service"
	salesserviceapiutils "after-sales/api/utils/sales-service"
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

func (*CarWashImpl) GetAll(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var responses []transactionjpcbpayloads.CarWashGetAllResponse

	baseModelQuery := tx.Table("trx_car_wash").
		Select(`trx_work_order.work_order_system_number, trx_work_order.work_order_document_number, trx_work_order.model_id, trx_work_order.vehicle_id,
				trx_work_order.promise_time, trx_work_order.promise_date, trx_car_wash.car_wash_bay_id, mtr_car_wash_bay.car_wash_bay_description,
				trx_car_wash.car_wash_status_id, mtr_car_wash_status.car_wash_status_description, trx_car_wash.start_time, trx_car_wash.end_time, 
				trx_car_wash.car_wash_priority_id, mtr_car_wash_priority.car_wash_priority_description`).
		Joins("LEFT JOIN trx_work_order ON trx_car_wash.work_order_system_number = trx_work_order.work_order_system_number AND trx_car_wash.company_id = trx_work_order.company_id").
		Joins("LEFT JOIN mtr_car_wash_priority ON trx_car_wash.car_wash_priority_id = mtr_car_wash_priority.car_wash_priority_id").
		Joins("LEFT JOIN mtr_car_wash_status ON trx_car_wash.car_wash_status_id = mtr_car_wash_status.car_wash_status_id").
		Joins("LEFT JOIN mtr_car_wash_bay ON trx_car_wash.car_wash_bay_id = mtr_car_wash_bay.car_wash_bay_id")

	whereQuery := utils.ApplyFilter(baseModelQuery, filterCondition)

	err := whereQuery.Scopes(pagination.Paginate(&pages, whereQuery)).Find(&responses).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if len(responses) == 0 {
		pages.Rows = []map[string]interface{}{}
		return pages, nil
	}

	var results []map[string]interface{}
	for _, response := range responses {
		// Fetch external data for Model, Vehicle, and Color
		getModelResponse, errFetchModel := salesserviceapiutils.GetUnitModelById(response.ModelId)
		if errFetchModel != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch unit model data from external service",
				Err:        errFetchModel,
			}
		}

		// getVehicleResponse, errFetchVehicle := salesserviceapiutils.GetVehicleById(response.VehicleId)
		// if errFetchVehicle != nil {
		// 	return pages, &exceptions.BaseErrorResponse{
		// 		StatusCode: http.StatusInternalServerError,
		// 		Message:    "Failed to fetch vehicle data from external service",
		// 		Err:        errFetchVehicle,
		// 	}
		// }

		result := map[string]interface{}{
			"work_order_system_number":   response.WorkOrderSystemNumber,
			"work_order_document_number": response.WorkOrderDocumentNumber,
			"model":                      getModelResponse.ModelName,
			// "color":                      getVehicleResponse.ColourCommercialName,
			"color": "",
			// "tnkb":                          getVehicleResponse.VehicleRegistrationCertificateTNKB,
			"tnkb":                          "",
			"promise_time":                  response.PromiseTime,
			"promise_date":                  response.PromiseDate,
			"car_wash_bay_id":               response.CarWashBayId,
			"car_wash_bay_description":      response.CarWashBayDescription,
			"car_wash_status_id":            response.CarWashStatusId,
			"car_wash_status_description":   response.CarWashStatusDescription,
			"start_time":                    response.StartTime,
			"end_time":                      response.EndTime,
			"car_wash_priority_id":          response.CarWashPriorityId,
			"car_wash_priority_description": response.CarWashPriorityDescription,
		}

		results = append(results, result)
	}

	pages.Rows = results

	return pages, nil
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
				Message:    updateQuery.Error.Error(),
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
	joinQuery := tx.Table("trx_car_wash").
		Joins("LEFT JOIN trx_work_order ON trx_car_wash.work_order_system_number = trx_work_order.work_order_system_number AND trx_car_wash.company_id = trx_work_order.company_id").
		Joins("LEFT JOIN mtr_car_wash_status ON trx_car_wash.car_wash_status_id = mtr_car_wash_status.car_wash_status_id").
		Joins("LEFT JOIN mtr_car_wash_bay ON trx_car_wash.car_wash_bay_id = mtr_car_wash_bay.car_wash_bay_id")

	keyAttributes := []string{
		"trx_work_order.work_order_system_number",
		"trx_work_order.work_order_document_number",
		"mtr_car_wash_bay.car_wash_bay_description",
		"trx_car_wash.car_wash_status_id",
		"mtr_car_wash_status.car_wash_status_description",
	}

	var result transactionjpcbpayloads.CarWashErrorDetail
	joinQuery = joinQuery.Select(keyAttributes).Where("trx_car_wash.work_order_system_number = ?", workOrderSystemNumber).
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
	var workOrderEntity transactionworkshopentities.WorkOrder
	var workOrderResponse transactionjpcbpayloads.CarWashWorkOrder

	err := tx.Model(&workOrderEntity).Select("car_wash, company_id, work_order_status_id").Where("work_order_system_number = ?", workOrderSystemNumber).Scan(&workOrderResponse).Error
	if err != nil {
		errorHelperInternalServerError(err)
	}

	getCompanyResponse, errFetchCompany := generalserviceapiutils.GetCompanyReferenceById(workOrderResponse.CompanyId)
	if errFetchCompany != nil {
		return transactionjpcbpayloads.CarWashPostResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Company Data from external service",
			Err:        errFetchCompany,
		}
	}

	if getCompanyResponse.UseJpcb {
		if workOrderResponse.WorkOrderStatusId == utils.WoStatQC {
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
						StatusId:              utils.CarWashStatDraft,
						PriorityId:            utils.CarWashPriorityNormal,
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
								StatusId:              utils.CarWashStatDraft,
								PriorityId:            utils.CarWashPriorityNormal,
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

func (*CarWashImpl) GetAllCarWashScreen(tx *gorm.DB, companyId int, carWashStatusId int) ([]transactionjpcbpayloads.CarWashScreenGetAllResponse, *exceptions.BaseErrorResponse) {
	var responses []transactionjpcbpayloads.CarWashScreenGetAllResponse

	keyAttributes := []string{
		"trx_car_wash.work_order_system_number", "trx_car_wash.car_wash_bay_id", "mtr_car_wash_bay.order_number", "trx_car_wash.car_wash_status_id",
		"mtr_car_wash_status.car_wash_status_description", "trx_work_order.model_id", "trx_work_order.vehicle_id",
	}

	rows, err := tx.Model(&transactionjpcbentities.CarWash{}).Select(keyAttributes).
		Order("trx_car_wash.car_wash_status_id desc, trx_car_wash.car_wash_bay_id asc, trx_car_wash.car_wash_priority_id desc").
		Order("trx_work_order.promise_date desc, trx_work_order.promise_time asc").
		Where("trx_work_order.company_id = ? AND trx_work_order.car_wash = ? AND trx_car_wash.car_wash_status_id <> ?", companyId, 1, utils.CarWashStatStop).
		Where("trx_car_wash.car_wash_status_id = ?", carWashStatusId).
		Joins("LEFT JOIN mtr_car_wash_bay on mtr_car_wash_bay.car_wash_bay_id = trx_car_wash.car_wash_bay_id AND mtr_car_wash_bay.company_id =  trx_car_wash.company_id").
		Joins("LEFT JOIN mtr_car_wash_status on mtr_car_wash_status.car_wash_status_id = trx_car_wash.car_wash_status_id").
		Joins("LEFT JOIN trx_work_order on trx_work_order.work_order_system_number = trx_car_wash.work_order_system_number").
		Rows()
	if err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        nil,
		}
	}
	defer rows.Close()

	for rows.Next() {
		var workOrderSystemNumber, carWashStatusId, modelId, vehicleId int
		var carWashBayId, orderNumber *int
		var carWashStatusDescription string

		err := rows.Scan(
			&workOrderSystemNumber, &carWashBayId, &orderNumber, &carWashStatusId,
			&carWashStatusDescription, &modelId, &vehicleId,
		)
		if err != nil {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		//Fetch data Model from external services
		getModelResponse, errFetchModel := salesserviceapiutils.GetUnitModelById(modelId)
		if errFetchModel != nil {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch unit model data grom external service",
				Err:        errFetchModel,
			}
		}

		// getVehicleResponse, errFetchVehicle := salesserviceapiutils.GetVehicleById(vehicleId)
		// if errFetchVehicle != nil {
		// 	return nil, &exceptions.BaseErrorResponse{
		// 		StatusCode: http.StatusInternalServerError,
		// 		Message:    "Failed to fetch vehicle data from external service",
		// 		Err:        errFetchVehicle,
		// 	}
		// }

		carWashScreen := transactionjpcbpayloads.CarWashScreenGetAllResponse{
			WorkOrderSystemNumber:    workOrderSystemNumber,
			CarWashBayId:             carWashBayId,
			OrderNumber:              orderNumber,
			CarWashStatusId:          carWashStatusId,
			CarWashStatusDescription: carWashStatusDescription,
			ModelId:                  modelId,
			ModelDescription:         getModelResponse.ModelName,
			VehicleId:                vehicleId,
			// ColourCommercialName:     getVehicleResponse.ColourCommercialName,
			ColourCommercialName: "",
		}

		fmt.Print(carWashScreen)
		responses = append(responses, carWashScreen)
	}

	return responses, nil
}

func (r *CarWashImpl) GetCarWashScreenDataByWorkOrderSystemNumber(tx *gorm.DB, workOrderSystemNumber int) (transactionjpcbpayloads.CarWashScreenGetAllResponse, *exceptions.BaseErrorResponse) {
	var response transactionjpcbpayloads.CarWashScreenGetAllResponse

	keyAttributes := []string{
		"trx_car_wash.work_order_system_number", "trx_car_wash.car_wash_bay_id", "mtr_car_wash_bay.order_number", "trx_car_wash.car_wash_status_id",
		"mtr_car_wash_status.car_wash_status_description", "trx_work_order.model_id", "trx_work_order.vehicle_id",
	}

	err := tx.Model(&transactionjpcbentities.CarWash{}).Select(keyAttributes).
		Where("trx_car_wash.work_order_system_number", workOrderSystemNumber).
		Joins("LEFT JOIN mtr_car_wash_bay on mtr_car_wash_bay.car_wash_bay_id = trx_car_wash.car_wash_bay_id AND mtr_car_wash_bay.company_id =  trx_car_wash.company_id").
		Joins("LEFT JOIN mtr_car_wash_status on mtr_car_wash_status.car_wash_status_id = trx_car_wash.car_wash_status_id").
		Joins("LEFT JOIN trx_work_order on trx_work_order.work_order_system_number = trx_car_wash.work_order_system_number").
		First(&response).Error
	if err != nil {
		return transactionjpcbpayloads.CarWashScreenGetAllResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return response, nil
}

func (r *CarWashImpl) UpdateBayNumberCarWashScreen(tx *gorm.DB, bayNumber int, workOrderSystemNumber int) (transactionjpcbpayloads.CarWashScreenGetAllResponse, *exceptions.BaseErrorResponse) {
	var carWashBay int
	err := tx.Model(&transactionjpcbentities.BayMaster{}).Select("car_wash_bay_id").Where(&transactionjpcbentities.BayMaster{
		CarWashBayId: bayNumber,
	}).First(&carWashBay).Error

	if err != nil {
		return transactionjpcbpayloads.CarWashScreenGetAllResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
			Message:    "Car wash bay doesn't exist",
		}
	}

	err = tx.Model(&transactionjpcbentities.CarWash{}).Where("work_order_system_number = ?", workOrderSystemNumber).Update("car_wash_bay_id", bayNumber).Error
	if err != nil {
		return transactionjpcbpayloads.CarWashScreenGetAllResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	response, getCarWashDataError := r.GetCarWashScreenDataByWorkOrderSystemNumber(tx, workOrderSystemNumber)
	if getCarWashDataError != nil {
		return transactionjpcbpayloads.CarWashScreenGetAllResponse{}, getCarWashDataError
	}

	return response, nil
}

func (r *CarWashImpl) StartCarWash(tx *gorm.DB, workOrderSystemNumber, carWashBayId int) (transactionjpcbpayloads.CarWashScreenGetAllResponse, *exceptions.BaseErrorResponse) {
	var carWashStatusId int
	err := tx.Model(&transactionjpcbentities.CarWash{}).Select("car_wash_status_id").Where("work_order_system_number = ? AND car_wash_bay_id = ?", workOrderSystemNumber, carWashBayId).Scan(&carWashStatusId).Error
	if err != nil {
		return transactionjpcbpayloads.CarWashScreenGetAllResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	if carWashStatusId == 0 {
		return transactionjpcbpayloads.CarWashScreenGetAllResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        fmt.Errorf("bad request"),
			Message:    "please check if bay is correct" + strconv.Itoa(carWashStatusId),
		}
	}

	if carWashStatusId != utils.CarWashStatStart {

		type carWashModel struct {
			CarWashBayId int `json:"car_wash_bay_id"`
			CompanyId    int `json:"company_id"`
		}

		var carWash carWashModel
		err := tx.Table("trx_car_wash AS carwash").
			Select("bay.car_wash_bay_id, carwash.company_id").
			Joins("LEFT JOIN mtr_car_wash_bay AS bay ON carwash.car_wash_bay_id = bay.car_wash_bay_id").
			Where("carwash.work_order_system_number = ?", workOrderSystemNumber).
			Scan(&carWash).Error
		if err != nil {
			return transactionjpcbpayloads.CarWashScreenGetAllResponse{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		if carWash.CarWashBayId == 0 {
			return transactionjpcbpayloads.CarWashScreenGetAllResponse{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusUnprocessableEntity,
				Err:        fmt.Errorf("update failed, bay already removed. please refresh page"),
				Message:    "update failed, bay already removed. please refresh page",
			}
		}

		getCompanyReferenceResponse, errFetchCompanyReference := generalserviceapiutils.GetCompanyReferenceById(carWash.CompanyId)
		if errFetchCompanyReference != nil {
			return transactionjpcbpayloads.CarWashScreenGetAllResponse{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
				Message:    "Failed to fetch company reference data from extternal service",
			}
		}

		err = tx.Model(&transactionjpcbentities.CarWash{}).Where(transactionjpcbentities.CarWash{
			WorkOrderSystemNumber: workOrderSystemNumber,
		}).Updates(&transactionjpcbpayloads.StartCarWashUpdates{
			CarWashStatusId: utils.CarWashStatStart,
			CarWashDate:     time.Now(),
			CarWashBayId:    carWash.CarWashBayId,
			StartTime:       createCurrentTime(getCompanyReferenceResponse.TimeDifference),
		}).Error
		if err != nil {
			return transactionjpcbpayloads.CarWashScreenGetAllResponse{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		result, getCarWashDataError := r.GetCarWashScreenDataByWorkOrderSystemNumber(tx, workOrderSystemNumber)
		if getCarWashDataError != nil {
			return transactionjpcbpayloads.CarWashScreenGetAllResponse{}, getCarWashDataError
		}

		//Fetch data Model from external services
		getModelResponse, errFetchModel := salesserviceapiutils.GetUnitModelById(result.ModelId)
		if errFetchModel != nil {
			return transactionjpcbpayloads.CarWashScreenGetAllResponse{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch unit model data grom external service",
				Err:        errFetchModel,
			}
		}

		// getVehicleResponse, errFetchVehicle := salesserviceapiutils.GetVehicleById(result.VehicleId)
		// if errFetchVehicle != nil {
		// 	return transactionjpcbpayloads.CarWashScreenGetAllResponse{}, &exceptions.BaseErrorResponse{
		// 		StatusCode: http.StatusInternalServerError,
		// 		Message:    "Failed to fetch vehicle data from external service",
		// 		Err:        errFetchVehicle,
		// 	}
		// }

		// TODO Exec uspg_wtWorkOrderLog_Insert

		return transactionjpcbpayloads.CarWashScreenGetAllResponse{
			WorkOrderSystemNumber:    result.WorkOrderSystemNumber,
			CarWashBayId:             result.CarWashBayId,
			OrderNumber:              result.OrderNumber,
			CarWashStatusId:          result.CarWashStatusId,
			CarWashStatusDescription: result.CarWashStatusDescription,
			ModelId:                  result.ModelId,
			ModelDescription:         getModelResponse.ModelName,
			VehicleId:                result.VehicleId,
			// ColourCommercialName:     getVehicleResponse.ColourCommercialName,
			ColourCommercialName: "",
		}, nil
	}
	return transactionjpcbpayloads.CarWashScreenGetAllResponse{}, &exceptions.BaseErrorResponse{
		StatusCode: http.StatusBadRequest,
		Err:        fmt.Errorf("carwash already started"),
		Message:    "Work Order carwash already started",
	}
}

func (r *CarWashImpl) StopCarWash(tx *gorm.DB, workOrderSystemNumber int) (transactionjpcbpayloads.CarWashScreenGetAllResponse, *exceptions.BaseErrorResponse) {
	var carWashStatus int

	err := tx.Model(&transactionjpcbentities.CarWash{}).Select("car_wash_status_id").Where(&transactionjpcbentities.CarWash{
		WorkOrderSystemNumber: workOrderSystemNumber,
	}).First(&carWashStatus).Error

	if err != nil {
		return transactionjpcbpayloads.CarWashScreenGetAllResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
		}
	}

	if carWashStatus == utils.CarWashStatStart {
		var startTime float32
		err := tx.Model(&transactionjpcbentities.CarWash{}).Select("start_time").Where(&transactionjpcbentities.CarWash{
			WorkOrderSystemNumber: workOrderSystemNumber,
		}).First(&startTime).Error
		if err != nil {
			return transactionjpcbpayloads.CarWashScreenGetAllResponse{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		err = tx.Table("trx_car_wash").
			Where("trx_car_wash.work_order_system_number = ?", workOrderSystemNumber).
			Updates(transactionjpcbentities.CarWash{
				StatusId:   utils.CarWashStatStop,
				EndTime:    createCurrentTime(0),
				ActualTime: createCurrentTime(0) - startTime,
			}).Error
		if err != nil {
			return transactionjpcbpayloads.CarWashScreenGetAllResponse{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		result, getCarWashDataError := r.GetCarWashScreenDataByWorkOrderSystemNumber(tx, workOrderSystemNumber)
		if getCarWashDataError != nil {
			return transactionjpcbpayloads.CarWashScreenGetAllResponse{}, getCarWashDataError
		}

		//Fetch data Model from external services
		getModelResponse, errFetchModel := salesserviceapiutils.GetUnitModelById(result.ModelId)
		if errFetchModel != nil {
			return transactionjpcbpayloads.CarWashScreenGetAllResponse{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch unit model data grom external service",
				Err:        errFetchModel,
			}
		}

		// getVehicleResponse, errFetchVehicle := salesserviceapiutils.GetVehicleById(result.VehicleId)
		// if errFetchVehicle != nil {
		// 	return transactionjpcbpayloads.CarWashScreenGetAllResponse{}, &exceptions.BaseErrorResponse{
		// 		StatusCode: http.StatusInternalServerError,
		// 		Message:    "Failed to fetch vehicle data from external service",
		// 		Err:        errFetchVehicle,
		// 	}
		// }

		//TODO uspg_wtWorkOrderLog_Insert

		return transactionjpcbpayloads.CarWashScreenGetAllResponse{
			WorkOrderSystemNumber:    result.WorkOrderSystemNumber,
			CarWashBayId:             result.CarWashBayId,
			OrderNumber:              result.OrderNumber,
			CarWashStatusId:          result.CarWashStatusId,
			CarWashStatusDescription: result.CarWashStatusDescription,
			ModelId:                  result.ModelId,
			ModelDescription:         getModelResponse.ModelName,
			VehicleId:                result.VehicleId,
			// ColourCommercialName:     getVehicleResponse.ColourCommercialName,
			ColourCommercialName: "",
		}, nil
	}

	return transactionjpcbpayloads.CarWashScreenGetAllResponse{}, &exceptions.BaseErrorResponse{
		StatusCode: http.StatusBadRequest,
		Err:        fmt.Errorf("work order carwash has not started"),
		Message:    "Work order carwash has not started",
	}
}

func (r *CarWashImpl) CancelCarWash(tx *gorm.DB, workOrderSystemNumber int) (transactionjpcbpayloads.CarWashScreenGetAllResponse, *exceptions.BaseErrorResponse) {
	var carWashStatus int

	err := tx.Model(&transactionjpcbentities.CarWash{}).Select("car_wash_status_id").Where(&transactionjpcbentities.CarWash{
		WorkOrderSystemNumber: workOrderSystemNumber,
	}).First(&carWashStatus).Error
	if err != nil {
		return transactionjpcbpayloads.CarWashScreenGetAllResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}
	fmt.Print(carWashStatus)

	if carWashStatus != utils.CarWashStatStart {
		err := tx.Model(&transactionjpcbentities.CarWash{}).Where(&transactionjpcbentities.CarWash{
			WorkOrderSystemNumber: workOrderSystemNumber,
		}).Update("car_wash_bay_id", nil).Error
		if err != nil {
			return transactionjpcbpayloads.CarWashScreenGetAllResponse{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}

		result, getCarWashError := r.GetCarWashScreenDataByWorkOrderSystemNumber(tx, workOrderSystemNumber)
		if getCarWashError != nil {
			return transactionjpcbpayloads.CarWashScreenGetAllResponse{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}

		return result, nil
	}

	return transactionjpcbpayloads.CarWashScreenGetAllResponse{}, &exceptions.BaseErrorResponse{
		StatusCode: http.StatusBadRequest,
		Err:        fmt.Errorf("already start"),
		Message:    "Work order carwash already started",
	}
}

func (*CarWashImpl) GetCarWashByWorkOrderSystemNumber(tx *gorm.DB, workOrderSystemNumber int) (transactionjpcbpayloads.CarWashGetAllResponse, *exceptions.BaseErrorResponse) {
	carWashPayload := transactionjpcbpayloads.CarWashGetAllResponse{}
	type selectCarWashResponse struct {
		WorkOrderSystemNumber      int
		WorkOrderDocumentNumber    string
		ModelId                    int
		VehicleId                  int
		PromiseTime                time.Time
		PromiseDate                time.Time
		CarWashBayId               int
		CarWashBayDescription      string
		CarWashStatusId            int
		CarWashStatusDescription   string
		StartTime                  float32
		EndTime                    float32
		CarWashPriorityId          int
		CarWashPriorityDescription string
	}
	var result selectCarWashResponse

	selectAttributes := []string{
		"trx_work_order.work_order_system_number", "trx_work_order.work_order_document_number",
		"trx_work_order.model_id", "trx_work_order.vehicle_id",
		"trx_work_order.promise_time", "trx_work_order.promise_date",
		"trx_car_wash.car_wash_bay_id", "mtr_car_wash_bay.car_wash_bay_description",
		"trx_car_wash.car_wash_status_id", "mtr_car_wash_status.car_wash_status_description",
		"trx_car_wash.start_time", "trx_car_wash.end_time",
		"trx_car_wash.car_wash_priority_id", "mtr_car_wash_priority.car_wash_priority_description",
	}
	query := tx.Model(&transactionjpcbentities.CarWash{}).Where("trx_work_order.work_order_system_number = ?", workOrderSystemNumber).
		Select(selectAttributes).
		Joins("LEFT JOIN trx_work_order ON trx_car_wash.work_order_system_number = trx_work_order.work_order_system_number AND trx_car_wash.company_id = trx_work_order.company_id").
		Joins("LEFT JOIN mtr_car_wash_priority ON trx_car_wash.car_wash_priority_id = mtr_car_wash_priority.car_wash_priority_id").
		Joins("LEFT JOIN mtr_car_wash_status ON trx_car_wash.car_wash_status_id = mtr_car_wash_status.car_wash_status_id").
		Joins("LEFT JOIN mtr_car_wash_bay ON trx_car_wash.car_wash_bay_id = mtr_car_wash_bay.car_wash_bay_id").First(&transactionjpcbentities.CarWash{}).Scan(&result).Error

	if query != nil {
		return transactionjpcbpayloads.CarWashGetAllResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Data not found",
			Err:        query,
		}
	}

	//Fetch data Model from external services
	getModelResponse, errFetchModel := salesserviceapiutils.GetUnitModelById(result.ModelId)
	if errFetchModel != nil {
		return transactionjpcbpayloads.CarWashGetAllResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch unit model data grom external service",
			Err:        errFetchModel,
		}
	}

	// getVehicleResponse, errFetchVehicle := salesserviceapiutils.GetVehicleById(result.VehicleId)
	// if errFetchVehicle != nil {
	// 	return transactionjpcbpayloads.CarWashGetAllResponse{}, &exceptions.BaseErrorResponse{
	// 		StatusCode: http.StatusInternalServerError,
	// 		Message:    "Failed to fetch vehicle data from external service",
	// 		Err:        errFetchVehicle,
	// 	}
	// }

	carWashPayload = transactionjpcbpayloads.CarWashGetAllResponse{
		WorkOrderSystemNumber:   result.WorkOrderSystemNumber,
		WorkOrderDocumentNumber: result.WorkOrderDocumentNumber,
		Model:                   getModelResponse.ModelName,
		// Color:                      getVehicleResponse.ColourCommercialName,
		Color: "",
		// Tnkb:  getVehicleResponse.VehicleRegistrationCertificateTNKB,
		Tnkb:                       "",
		PromiseTime:                &result.PromiseTime,
		PromiseDate:                &result.PromiseDate,
		CarWashBayId:               &result.CarWashBayId,
		CarWashBayDescription:      &result.CarWashBayDescription,
		CarWashStatusId:            result.CarWashStatusId,
		CarWashStatusDescription:   result.CarWashStatusDescription,
		StartTime:                  result.StartTime,
		EndTime:                    result.EndTime,
		CarWashPriorityId:          result.CarWashPriorityId,
		CarWashPriorityDescription: result.CarWashPriorityDescription,
	}

	return carWashPayload, nil
}

func createCurrentTime(timeDiff float32) float32 {
	// Split the float into whole hours and fractional minutes
	wholeHours := int(timeDiff)
	fractionalMinutes := int((timeDiff - float32(wholeHours)) * 100)
	// Get the current time
	now := time.Now()

	// Add the whole hours and fractional minutes
	result := now.Add(time.Duration(wholeHours)*time.Hour + time.Duration(fractionalMinutes)*time.Minute)
	fmt.Println(float32(result.Minute()) / 100)
	fmt.Println(result.Hour())
	return float32(result.Hour()) + (float32(result.Minute()) / 100)
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
