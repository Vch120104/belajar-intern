package masterrepositoryimpl

import (
	"after-sales/api/config"
	masterentities "after-sales/api/entities/master"
	exceptionsss_test "after-sales/api/expectionsss"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	"after-sales/api/utils"
	"fmt"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type AgreementRepositoryImpl struct {
}

func StartAgreementRepositoryImpl() masterrepository.AgreementRepository {
	return &AgreementRepositoryImpl{}
}

func (r *AgreementRepositoryImpl) GetAgreementById(tx *gorm.DB, AgreementId int) (masterpayloads.AgreementRequest, *exceptionsss_test.BaseErrorResponse) {
	entities := masterentities.Agreement{}
	response := masterpayloads.AgreementRequest{}

	err := tx.Model(&entities).
		Where(masterentities.Agreement{
			AgreementId: AgreementId,
		}).
		First(&entities).
		Error

	if err != nil {
		return response, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	// Copying values from entities to response
	response.AgreementId = entities.AgreementId
	response.AgreementCode = entities.AgreementCode
	response.IsActive = entities.IsActive
	response.BrandId = entities.BrandId
	response.CustomerId = entities.CustomerId
	response.ProfitCenterId = entities.ProfitCenterId
	response.AgreementDateFrom = entities.AgreementDateFrom
	response.AgreementDateTo = entities.AgreementDateTo
	response.DealerId = entities.DealerId
	response.TopId = entities.TopId
	response.AgreementRemark = entities.AgreementRemark

	return response, nil
}

func (r *AgreementRepositoryImpl) SaveAgreement(tx *gorm.DB, req masterpayloads.AgreementResponse) (bool, *exceptionsss_test.BaseErrorResponse) {
	entities := masterentities.Agreement{
		AgreementCode:     req.AgreementCode,
		BrandId:           req.BrandId,
		DealerId:          req.DealerId,
		TopId:             req.TopId,
		AgreementDateFrom: req.AgreementDateFrom,
		AgreementDateTo:   req.AgreementDateTo,
		AgreementRemark:   req.AgreementRemark,
		ProfitCenterId:    req.ProfitCenterId,
		IsActive:          req.IsActive,
		AgreementId:       req.AgreementId,
		CustomerId:        req.CustomerId,
	}

	err := tx.Save(&entities).Error

	if err != nil {
		return false, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return true, nil
}

func (r *AgreementRepositoryImpl) ChangeStatusAgreement(tx *gorm.DB, Id int) (bool, *exceptionsss_test.BaseErrorResponse) {
	var entities masterentities.Agreement

	result := tx.Model(&entities).
		Where("agreement_id = ?", Id).
		First(&entities)

	if result.Error != nil {
		return false, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	if entities.IsActive {
		entities.IsActive = false
	} else {
		entities.IsActive = true
	}

	result = tx.Save(&entities)

	if result.Error != nil {
		return false, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	return true, nil
}

func (r *AgreementRepositoryImpl) GetAllAgreement(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse) {
	// Define a slice to hold Agreement responses
	var responses []masterpayloads.AgreementRequest

	// Define table struct
	tableStruct := masterpayloads.AgreementRequest{}

	// Define join table
	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)

	// Apply filters
	whereQuery := utils.ApplyFilter(joinTable, filterCondition)

	// Execute query
	rows, err := whereQuery.Find(&responses).Rows()
	if err != nil {
		return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}
	defer rows.Close()

	// Define a slice to hold Agreement responses
	var convertedResponses []masterpayloads.AgreementResponse

	// Iterate over rows
	for rows.Next() {
		// Define variables to hold row data
		var (
			AgreementReq masterpayloads.AgreementRequest
			AgreementRes masterpayloads.AgreementResponse
		)

		// Scan the row into PurchasePriceRequest struct
		if err := rows.Scan(
			&AgreementReq.AgreementId,
			&AgreementReq.AgreementCode,
			&AgreementReq.IsActive,
			&AgreementReq.BrandId,
			&AgreementReq.CustomerId,
			&AgreementReq.ProfitCenterId,
			&AgreementReq.AgreementDateFrom,
			&AgreementReq.AgreementDateTo,
			&AgreementReq.DealerId,
			&AgreementReq.TopId,
			&AgreementReq.AgreementRemark); err != nil {
			return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		// Fetch Customer data from external service
		CustomerURL := config.EnvConfigs.GeneralServiceUrl + "api/general/customer/" + strconv.Itoa(AgreementReq.CustomerId)
		fmt.Println("Fetching Customer data from:", CustomerURL)
		var getCustomerResponse masterpayloads.AgreementCustomerResponse
		if err := utils.Get(CustomerURL, &getCustomerResponse, nil); err != nil {
			return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		// Fetch Company data from external service
		CompanyURL := config.EnvConfigs.GeneralServiceUrl + "api/general/company/" + strconv.Itoa(AgreementReq.DealerId)
		fmt.Println("Fetching Company data from:", CompanyURL)
		var getCompanyResponse masterpayloads.AgreementCompanyResponse
		if err := utils.Get(CompanyURL, &getCompanyResponse, nil); err != nil {
			return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
		// Create AgreementResponse
		AgreementRes = masterpayloads.AgreementResponse{
			AgreementId:       AgreementReq.AgreementId,
			AgreementCode:     AgreementReq.AgreementCode,
			IsActive:          AgreementReq.IsActive,
			BrandId:           AgreementReq.BrandId,
			CustomerId:        AgreementReq.CustomerId,
			CustomerName:      getCustomerResponse.CustomerName,
			CustomerCode:      getCustomerResponse.CustomerCode,
			ProfitCenterId:    AgreementReq.ProfitCenterId,
			AgreementDateFrom: AgreementReq.AgreementDateFrom,
			AgreementDateTo:   AgreementReq.AgreementDateTo,
			DealerId:          AgreementReq.DealerId,
			DealerName:        getCompanyResponse.CompanyName,
			DealerCode:        getCompanyResponse.CompanyCode,
			TopId:             AgreementReq.TopId,
			AgreementRemark:   AgreementReq.AgreementRemark,
		}

		// Append PurchasePriceResponse to the slice
		convertedResponses = append(convertedResponses, AgreementRes)
	}

	// Define a slice to hold map responses
	var mapResponses []map[string]interface{}

	// Iterate over convertedResponses and convert them to maps
	for _, response := range convertedResponses {
		responseMap := map[string]interface{}{
			"agreement_id":        response.AgreementId,
			"agreement_code":      response.AgreementCode,
			"customer_id":         response.CustomerId,
			"customer_name":       response.CustomerName,
			"customer_code":       response.CustomerCode,
			"profit_center_id":    response.ProfitCenterId,
			"agreement_date_from": response.AgreementDateFrom,
			"agreement_date_to":   response.AgreementDateTo,
			"dealer_id":           response.DealerId,
			"dealer_name":         response.DealerName,
			"dealer_code":         response.DealerCode,
			"top_id":              response.TopId,
			"agreement_remark":    response.AgreementRemark,
			"brand_id":            response.BrandId,
			"is_active":           response.IsActive,
		}
		mapResponses = append(mapResponses, responseMap)
	}

	// Paginate the response data
	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(mapResponses, &pages)

	return paginatedData, totalPages, totalRows, nil
}

func (r *AgreementRepositoryImpl) AddDiscountGroup(tx *gorm.DB, AgreementId int, req masterpayloads.DiscountGroupRequest) *exceptionsss_test.BaseErrorResponse {
	entities := masterentities.AgreementDiscountGroupDetail{
		AgreementId:                    AgreementId,
		AgreementDiscountGroupDetailId: req.AgreementDiscountGroupId,
		AgreementSelection:             req.AgreementSelection,
		AgreementOrderType:             req.AgreementLineTypeId,
		AgreementDiscountMarkupId:      req.AgreementDiscountMarkup,
		AgreementDiscount:              req.AgreementDiscount,
		AgreementDetailRemarks:         req.AgreementDetailRemaks,
	}

	err := tx.Save(&entities).Error

	if err != nil {
		return &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return nil
}

func (r *AgreementRepositoryImpl) DeleteDiscountGroup(tx *gorm.DB, AgreementId int, DiscountGroupId int) *exceptionsss_test.BaseErrorResponse {
	var entities masterentities.AgreementDiscountGroupDetail

	result := tx.Model(&entities).
		Where("agreement_id = ? AND agreement_discount_group_id = ?", AgreementId, DiscountGroupId).
		Delete(&entities)

	if result.Error != nil {
		return &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	return nil
}

func (r *AgreementRepositoryImpl) AddItemDiscount(tx *gorm.DB, AgreementId int, req masterpayloads.ItemDiscountRequest) *exceptionsss_test.BaseErrorResponse {
	entities := masterentities.AgreementItemDetail{
		AgreementId: AgreementId,
	}

	err := tx.Save(&entities).Error

	if err != nil {
		return &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return nil
}

func (r *AgreementRepositoryImpl) DeleteItemDiscount(tx *gorm.DB, AgreementId int, ItemDiscountId int) *exceptionsss_test.BaseErrorResponse {
	var entities masterentities.AgreementItemDetail

	result := tx.Model(&entities).
		Where("agreement_id = ? AND agreement_item_id = ?", AgreementId, ItemDiscountId).
		Delete(&entities)

	if result.Error != nil {
		return &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	return nil
}

func (r *AgreementRepositoryImpl) AddDiscountValue(tx *gorm.DB, AgreementId int, req masterpayloads.DiscountValueRequest) *exceptionsss_test.BaseErrorResponse {
	entities := masterentities.AgreementDiscount{
		AgreementId: AgreementId,
	}

	err := tx.Save(&entities).Error

	if err != nil {
		return &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return nil
}

func (r *AgreementRepositoryImpl) DeleteDiscountValue(tx *gorm.DB, AgreementId int, DiscountValueId int) *exceptionsss_test.BaseErrorResponse {
	var entities masterentities.AgreementDiscount

	result := tx.Model(&entities).
		Where("agreement_id = ? AND agreement_discount_id = ?", AgreementId, DiscountValueId).
		Delete(&entities)

	if result.Error != nil {
		return &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	return nil
}
