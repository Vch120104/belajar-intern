package transactionworkshoprepositoryimpl

import (
	"after-sales/api/config"
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshoprepository "after-sales/api/repositories/transaction/workshop"
	"after-sales/api/utils"
	salesserviceapiutils "after-sales/api/utils/sales-service"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type ContractServiceRepositoryImpl struct {
}

func OpenContractServicelRepositoryImpl() transactionworkshoprepository.ContractServiceRepository {
	return &ContractServiceRepositoryImpl{}
}

func (r *ContractServiceRepositoryImpl) GetAll(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var responses []transactionworkshoppayloads.ContractServiceRequest

	baseModelQuery := tx.Model(&transactionworkshopentities.ContractService{})

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
		// Fetch external data for brand, model, and vehicle
		brandResponse, errBrand := salesserviceapiutils.GetUnitBrandById(response.BrandId)
		if errBrand != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: errBrand.StatusCode,
				Err:        errBrand.Err,
			}
		}

		modelResponse, errModel := salesserviceapiutils.GetUnitModelById(response.ModelId)
		if errModel != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: errModel.StatusCode,
				Err:        errModel.Err,
			}
		}

		vehicleResponse, errVehicle := salesserviceapiutils.GetVehicleById(response.VehicleId)
		if errVehicle != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: errVehicle.StatusCode,
				Err:        errVehicle.Err,
			}
		}

		// Construct the final result map
		result := map[string]interface{}{
			"company_id":                       response.CompanyId,
			"contract_service_system_number":   response.ContractServiceSystemNumber,
			"contract_service_document_number": response.ContractServiceDocumentNumber,
			"contract_service_from":            response.ContractServiceFrom,
			"contract_service_to":              response.ContractServiceTo,
			"brand_id":                         brandResponse.BrandId,
			"brand_name":                       brandResponse.BrandName,
			"brand_code":                       brandResponse.BrandCode,
			"model_id":                         modelResponse.ModelId,
			"model_name":                       modelResponse.ModelName,
			"model_code":                       modelResponse.ModelCode,
			"vehicle_id":                       vehicleResponse.VehicleID,
			"vehicle_tnkb":                     vehicleResponse.VehicleRegistrationCertificateTNKB,
			"vehicle_code":                     vehicleResponse.VehicleChassisNumber,
			"contract_service_status_id":       response.ContractServiceStatusId,
		}

		results = append(results, result)
	}

	pages.Rows = results

	return pages, nil
}

// GetById implements transactionworkshoprepository.ContractServiceRepository.
func (r *ContractServiceRepositoryImpl) GetById(tx *gorm.DB, Id int, filterCondition []utils.FilterCondition, pages pagination.Pagination) (transactionworkshoppayloads.ContractServiceResponseId, *exceptions.BaseErrorResponse) {

	var entity transactionworkshopentities.ContractService
	err := tx.Model(&transactionworkshopentities.ContractService{}).
		Where("contract_service_system_number = ?", Id).
		First(&entity).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return transactionworkshoppayloads.ContractServiceResponseId{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Contract service not found",
			}
		}
		return transactionworkshoppayloads.ContractServiceResponseId{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	BrandUrl := config.EnvConfigs.SalesServiceUrl + "unit-brand/" + strconv.Itoa(entity.BrandId)
	var brandResponse transactionworkshoppayloads.ContractServiceBrand
	errBrand := utils.Get(BrandUrl, &brandResponse, nil)
	if errBrand != nil {
		return transactionworkshoppayloads.ContractServiceResponseId{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve brand data from external API",
			Err:        errBrand,
		}
	}

	ModelUrl := config.EnvConfigs.SalesServiceUrl + "unit-model/" + strconv.Itoa(entity.ModelId)
	var modelResponse transactionworkshoppayloads.ContractServiceModel
	errModel := utils.Get(ModelUrl, &modelResponse, nil)
	if errModel != nil {
		return transactionworkshoppayloads.ContractServiceResponseId{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve model data from external API",
			Err:        errModel,
		}
	}

	VehicleUrl := config.EnvConfigs.SalesServiceUrl + "vehicle-master/" + strconv.Itoa(entity.VehicleId)
	var vehicleResponses transactionworkshoppayloads.ContractServiceVehicleResponse
	errVehicle := utils.Get(VehicleUrl, &vehicleResponses, nil)
	if errVehicle != nil {
		return transactionworkshoppayloads.ContractServiceResponseId{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve vehicle data from external API",
			Err:        errVehicle,
		}
	}

	vehicleTnkb := vehicleResponses.STNK.VehicleRegistrationCertificateTnkb
	vehicleCode := vehicleResponses.Master.VehicleCode // Pastikan nama field benar
	vehicleOwner := vehicleResponses.STNK.VehicleRegistrationCertificateOwnerName
	vehicleEngineNumber := vehicleResponses.Master.VehicleEngineNumber

	payload := transactionworkshoppayloads.ContractServiceResponseId{
		CompanyId:                     entity.CompanyId,
		ContractServiceSystemNumber:   entity.ContractServiceSystemNumber,
		ContractServiceDocumentNumber: entity.ContractServiceDocumentNumber,
		ContractServiceFrom:           entity.ContractServiceFrom,
		ContractServiceTo:             entity.ContractServiceTo,
		BrandId:                       entity.BrandId,
		BrandName:                     brandResponse.BrandName,
		BrandCode:                     brandResponse.BrandCode,
		ModelId:                       entity.ModelId,
		ModelName:                     modelResponse.ModelName,
		ModelCode:                     modelResponse.ModelCode,
		VehicleId:                     entity.VehicleId,
		VehicleTnkb:                   vehicleTnkb,
		VehicleCode:                   vehicleCode,
		VehicleOwner:                  vehicleOwner,
		VehicleEngineNumber:           vehicleEngineNumber,
		ContractServiceStatusId:       entity.ContractServiceStatusId,
	}

	return payload, nil
}

// Save implements transactionworkshoprepository.ContractServiceRepository.
func (r *ContractServiceRepositoryImpl) Save(tx *gorm.DB, payload transactionworkshoppayloads.ContractServiceInsert) (transactionworkshoppayloads.ContractServiceInsert, *exceptions.BaseErrorResponse) {
	BrandUrl := config.EnvConfigs.SalesServiceUrl + "unit-brand/" + strconv.Itoa(payload.BrandId)
	var brandResponse transactionworkshoppayloads.ContractServiceBrand
	errBrand := utils.Get(BrandUrl, &brandResponse, nil)
	if errBrand != nil {
		return payload, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve brand data from external API",
			Err:        errBrand,
		}
	}

	ModelUrl := config.EnvConfigs.SalesServiceUrl + "unit-model/" + strconv.Itoa(payload.ModelId)
	var modelResponse transactionworkshoppayloads.ContractServiceModel
	errModel := utils.Get(ModelUrl, &modelResponse, nil)
	if errModel != nil {
		return payload, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve model data from external API",
			Err:        errModel,
		}
	}

	VehicleUrl := config.EnvConfigs.SalesServiceUrl + "vehicle-master/" + strconv.Itoa(payload.VehicleId)
	var vehicleResponses transactionworkshoppayloads.ContractServiceVehicleResponse
	errVehicle := utils.Get(VehicleUrl, &vehicleResponses, nil)
	if errVehicle != nil {
		return payload, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve vehicle data from external API",
			Err:        errVehicle,
		}
	}

	payload.Total = 0.0
	payload.Vat = 0.0
	payload.GrandTotal = 0.0

	contractService := transactionworkshopentities.ContractService{
		CompanyId:                     payload.CompanyId,
		ContractServiceDocumentNumber: payload.ContractServiceDocumentNumber,
		ContractServiceDate:           payload.ContractServiceDate,
		ContractServiceFrom:           payload.ContractServiceFrom,
		ContractServiceTo:             payload.ContractServiceTo,
		BrandId:                       payload.BrandId,
		ModelId:                       payload.ModelId,
		VehicleId:                     payload.VehicleId,
		RegisteredMileage:             payload.RegisteredMileage,
		Remark:                        payload.Remark,
		ContractServiceStatusId:       payload.ContractServiceStatusId,
		Total:                         payload.Total,
		TotalValueAfterTax:            payload.Vat,
		ValueAfterTaxrate:             payload.GrandTotal,
	}

	err := tx.Create(&contractService).Error
	if err != nil {
		return payload, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to insert contract service data",
			Err:        err,
		}
	}

	payload.BrandName = brandResponse.BrandName
	payload.ModelName = modelResponse.ModelName
	payload.VehicleTnkb = vehicleResponses.STNK.VehicleRegistrationCertificateTnkb
	payload.VehicleOwner = vehicleResponses.STNK.VehicleRegistrationCertificateOwnerName

	return payload, nil
}

// Void implements transactionworkshoprepository.ContractServiceRepository.
func (r *ContractServiceRepositoryImpl) Void(tx *gorm.DB, Id int) (bool, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.ContractService

	err := tx.Model(transactionworkshopentities.ContractService{}).Where("contract_service_system_number", Id).First(&entity).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Data not found",
			}
		}
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	err = tx.Delete(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to delete the work order",
			Err:        err,
		}
	}
	return true, nil
}

func (r *ContractServiceRepositoryImpl) GenerateDocumentNumber(tx *gorm.DB, Id int) (string, *exceptions.BaseErrorResponse) {
	var supplySlip transactionsparepartentities.SupplySlip

	err1 := tx.Model(&transactionsparepartentities.SupplySlip{}).
		Where("contract_service_system_number = ?", Id).
		First(&supplySlip).
		Error
	if err1 != nil {
		return "", &exceptions.BaseErrorResponse{Message: fmt.Sprintf("Failed to retrieve contract service from the database: %v", err1)}
	}

	var workOrder transactionworkshopentities.WorkOrder
	var brandResponse transactionworkshoppayloads.BrandDocResponse

	workOrderId := supplySlip.WorkOrderSystemNumber

	// Get the work order based on the work order system number
	err := tx.Model(&transactionworkshopentities.WorkOrder{}).Where("work_order_system_number = ?", workOrderId).First(&workOrder).Error
	if err != nil {

		return "", &exceptions.BaseErrorResponse{Message: fmt.Sprintf("Failed to retrieve work order from the database: %v", err)}
	}

	if workOrder.BrandId == 0 {

		return "", &exceptions.BaseErrorResponse{Message: "brand_id is missing in the work order. Please ensure the work order has a valid brand_id before generating document number."}
	}

	// Get the last work order based on the work order system number
	var lastWorkOrder transactionworkshopentities.WorkOrder
	err = tx.Model(&transactionworkshopentities.WorkOrder{}).
		Where("brand_id = ?", workOrder.BrandId).
		Order("work_order_document_number desc").
		First(&lastWorkOrder).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {

		return "", &exceptions.BaseErrorResponse{Message: fmt.Sprintf("Failed to retrieve last work order: %v", err)}
	}

	currentTime := time.Now()
	month := int(currentTime.Month())
	year := currentTime.Year() % 100 // Use last two digits of the year

	// fetch data brand from external api
	brandUrl := config.EnvConfigs.SalesServiceUrl + "unit-brand/" + strconv.Itoa(workOrder.BrandId)
	errUrl := utils.Get(brandUrl, &brandResponse, nil)
	if errUrl != nil {
		return "", &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errUrl,
		}
	}

	// Check if BrandCode is not empty before using it
	if brandResponse.BrandCode == "" {
		return "", &exceptions.BaseErrorResponse{StatusCode: http.StatusInternalServerError, Message: "Brand code is empty"}
	}

	// Get the initial of the brand code
	brandInitial := brandResponse.BrandCode[0]

	// Handle the case when there is no last work order or the format is invalid
	newDocumentNumber := fmt.Sprintf("WSCS/%c/%02d/%02d/00001", brandInitial, month, year)
	if lastWorkOrder.WorkOrderSystemNumber != 0 {
		lastWorkOrderDate := lastWorkOrder.WorkOrderDate
		lastWorkOrderYear := lastWorkOrderDate.Year() % 100

		// Check if the last work order is from the same year
		if lastWorkOrderYear == year {
			lastWorkOrderCode := lastWorkOrder.WorkOrderDocumentNumber
			codeParts := strings.Split(lastWorkOrderCode, "/")
			if len(codeParts) == 5 {
				lastWorkOrderNumber, err := strconv.Atoi(codeParts[4])
				if err == nil {
					newWorkOrderNumber := lastWorkOrderNumber + 1
					newDocumentNumber = fmt.Sprintf("SPSS/%c/%02d/%02d/%05d", brandInitial, month, year, newWorkOrderNumber)
				} else {
					log.Printf("Failed to parse last work order code: %v", err)
				}
			} else {
				log.Println("Invalid last work order code format")
			}
		}
	}

	log.Printf("New document number: %s", newDocumentNumber)
	return newDocumentNumber, nil
}

func (r *ContractServiceRepositoryImpl) Submit(tx *gorm.DB, Id int) (bool, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.ContractService

	err := tx.Model(&transactionworkshopentities.ContractService{}).Where("contract_service_system_number = ?", Id).First(&entity).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, &exceptions.BaseErrorResponse{
				Message: "No COntract Service Data Found",
			}
		}
		return false, &exceptions.BaseErrorResponse{
			Message: fmt.Sprintf("Failed to retrive contract service from the database: %v", err),
		}
	}

	if entity.ContractServiceDocumentNumber == " " {
		newDocumentNumber, genErr := r.GenerateDocumentNumber(tx, entity.ContractServiceSystemNumber)
		if genErr != nil {
			return false, genErr
		}

		entity.ContractServiceDocumentNumber = newDocumentNumber

		err = tx.Save(&entity).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				Message: "Failed to submit contract service",
			}
		}
		return true, nil
	} else {
		return false, &exceptions.BaseErrorResponse{
			Message: "Document number has already been generated",
		}
	}
}
