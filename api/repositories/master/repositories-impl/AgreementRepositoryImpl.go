package masterrepositoryimpl

import (
	masterentities "after-sales/api/entities/master"
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	"after-sales/api/utils"
	generalserviceapiutils "after-sales/api/utils/general-service"
	"errors"
	"fmt"
	"net/http"

	"gorm.io/gorm"
)

type AgreementRepositoryImpl struct {
}

func StartAgreementRepositoryImpl() masterrepository.AgreementRepository {
	return &AgreementRepositoryImpl{}
}

func (r *AgreementRepositoryImpl) GetAgreementById(tx *gorm.DB, AgreementId int) (masterpayloads.AgreementRequest, *exceptions.BaseErrorResponse) {
	entities := masterentities.Agreement{}
	response := masterpayloads.AgreementRequest{}

	err := tx.Model(&entities).
		Where(masterentities.Agreement{
			AgreementId: AgreementId,
		}).
		First(&entities).
		Error

	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	response.AgreementId = entities.AgreementId
	response.AgreementCode = entities.AgreementCode
	response.IsActive = entities.IsActive
	response.BrandId = entities.BrandId
	response.CustomerId = entities.CustomerId
	response.ProfitCenterId = entities.ProfitCenterId
	response.AgreementDateFrom = entities.AgreementDateFrom
	response.AgreementDateTo = entities.AgreementDateTo
	response.CompanyId = entities.DealerId
	response.TopId = entities.TopId
	response.AgreementRemark = entities.AgreementRemark

	return response, nil
}

func (r *AgreementRepositoryImpl) SaveAgreement(tx *gorm.DB, req masterpayloads.AgreementRequest) (masterentities.Agreement, *exceptions.BaseErrorResponse) {
	entities := masterentities.Agreement{
		AgreementCode:     req.AgreementCode,
		BrandId:           req.BrandId,
		DealerId:          req.CompanyId,
		TopId:             req.TopId,
		AgreementDateFrom: req.AgreementDateFrom,
		AgreementDateTo:   req.AgreementDateTo,
		AgreementRemark:   req.AgreementRemark,
		ProfitCenterId:    req.ProfitCenterId,
		IsActive:          req.IsActive,
		CustomerId:        req.CustomerId,
	}

	err := tx.Save(&entities).Error

	if err != nil {
		return masterentities.Agreement{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return entities, nil
}

func (r *AgreementRepositoryImpl) UpdateAgreement(tx *gorm.DB, Id int, req masterpayloads.AgreementRequest) (masterentities.Agreement, *exceptions.BaseErrorResponse) {
	var entities masterentities.Agreement

	result := tx.Model(&entities).
		Where("agreement_id = ?", Id).
		Updates(map[string]interface{}{
			"aggreement_id":       Id,
			"agreement_code":      req.AgreementCode,
			"brand_id":            req.BrandId,
			"company_id":          req.CompanyId,
			"top_id":              req.TopId,
			"agreement_date_from": req.AgreementDateFrom,
			"agreement_date_to":   req.AgreementDateTo,
			"agreement_remark":    req.AgreementRemark,
			"profit_center_id":    req.ProfitCenterId,
			"is_active":           req.IsActive,
			"customer_id":         req.CustomerId,
		})

	if result.Error != nil {
		return masterentities.Agreement{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	return entities, nil
}

func (r *AgreementRepositoryImpl) ChangeStatusAgreement(tx *gorm.DB, Id int) (masterentities.Agreement, *exceptions.BaseErrorResponse) {
	var entities masterentities.Agreement

	result := tx.Model(&entities).
		Where("agreement_id = ?", Id).
		First(&entities)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return masterentities.Agreement{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        fmt.Errorf("agreement with ID %d not found", Id),
			}
		}

		return masterentities.Agreement{}, &exceptions.BaseErrorResponse{
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
		return masterentities.Agreement{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	return entities, nil
}

func (r *AgreementRepositoryImpl) GetAllAgreement(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var responses []masterentities.Agreement
	var entities masterentities.Agreement

	baseModelQuery := tx.Model(&entities)

	baseModelQuery = utils.ApplyFilter(baseModelQuery, filterCondition)

	var customerCode, customerName, profitCenterName string

	for _, filter := range filterCondition {
		switch filter.ColumnField {
		case "customer_code":
			customerCode = filter.ColumnValue
		case "customer_name":
			customerName = filter.ColumnValue
		case "profit_center_name":
			profitCenterName = filter.ColumnValue
		}
	}

	// External filter untuk Customer
	if customerCode != "" || customerName != "" {
		customerParams := generalserviceapiutils.CustomerMasterParams{
			CustomerCode: customerCode,
			CustomerName: customerName,
		}

		customerResponse, customerError := generalserviceapiutils.GetAllCustomerMaster(customerParams)
		if customerError != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: customerError.StatusCode,
				Message:    "Error fetching customer data",
				Err:        customerError.Err,
			}
		}

		var customerIds []int
		for _, customer := range customerResponse {
			customerIds = append(customerIds, customer.CustomerId)
		}

		if len(customerIds) > 0 {
			baseModelQuery = baseModelQuery.Where("customer_id IN ?", customerIds)
		} else {
			pages.Rows = []map[string]interface{}{}
			return pages, nil
		}
	}

	// External filter untuk Profit Center
	if profitCenterName != "" {
		profitCenterParams := generalserviceapiutils.ProfitCenterParams{
			ProfitCenterName: profitCenterName,
		}

		profitCenterResponse, profitCenterError := generalserviceapiutils.GetAllProfitCenter(profitCenterParams)
		if profitCenterError != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: profitCenterError.StatusCode,
				Message:    "Error fetching profit center data",
				Err:        profitCenterError.Err,
			}
		}

		var profitCenterIds []int
		for _, profitCenter := range profitCenterResponse {
			profitCenterIds = append(profitCenterIds, profitCenter.ProfitCenterId)
		}

		if len(profitCenterIds) > 0 {
			baseModelQuery = baseModelQuery.Where("profit_center_id IN ?", profitCenterIds)
		} else {
			pages.Rows = []map[string]interface{}{}
			return pages, nil
		}
	}

	// Apply pagination scope and execute the query
	err := baseModelQuery.Scopes(pagination.Paginate(&pages, baseModelQuery)).Find(&responses).Error
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

	// Map the query results to the response structure
	var results []map[string]interface{}
	for _, response := range responses {
		// Fetch customer data
		getCustomerResponse, custErr := generalserviceapiutils.GetCustomerMasterByID(response.CustomerId)
		if custErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        custErr.Err,
			}
		}

		// Fetch company data
		getCompanyResponse, compErr := generalserviceapiutils.GetCompanyDataById(response.DealerId)
		if compErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        compErr.Err,
			}
		}

		// Fetch profit center data
		profitCenterResponse, profitCenterErr := generalserviceapiutils.GetProfitCenterById(response.ProfitCenterId)
		if profitCenterErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        profitCenterErr.Err,
			}
		}

		termOfPaymentResponse, termOfPaymentErr := generalserviceapiutils.GetTermOfPaymentById(response.TopId)
		if termOfPaymentErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        termOfPaymentErr.Err,
			}
		}

		result := map[string]interface{}{
			"agreement_id":        response.AgreementId,
			"agreement_code":      response.AgreementCode,
			"is_active":           response.IsActive,
			"brand_id":            response.BrandId,
			"customer_id":         response.CustomerId,
			"customer_name":       getCustomerResponse.CustomerName,
			"customer_code":       getCustomerResponse.CustomerCode,
			"address_street_1":    getCustomerResponse.AddressStreet1,
			"address_street_2":    getCustomerResponse.AddressStreet2,
			"address_street_3":    getCustomerResponse.AddressStreet3,
			"village_name":        getCustomerResponse.VillageName,
			"village_zip_code":    getCustomerResponse.VillageZipCode,
			"district_name":       getCustomerResponse.DistrictName,
			"city_name":           getCustomerResponse.CityName,
			"city_phone_area":     getCustomerResponse.CityPhoneArea,
			"province_name":       getCustomerResponse.ProvinceName,
			"country_name":        getCustomerResponse.CountryName,
			"profit_center_id":    response.ProfitCenterId,
			"profit_center_name":  profitCenterResponse.ProfitCenterName,
			"agreement_date_from": response.AgreementDateFrom,
			"agreement_date_to":   response.AgreementDateTo,
			"dealer_id":           response.DealerId,
			"dealer_name":         getCompanyResponse.CompanyName,
			"dealer_code":         getCompanyResponse.CompanyCode,
			"top_id":              response.TopId,
			"top_code":            termOfPaymentResponse.TermOfPaymentCode,
			"top_description":     termOfPaymentResponse.TermOfPaymentName,
		}
		results = append(results, result)
	}

	pages.Rows = results
	return pages, nil
}

func (r *AgreementRepositoryImpl) AddDiscountGroup(tx *gorm.DB, AgreementId int, req masterpayloads.DiscountGroupRequest) (masterentities.AgreementDiscountGroupDetail, *exceptions.BaseErrorResponse) {

	// Validasi AgreementSelection
	if req.AgreementSelection != 0 && req.AgreementSelection != 1 {
		return masterentities.AgreementDiscountGroupDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("agreement Selection not valid"),
		}
	}

	entities := masterentities.AgreementDiscountGroupDetail{
		AgreementId:               AgreementId,
		AgreementSelection:        req.AgreementSelection,
		AgreementOrderType:        req.AgreementOrderTypeId,
		AgreementDiscountMarkupId: req.AgreementDiscountMarkup,
		AgreementDiscount:         req.AgreementDiscount,
		AgreementDetailRemarks:    req.AgreementDetailRemaks,
	}

	err := tx.Save(&entities).Error

	if err != nil {
		return masterentities.AgreementDiscountGroupDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return entities, nil
}

func (r *AgreementRepositoryImpl) UpdateDiscountGroup(tx *gorm.DB, AgreementId int, DiscountGroupId int, req masterpayloads.DiscountGroupRequest) (masterentities.AgreementDiscountGroupDetail, *exceptions.BaseErrorResponse) {
	var entities masterentities.AgreementDiscountGroupDetail

	result := tx.Model(&entities).
		Where("agreement_id = ? AND agreement_discount_group_id = ?", AgreementId, DiscountGroupId).
		Updates(map[string]interface{}{
			"agreement_selection_id":       req.AgreementSelection,
			"agreement_order_type_id":      req.AgreementOrderTypeId,
			"agreement_discount_markup_id": req.AgreementDiscountMarkup,
			"agreement_discount":           req.AgreementDiscount,
			"agreement_detail_remarks":     req.AgreementDetailRemaks,
			"agreement_id":                 AgreementId,
			"agreement_discount_group_id":  DiscountGroupId,
		})

	if result.Error != nil {
		return masterentities.AgreementDiscountGroupDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	return entities, nil
}

func (r *AgreementRepositoryImpl) DeleteDiscountGroup(tx *gorm.DB, AgreementId int, DiscountGroupId int) *exceptions.BaseErrorResponse {
	var entities masterentities.AgreementDiscountGroupDetail

	result := tx.Model(&entities).
		Where("agreement_id = ? AND agreement_discount_group_id = ?", AgreementId, DiscountGroupId).
		Delete(&entities)

	if result.Error != nil {
		return &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	return nil
}

func (r *AgreementRepositoryImpl) AddItemDiscount(tx *gorm.DB, AgreementId int, req masterpayloads.ItemDiscountRequest) (masterentities.AgreementItemDetail, *exceptions.BaseErrorResponse) {
	entities := masterentities.AgreementItemDetail{
		AgreementId:              AgreementId,
		LineTypeId:               req.LineTypeId,
		AgreementItemOperationId: req.AgreementItemOperationId,
		DiscountPercent:          req.DiscountPercent,
		MinValue:                 req.MinValue,
		AgreementRemark:          req.AgreementRemark,
	}

	err := tx.Save(&entities).Error

	if err != nil {
		return masterentities.AgreementItemDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return entities, nil
}

func (r *AgreementRepositoryImpl) UpdateItemDiscount(tx *gorm.DB, AgreementId int, ItemDiscountId int, req masterpayloads.ItemDiscountRequest) (masterentities.AgreementItemDetail, *exceptions.BaseErrorResponse) {
	var entities masterentities.AgreementItemDetail

	result := tx.Model(&entities).
		Where("agreement_id = ? AND agreement_item_id = ?", AgreementId, ItemDiscountId).
		Updates(map[string]interface{}{
			"line_type_id":                req.LineTypeId,
			"agreement_item_operation_id": req.AgreementItemOperationId,
			"discount_percent":            req.DiscountPercent,
			"min_value":                   req.MinValue,
			"agreement_remark":            req.AgreementRemark,
			"agreement_id":                AgreementId,
		})

	if result.Error != nil {
		return masterentities.AgreementItemDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	return entities, nil
}

func (r *AgreementRepositoryImpl) DeleteItemDiscount(tx *gorm.DB, AgreementId int, ItemDiscountId int) *exceptions.BaseErrorResponse {
	var entities masterentities.AgreementItemDetail

	result := tx.Model(&entities).
		Where("agreement_id = ? AND agreement_item_id = ?", AgreementId, ItemDiscountId).
		Delete(&entities)

	if result.Error != nil {
		return &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	return nil
}

func (r *AgreementRepositoryImpl) AddDiscountValue(tx *gorm.DB, AgreementId int, req masterpayloads.DiscountValueRequest) (masterentities.AgreementDiscount, *exceptions.BaseErrorResponse) {
	entities := masterentities.AgreementDiscount{
		AgreementId:     AgreementId,
		LineTypeId:      req.LineTypeId,
		MinValue:        req.MinValue,
		DiscountPercent: req.DiscountPercent,
		DiscountRemarks: req.DiscountRemarks,
	}

	err := tx.Save(&entities).Error

	if err != nil {
		return masterentities.AgreementDiscount{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return entities, nil
}

func (r *AgreementRepositoryImpl) UpdateDiscountValue(tx *gorm.DB, AgreementId int, DiscountValueId int, req masterpayloads.DiscountValueRequest) (masterentities.AgreementDiscount, *exceptions.BaseErrorResponse) {
	var entities masterentities.AgreementDiscount

	result := tx.Model(&entities).
		Where("agreement_id = ? AND agreement_discount_id = ?", AgreementId, DiscountValueId).
		Updates(map[string]interface{}{
			"line_type_id":     req.LineTypeId,
			"min_value":        req.MinValue,
			"discount_percent": req.DiscountPercent,
			"discount_remarks": req.DiscountRemarks,
		})

	if result.Error != nil {
		return masterentities.AgreementDiscount{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	return entities, nil
}

func (r *AgreementRepositoryImpl) DeleteDiscountValue(tx *gorm.DB, AgreementId int, DiscountValueId int) *exceptions.BaseErrorResponse {
	var entities masterentities.AgreementDiscount

	result := tx.Model(&entities).
		Where("agreement_id = ? AND agreement_discount_id = ?", AgreementId, DiscountValueId).
		Delete(&entities)

	if result.Error != nil {
		return &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	return nil
}

func (r *AgreementRepositoryImpl) GetAllDiscountGroup(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {

	entities := []masterentities.AgreementDiscountGroupDetail{}

	baseModelQuery := tx.Model(&masterentities.AgreementDiscountGroupDetail{})

	whereQuery := utils.ApplyFilter(baseModelQuery, filterCondition)

	err := whereQuery.Scopes(pagination.Paginate(&pages, whereQuery)).Find(&entities).Error
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

	var results []map[string]interface{}
	for _, entity := range entities {
		selectionName := "unknown"
		if entity.AgreementSelection == 0 {
			selectionName = "discount"
		} else if entity.AgreementSelection == 1 {
			selectionName = "markup"
		}

		result := map[string]interface{}{
			"agreement_discount_group_id":  entity.AgreementDiscountGroupId,
			"agreement_id":                 entity.AgreementId,
			"agreement_selection":          entity.AgreementSelection,
			"agreement_order_type":         entity.AgreementOrderType,
			"agreement_discount_markup_id": entity.AgreementDiscountMarkupId,
			"agreement_discount":           entity.AgreementDiscount,
			"agreement_detail_remarks":     entity.AgreementDetailRemarks,
			"selection_name":               selectionName,
		}
		results = append(results, result)
	}

	pages.Rows = results

	return pages, nil
}

func (r *AgreementRepositoryImpl) GetAllItemDiscount(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	entities := []masterentities.AgreementItemDetail{}

	baseModelQuery := tx.Model(&masterentities.AgreementItemDetail{})

	whereQuery := utils.ApplyFilter(baseModelQuery, filterCondition)

	err := whereQuery.Scopes(pagination.Paginate(&pages, whereQuery)).Find(&entities).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if len(entities) == 0 {
		pages.Rows = []masterentities.AgreementItemDetail{}
		return pages, nil
	}

	pages.Rows = entities

	return pages, nil
}

func (r *AgreementRepositoryImpl) GetAllDiscountValue(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {

	entities := []masterentities.AgreementDiscount{}

	baseModelQuery := tx.Model(&masterentities.AgreementDiscount{})

	whereQuery := utils.ApplyFilter(baseModelQuery, filterCondition)

	err := whereQuery.Scopes(pagination.Paginate(&pages, whereQuery)).Find(&entities).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if len(entities) == 0 {
		pages.Rows = []masterentities.AgreementDiscount{}
		return pages, nil
	}

	pages.Rows = entities

	return pages, nil
}

func (r *AgreementRepositoryImpl) GetDiscountGroupAgreementById(tx *gorm.DB, DiscountGroupId, AgreementId int) (masterpayloads.DiscountGroupRequest, *exceptions.BaseErrorResponse) {
	entities := masterentities.AgreementDiscountGroupDetail{}
	response := masterpayloads.DiscountGroupRequest{}

	err := tx.Model(&entities).
		Where(masterentities.AgreementDiscountGroupDetail{
			AgreementDiscountGroupId: DiscountGroupId,
			AgreementId:              AgreementId,
		}).
		First(&entities).
		Error

	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	response.AgreementId = entities.AgreementId
	response.AgreementSelection = entities.AgreementSelection
	response.AgreementOrderTypeId = entities.AgreementOrderType
	response.AgreementDiscountMarkup = entities.AgreementDiscountMarkupId
	response.AgreementDiscount = entities.AgreementDiscount
	response.AgreementDetailRemaks = entities.AgreementDetailRemarks

	return response, nil
}

func (r *AgreementRepositoryImpl) GetDiscountItemAgreementById(tx *gorm.DB, ItemDiscountId, AgreementId int) (masterpayloads.ItemDiscountRequest, *exceptions.BaseErrorResponse) {
	entities := masterentities.AgreementItemDetail{}
	response := masterpayloads.ItemDiscountRequest{}

	err := tx.Model(&entities).
		Where(masterentities.AgreementItemDetail{
			AgreementItemId: ItemDiscountId,
			AgreementId:     AgreementId,
		}).
		First(&entities).
		Error

	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	response.AgreementItemId = entities.AgreementItemId
	response.AgreementId = entities.AgreementId
	response.LineTypeId = entities.LineTypeId
	response.AgreementItemOperationId = entities.AgreementItemOperationId
	response.DiscountPercent = entities.DiscountPercent
	response.MinValue = entities.MinValue
	response.AgreementRemark = entities.AgreementRemark

	return response, nil
}

func (r *AgreementRepositoryImpl) GetDiscountValueAgreementById(tx *gorm.DB, DiscountValueId, AgreementId int) (masterpayloads.DiscountValueRequest, *exceptions.BaseErrorResponse) {
	entities := masterentities.AgreementDiscount{}
	response := masterpayloads.DiscountValueRequest{}

	err := tx.Model(&entities).
		Where(masterentities.AgreementDiscount{
			AgreementDiscountId: DiscountValueId,
			AgreementId:         AgreementId,
		}).
		First(&entities).
		Error

	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	response.AgreementDiscountId = entities.AgreementDiscountId
	response.AgreementId = entities.AgreementId
	response.LineTypeId = entities.LineTypeId
	response.MinValue = entities.MinValue
	response.DiscountPercent = entities.DiscountPercent
	response.DiscountRemarks = entities.DiscountRemarks

	return response, nil
}

func (r *AgreementRepositoryImpl) GetAgreementByCode(tx *gorm.DB, AgreementCode string) (masterpayloads.AgreementResponse, *exceptions.BaseErrorResponse) {
	entities := masterentities.Agreement{}
	response := masterpayloads.AgreementResponse{}

	err := tx.Model(&entities).
		Where(masterentities.Agreement{
			AgreementCode: AgreementCode,
		}).
		First(&entities).
		Error

	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	// fetch data customer from utils cross service
	customerResponse, custErr := generalserviceapiutils.GetCustomerMasterByID(entities.CustomerId)
	if custErr != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        custErr.Err,
		}
	}

	// fetch data company from utils cross service
	companyResponse, compErr := generalserviceapiutils.GetCompanyDataById(entities.DealerId)
	if compErr != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        compErr.Err,
		}
	}

	response.AgreementId = entities.AgreementId
	response.AgreementCode = entities.AgreementCode
	response.IsActive = entities.IsActive
	response.BrandId = entities.BrandId
	response.CustomerId = entities.CustomerId
	response.CustomerName = customerResponse.CustomerName
	response.CustomerCode = customerResponse.CustomerCode
	response.AddressStreet1 = customerResponse.AddressStreet1
	response.AddressStreet2 = customerResponse.AddressStreet2
	response.AddressStreet3 = customerResponse.AddressStreet3
	response.VillageName = customerResponse.VillageName
	response.VillageZipCode = customerResponse.VillageZipCode
	response.DistrictName = customerResponse.DistrictName
	response.CityName = customerResponse.CityName
	response.ProvinceName = customerResponse.ProvinceName
	response.CountryName = customerResponse.CountryName
	response.ProfitCenterId = entities.ProfitCenterId
	response.AgreementDateFrom = entities.AgreementDateFrom
	response.AgreementDateTo = entities.AgreementDateTo
	response.DealerId = entities.DealerId
	response.DealerName = companyResponse.CompanyName
	response.DealerCode = companyResponse.CompanyCode
	response.TopId = entities.TopId
	response.AgreementRemark = entities.AgreementRemark

	return response, nil
}
