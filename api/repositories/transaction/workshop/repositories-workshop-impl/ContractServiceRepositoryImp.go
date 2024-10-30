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

func (r *ContractServiceRepositoryImpl) GetAll(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var entities []transactionworkshoppayloads.ContractServiceRequest

	// Buat query dasar
	joinTable := utils.CreateJoinSelectStatement(tx, transactionworkshoppayloads.ContractServiceRequest{})
	whereQuery := utils.ApplyFilter(joinTable, filterCondition)

	// Eksekusi query untuk mendapatkan data contract service
	if err := whereQuery.Find(&entities).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Contract service not found",
				Err:        err,
			}
		}
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch contract service entity",
			Err:        err,
		}
	}

	if len(entities) == 0 {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "No contract service entities found",
		}
	}

	// Konversi ke response yang dibutuhkan, termasuk fetch data dari external API jika perlu
	var convertedResponses []transactionworkshoppayloads.ContractServiceResponse
	for _, entity := range entities {
		// Fetch data external (contoh: brand, model, TNKB)
		BrandUrl := config.EnvConfigs.SalesServiceUrl + "unit-brand/" + strconv.Itoa(entity.BrandId)
		var brandResponse transactionworkshoppayloads.ContractServiceBrand
		errBrand := utils.Get(BrandUrl, &brandResponse, nil)
		if errBrand != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve brand data from external API",
				Err:        errBrand,
			}
		}

		ModelUrl := config.EnvConfigs.SalesServiceUrl + "unit-model/" + strconv.Itoa(entity.ModelId)
		var modelResponse transactionworkshoppayloads.ContractServiceModel
		errModel := utils.Get(ModelUrl, &modelResponse, nil)
		if errModel != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve model data from external API",
				Err:        errModel,
			}
		}

		// Contoh tambahan untuk TNKB/Vehicle data dari API eksternal jika perlu
		VehicleUrl := config.EnvConfigs.SalesServiceUrl + "vehicle-master/" + strconv.Itoa(entity.VehicleId)
		var vehicleResponses transactionworkshoppayloads.VehicleResponse
		errVehicle := utils.Get(VehicleUrl, &vehicleResponses, nil)
		if errVehicle != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve vehicle data from the external API",
				Err:        errVehicle,
			}
		}

		// Buat response contract service
		convertedResponses = append(convertedResponses, transactionworkshoppayloads.ContractServiceResponse{
			ContractServiceDocumentNumber: entity.ContractServiceDocumentNumber,
			ContractServiceFrom:           entity.ContractServiceFrom,
			ContractServiceTo:             entity.ContractServiceTo,
			BrandId:                       brandResponse.BrandId,
			BrandName:                     brandResponse.BrandName,
			BrandCode:                     brandResponse.BrandCode, // Mengambil nama brand dari response API
			ModelId:                       modelResponse.ModelId,
			ModelName:                     modelResponse.ModelName,
			ModelCode:                     modelResponse.ModelCode, // Mengambil nama model dari response API
			VehicleTnkb:                   vehicleResponses.VehicleTnkb,
			VehicleCode:                   vehicleResponses.VehicleCode,   // Mengambil TNKB kendaraan dari response API
			ContractServiceStatusId:       entity.ContractServiceStatusId, // Mengambil status dari entitas contract service
			ContractServiceSystemNumber:   entity.ContractServiceSystemNumber,
		})

	}

	// Konversi hasil ke dalam format yang diinginkan
	var mapResponses []map[string]interface{}
	for _, response := range convertedResponses {
		responseMap := map[string]interface{}{
			"contract_service_document_number": response.ContractServiceDocumentNumber,
			"contract_service_from":            response.ContractServiceFrom,
			"contract_service_to":              response.ContractServiceTo,
			"brand_id":                         response.BrandId,
			"brand_name":                       response.BrandName,
			"brand_code":                       response.BrandCode,
			"model_id":                         response.ModelId,
			"model_name":                       response.ModelName,
			"model_code":                       response.ModelCode,
			"vehicle_tnkb":                     response.VehicleTnkb,
			"vehicle_code":                     response.VehicleCode,
			// "customer_name":                    response.ContractServiceStatusId,
			"contract_service_system_number": response.ContractServiceSystemNumber,
			"status":                         response.ContractServiceStatusId,
		}
		mapResponses = append(mapResponses, responseMap)
	}

	// Lakukan pagination
	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(mapResponses, &pages)

	return paginatedData, totalPages, totalRows, nil
}

// GetById implements transactionworkshoprepository.ContractServiceRepository.
func (r *ContractServiceRepositoryImpl) GetById(tx *gorm.DB, Id int, filterCondition []utils.FilterCondition, pages pagination.Pagination) (transactionworkshoppayloads.ContractServiceResponseId, *exceptions.BaseErrorResponse) {

	var entity transactionworkshopentities.ContractService
	// Fetch Contract Service by Id
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

	// Fetch data brand from external API
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

	// Fetch data model from external API
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

	// Fetch data vehicle from external API
	VehicleUrl := config.EnvConfigs.SalesServiceUrl + "vehicle-master?page=0&limit=100&vehicle_id=" + strconv.Itoa(entity.VehicleId)
	var vehicleResponses []transactionworkshoppayloads.ContractServiceVehicleResponse
	errVehicle := utils.GetArray(VehicleUrl, &vehicleResponses, nil)
	if errVehicle != nil {
		return transactionworkshoppayloads.ContractServiceResponseId{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve vehicle data from external API",
			Err:        errVehicle,
		}
	}

	// Handle case where vehicle data is not found
	var vehicleTnkb, vehicleCode, vehicleOwner string
	if len(vehicleResponses) > 0 {
		vehicleTnkb = vehicleResponses[0].VehicleTnkb
		vehicleCode = vehicleResponses[0].VehicleCode   // Mengambil VehicleCode dari respons API
		vehicleOwner = vehicleResponses[0].VehicleOwner // Mengambil VehicleOwner dari respons API
	} else {
		vehicleTnkb = "Unknown"
		vehicleCode = "Unknown"  // Memberikan nilai default jika tidak ditemukan
		vehicleOwner = "Unknown" // Memberikan nilai default jika tidak ditemukan
	}

	// Prepare the response payload
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
		ContractServiceStatusId:       entity.ContractServiceStatusId,
	}

	return payload, nil
}

// Save implements transactionworkshoprepository.ContractServiceRepository.
func (r *ContractServiceRepositoryImpl) Save(tx *gorm.DB, payload transactionworkshoppayloads.ContractServiceInsert) (transactionworkshoppayloads.ContractServiceInsert, *exceptions.BaseErrorResponse) {
	// Fetch data eksternal dari API
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

	// Set default nilai total, vat, dan grand_total
	payload.Total = 0.0
	payload.Vat = 0.0
	payload.GrandTotal = 0.0

	// Prepare entity for insertion (only IDs and finance values set to 0 initially)
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

	// Simpan ke database
	err := tx.Create(&contractService).Error
	if err != nil {
		return payload, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to insert contract service data",
			Err:        err,
		}
	}

	// Update response payload dengan data dari external API
	payload.BrandName = brandResponse.BrandName
	payload.ModelName = modelResponse.ModelName
	payload.VehicleTnkb = vehicleResponses.VehicleTnkb
	payload.VehicleOwner = vehicleResponses.VehicleOwner

	// Return updated payload
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

// Submit implements transactionworkshoprepository.ContractServiceRepository.
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
