package masteritemrepositoryimpl

import (
	"after-sales/api/config"
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	"after-sales/api/utils"
	"errors"
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

	if err := query.Scopes(pagination.Paginate(model, &pages, query)).Scan(&result).Error; err != nil {
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

	if err := tx.Model(model).Where(masteritementities.ItemPriceList{BrandId: brandId, ItemId: itemId}).
		Where("mtr_item_price_list.company_id = ?", companyId).
		Where("CONVERT(DATE, mtr_item_price_list.effective_date) like ?", date).First(&model).Error; err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
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
	brandpayloads := masteritempayloads.UnitBrandResponses{}
	itemgrouppayloads := masteritempayloads.ItemGroupResponse{}
	currencypayloads := masteritempayloads.CurrencyResponse{}

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

	ErrUrlBrand := utils.Get(config.EnvConfigs.SalesServiceUrl+"unit-brand/"+strconv.Itoa(response.BrandId), &brandpayloads, nil)
	if ErrUrlBrand != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        ErrUrlBrand,
		}
	}

	if brandpayloads != (masteritempayloads.UnitBrandResponses{}) {
		response.BrandId = brandpayloads.BrandId
		response.BrandName = brandpayloads.BrandName
	}

	ErrUrlItemGroup := utils.Get(config.EnvConfigs.GeneralServiceUrl+"item-group/"+strconv.Itoa(response.ItemGroupId), &itemgrouppayloads, nil)
	if ErrUrlItemGroup != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        ErrUrlItemGroup,
		}
	}

	if itemgrouppayloads != (masteritempayloads.ItemGroupResponse{}) {
		response.ItemGroupId = itemgrouppayloads.ItemGroupId
		response.ItemGroupName = itemgrouppayloads.ItemGroupName
	}

	ErrUrlCurrency := utils.Get(config.EnvConfigs.FinanceServiceUrl+"currency-code/"+strconv.Itoa(response.CurrencyId), &currencypayloads, nil)
	if ErrUrlCurrency != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        ErrUrlCurrency,
		}
	}

	if currencypayloads != (masteritempayloads.CurrencyResponse{}) {
		response.CurrencyId = currencypayloads.CurrencyId
		response.CurrencyCode = currencypayloads.CurrencyCode
	}

	return response, nil
}

func (r *PriceListRepositoryImpl) SavePriceList(tx *gorm.DB, request masteritempayloads.SavePriceListMultiple) (int, *exceptions.BaseErrorResponse) {
	// dateParse, _ := time.Parse("2006-01-02", request.EffectiveDate)

	//NOTE!! MUST CHECK PRICELISTCODEID IS EXIST, PRICE LIST CODE (COMMON) STILL ON DEVELOPMENT - 9/AUG/2024 last status

	PriceListId := -1

	for _, value := range request.Detail {

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

		err := tx.Save(&entities).Where(entities).Select("mtr_item_price_list.price_list_id").First(&PriceListId).Error

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

func (r *PriceListRepositoryImpl) GetAllPriceListNew(tx *gorm.DB, filtercondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]any, int, int, *exceptions.BaseErrorResponse) {
	var payloads []masteritempayloads.PriceListGetAllResponse
	var brandpayloads []masterpayloads.BrandResponse
	var itemgrouppayloads []masteritempayloads.ItemGroupResponse
	var currencypayloads []masteritempayloads.CurrencyResponse

	model := masteritementities.ItemPriceList{}

	query := tx.Model(model).
		Select("mtr_item.*,mtr_item_class.*,mtr_item_price_list.*,mtr_item_price_code.item_price_code").
		Joins("JOIN mtr_item on mtr_item.item_id=mtr_item_price_list.item_id").
		Joins("JOIN mtr_item_class on mtr_item_class.item_class_id = mtr_item_price_list.item_class_id").
		Joins("LEFT JOIN mtr_item_price_code on mtr_item_price_code.item_price_code_id = mtr_item_price_list.price_list_code_id")

	//apply where query
	whereQuery := utils.ApplyFilterExact(query, filtercondition)
	//execute
	err := whereQuery.Scan(&payloads).Error

	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	errBrandUrl := utils.Get(config.EnvConfigs.SalesServiceUrl+"unit-brand-dropdown", &brandpayloads, nil)
	if errBrandUrl != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New("failed to fetch brand data"),
		}
	}

	joinedData := utils.DataFrameLeftJoin(payloads, brandpayloads, "BrandId")

	errItemGroupUrl := utils.Get(config.EnvConfigs.GeneralServiceUrl+"item-group", &itemgrouppayloads, nil)
	if errItemGroupUrl != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New("failed to fetch item group data"),
		}
	}

	joinedData1 := utils.DataFrameLeftJoin(joinedData, itemgrouppayloads, "ItemGroupId")

	errCurrencyUrl := utils.Get(config.EnvConfigs.FinanceServiceUrl+"currency-code/", &currencypayloads, nil)
	if errCurrencyUrl != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New("failed to fetch currency data"),
		}
	}

	joinedData2 := utils.DataFrameLeftJoin(joinedData1, currencypayloads, "CurrencyId")

	dataPaginate, totalPages, totalRows := pagination.NewDataFramePaginate(joinedData2, &pages)

	return dataPaginate, totalPages, totalRows, nil
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
