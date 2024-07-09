package masterrepositoryimpl

import (
	"after-sales/api/config"
	masterentities "after-sales/api/entities/master"
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	"after-sales/api/utils"
	"database/sql"
	"errors"
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
	response.DealerId = entities.DealerId
	response.TopId = entities.TopId
	response.AgreementRemark = entities.AgreementRemark

	return response, nil
}

func (r *AgreementRepositoryImpl) SaveAgreement(tx *gorm.DB, req masterpayloads.AgreementRequest) (bool, *exceptions.BaseErrorResponse) {
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
		CustomerId:        req.CustomerId,
	}

	err := tx.Save(&entities).Error

	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return true, nil
}

func (r *AgreementRepositoryImpl) UpdateAgreement(tx *gorm.DB, Id int, req masterpayloads.AgreementRequest) (bool, *exceptions.BaseErrorResponse) {
	var entities masterentities.Agreement

	result := tx.Model(&entities).
		Where("agreement_id = ?", Id).
		Updates(map[string]interface{}{
			"agreement_code":      req.AgreementCode,
			"brand_id":            req.BrandId,
			"company_id":          req.DealerId,
			"top_id":              req.TopId,
			"agreement_date_from": req.AgreementDateFrom,
			"agreement_date_to":   req.AgreementDateTo,
			"agreement_remark":    req.AgreementRemark,
			"profit_center_id":    req.ProfitCenterId,
			"is_active":           req.IsActive,
			"customer_id":         req.CustomerId,
		})

	if result.Error != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	return true, nil
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

func (r *AgreementRepositoryImpl) GetAllAgreement(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {

	var responses []masterpayloads.AgreementRequest

	tableStruct := masterpayloads.AgreementRequest{}

	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)

	whereQuery := utils.ApplyFilter(joinTable, filterCondition)

	rows, err := whereQuery.Find(&responses).Rows()
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}
	defer rows.Close()

	var convertedResponses []masterpayloads.AgreementResponse

	for rows.Next() {

		var (
			AgreementReq masterpayloads.AgreementRequest
			AgreementRes masterpayloads.AgreementResponse
		)

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
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		// Fetch Customer data from external service
		CustomerURL := config.EnvConfigs.GeneralServiceUrl + "customer/" + strconv.Itoa(AgreementReq.CustomerId)
		fmt.Println("Fetching Customer data from:", CustomerURL)
		var getCustomerResponse masterpayloads.AgreementCustomerResponse
		if err := utils.Get(CustomerURL, &getCustomerResponse, nil); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		// Fetch Company data from external service
		CompanyURL := config.EnvConfigs.GeneralServiceUrl + "company/" + strconv.Itoa(AgreementReq.DealerId)
		fmt.Println("Fetching Company data from:", CompanyURL)
		var getCompanyResponse masterpayloads.AgreementCompanyResponse
		if err := utils.Get(CompanyURL, &getCompanyResponse, nil); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

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

		convertedResponses = append(convertedResponses, AgreementRes)
	}

	var mapResponses []map[string]interface{}

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

	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(mapResponses, &pages)

	return paginatedData, totalPages, totalRows, nil
}

func (r *AgreementRepositoryImpl) AddDiscountGroup(tx *gorm.DB, AgreementId int, req masterpayloads.DiscountGroupRequest) *exceptions.BaseErrorResponse {
	entities := masterentities.AgreementDiscountGroupDetail{
		AgreementId:               AgreementId,
		AgreementSelection:        req.AgreementSelection,
		AgreementOrderType:        req.AgreementLineTypeId,
		AgreementDiscountMarkupId: req.AgreementDiscountMarkup,
		AgreementDiscount:         req.AgreementDiscount,
		AgreementDetailRemarks:    req.AgreementDetailRemaks,
	}

	err := tx.Save(&entities).Error

	if err != nil {
		return &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return nil
}

func (r *AgreementRepositoryImpl) UpdateDiscountGroup(tx *gorm.DB, AgreementId int, DiscountGroupId int, req masterpayloads.DiscountGroupRequest) *exceptions.BaseErrorResponse {
	var entities masterentities.AgreementDiscountGroupDetail

	result := tx.Model(&entities).
		Where("agreement_id = ? AND agreement_discount_group_id = ?", AgreementId, DiscountGroupId).
		Updates(map[string]interface{}{
			"agreement_selection":          req.AgreementSelection,
			"agreement_line_type_id":       req.AgreementLineTypeId,
			"agreement_discount_markup_id": req.AgreementDiscountMarkup,
			"agreement_discount":           req.AgreementDiscount,
			"agreement_detail_remarks":     req.AgreementDetailRemaks,
		})

	if result.Error != nil {
		return &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	return nil
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

func (r *AgreementRepositoryImpl) AddItemDiscount(tx *gorm.DB, AgreementId int, req masterpayloads.ItemDiscountRequest) *exceptions.BaseErrorResponse {
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
		return &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return nil
}

func (r *AgreementRepositoryImpl) UpdateItemDiscount(tx *gorm.DB, AgreementId int, ItemDiscountId int, req masterpayloads.ItemDiscountRequest) *exceptions.BaseErrorResponse {
	var entities masterentities.AgreementItemDetail

	result := tx.Model(&entities).
		Where("agreement_id = ? AND agreement_item_id = ?", AgreementId, ItemDiscountId).
		Updates(map[string]interface{}{
			"line_type_id":                req.LineTypeId,
			"agreement_item_operation_id": req.AgreementItemOperationId,
			"discount_percent":            req.DiscountPercent,
			"min_value":                   req.MinValue,
			"agreement_remark":            req.AgreementRemark,
		})

	if result.Error != nil {
		return &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	return nil
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

func (r *AgreementRepositoryImpl) AddDiscountValue(tx *gorm.DB, AgreementId int, req masterpayloads.DiscountValueRequest) *exceptions.BaseErrorResponse {
	entities := masterentities.AgreementDiscount{
		AgreementId:     AgreementId,
		LineTypeId:      req.LineTypeId,
		MinValue:        req.MinValue,
		DiscountPercent: req.DiscountPercent,
		DiscountRemarks: req.DiscountRemarks,
	}

	err := tx.Save(&entities).Error

	if err != nil {
		return &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return nil
}

func (r *AgreementRepositoryImpl) UpdateDiscountValue(tx *gorm.DB, AgreementId int, DiscountValueId int, req masterpayloads.DiscountValueRequest) *exceptions.BaseErrorResponse {
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
		return &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	return nil
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

func (r *AgreementRepositoryImpl) GetAllDiscountGroup(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {

	var responses []masterpayloads.DiscountGroupRequest

	tableStruct := masterpayloads.DiscountGroupRequest{}

	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)

	whereQuery := utils.ApplyFilter(joinTable, filterCondition)

	rows, err := whereQuery.Find(&responses).Rows()
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}
	defer rows.Close()

	var mapResponses []map[string]interface{}

	for rows.Next() {

		var DiscountGroupRes masterpayloads.DiscountGroupResponse

		if err := rows.Scan(
			&DiscountGroupRes.AgreementDiscountGroupId,
			&DiscountGroupRes.AgreementId,
			&DiscountGroupRes.AgreementSelection,
			&DiscountGroupRes.AgreementLineTypeId,
			&DiscountGroupRes.AgreementDiscountMarkup,
			&DiscountGroupRes.AgreementDiscount,
			&DiscountGroupRes.AgreementDetailRemaks); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		responseMap := map[string]interface{}{
			"agreement_discount_group_id": DiscountGroupRes.AgreementDiscountGroupId,
			"agreement_id":                DiscountGroupRes.AgreementId,
			"agreement_selection":         DiscountGroupRes.AgreementSelection,
			"agreement_line_type_id":      DiscountGroupRes.AgreementLineTypeId,
			"agreement_discount_markup":   DiscountGroupRes.AgreementDiscountMarkup,
			"agreement_discount":          DiscountGroupRes.AgreementDiscount,
			"agreement_detail_remarks":    DiscountGroupRes.AgreementDetailRemaks,
		}

		mapResponses = append(mapResponses, responseMap)
	}

	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(mapResponses, &pages)

	return paginatedData, totalPages, totalRows, nil
}

func (r *AgreementRepositoryImpl) GetAllItemDiscount(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var responses []masterpayloads.ItemDiscountRequest

	tableStruct := masterpayloads.ItemDiscountRequest{}

	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)

	whereQuery := utils.ApplyFilter(joinTable, filterCondition)

	rows, err := whereQuery.Find(&responses).Rows()
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}
	defer rows.Close()

	var convertedResponses []masterpayloads.ItemDiscountResponse

	for rows.Next() {

		var ItemDiscountRes masterpayloads.ItemDiscountResponse
		var discountPercent, minValue sql.NullFloat64

		if err := rows.Scan(
			&ItemDiscountRes.AgreementItemId,
			&ItemDiscountRes.AgreementId,
			&ItemDiscountRes.LineTypeId,
			&ItemDiscountRes.AgreementItemOperationId,
			&discountPercent,
			&minValue,
			&ItemDiscountRes.AgreementRemark); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		if discountPercent.Valid {
			ItemDiscountRes.DiscountPercent = float32(discountPercent.Float64)
		} else {
			ItemDiscountRes.DiscountPercent = 0
		}
		if minValue.Valid {
			ItemDiscountRes.MinValue = int(minValue.Float64)
		} else {
			ItemDiscountRes.MinValue = 0
		}

		convertedResponses = append(convertedResponses, ItemDiscountRes)
	}

	var mapResponses []map[string]interface{}

	for _, response := range convertedResponses {
		responseMap := map[string]interface{}{
			"agreement_item_id":           response.AgreementItemId,
			"agreement_id":                response.AgreementId,
			"line_type_id":                response.LineTypeId,
			"agreement_item_operation_id": response.AgreementItemOperationId,
			"discount_percent":            response.DiscountPercent,
			"min_value":                   response.MinValue,
			"agreement_remark":            response.AgreementRemark,
		}
		mapResponses = append(mapResponses, responseMap)
	}

	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(mapResponses, &pages)

	return paginatedData, totalPages, totalRows, nil
}

func (r *AgreementRepositoryImpl) GetAllDiscountValue(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {

	var responses []masterpayloads.DiscountValueRequest

	tableStruct := masterpayloads.DiscountValueRequest{}

	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)

	whereQuery := utils.ApplyFilter(joinTable, filterCondition)

	rows, err := whereQuery.Find(&responses).Rows()
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}
	defer rows.Close()

	var convertedResponses []masterpayloads.DiscountValueResponse

	for rows.Next() {

		var DiscountValueRes masterpayloads.DiscountValueResponse

		if err := rows.Scan(
			&DiscountValueRes.AgreementDiscountId,
			&DiscountValueRes.AgreementId,
			&DiscountValueRes.LineTypeId,
			&DiscountValueRes.MinValue,
			&DiscountValueRes.DiscountPercent,
			&DiscountValueRes.DiscountRemarks); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		convertedResponses = append(convertedResponses, DiscountValueRes)
	}

	var mapResponses []map[string]interface{}

	for _, response := range convertedResponses {
		responseMap := map[string]interface{}{
			"agreement_discount_id": response.AgreementDiscountId,
			"agreement_id":          response.AgreementId,
			"line_type_id":          response.LineTypeId,
			"min_value":             response.MinValue,
			"discount_percent":      response.DiscountPercent,
			"discount_remarks":      response.DiscountRemarks,
		}
		mapResponses = append(mapResponses, responseMap)
	}

	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(mapResponses, &pages)

	return paginatedData, totalPages, totalRows, nil
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
	response.AgreementLineTypeId = entities.AgreementOrderType
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
