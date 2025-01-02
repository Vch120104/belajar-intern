package masteritemrepositoryimpl

import (
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	"after-sales/api/utils"
	generalserviceapiutils "after-sales/api/utils/general-service"
	"errors"
	"net/http"
	"strings"
	"time"

	"gorm.io/gorm"
)

type ItemSubstituteRepositoryImpl struct {
}

func StartItemSubstituteRepositoryImpl() masteritemrepository.ItemSubstituteRepository {
	return &ItemSubstituteRepositoryImpl{}
}

func (r *ItemSubstituteRepositoryImpl) GetAllItemSubstitute(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination, from time.Time, to time.Time) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var payloads []masteritempayloads.ItemSubstitutePayloads

	query := tx.Model(&masteritementities.ItemSubstitute{}).
		Select("mtr_item_substitute.*, Item.item_code, Item.item_name").
		Joins("JOIN mtr_item Item ON mtr_item_substitute.item_id = Item.item_id")

	whereQuery := utils.ApplyFilter(query, filterCondition)
	if !from.IsZero() {
		fromFormatted := from.Format("2006-01-02") + " 00:00:00.000"
		whereQuery = whereQuery.Where("effective_date >= ?", fromFormatted)
	}
	if !to.IsZero() {
		toFormatted := to.Format("2006-01-02") + " 23:59:59.999"
		whereQuery = whereQuery.Where("effective_date <= ?", toFormatted)
	}

	err := whereQuery.Scopes(pagination.Paginate(&pages, whereQuery)).Find(&payloads).Error
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

	typeResponse, typeError := generalserviceapiutils.GetAllSubstituteType()
	if typeError != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: typeError.StatusCode,
			Err:        typeError.Err,
		}
	}

	joinedData := utils.DataFrameLeftJoin(payloads, typeResponse, "SubstituteTypeId")

	var results []map[string]interface{}
	for _, data := range joinedData {
		result := map[string]interface{}{
			"effective_date":       data["EffectiveDate"],
			"is_active":            data["IsActive"],
			"item_class_code":      data["ItemClassCode"],
			"item_class_id":        data["ItemClassId"],
			"item_code":            data["ItemCode"],
			"item_group_id":        data["ItemGroupId"],
			"item_id":              data["ItemId"],
			"item_name":            data["ItemName"],
			"item_substitute_id":   data["ItemSubstituteId"],
			"substitute_type_id":   data["SubstituteTypeId"],
			"substitute_type_name": data["SubstituteTypeName"],
		}
		results = append(results, result)
	}

	pages.Rows = results
	return pages, nil
}

func (r *ItemSubstituteRepositoryImpl) GetByIdItemSubstitute(tx *gorm.DB, id int) (map[string]interface{}, *exceptions.BaseErrorResponse) {
	var entity masteritementities.ItemSubstitute
	var response masteritempayloads.ItemSubstitutePayloads

	err := tx.Model(entity).
		Select("mtr_item_substitute.*, Item.item_code, Item.item_name, Item.item_class_id, Item.item_group_id").
		Where(masteritementities.ItemSubstitute{ItemSubstituteId: id}).
		Joins("JOIN mtr_item Item ON mtr_item_substitute.item_id = Item.item_id").
		Joins("JOIN mtr_item_class ON Item.item_class_id = mtr_item_class.item_class_id").
		First(&response).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        errors.New("item substitute not found"),
			}
		}
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	typeResponse, errSubType := generalserviceapiutils.GetSubstituteTypeById(response.SubstituteTypeId)
	if errSubType != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: errSubType.StatusCode,
			Err:        errSubType.Err,
		}
	}

	// Construct response map
	result := map[string]interface{}{
		"effective_date":       response.EffectiveDate,
		"is_active":            response.IsActive,
		"item_class_code":      response.ItemClassCode,
		"item_class_id":        response.ItemClassId,
		"item_code":            response.ItemCode,
		"item_group_id":        response.ItemGroupId,
		"item_id":              response.ItemId,
		"item_name":            response.ItemName,
		"item_substitute_id":   response.ItemSubstituteId,
		"description":          response.Description,
		"substitute_type_id":   typeResponse.SubstituteTypeId,
		"substitute_type_name": typeResponse.SubstituteTypeName,
	}

	return result, nil

}

func (r *ItemSubstituteRepositoryImpl) GetAllItemSubstituteDetail(tx *gorm.DB, pages pagination.Pagination, id int) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	entities := []masteritementities.ItemSubstituteDetail{}
	payloads := []masteritempayloads.ItemSubstituteDetailPayloads{}

	query := tx.Model(entities).Select("mtr_item_substitute_detail.*, Item.item_code, Item.item_name").
		Joins("Item", tx.Select("")).Where("mtr_item_substitute_detail.item_substitute_id = ?", id)

	err := query.Scopes(pagination.Paginate(&pages, query)).Scan(&payloads).Error

	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	pages.Rows = payloads

	return pages, nil
}

func (r *ItemSubstituteRepositoryImpl) GetByIdItemSubstituteDetail(tx *gorm.DB, id int) (masteritempayloads.ItemSubstituteDetailGetPayloads, *exceptions.BaseErrorResponse) {
	response := masteritempayloads.ItemSubstituteDetailGetPayloads{}
	tableStruct := masteritempayloads.ItemSubstituteDetailPayloads{}
	baseModelQuery := utils.CreateJoinSelectStatement(tx, tableStruct).Where(masteritementities.ItemSubstituteDetail{ItemSubstituteDetailId: id})

	rows, err := baseModelQuery.First(&response).Rows()

	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	defer rows.Close()
	return response, nil
}

func (r *ItemSubstituteRepositoryImpl) SaveItemSubstitute(tx *gorm.DB, req masteritempayloads.ItemSubstitutePostPayloads) (masteritementities.ItemSubstitute, *exceptions.BaseErrorResponse) {

	// parseEffectiveDate, _ := time.Parse("2006-01-02T15:04:05.000Z", req.EffectiveDate)

	entities := masteritementities.ItemSubstitute{
		IsActive:         req.IsActive,
		SubstituteTypeId: req.SubstituteTypeId,
		ItemSubstituteId: req.ItemSubstituteId,
		EffectiveDate:    req.EffectiveDate,
		ItemId:           req.ItemId,
		Description:      req.Description,
	}
	err := tx.Save(&entities).Error

	if err != nil {
		return entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err,
		}
	}

	return entities, nil
}

func (r *ItemSubstituteRepositoryImpl) SaveItemSubstituteDetail(tx *gorm.DB, req masteritempayloads.ItemSubstituteDetailPostPayloads, id int) (masteritementities.ItemSubstituteDetail, *exceptions.BaseErrorResponse) {
	var existing masteritementities.ItemSubstituteDetail
	if err := tx.Where("item_id = ? AND item_substitute_id = ?", req.ItemId, id).First(&existing).Error; err == nil {
		return masteritementities.ItemSubstituteDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Message:    "duplicate item in substitute detail",
			Err:        errors.New("duplicate item in substitute detail"),
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return masteritementities.ItemSubstituteDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Fail to check existing data",
			Err:        err,
		}
	}

	var header masteritementities.ItemSubstitute
	if err := tx.Where("item_id = ? AND item_substitute_id = ?", req.ItemId, id).First(&header).Error; err == nil {
		return masteritementities.ItemSubstituteDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Message:    "duplicate item with substitute header",
			Err:        errors.New("duplicate item with substitute header"),
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return masteritementities.ItemSubstituteDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Fail to check existing data",
			Err:        err,
		}
	}

	entities := masteritementities.ItemSubstituteDetail{
		IsActive:               req.IsActive,
		ItemSubstituteDetailId: req.ItemSubstituteDetailId,
		ItemId:                 req.ItemId,
		ItemSubstituteId:       id,
		Quantity:               req.Quantity,
		Sequence:               req.Sequence,
	}

	err := tx.Save(&entities).Error
	if err != nil {
		return masteritementities.ItemSubstituteDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	return entities, nil
}

func (r *ItemSubstituteRepositoryImpl) UpdateItemSubstituteDetail(tx *gorm.DB, req masteritempayloads.ItemSubstituteDetailUpdatePayloads) (masteritementities.ItemSubstituteDetail, *exceptions.BaseErrorResponse) {
	var err error
	entities := masteritementities.ItemSubstituteDetail{}

	err = tx.Model(&entities).
		Where(masteritementities.ItemSubstituteDetail{ItemSubstituteDetailId: req.ItemSubstituteDetailId}).
		First(&entities).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return masteritementities.ItemSubstituteDetail{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "item substitute detail not found",
				Err:        err,
			}
		}
		return masteritementities.ItemSubstituteDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error when fetching item substitute detail data",
			Err:        err,
		}
	}

	entities.Quantity = req.Quantity

	err = tx.Save(&entities).Error
	if err != nil {
		return masteritementities.ItemSubstituteDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	return entities, nil
}

func (r *ItemSubstituteRepositoryImpl) ChangeStatusItemSubstitute(tx *gorm.DB, id int) (bool, *exceptions.BaseErrorResponse) {
	var entities masteritementities.ItemSubstitute

	result := tx.Model(&entities).
		Where("item_substitute_id = ?", id).
		First(&entities)

	if result.Error != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
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
			StatusCode: http.StatusNotFound,
			Err:        result.Error,
		}
	}

	return true, nil
}

func (r *ItemSubstituteRepositoryImpl) DeactivateItemSubstituteDetail(tx *gorm.DB, id string) (bool, *exceptions.BaseErrorResponse) {
	idSlice := strings.Split(id, ",")

	for _, Ids := range idSlice {
		var entityToUpdate masteritementities.ItemSubstituteDetail
		err := tx.Model(&entityToUpdate).Where("item_substitute_detail_id = ?", Ids).First(&entityToUpdate).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}

		entityToUpdate.IsActive = false
		result := tx.Save(&entityToUpdate)
		if result.Error != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        result.Error,
			}
		}
	}

	return true, nil
}

func (r *ItemSubstituteRepositoryImpl) ActivateItemSubstituteDetail(tx *gorm.DB, id string) (bool, *exceptions.BaseErrorResponse) {
	idSlice := strings.Split(id, ",")

	for _, Ids := range idSlice {
		var entityToUpdate masteritementities.ItemSubstituteDetail
		err := tx.Model(&entityToUpdate).Where("item_substitute_detail_id = ?", Ids).First(&entityToUpdate).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}

		entityToUpdate.IsActive = true
		result := tx.Save(&entityToUpdate)
		if result.Error != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        result.Error,
			}
		}
	}

	return true, nil
}

func (r *ItemSubstituteRepositoryImpl) GetallItemForFilter(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	entities := masteritementities.Item{}
	payloads := []masteritempayloads.Itemforfilter{}

	// Build the query, ensure Select is used properly
	query := tx.Model(&entities).
		Select(`
			mtr_item.item_id,
			mtr_item.item_code,
			mtr_item.item_name,
			mil1.item_level_1_code,
			mil2.item_level_2_code,
			mil3.item_level_3_code,
			mil4.item_level_4_code,
			mtr_item_class.item_class_code,
			mit.item_type_code
		`).
		Joins("JOIN mtr_item_class ON mtr_item_class.item_class_id = mtr_item.item_class_id").
		Joins("INNER JOIN mtr_item_level_1 mil1 ON mil1.item_level_1_id = mtr_item.item_level_1_id").
		Joins("LEFT JOIN mtr_item_level_2 mil2 ON mil2.item_level_2_id = mtr_item.item_level_2_id").
		Joins("LEFT JOIN mtr_item_level_3 mil3 ON mil3.item_level_3_id = mtr_item.item_level_3_id").
		Joins("LEFT JOIN mtr_item_level_4 mil4 ON mil4.item_level_4_id = mtr_item.item_level_4_id").
		Joins("INNER JOIN mtr_item_type mit ON mit.item_type_id = mtr_item.item_type_id")

	// Apply filters
	whereQuery := utils.ApplyFilter(query, filterCondition)

	// Apply pagination and execute the query
	err := whereQuery.Scopes(pagination.Paginate(&pages, whereQuery)).Scan(&payloads).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	// Check if the result set is empty
	if len(payloads) == 0 {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New("no data found"),
		}
	}

	// Set the result rows to the pagination object
	pages.Rows = payloads

	return pages, nil
}

func (r *ItemSubstituteRepositoryImpl) GetItemSubstituteDetailLastSequence(tx *gorm.DB, id int) (map[string]interface{}, *exceptions.BaseErrorResponse) {
	entities := masteritementities.ItemSubstituteDetail{}
	response := map[string]interface{}{}

	err := tx.Model(&entities).Where(masteritementities.ItemSubstituteDetail{ItemSubstituteId: id}).Last(&entities).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error when fetching item substitute detail last sequence number",
			Err:        err,
		}
	}

	if err == nil {
		response["last_sequence"] = entities.Sequence
	} else {
		response["last_sequence"] = 0
	}

	return response, nil
}
