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

func (r *PriceListRepositoryImpl) GetPriceListLookup(tx *gorm.DB, request masteritempayloads.PriceListGetAllRequest) ([]masteritempayloads.PriceListResponse, *exceptions.BaseErrorResponse) {
	var responses []masteritempayloads.PriceListResponse

	tempRows := tx.
		Model(&masteritementities.PriceList{})

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
		Model(&masteritementities.PriceList{})

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

func (r *PriceListRepositoryImpl) GetPriceListById(tx *gorm.DB, Id int) (map[string]interface{}, *exceptions.BaseErrorResponse) {
	entities := masteritementities.PriceList{}
	response := masteritempayloads.PriceListResponse{}
	brandpayloads := masteritempayloads.UnitBrandResponses{}
	itemgrouppayloads := masteritempayloads.ItemGroupResponse{}
	currencypayloads := masteritempayloads.CurrencyResponse{}

	err := tx.Model(&entities).
		Where("price_list_id = ?", Id).
		First(&response).Error

	if err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if response.PriceListId != 0 {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	ErrUrlBrand := utils.Get(config.EnvConfigs.SalesServiceUrl+"/unit-brand/"+strconv.Itoa(response.BrandId), &brandpayloads, nil)
	if ErrUrlBrand != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        ErrUrlBrand,
		}
	}

	joinedData := utils.DataFrameInnerJoin([]masteritempayloads.PriceListResponse{response}, []masteritempayloads.UnitBrandResponses{brandpayloads}, "BrandId")

	ErrUrlItemGroup := utils.Get(config.EnvConfigs.GeneralServiceUrl+"/item-group/"+strconv.Itoa(response.ItemGroupId), &itemgrouppayloads, nil)
	if ErrUrlItemGroup != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        ErrUrlItemGroup,
		}
	}

	joineddata2 := utils.DataFrameInnerJoin(joinedData, []masteritempayloads.ItemGroupResponse{itemgrouppayloads}, "ItemGroupId")

	ErrUrlCurrency := utils.Get(config.EnvConfigs.FinanceServiceUrl+"/currency/"+strconv.Itoa(response.CurrencyId), &currencypayloads, nil)
	if ErrUrlCurrency != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        ErrUrlCurrency,
		}
	}
	joineddata3 := utils.DataFrameInnerJoin(joineddata2, []masteritempayloads.CurrencyResponse{currencypayloads}, "CurrencyId")

	return joineddata3[0], nil
}

func (r *PriceListRepositoryImpl) SavePriceList(tx *gorm.DB, request masteritempayloads.PriceListResponse) (bool, *exceptions.BaseErrorResponse) {
	entities := masteritementities.PriceList{
		IsActive:            request.IsActive,
		PriceListId:         request.PriceListId,
		PriceListCode:       request.PriceListCode,
		CompanyId:           request.CompanyId,
		BrandId:             request.BrandId,
		CurrencyId:          request.CurrencyId,
		EffectiveDate:       request.EffectiveDate,
		ItemId:              request.ItemId,
		ItemGroupId:         request.ItemGroupId,
		ItemClassId:         request.ItemClassId,
		PriceListAmount:     request.PriceListAmount,
		PriceListModifiable: request.PriceListModifiable,
		AtpmSyncronize:      request.AtpmSyncronize,
		AtpmSyncronizeTime:  request.AtpmSyncronizeTime,
	}

	err := tx.Save(&entities).Error

	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Err:        err,
		}
	}

	return true, nil
}

func (r *PriceListRepositoryImpl) ChangeStatusPriceList(tx *gorm.DB, Id int) (bool, *exceptions.BaseErrorResponse) {
	var entities masteritementities.PriceList

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

func (r *PriceListRepositoryImpl) GetAllPriceListNew(tx *gorm.DB, filtercondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var payloads []masteritempayloads.PriceListGetAllResponse
	var brandpayloads []masterpayloads.BrandResponse
	var itemgrouppayloads []masteritempayloads.ItemGroupResponse
	var currencypayloads []masteritempayloads.CurrencyResponse

	err := tx.Table("mtr_price_list").
		Select("mtr_price_list.*,mtr_item.*,mtr_item_class.*").
		Joins("JOIN mtr_item on mtr_item.item_id=mtr_price_list.item_id").
		Joins("JOIN mtr_item_class on mtr_item_class.item_class_id = mtr_price_list.item_class_id").
		Scan(&payloads).Error
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	errBrandUrl := utils.Get(config.EnvConfigs.SalesServiceUrl+"/unit-brand?page=0&limit=10000", &brandpayloads, nil)
	if errBrandUrl != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New("failed to fetch brand data"),
		}
	}

	joinedData := utils.DataFrameInnerJoin(payloads,brandpayloads,"BrandId")

	errItemGroupUrl := utils.Get(config.EnvConfigs.GeneralServiceUrl+"/item-group", &itemgrouppayloads, nil)
	if errItemGroupUrl != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New("failed to fetch item group data"),
		}
	}

	joinedData1 := utils.DataFrameInnerJoin(joinedData, itemgrouppayloads, "ItemGroupId")

	errCurrencyUrl := utils.Get(config.EnvConfigs.FinanceServiceUrl+"/currency-code/", &currencypayloads, nil)
	if errCurrencyUrl != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New("failed to fetch currency data"),
		}
	}

	joinedData2 := utils.DataFrameInnerJoin(joinedData1, currencypayloads, "CurrencyId")

	dataPaginate, totalPages, totalRows := pagination.NewDataFramePaginate(joinedData2, &pages)
	return dataPaginate, totalPages, totalRows, nil
}

func (r *PriceListRepositoryImpl) DeletePriceList(tx *gorm.DB, id string) (bool, *exceptions.BaseErrorResponse) {
	idslice := strings.Split(id, ",")

	for _, ids := range idslice {
		var entityToUpdate masteritementities.PriceList
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
		var entityToUpdate masteritementities.PriceList
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
		var entityToUpdate masteritementities.PriceList
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
