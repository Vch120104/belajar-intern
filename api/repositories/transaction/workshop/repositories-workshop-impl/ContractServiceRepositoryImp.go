package transactionworkshoprepositoryimpl

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshoprepository "after-sales/api/repositories/transaction/workshop"
	"after-sales/api/utils"
	"errors"
	"fmt"
	"net/http"
	"strconv"

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
	fmt.Println("Generated Query:", joinTable)
	whereQuery := utils.ApplyFilter(joinTable, filterCondition)
	fmt.Println("Where Query Generated:", whereQuery)

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
			BrandName:                     brandResponse.BrandName,        // Mengambil nama brand dari response API
			ModelName:                     modelResponse.ModelName,        // Mengambil nama model dari response API
			VehicleTnkb:                   vehicleResponses.VehicleTnkb,   // Mengambil TNKB kendaraan dari response API
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
			"brand_name":                       response.BrandName,
			"model_name":                       response.ModelName,
			"tnkb":                             response.VehicleTnkb,
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
