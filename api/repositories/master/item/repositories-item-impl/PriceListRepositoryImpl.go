package masteritemrepositoryimpl

import (
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	"after-sales/api/utils"
	aftersalesserviceapiutils "after-sales/api/utils/aftersales-service"
	financeserviceapiutils "after-sales/api/utils/finance-service"
	salesserviceapiutils "after-sales/api/utils/sales-service"
	"errors"
	"math"
	"net/http"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type PriceListRepositoryImpl struct {
}

func StartPriceListRepositoryImpl() masteritemrepository.PriceListRepository {
	return &PriceListRepositoryImpl{}
}

// Duplicate implements masteritemrepository.PriceListRepository.
func (r *PriceListRepositoryImpl) Duplicate(tx *gorm.DB, itemGroupId int, brandId int, currencyId int, date string) ([]masteritempayloads.PriceListItemResponses, *exceptions.BaseErrorResponse) {
	model := masteritementities.ItemPriceList{}

	result := []masteritempayloads.PriceListItemResponses{}

	if err := tx.Model(model).Select("mtr_item.item_code, mtr_item.item_name,mtr_item_price_list.price_list_amount,mtr_item_price_list.is_active,mtr_item.item_id,mtr_item.item_class_id").
		Joins("LEFT JOIN mtr_item ON mtr_item_price_list.item_id = mtr_item.item_id").
		Where(masteritementities.ItemPriceList{ItemGroupId: itemGroupId, BrandId: brandId, CurrencyId: currencyId}).
		Where("CONVERT(DATE, mtr_item_price_list.effective_date) like ?", date).
		Scan(&result).Error; err != nil {
		return result, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New("price list item not found"),
		}
	}

	return result, nil
}

// CheckPriceListItem implements masteritemrepository.PriceListRepository.
func (r *PriceListRepositoryImpl) CheckPriceListItem(tx *gorm.DB, itemGroupId int, brandId int, currencyId int, date string, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	model := masteritementities.ItemPriceList{}

	result := []masteritempayloads.PriceListItemResponses{}

	query := tx.Model(model).Select("mtr_item.item_code, mtr_item.item_name,mtr_item_price_list.price_list_amount,mtr_item_price_list.is_active,mtr_item.item_id,mtr_item.item_class_id,mtr_item_price_list.price_list_id").
		Joins("LEFT JOIN mtr_item ON mtr_item_price_list.item_id = mtr_item.item_id").
		Where(masteritementities.ItemPriceList{ItemGroupId: itemGroupId, BrandId: brandId, CurrencyId: currencyId}).
		Where("CONVERT(DATE, mtr_item_price_list.effective_date) like ?", date)

	if err := query.Scopes(pagination.Paginate(&pages, query)).Scan(&result).Error; err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New("price list item not found"),
		}
	}

	pages.Rows = result

	return pages, nil
}

// CheckPriceListAlreadyExist implements masteritemrepository.PriceListRepository.
func (r *PriceListRepositoryImpl) CheckPriceListExist(tx *gorm.DB, itemId int, brandId int, currencyId int, date string, companyId int) (bool, *exceptions.BaseErrorResponse) {
	model := masteritementities.ItemPriceList{}

	if err := tx.Model(&model).Where(masteritementities.ItemPriceList{BrandId: brandId, ItemId: itemId}).
		Where("mtr_item_price_list.company_id = ?", companyId).
		Where("CONVERT(DATE, mtr_item_price_list.effective_date) = ?", date).First(&model).Error; err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
			Message:    "failed to check price list is exist",
		}
	}

	return true, nil
}

func (r *PriceListRepositoryImpl) GetPriceListLookup(tx *gorm.DB, request masteritempayloads.PriceListGetAllRequest) ([]masteritempayloads.PriceListResponse, *exceptions.BaseErrorResponse) {
	var responses []masteritempayloads.PriceListResponse

	tempRows := tx.
		Model(&masteritementities.ItemPriceList{})

	if request.CompanyId != 0 {
		tempRows = tempRows.Where("company_id = ?", request.CompanyId)
	}

	if request.PriceListCode != "" {
		tempRows = tempRows.Where("price_list_code_id like ?", request.PriceListCode)
	}

	if request.BrandId != 0 {
		tempRows = tempRows.Where("brand_id = ?", request.BrandId)
	}

	if request.CurrencyId != 0 {
		tempRows = tempRows.Where("currency_id = ?", request.CurrencyId)
	}

	if !request.EffectiveDate.IsZero() {
		tempRows = tempRows.Where("effective_date >= ?", request.EffectiveDate)
	}

	if request.ItemGroupId != 0 {
		tempRows = tempRows.Where("item_group_id = ?", request.ItemGroupId)
	}

	if request.ItemClassId != 0 {
		tempRows = tempRows.Where("item_class_id = ?", request.ItemClassId)
	}

	rows, err := tempRows.
		Scan(&responses).
		Rows()

	if err != nil {
		return responses, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Err:        err,
		}
	}

	if len(responses) == 0 {
		return responses, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New(""),
		}
	}

	defer rows.Close()

	return responses, nil
}

func (r *PriceListRepositoryImpl) GetPriceList(tx *gorm.DB, request masteritempayloads.PriceListGetAllRequest) ([]masteritempayloads.PriceListResponse, *exceptions.BaseErrorResponse) {
	var responses []masteritempayloads.PriceListResponse
	var idMaps = make(map[string][]string)

	tempRows := tx.
		Model(&masteritementities.ItemPriceList{})

	if request.CompanyId != 0 {
		tempRows = tempRows.Where("company_id = ?", request.CompanyId)
	}

	if request.PriceListCode != "" {
		tempRows = tempRows.Where("price_list_code like ?", "%"+request.PriceListCode+"%")
	}

	if request.BrandId != 0 {
		tempRows = tempRows.Where("brand_id = ?", request.BrandId)
	}

	if request.CurrencyId != 0 {
		tempRows = tempRows.Where("currency_id = ?", request.CurrencyId)
	}

	if !request.EffectiveDate.IsZero() {
		tempRows = tempRows.Where("effective_date >= ?", request.EffectiveDate)
	}

	if request.ItemId != 0 {
		tempRows = tempRows.Where("item_id = ?", request.ItemId)
	}

	if request.ItemGroupId != 0 {
		tempRows = tempRows.Where("item_group_id = ?", request.ItemGroupId)
	}

	if request.ItemClassId != 0 {
		tempRows = tempRows.Where("item_class_id = ?", request.ItemClassId)
	}

	if request.PriceListAmount != 0 {
		tempRows = tempRows.Where("price_list_amount = ?", request.PriceListAmount)
	}

	if request.PriceListModifiable != "" {
		tempRows = tempRows.Where("price_list_modifiable = ?", request.PriceListModifiable)
	}

	if request.AtpmSyncronize != "" {
		tempRows = tempRows.Where("atpm_syncronize = ?", request.AtpmSyncronize)
	}

	if !request.AtpmSyncronizeTime.IsZero() {
		tempRows = tempRows.Where("atpm_syncronize_time >= ?", request.AtpmSyncronizeTime)
	}

	rows, err := tempRows.
		Scan(&responses).
		Rows()

	if err != nil {
		return responses, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	if len(responses) == 0 {
		return responses, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New(""),
		}
	}

	for _, response := range responses {
		idMaps["BrandId"] = append(idMaps["BrandId"], strconv.Itoa(int(response.BrandId)))
		idMaps["ItemGroupId"] = append(idMaps["ItemGroupId"], strconv.Itoa(int(response.ItemGroupId)))
		idMaps["ItemClassId"] = append(idMaps["ItemClassId"], strconv.Itoa(int(response.ItemClassId)))
		idMaps["CurrencyId"] = append(idMaps["CurrencyId"], strconv.Itoa(int(response.CurrencyId)))
	}

	defer rows.Close()

	return responses, nil
}

func (r *PriceListRepositoryImpl) GetPriceListById(tx *gorm.DB, Id int) (masteritempayloads.PriceListGetbyId, *exceptions.BaseErrorResponse) {
	entities := masteritementities.ItemPriceList{}
	response := masteritempayloads.PriceListGetbyId{}

	err := tx.Model(&entities).Select("mtr_item.*,mtr_item_class.*,mtr_item_price_list.*").
		Joins("JOIN mtr_item on mtr_item.item_id=mtr_item_price_list.item_id").
		Joins("JOIN mtr_item_class on mtr_item_class.item_class_id = mtr_item_price_list.item_class_id").
		Where(masteritementities.ItemPriceList{PriceListId: Id}).
		First(&response).Error

	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	// Fetch Brand
	brandResponse, errBrand := salesserviceapiutils.GetUnitBrandById(response.BrandId)
	if errBrand != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: errBrand.StatusCode,
			Err:        errBrand.Err,
		}
	}
	response.BrandId = brandResponse.BrandId
	response.BrandName = brandResponse.BrandName
	response.BrandCode = brandResponse.BrandCode

	// Fetch Item Group
	itemGroupResponse, errItemGroup := aftersalesserviceapiutils.GetItemGroupById(response.ItemGroupId)
	if errItemGroup != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: errItemGroup.StatusCode,
			Err:        errItemGroup.Err,
		}
	}
	response.ItemGroupId = itemGroupResponse.ItemGroupId
	response.ItemGroupName = itemGroupResponse.ItemGroupName

	// Fetch Currency
	currencyResponse, errCurrency := financeserviceapiutils.GetCurrencyId(response.CurrencyId)
	if errCurrency != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: errCurrency.StatusCode,
			Err:        errCurrency.Err,
		}
	}
	response.CurrencyId = currencyResponse.CurrencyId
	response.CurrencyCode = currencyResponse.CurrencyCode

	return response, nil
}

func (r *PriceListRepositoryImpl) SavePriceList(tx *gorm.DB, request masteritempayloads.SavePriceListMultiple) (int, *exceptions.BaseErrorResponse) {
	// dateParse, _ := time.Parse("2006-01-02", request.EffectiveDate)

	//NOTE!! MUST CHECK PRICELISTCODEID IS EXIST, PRICE LIST CODE (COMMON) STILL ON DEVELOPMENT - 9/AUG/2024 last status

	PriceListId := -1

	for _, value := range request.Detail {
		isExist := 0
		err := tx.Model(&masteritementities.ItemPriceList{}).
			Where("CONVERT(DATE, effective_date) = ?", request.EffectiveDate.Format("2006-01-02")).
			Where(masteritementities.ItemPriceList{
				ItemId:          value.ItemId,
				BrandId:         request.BrandId,
				ItemGroupId:     request.ItemGroupId,
				PriceListCodeId: request.PriceListCodeId,
				CurrencyId:      request.CurrencyId,
				CompanyId:       request.CompanyId,
			}).Select("1").Scan(&isExist).Error

		if err != nil {
			return PriceListId, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
				Message:    "error on check data price list",
			}
		}
		if isExist == 1 {
			return PriceListId, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "cannot insert duplicate item",
				Err:        err,
			}
		}
		entities := masteritementities.ItemPriceList{
			IsActive:            value.IsActive,
			PriceListCodeId:     request.PriceListCodeId,
			CompanyId:           request.CompanyId,
			BrandId:             request.BrandId,
			CurrencyId:          request.CurrencyId,
			EffectiveDate:       request.EffectiveDate,
			ItemId:              value.ItemId,
			ItemGroupId:         request.ItemGroupId,
			ItemClassId:         value.ItemClassId,
			PriceListAmount:     value.PriceListAmount,
			PriceListModifiable: true,
		}

		err = tx.Save(&entities).Where(entities).Select("mtr_item_price_list.price_list_id").First(&PriceListId).Error

		if err != nil {
			return PriceListId, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusConflict,
				Err:        err,
			}
		}
	}

	return PriceListId, nil
}

func (r *PriceListRepositoryImpl) ChangeStatusPriceList(tx *gorm.DB, Id int) (bool, *exceptions.BaseErrorResponse) {
	var entities masteritementities.ItemPriceList

	result := tx.Model(&entities).
		Where("price_list_id = ?", Id).
		First(&entities)

	if result.Error != nil {
		return false, &exceptions.BaseErrorResponse{
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
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	return true, nil
}

func RemoveDuplicates(input []string) []string {
	var result []string
	encountered := make(map[string]bool)

	for _, value := range input {
		if !encountered[value] {
			encountered[value] = true
			result = append(result, value)
		}
	}

	return result
}

func (r *PriceListRepositoryImpl) GetAllPriceListNew(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var payloads []masteritempayloads.PriceListGetAllResponse

	baseModelQuery := tx.Model(&masteritementities.ItemPriceList{}).
		Select(`
			mtr_item_price_list.brand_id,
			mtr_item_price_list.item_group_id,
			mtr_item_price_list.item_class_id,
			mtr_item_class.item_class_name,
			mtr_item_price_list.currency_id,
			CAST(effective_date AS DATE) AS effective_date,
			mtr_item_price_code.item_price_code,
			MIN(mtr_item_price_list.price_list_id) AS price_list_id,
			mtr_item_price_list.is_active
		`).
		Joins("JOIN mtr_item ON mtr_item.item_id = mtr_item_price_list.item_id").
		Joins("JOIN mtr_item_class ON mtr_item_class.item_class_id = mtr_item_price_list.item_class_id").
		Joins("LEFT JOIN mtr_item_price_code ON mtr_item_price_code.item_price_code_id = mtr_item_price_list.price_list_code_id").
		Group(`
			mtr_item_price_list.brand_id,
			mtr_item_price_list.item_group_id,
			mtr_item_price_list.item_class_id,
			mtr_item_class.item_class_name,
			mtr_item_price_list.currency_id,
			CAST(effective_date AS DATE),
			mtr_item_price_code.item_price_code,
			mtr_item_price_list.is_active
		`).
		Order("CAST(effective_date AS DATE) DESC")

	whereQuery := utils.ApplyFilter(baseModelQuery, filterCondition)

	var totalRows int64
	if err := whereQuery.Count(&totalRows).Error; err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	pages.TotalRows = totalRows

	totalPages := int(math.Ceil(float64(totalRows) / float64(pages.GetLimit())))
	pages.TotalPages = totalPages

	err := whereQuery.Offset(pages.GetOffset()).Limit(pages.GetLimit()).Find(&payloads).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if len(payloads) == 0 {
		pages.Rows = []map[string]interface{}{}
		return pages, nil
	}

	var results []map[string]interface{}
	for _, payload := range payloads {
		// Fetch data eksternal seperti brand, item group, currency
		brandResponse, brandErr := salesserviceapiutils.GetUnitBrandById(payload.BrandId)
		if brandErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        brandErr.Err,
			}
		}

		itemGroupResponse, itemGroupErr := aftersalesserviceapiutils.GetItemGroupById(payload.ItemGroupId)
		if itemGroupErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        itemGroupErr.Err,
			}
		}

		currencyResponse, currencyErr := financeserviceapiutils.GetCurrencyId(payload.CurrencyId)
		if currencyErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        currencyErr.Err,
			}
		}

		result := map[string]interface{}{
			"price_list_id":   payload.PriceListId,
			"brand_id":        payload.BrandId,
			"brand_name":      brandResponse.BrandName,
			"item_group_id":   payload.ItemGroupId,
			"item_group_name": itemGroupResponse.ItemGroupName,
			"item_class_id":   payload.ItemClassId,
			"item_class_name": payload.ItemClassName,
			"currency_id":     payload.CurrencyId,
			"currency_code":   currencyResponse.CurrencyCode,
			"effective_date":  payload.EffectiveDate,
			"item_price_code": payload.ItemPriceCode,
			"is_active":       payload.IsActive,
		}

		results = append(results, result)
	}

	pages.Rows = results

	return pages, nil
}

func (r *PriceListRepositoryImpl) DeletePriceList(tx *gorm.DB, id string) (bool, *exceptions.BaseErrorResponse) {
	idslice := strings.Split(id, ",")

	for _, ids := range idslice {
		var entityToUpdate masteritementities.ItemPriceList
		err := tx.Model(&entityToUpdate).Where("price_list_id = ?", ids).Delete(&entityToUpdate).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Err:        err,
			}
		}
	}
	return true, nil
}

func (r *PriceListRepositoryImpl) ActivatePriceList(tx *gorm.DB, id string) (bool, *exceptions.BaseErrorResponse) {
	idslice := strings.Split(id, ",")

	for _, ids := range idslice {
		var entityToUpdate masteritementities.ItemPriceList
		err := tx.Model(&entityToUpdate).Where("price_list_id = ?", ids).First(&entityToUpdate).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Err:        err,
			}
		}
		entityToUpdate.IsActive = true
		result := tx.Save(&entityToUpdate)
		if result.Error != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusConflict,
				Err:        result.Error,
			}
		}
	}
	return true, nil
}

func (r *PriceListRepositoryImpl) DeactivatePriceList(tx *gorm.DB, id string) (bool, *exceptions.BaseErrorResponse) {
	idslice := strings.Split(id, ",")

	for _, ids := range idslice {
		var entityToUpdate masteritementities.ItemPriceList
		err := tx.Model(&entityToUpdate).Where("price_list_id = ?", ids).First(&entityToUpdate).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Err:        err,
			}
		}
		entityToUpdate.IsActive = false
		result := tx.Save(&entityToUpdate)
		if result.Error != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusConflict,
				Err:        result.Error,
			}
		}
	}
	return true, nil
}

func (r *PriceListRepositoryImpl) GetPriceListByCodeId(tx *gorm.DB, CodeId string) (masteritempayloads.PriceListGetbyCode, *exceptions.BaseErrorResponse) {
	entities := masteritementities.ItemPriceList{}
	response := masteritempayloads.PriceListGetbyCode{}

	err := tx.Model(&entities).Select("mtr_item.*,mtr_item_class.*,mtr_item_price_list.*").
		Joins("JOIN mtr_item on mtr_item.item_id=mtr_item_price_list.item_id").
		Joins("JOIN mtr_item_class on mtr_item_class.item_class_id = mtr_item_price_list.item_class_id").
		Where("mtr_item_price_list.price_list_code_id = ?", CodeId).
		First(&response).Error

	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error on get price list by code id",
			Err:        err,
		}
	}

	// Fetch Brand
	brandResponse, errBrand := salesserviceapiutils.GetUnitBrandById(response.BrandId)
	if errBrand != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: errBrand.StatusCode,
			Err:        errBrand.Err,
		}
	}
	response.BrandId = brandResponse.BrandId
	response.BrandName = brandResponse.BrandName
	response.BrandCode = brandResponse.BrandCode

	// Fetch Item Group
	itemGroupResponse, errItemGroup := aftersalesserviceapiutils.GetItemGroupById(response.ItemGroupId)
	if errItemGroup != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: errItemGroup.StatusCode,
			Err:        errItemGroup.Err,
		}
	}
	response.ItemGroupId = itemGroupResponse.ItemGroupId
	response.ItemGroupName = itemGroupResponse.ItemGroupName

	// Fetch Currency
	currencyResponse, errCurrency := financeserviceapiutils.GetCurrencyId(response.CurrencyId)
	if errCurrency != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: errCurrency.StatusCode,
			Err:        errCurrency.Err,
		}
	}
	response.CurrencyId = currencyResponse.CurrencyId
	response.CurrencyCode = currencyResponse.CurrencyCode

	return response, nil
}
