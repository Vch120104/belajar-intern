package transactionworkshoprepositoryimpl

import (
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionworkshoprepository "after-sales/api/repositories/transaction/workshop"
	"after-sales/api/utils"
	salesserviceapiutils "after-sales/api/utils/sales-service"
	"net/http"

	"gorm.io/gorm"
)

type LicenseOwncerChangeRepository struct {
}

func OpenLicenseOwnerChangeRepositoryImpl() transactionworkshoprepository.LicenseOwncerChangeRepository {
	return &LicenseOwncerChangeRepository{}
}

// GetAll implements transactionworkshoprepository.LicenseOwncerChangeRepository.
func (l *LicenseOwncerChangeRepository) GetAll(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var entities []transactionworkshopentities.LicenseOwnerChange
	combinePayloads := make([]map[string]interface{}, 0)

	query := tx.Model(&transactionworkshopentities.LicenseOwnerChange{})

	for _, condition := range filterCondition {
		query = query.Where(condition.ColumnField+"= ?", condition.ColumnValue)
	}

	if err := query.Find(&entities).Error; err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	for _, entity := range entities {
		brandDetails, brandErr := salesserviceapiutils.GetUnitBrandById(entity.BrandId)
		if brandErr != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        brandErr.Err,
			}
		}

		modelDetails, modelErr := salesserviceapiutils.GetUnitModelById(entity.ModelId)
		if modelErr != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        modelErr.Err,
			}
		}

		variantDetails, variantErr := salesserviceapiutils.GetUnitVariantById(entity.VehicleId)
		if variantErr != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        variantErr.Err,
			}
		}

		response := map[string]interface{}{
			"license_owner_change_id": entity.LicenseOwnerChangeId,
			"brand_id":                entity.BrandId,
			"brand_name":              brandDetails.BrandName,
			"model_id":                entity.ModelId,
			"model_name":              modelDetails.ModelName,
			"variant_id":              variantDetails.VariantId,
			"variant_description":     variantDetails.VariantDescription,
			"vehicle_id":              entity.VehicleId,
			"change_date":             entity.ChangeDate,
			"change_type":             entity.ChangeType,
			"tnkb_old":                entity.TnkbOld,
			"tnkb_new":                entity.TnkbNew,
			"owner_name_old":          entity.OwnerNameOld,
			"owner_name_new":          entity.OwnerNameNew,
			"owner_address_old":       entity.OwnerAddressOld,
			"owner_address_new":       entity.OwnerAddressNew,
		}

		combinePayloads = append(combinePayloads, response)
	}

	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(combinePayloads, &pages)
	return paginatedData, totalPages, totalRows, nil
}

// GetAllHistory implements transactionworkshoprepository.LicenseOwncerChangeRepository.
func (l *LicenseOwncerChangeRepository) GetHistoryByChassisNumber(chassisNumber string, tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var entities []transactionworkshopentities.LicenseOwnerChange
	combinePayloads := make([]map[string]interface{}, 0)

	vehicleResponse, vehicleErr := salesserviceapiutils.GetVehicleByChassisNumber(chassisNumber)
	if vehicleErr != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid chassis number or vehicle not found",
			Err:        vehicleErr.Err,
		}
	}

	vehicleID := vehicleResponse.Data.Master.VehicleID

	query := tx.Model(&transactionworkshopentities.LicenseOwnerChange{}).
		Select("change_date, change_type, tnkb_old, tnkb_new, owner_name_old, owner_name_new").
		Where("vehicle_id = ?", vehicleID).
		Order("change_date DESC")

	if err := query.Find(&entities).Error; err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "No history data found for the given chassis number",
			Err:        err,
		}
	}

	for _, entity := range entities {
		response := map[string]interface{}{
			"change_date":    entity.ChangeDate,
			"change_type":    entity.ChangeType,
			"tnkb_old":       entity.TnkbOld,
			"tnkb_new":       entity.TnkbNew,
			"owner_name_old": entity.OwnerNameOld,
			"owner_name_new": entity.OwnerNameNew,
		}
		combinePayloads = append(combinePayloads, response)
	}

	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(combinePayloads, &pages)
	return paginatedData, totalPages, totalRows, nil
}
