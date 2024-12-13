package masterrepositoryimpl

import (
	masterentities "after-sales/api/entities/master"
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	"after-sales/api/utils"
	aftersalesserviceapiutils "after-sales/api/utils/aftersales-service"
	generalserviceapiutils "after-sales/api/utils/general-service"
	"errors"
	"fmt"
	"math"
	"net/http"
	"time"

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
	response.CompanyId = entities.CompanyId
	response.TopId = entities.TopId
	response.AgreementRemark = entities.AgreementRemark

	return response, nil
}

func (r *AgreementRepositoryImpl) SaveAgreement(tx *gorm.DB, req masterpayloads.AgreementRequest) (masterentities.Agreement, *exceptions.BaseErrorResponse) {
	// Validasi: tidak ada AgreementCode yang sama dalam satu CompanyId
	var existingAgreement masterentities.Agreement
	err := tx.Where("agreement_code = ? AND company_id = ?", req.AgreementCode, req.CompanyId).First(&existingAgreement).Error

	if err == nil {
		return masterentities.Agreement{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Err:        errors.New("agreement code already exists in the same company"),
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return masterentities.Agreement{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	entities := masterentities.Agreement{
		AgreementCode:     req.AgreementCode,
		BrandId:           req.BrandId,
		CompanyId:         req.CompanyId,
		TopId:             req.TopId,
		AgreementDateFrom: req.AgreementDateFrom,
		AgreementDateTo:   req.AgreementDateTo,
		AgreementRemark:   req.AgreementRemark,
		ProfitCenterId:    req.ProfitCenterId,
		IsActive:          req.IsActive,
		CustomerId:        req.CustomerId,
	}

	err = tx.Save(&entities).Error
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

func (r *AgreementRepositoryImpl) GetAllAgreement(tx *gorm.DB, internalFilter []utils.FilterCondition, externalFilter []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	model := masterentities.Agreement{}
	var responses []masterpayloads.AgreementResponse

	var customerCode, customerName, profitCenterName string
	var dateFrom, dateTo string
	layoutDate := "2006-01-02" // Format parsing date

	// External filter processing for customer and profit center
	for _, filter := range externalFilter {
		switch filter.ColumnField {
		case "customer_code":
			customerCode = filter.ColumnValue
		case "customer_name":
			customerName = filter.ColumnValue
		case "profit_center_name":
			profitCenterName = filter.ColumnValue
		case "agreement_date_from":
			dateFrom = filter.ColumnValue
		case "agreement_date_to":
			dateTo = filter.ColumnValue
		}
	}

	query := tx.Model(&model).
		Select("mtr_agreement.*, cust.customer_name, cust.customer_code, comp.company_name, comp.company_code, pc.profit_center_name, tops.term_of_payment_code, tops.term_of_payment_name").
		Joins("JOIN dms_microservices_general_dev.dbo.mtr_customer cust ON mtr_agreement.customer_id = cust.customer_id").
		Joins("JOIN dms_microservices_general_dev.dbo.mtr_company comp ON mtr_agreement.company_id = comp.company_id").
		Joins("JOIN dms_microservices_general_dev.dbo.mtr_profit_center pc ON mtr_agreement.profit_center_id = pc.profit_center_id").
		Joins("JOIN dms_microservices_general_dev.dbo.mtr_term_of_payment tops ON mtr_agreement.top_id = tops.term_of_payment_id")

	// External filters processing
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

		if len(customerResponse) == 0 {

			internalFilter = append(internalFilter, utils.FilterCondition{
				ColumnField: "mtr_agreement.customer_id",
				ColumnValue: "-1",
			})
		} else if len(customerResponse) > 0 {
			if len(customerCode) > 0 && len(customerResponse) > 1 {
				query = query.Where("cust.customer_code LIKE ?", fmt.Sprintf("%%%s%%", customerCode))
			} else if len(customerName) > 0 && len(customerResponse) > 1 {
				query = query.Where("cust.customer_name LIKE ?", fmt.Sprintf("%%%s%%", customerName))
			} else {
				var customerIds []int
				for _, customer := range customerResponse {
					customerIds = append(customerIds, customer.CustomerId)
				}
				query = query.Where("mtr_agreement.customer_id IN (?)", customerIds)
			}
		} else {
			// Return empty rows if no customers found
			pages.Rows = []map[string]interface{}{}
			return pages, nil
		}
	}

	// External filter for Profit Center
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

		if len(profitCenterResponse) == 0 {

			internalFilter = append(internalFilter, utils.FilterCondition{
				ColumnField: "mtr_agreement.profit_center_id",
				ColumnValue: "-1",
			})
		} else if len(profitCenterResponse) > 0 {

			if len(profitCenterName) > 0 && len(profitCenterResponse) > 1 {
				query = query.Where("pc.profit_center_name LIKE ?", fmt.Sprintf("%%%s%%", profitCenterName))
			} else {

				var profitCenterIds []int
				for _, profitCenter := range profitCenterResponse {
					profitCenterIds = append(profitCenterIds, profitCenter.ProfitCenterId)
				}
				query = query.Where("mtr_agreement.profit_center_id IN (?)", profitCenterIds)
			}
		} else {

			pages.Rows = []map[string]interface{}{}
			return pages, nil
		}
	}

	if dateTo != "" {
		parsedDateTo, err := time.Parse(layoutDate, dateTo)
		if err != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    fmt.Sprintf("Invalid dateTo format: %v", dateTo),
				Err:        err,
			}
		}
		dateTo = parsedDateTo.Format(layoutDate)
		query = query.Where("mtr_agreement.agreement_date_to <= ?", dateTo)
	}

	if dateFrom != "" {
		parsedDateFrom, err := time.Parse(layoutDate, dateFrom)
		if err != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    fmt.Sprintf("Invalid dateFrom format: %v", dateFrom),
				Err:        err,
			}
		}
		dateFrom = parsedDateFrom.Format(layoutDate)
		query = query.Where("mtr_agreement.agreement_date_from >= ?", dateFrom)
	}

	whereQuery := utils.ApplyFilter(query, internalFilter)

	// Manually calculate total rows for pagination
	var totalRows int64
	err := whereQuery.Model(&model).Count(&totalRows).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	// Calculate pagination parameters
	offset := pages.GetOffset()
	limit := pages.GetLimit()

	// Apply pagination manually
	err = whereQuery.Offset(offset).Limit(limit).Order("mtr_agreement.agreement_id").Find(&responses).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(totalRows) / float64(limit)))
	pages.TotalRows = totalRows
	pages.TotalPages = totalPages

	// If no responses found, return empty rows
	if len(responses) == 0 {
		pages.Rows = []map[string]interface{}{}
		return pages, nil
	}

	var results []map[string]interface{}
	for _, response := range responses {
		result := map[string]interface{}{
			"agreement_id":        response.AgreementId,
			"agreement_code":      response.AgreementCode,
			"is_active":           response.IsActive,
			"brand_id":            response.BrandId,
			"customer_id":         response.CustomerId,
			"customer_name":       response.CustomerName,
			"customer_code":       response.CustomerCode,
			"profit_center_id":    response.ProfitCenterId,
			"profit_center_name":  response.ProfitCenterName,
			"agreement_date_from": response.AgreementDateFrom,
			"agreement_date_to":   response.AgreementDateTo,
			"company_id":          response.CompanyId,
			"top_id":              response.TopId,
			"top_code":            response.TermOfPaymentCode,
			"top_description":     response.TermOfPaymentName,
			"company_name":        response.CompanyName,
			"company_code":        response.CompanyCode,
		}

		results = append(results, result)
	}

	// Set the pagination info
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

		// fetch order type from utils cross service
		orderTypeResponse, orderTypeError := aftersalesserviceapiutils.GetOrderTypeById(entity.AgreementOrderType)
		if orderTypeError != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: orderTypeError.StatusCode,
				Err:        orderTypeError.Err,
			}
		}

		result := map[string]interface{}{
			"agreement_discount_group_id":  entity.AgreementDiscountGroupId,
			"agreement_id":                 entity.AgreementId,
			"agreement_selection":          entity.AgreementSelection,
			"agreement_selection_name":     selectionName,
			"agreement_order_type":         entity.AgreementOrderType,
			"agreement_order_type_name":    orderTypeResponse.OrderTypeName,
			"agreement_discount_markup_id": entity.AgreementDiscountMarkupId,
			"agreement_discount":           entity.AgreementDiscount,
			"agreement_detail_remarks":     entity.AgreementDetailRemarks,
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
		Where("agreement_discount_group_id = ? AND agreement_id = ?", DiscountGroupId, AgreementId).
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
		Where("agreement_item_id = ? AND agreement_id = ?", ItemDiscountId, AgreementId).
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
		Where("agreement_discount_id = ? AND agreement_id = ?", DiscountValueId, AgreementId).
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
	customerResponse, custErr := generalserviceapiutils.GetCustomerMasterById(entities.CustomerId)
	if custErr != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        custErr.Err,
		}
	}

	// fetch data company from utils cross service
	companyResponse, compErr := generalserviceapiutils.GetCompanyDataById(entities.CompanyId)
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
	response.CompanyId = entities.CompanyId
	response.CompanyName = companyResponse.CompanyName
	response.CompanyCode = companyResponse.CompanyCode
	response.TopId = entities.TopId
	response.AgreementRemark = entities.AgreementRemark

	return response, nil
}

func (r *AgreementRepositoryImpl) GetDiscountGroupAgreementByHeaderId(tx *gorm.DB, AgreementId int, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	entities := []masterentities.AgreementDiscountGroupDetail{}

	baseModelQuery := tx.Model(&masterentities.AgreementDiscountGroupDetail{})

	whereQuery := utils.ApplyFilter(baseModelQuery, filterCondition)

	whereQuery = whereQuery.Where("agreement_id = ?", AgreementId)

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

		// fetch order type from utils cross service
		orderTypeResponse, orderTypeError := aftersalesserviceapiutils.GetOrderTypeById(entity.AgreementOrderType)
		if orderTypeError != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: orderTypeError.StatusCode,
				Err:        orderTypeError.Err,
			}
		}

		result := map[string]interface{}{
			"agreement_discount_group_id":  entity.AgreementDiscountGroupId,
			"agreement_id":                 entity.AgreementId,
			"agreement_selection":          entity.AgreementSelection,
			"agreement_selection_name":     selectionName,
			"agreement_order_type":         entity.AgreementOrderType,
			"agreement_order_type_name":    orderTypeResponse.OrderTypeName,
			"agreement_discount_markup_id": entity.AgreementDiscountMarkupId,
			"agreement_discount":           entity.AgreementDiscount,
			"agreement_detail_remarks":     entity.AgreementDetailRemarks,
		}
		results = append(results, result)
	}

	pages.Rows = results

	return pages, nil
}

func (r *AgreementRepositoryImpl) GetDiscountItemAgreementByHeaderId(tx *gorm.DB, AgreementId int, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	entities := []masterentities.AgreementItemDetail{}

	baseModelQuery := tx.Model(&masterentities.AgreementItemDetail{})

	whereQuery := utils.ApplyFilter(baseModelQuery, filterCondition)

	whereQuery = whereQuery.Where("agreement_id = ?", AgreementId)

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

func (r *AgreementRepositoryImpl) GetDiscountValueAgreementByHeaderId(tx *gorm.DB, AgreementId int, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {

	entities := []masterentities.AgreementDiscount{}

	baseModelQuery := tx.Model(&masterentities.AgreementDiscount{})

	whereQuery := utils.ApplyFilter(baseModelQuery, filterCondition)

	whereQuery = whereQuery.Where("agreement_id = ?", AgreementId)

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

func (r *AgreementRepositoryImpl) DeleteMultiIdDiscountGroup(tx *gorm.DB, AgreementId int, DiscountGroupIds []int) (bool, *exceptions.BaseErrorResponse) {
	var entities masterentities.AgreementDiscountGroupDetail
	result := tx.Model(&entities).
		Where("agreement_id = ? AND agreement_discount_group_id IN (?)", AgreementId, DiscountGroupIds).
		Delete(&entities)

	if result.Error != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error deleting discount group",
			Err:        result.Error,
		}
	}

	return true, nil
}

func (r *AgreementRepositoryImpl) DeleteMultiIdItemDiscount(tx *gorm.DB, AgreementId int, ItemDiscountIds []int) (bool, *exceptions.BaseErrorResponse) {
	var entities masterentities.AgreementItemDetail
	result := tx.Model(&entities).
		Where("agreement_id = ? AND agreement_item_id IN (?)", AgreementId, ItemDiscountIds).
		Delete(&entities)

	if result.Error != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error deleting item discount",
			Err:        result.Error,
		}
	}

	return true, nil
}

func (r *AgreementRepositoryImpl) DeleteMultiIdDiscountValue(tx *gorm.DB, AgreementId int, DiscountValueIds []int) (bool, *exceptions.BaseErrorResponse) {
	var entities masterentities.AgreementDiscount
	result := tx.Model(&entities).
		Where("agreement_id = ? AND agreement_discount_id IN (?)", AgreementId, DiscountValueIds).
		Delete(&entities)

	if result.Error != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error deleting discount value",
			Err:        result.Error,
		}
	}

	return true, nil
}
