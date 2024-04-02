package masteritemrepositoryimpl

import (
	masteritementities "after-sales/api/entities/master/item"
	masteritempayloads "after-sales/api/payloads/master/item"
	masteritemrepository "after-sales/api/repositories/master/item"
	"strconv"

	"gorm.io/gorm"
)

type PriceListRepositoryImpl struct {
}

func StartPriceListRepositoryImpl() masteritemrepository.PriceListRepository {
	return &PriceListRepositoryImpl{}
}

func (r *PriceListRepositoryImpl) GetPriceListLookup(tx *gorm.DB, request masteritempayloads.PriceListGetAllRequest) ([]masteritempayloads.PriceListResponse, error) {
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
		return responses, err
	}

	defer rows.Close()

	return responses, nil
}

func (r *PriceListRepositoryImpl) GetPriceList(tx *gorm.DB, request masteritempayloads.PriceListGetAllRequest) ([]masteritempayloads.PriceListResponse, error) {
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
		return responses, err
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

func (r *PriceListRepositoryImpl) GetPriceListById(tx *gorm.DB, Id int) (masteritempayloads.PriceListResponse, error) {
	entities := masteritementities.PriceList{}
	response := masteritempayloads.PriceListResponse{}

	rows, err := tx.Model(&entities).
		Where(masteritementities.PriceList{
			PriceListId: int32(Id),
		}).
		First(&response).
		Rows()

	if err != nil {
		return response, err
	}

	defer rows.Close()

	return response, nil
}

func (r *PriceListRepositoryImpl) SavePriceList(tx *gorm.DB, request masteritempayloads.PriceListResponse) (bool, error) {
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
		return false, err
	}

	return true, nil
}

func (r *PriceListRepositoryImpl) ChangeStatusPriceList(tx *gorm.DB, Id int) (bool, error) {
	var entities masteritementities.PriceList

	result := tx.Model(&entities).
		Where("price_list_id = ?", Id).
		First(&entities)

	if result.Error != nil {
		return false, result.Error
	}

	if entities.IsActive {
		entities.IsActive = false
	} else {
		entities.IsActive = true
	}

	result = tx.Save(&entities)

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
