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
func (l *LicenseOwncerChangeRepository) GetAll(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var entities []transactionworkshopentities.LicenseOwnerChange
	combinePayloads := make([]map[string]interface{}, 0)

	query := tx.Model(&transactionworkshopentities.LicenseOwnerChange{})
	query = utils.ApplyFilter(query, filterCondition)

	err := query.Scopes(pagination.Paginate(&pages, query)).Find(&entities).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if len(entities) == 0 {
		pages.Rows = []map[string]interface{}{}
		return pages, nil
	}

	for _, entity := range entities {
		brandDetails, errBrand := salesserviceapiutils.GetUnitBrandById(entity.BrandId)
		if errBrand != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: errBrand.StatusCode,
				Err:        errBrand.Err,
			}
		}

		modelDetails, errModel := salesserviceapiutils.GetUnitModelById(entity.ModelId)
		if errModel != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: errModel.StatusCode,
				Err:        errModel.Err,
			}
		}

		variantDetails, errVariant := salesserviceapiutils.GetUnitVariantById(entity.VehicleId)
		if errVariant != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: errVariant.StatusCode,
				Err:        errVariant.Err,
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

	pages.Rows = combinePayloads
	return pages, nil
}

// GetAllHistory implements transactionworkshoprepository.LicenseOwncerChangeRepository.
func (l *LicenseOwncerChangeRepository) GetHistoryByChassisNumber(chassisNumber string, tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var entities []transactionworkshopentities.LicenseOwnerChange
	combinePayloads := make([]map[string]interface{}, 0)

	vehicleResponse, vehicleErr := salesserviceapiutils.GetVehicleByChassisNumber(chassisNumber)
	if vehicleErr != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid chassis number or vehicle not found",
			Err:        vehicleErr.Err,
		}
	}

	vehicleID := vehicleResponse.Data.Master.VehicleID

	query := tx.Model(&transactionworkshopentities.LicenseOwnerChange{}).
		Select("change_date, change_type, vehicle_stnk_tnkb_old AS tnkb_old, vehicle_stnk_tnkb_new AS tnkb_new, vehicle_owner_name_old AS owner_name_old, vehicle_owner_name_new AS owner_name_new").
		Where("vehicle_id = ?", vehicleID).
		Order("change_date DESC")

	query = utils.ApplyFilter(query, filterCondition)

	err := query.Scopes(pagination.Paginate(&pages, query)).Find(&entities).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "No history data found for the given chassis number",
			Err:        err,
		}
	}

	if len(entities) == 0 {
		pages.Rows = []map[string]interface{}{}
		return pages, nil
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
	pages.Rows = combinePayloads
	return pages, nil
}
