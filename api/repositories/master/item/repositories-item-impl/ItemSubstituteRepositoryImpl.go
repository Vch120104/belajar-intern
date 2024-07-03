package masteritemrepositoryimpl

import (
	"after-sales/api/config"
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type ItemSubstituteRepositoryImpl struct {
}

func StartItemSubstituteRepositoryImpl() masteritemrepository.ItemSubstituteRepository {
	return &ItemSubstituteRepositoryImpl{}
}

func (r *ItemSubstituteRepositoryImpl) GetAllItemSubstitute(tx *gorm.DB, filterCondition map[string]string, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
    var responses []masteritempayloads.ItemSubstituteById
    var getSubtituteCode []masteritempayloads.GetSubstitutePayloads
    is_active := filterCondition["is_active"]
    item_id := filterCondition["item_id"]
    brand_id := filterCondition["brand_id"]
    substitute_item_id := filterCondition["substitute_item_code"]
    effective_date_from := filterCondition["effective_date_from"]
    effective_date_to := filterCondition["effective_date_to"]

    layout := "2006-01-02T15:04:05.0000000-07:00"
    var startDate, endDate time.Time
    var err error
    if effective_date_from != "" {
        startDate, err = time.Parse(layout, effective_date_from)
        if err != nil {
            return nil, 0, 0, &exceptions.BaseErrorResponse{
                StatusCode: http.StatusBadRequest,
                Err:        err,
            }
        }
    }

    if effective_date_to != "" {
        endDate, err = time.Parse(layout, effective_date_to)
        if err != nil {
            return nil, 0, 0, &exceptions.BaseErrorResponse{
                StatusCode: http.StatusBadRequest,
                Err:        err,
            }
        }
    }

    query := tx.Table("mtr_item_substitute AS mis").
        Select("mis.*, mi.*,ic.*").
        Joins("JOIN mtr_item AS mi ON mi.item_id = mis.item_id").
		Joins("JOIN mtr_item_class as ic on ic.item_class_id = mis.item_class_id")

    if is_active != "" {
        query = query.Where("mis.is_active = ?", is_active)
    }
    if item_id != "" {
        query = query.Where("mis.item_id = ?", item_id)
    }
    if brand_id != "" {
        query = query.Where("mi.brand_id = ?", brand_id)
    }
    if substitute_item_id != "" {
        query = query.Where("mis.substitute_item_code = ?", substitute_item_id)
    }
    if !startDate.IsZero() && !endDate.IsZero() {
        query = query.Where("mis.effective_date >= ? AND mis.effective_date <= ?", startDate, endDate)
    }

    err = query.Scan(&responses).Error
    if err != nil {
        return nil, 0, 0, &exceptions.BaseErrorResponse{
            StatusCode: http.StatusNotFound,
            Err:        err,
        }
    }

    if len(responses) == 0 {
        return nil, 0, 0, &exceptions.BaseErrorResponse{
            StatusCode: http.StatusInternalServerError,
            Err:        err,
        }
    }

    ItemSubstituteCodeUrl := config.EnvConfigs.GeneralServiceUrl + "/substitute-type"

    errUrlSubstituteItem := utils.Get(ItemSubstituteCodeUrl, &getSubtituteCode, nil)
    if errUrlSubstituteItem != nil {
        return nil, 0, 0, &exceptions.BaseErrorResponse{
            StatusCode: http.StatusInternalServerError,
            Err:        err,
        }
    }

    joinedData := utils.DataFrameInnerJoin(responses, getSubtituteCode, "SubstituteTypeId")
    dataPaginate, totalPages, totalRows := pagination.NewDataFramePaginate(joinedData, &pages)

    return dataPaginate, totalPages, totalRows, nil
}


func (r *ItemSubstituteRepositoryImpl) GetByIdItemSubstitute(tx *gorm.DB, id int) (map[string]interface{}, *exceptions.BaseErrorResponse) {
	response := masteritempayloads.ItemSubstituteById{}
	itemgroupresponse := masteritempayloads.ItemGroupResponse{}
	err := tx.Table("mtr_item_substitute").Select("mtr_item_substitute.*,mtr_item.*,mtr_item_class.*").
		Joins("JOIN mtr_item on mtr_item.item_id = mtr_item_substitute.item_id").
		Joins("Join mtr_item_class on mtr_item_class.item_class_id = mtr_item_substitute.item_class_id").
		Scan(&response).Error
	if err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if errItemGroup := utils.Get(config.EnvConfigs.GeneralServiceUrl+"/item-group/"+strconv.Itoa(response.ItemGroupId),&itemgroupresponse,nil); errItemGroup != nil{
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err: errors.New(""),
		}
	}
	joineddata2 := utils.DataFrameInnerJoin([]masteritempayloads.ItemSubstituteById{response},[]masteritempayloads.ItemGroupResponse{itemgroupresponse},"ItemGroupId")
	if len(joineddata2)==0{
		return nil,&exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err: errors.New(""),
		}
	}
	return joineddata2[0], nil
}

func (r *ItemSubstituteRepositoryImpl) GetAllItemSubstituteDetail(tx *gorm.DB, pages pagination.Pagination, id int) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	entities := []masteritementities.ItemSubstituteDetail{}
	payloads := []masteritempayloads.ItemSubstituteDetailPayloads{}
	tableStruct := masteritempayloads.ItemSubstituteDetailPayloads{}

	baseModelQuery := utils.CreateJoinSelectStatement(tx, tableStruct).Where(masteritementities.ItemSubstituteDetail{ItemSubstituteId: id})

	rows, err := baseModelQuery.Scopes(pagination.Paginate(&entities, &pages, baseModelQuery)).Scan(&payloads).Rows()
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	if len(payloads) == 0 {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New(""),
		}
	}

	defer rows.Close()

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
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	defer rows.Close()
	return response, nil
}

func (r *ItemSubstituteRepositoryImpl) SaveItemSubstitute(tx *gorm.DB, req masteritempayloads.ItemSubstitutePostPayloads) (bool, *exceptions.BaseErrorResponse) {
	entities := masteritementities.ItemSubstitute{
		SubstituteTypeId: req.SubstituteTypeId,
		ItemSubstituteId: req.ItemSubstituteId,
		EffectiveDate:    req.EffectiveDate,
		ItemId:           req.ItemId,
		ItemGroupId:      req.ItemGroupId,
		ItemClassId:      req.ItemClassId,
		Description:      req.Description,
	}
	err := tx.Save(&entities).Error

	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	return true, nil
}

func (r *ItemSubstituteRepositoryImpl) SaveItemSubstituteDetail(tx *gorm.DB, req masteritempayloads.ItemSubstituteDetailPostPayloads, id int) (bool, *exceptions.BaseErrorResponse) {
	var count int64
	err := tx.Model(&masteritementities.ItemSubstituteDetail{}).
		Where("item_substitute_id = ?", id).
		Count(&count).Error

	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	entities := masteritementities.ItemSubstituteDetail{
		ItemSubstituteDetailId: req.ItemSubstituteDetailId,
		ItemId:                 req.ItemId,
		ItemSubstituteId:       id,
		Quantity:               req.Quantity,
		Sequence:               int(count) + 1,
	}
	err2 := tx.Save(&entities).Error

	if err2 != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err2,
		}
	}

	return true, nil
}

func (r *ItemSubstituteRepositoryImpl) ChangeStatusItemOperation(tx *gorm.DB, id int) (bool, *exceptions.BaseErrorResponse) {
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