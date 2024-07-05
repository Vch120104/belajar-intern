package masterrepositoryimpl

import (
	masterentities "after-sales/api/entities/master"
	"net/http"
	"strconv"
	"strings"

	// masterpayloads "after-sales/api/payloads/master"
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type FieldActionRepositoryImpl struct {
}

func StartFieldActionRepositoryImpl() masterrepository.FieldActionRepository {
	return &FieldActionRepositoryImpl{}
}

func (r *FieldActionRepositoryImpl) GetAllFieldAction(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var responses []masterpayloads.FieldActionResponse
	entities := masterentities.FieldAction{}
	JoinTable := tx.Table("mtr_field_action as fa").
		Select("fa.*,faev.*").
		Joins("Join mtr_field_action_eligible_vehicle as faev ON faev.field_action_system_number=fa.field_action_system_number")

	whereQuery := utils.ApplyFilter(JoinTable, filterCondition)
	err := whereQuery.Scopes(pagination.Paginate(&entities, &pages, JoinTable)).Order("fa.field_action_system_number").Scan(&responses).Error

	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if len(responses) == 0 {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	pages.Rows = responses
	return pages, nil
}

func (r *FieldActionRepositoryImpl) SaveFieldAction(tx *gorm.DB, req masterpayloads.FieldActionRequest) (bool, *exceptions.BaseErrorResponse) {
	entities := masterentities.FieldAction{
		IsActive:                  req.IsActive,
		FieldActionSystemNumber:   req.FieldActionSystemNumber,
		FieldActionDocumentNumber: req.FieldActionDocumentNumber,
		ApprovalStatusId:          req.ApprovalStatusId,
		BrandId:                   req.BrandId,
		FieldActionName:           req.FieldActionName,
		FieldActionPeriodFrom:     req.FieldActionPeriodFrom,
		FieldActionPeriodTo:       req.FieldActionPeriodTo,
		IsNeverExpired:            req.IsNeverExpired,
		RemarkPopup:               req.RemarkPopup,
		IsCritical:                req.IsCritical,
		RemarkInvoice:             req.RemarkInvoice,
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

func (r *FieldActionRepositoryImpl) GetFieldActionHeaderById(tx *gorm.DB, Id int) (masterpayloads.FieldActionResponse, *exceptions.BaseErrorResponse) {
	entities := masterentities.FieldAction{}
	response := masterpayloads.FieldActionResponse{}

	rows, err := tx.Model(&entities).
		Where(masterentities.FieldAction{
			FieldActionSystemNumber: Id,
		}).
		First(&response).
		Rows()

	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()

	return response, nil
}

func (r *FieldActionRepositoryImpl) GetAllFieldActionVehicleDetailById(tx *gorm.DB, Id int, pages pagination.Pagination, filterCondition []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	entities := []masterentities.FieldActionEligibleVehicle{}
	payloads := []masterpayloads.FieldActionDetailResponse{}
	// tableStruct := masterpayloads.FieldActionDetailResponse{}

	baseModelQuery := tx.Model(&entities).
		Where(masterentities.FieldActionEligibleVehicle{
			FieldActionSystemNumber: Id,
		})
	filterQuery := utils.ApplyFilter(baseModelQuery, filterCondition)
	rows, err := baseModelQuery.Scopes(pagination.Paginate(&entities, &pages, filterQuery)).Scan(&payloads).Rows()

	if len(payloads) == 0 {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()

	pages.Rows = payloads

	return pages, nil
}

func (r *FieldActionRepositoryImpl) GetFieldActionVehicleDetailById(tx *gorm.DB, Id int) (masterpayloads.FieldActionDetailResponse, *exceptions.BaseErrorResponse) {
	entities := masterentities.FieldActionEligibleVehicle{}
	response := masterpayloads.FieldActionDetailResponse{}

	rows, err := tx.Model(&entities).
		Where(masterentities.FieldActionEligibleVehicle{
			FieldActionEligibleVehicleSystemNumber: Id,
		}).
		First(&response).
		Rows()

	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()

	return response, nil
}

func (r *FieldActionRepositoryImpl) GetAllFieldActionVehicleItemDetailById(tx *gorm.DB, Id int, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	entitiesitem := []masterentities.FieldActionEligibleVehicleItem{}
	entitiesoperation := []masterentities.FieldActionEligibleVehicleOperation{}
	payloadsoperation := []masterpayloads.FieldActionEligibleVehicleOperation{}
	payloadsitem := []masterpayloads.FieldActionEligibleVehicleItem{}
	combinedpayloads := make([]map[string]interface{}, 0)
	// tableStruct := masterpayloads.FieldActionItemDetailResponse{}

	// baseModelQuery := utils.CreateJoinSelectStatement(tx, tableStruct).Where(masterentities.FieldActionEligibleVehicle{FieldActionEligibleVehicleSystemNumber: Id})
	err := tx.Model(&entitiesitem).
		Where(masterentities.FieldActionEligibleVehicleItem{
			FieldActionEligibleVehicleSystemNumber: Id,
		}).Joins("JOIN mtr_item ON mtr_item.item_id=mtr_field_action_eligible_vehicle_item.item_id").Select("mtr_field_action_eligible_vehicle_item.*,mtr_item.*").
		Scan(&payloadsitem).Error

	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}
	for _, it := range payloadsitem {
		combinedpayloads = append(combinedpayloads, map[string]interface{}{
			"is_active": it.IsActive,
			"field_action_eligible_vehicle_item_system_number": it.FieldActionEligibleVehicleItemSystemNumber,
			"field_action_eligible_vehicle_system_number":      it.FieldActionEligibleVehicleSystemNumber,
			"line_type_id": it.LineTypeId,
			"field_action_eligible_vehicle_line_number": it.FieldActionEligibleVehicleItemLineNumber,
			"item_id":          it.ItemId,
			"item_name":        it.ItemName,
			"field_action_frt": it.FieldActionFrt,
		})
	}

	err2 := tx.Model(&entitiesoperation).
		Where("field_action_eligible_vehicle_system_number =?",Id).Joins("JOIN mtr_operation_model_mapping on mtr_operation_model_mapping.operation_model_mapping_id=mtr_field_action_eligible_vehicle_operation.operation_model_mapping_id").Select("mtr_field_action_eligible_vehicle_operation.*,mtr_operation_model_mapping.*").
		Scan(&payloadsoperation).Error

	if err2 != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err2,
		}
	}
	for _, it := range payloadsoperation {
		combinedpayloads = append(combinedpayloads, map[string]interface{}{
			"is_active": it.IsActive,
			"field_action_eligible_vehicle_item_system_number": it.FieldActionEligibleVehicleItemSystemNumber,
			"field_action_eligible_vehicle_system_number":      it.FieldActionEligibleVehicleSystemNumber,
			"line_type_id": it.LineTypeId,
			"field_action_eligible_vehicle_line_number": it.FieldActionEligibleVehicleItemLineNumber,
			"operation_id":     it.OperationModelMappingId,
			"operation_name":   it.OperationName,
			"field_action_frt": it.FieldActionFrt,
		})
	}
	dataPaginate, totalPages, totalRows := pagination.NewDataFramePaginate(combinedpayloads, &pages)
	return dataPaginate, totalPages, totalRows, nil
}

func (r *FieldActionRepositoryImpl) GetFieldActionVehicleItemDetailById(tx *gorm.DB, Id int, LineTypeId int) (map[string]interface{}, *exceptions.BaseErrorResponse) {
	entitiesitem := masterentities.FieldActionEligibleVehicleItem{}
	entitiesoperation := masterentities.FieldActionEligibleVehicleOperation{}
	responseitem := masterpayloads.FieldActionEligibleVehicleItem{}
	responseoperation := masterpayloads.FieldActionEligibleVehicleOperation{}
	if LineTypeId == 5 {
		err := tx.Model(&entitiesoperation).
			Where(masterentities.FieldActionEligibleVehicleOperation{
				FieldActionEligibleVehicleOperationSystemNumber: Id,
			}).Joins("JOIN mtr_operation_model_mapping on mtr_operation_model_mapping.operation_model_mapping_id = mtr_field_action_eligible_vehicle_operation.operation_model_mapping_id").Select("mtr_field_action_eligible_vehicle_operation.*,mtr_operation_model_mapping.*").
			First(&responseoperation).Error

		if err != nil {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
		responsepayloads := map[string]interface{}{
			"is_active": responseoperation.IsActive,
			"field_action_eligible_vehicle_item_system_number": responseoperation.FieldActionEligibleVehicleItemSystemNumber,
			"field_action_eligible_vehicle_system_number":      responseoperation.FieldActionEligibleVehicleSystemNumber,
			"line_type_id": responseoperation.LineTypeId,
			"field_action_eligible_vehicle_line_number": responseoperation.FieldActionEligibleVehicleItemLineNumber,
			"operation_id":     responseoperation.OperationModelMappingId,
			"operation_name":   responseoperation.OperationName,
			"field_action_frt": responseoperation.FieldActionFrt,
		}
		return responsepayloads, nil
	} else {
		err := tx.Model(&entitiesitem).
			Where(masterentities.FieldActionEligibleVehicleItem{
				FieldActionEligibleVehicleItemSystemNumber: Id,
			}).Joins("JOIN mtr_item on mtr_item.item_id = mtr_field_action_eligible_vehicle_item.item_id").Select("mtr_field_action_eligible_vehicle_item.*,mtr_item.*").
			First(&responseitem).Error

		responsepayloads := map[string]interface{}{
			"is_active": responseitem.IsActive,
			"field_action_eligible_vehicle_item_system_number": responseitem.FieldActionEligibleVehicleItemSystemNumber,
			"field_action_eligible_vehicle_system_number":      responseitem.FieldActionEligibleVehicleSystemNumber,
			"line_type_id": responseitem.LineTypeId,
			"field_action_eligible_vehicle_line_number": responseitem.FieldActionEligibleVehicleItemLineNumber,
			"item_id":          responseitem.ItemId,
			"item_name":        responseitem.ItemName,
			"field_action_frt": responseitem.FieldActionFrt,
		}
		if err != nil {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
		return responsepayloads, nil
	}

}

func (r *FieldActionRepositoryImpl) PostFieldActionVehicleItemDetail(tx *gorm.DB, req masterpayloads.FieldActionItemDetailResponse, id int) (bool, *exceptions.BaseErrorResponse) {
	if req.LineTypeId == 5 {
		entities := masterentities.FieldActionEligibleVehicleOperation{
			FieldActionEligibleVehicleOperationSystemNumber: req.FieldActionEligibleVehicleItemSystemNumber,
			FieldActionEligibleVehicleSystemNumber:     id,
			LineTypeId:                                 req.LineTypeId,
			FieldActionEligibleVehicleItemLineNumber:   req.FieldActionEligibleVehicleItemLineNumber,
			OperationModelMappingId:                    req.ItemOperationCode,
			FieldActionFrt:                             req.FieldActionFrt,
		}

		err := tx.Save(&entities).Error

		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
	} else {
		entities := masterentities.FieldActionEligibleVehicleItem{
			FieldActionEligibleVehicleItemSystemNumber: req.FieldActionEligibleVehicleItemSystemNumber,
			FieldActionEligibleVehicleSystemNumber:     id,
			LineTypeId:                                 req.LineTypeId,
			FieldActionEligibleVehicleItemLineNumber:   req.FieldActionEligibleVehicleItemLineNumber,
			ItemId:                                     req.ItemOperationCode,
			FieldActionFrt:                             req.FieldActionFrt,
		}

		err := tx.Save(&entities).Error

		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
	}

	return true, nil
}

func (r *FieldActionRepositoryImpl) PostFieldActionVehicleDetail(tx *gorm.DB, req masterpayloads.FieldActionDetailResponse, id int) (bool, *exceptions.BaseErrorResponse) {
	entities := masterentities.FieldActionEligibleVehicle{
		FieldActionEligibleVehicleSystemNumber: req.FieldActionEligibleVehicleSystemNumber,
		FieldActionRecallLineNumber:            req.FieldActionRecallLineNumber,
		FieldActionSystemNumber:                id,
		VehicleId:                              req.VehicleId,
		CompanyId:                              req.CompanyId,
		FieldActionDate:                        req.FieldActionDate,
		FieldActionHasTaken:                    req.FieldActionHasTaken,
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

func (r *FieldActionRepositoryImpl) PostMultipleVehicleDetail(tx *gorm.DB, headerId int, id string) (bool, *exceptions.BaseErrorResponse) {

	var entities masterentities.FieldActionEligibleVehicle
	var entityToUpdate []masterentities.FieldActionEligibleVehicle
	strid := strings.Split(id, ",")

	var strids []int

	for _, numid := range strid {
		num, err := strconv.Atoi(numid)
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
		strids = append(strids, num)

	}

	err := tx.Model(&entities).Where("vehicle_id in (?) AND field_action_system_number == ?", strids, headerId).Scan(&entityToUpdate).Error

	if len(entityToUpdate) != 0 {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	for _, value := range strids {
		data := masterentities.FieldActionEligibleVehicle{
			IsActive:                true,
			FieldActionSystemNumber: headerId,
			// CompanyId:               companyId,
			VehicleId:           value,
			FieldActionHasTaken: false,
		}

		err := tx.Save(&data).Error

		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

	}

	return true, nil

}

func (r *FieldActionRepositoryImpl) PostVehicleItemIntoAllVehicleDetail(tx *gorm.DB, headerId int, req masterpayloads.FieldActionItemDetailResponse) (bool, *exceptions.BaseErrorResponse) {

	// var entities masterentities.FieldActionEligibleVehicleItem
	var headerEntities masterentities.FieldActionEligibleVehicle

	var listId []masterentities.FieldActionEligibleVehicle
	// strid := strings.Split(headerId, ",")

	tx.Model(&headerEntities).Where("field_action_system_number = ?", headerId).Scan(&listId)

	for _, value := range listId {
		if req.LineTypeId == 5 {
			entities := masterentities.FieldActionEligibleVehicleOperation{
				FieldActionEligibleVehicleOperationSystemNumber: req.FieldActionEligibleVehicleItemSystemNumber,
				FieldActionEligibleVehicleSystemNumber:     value.FieldActionSystemNumber,
				LineTypeId:                                 req.LineTypeId,
				FieldActionEligibleVehicleItemLineNumber:   req.FieldActionEligibleVehicleItemLineNumber,
				OperationModelMappingId:                    req.ItemOperationCode,
				FieldActionFrt:                             req.FieldActionFrt,
			}

			err := tx.Save(&entities).Error

			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        err,
				}
			}
		} else {
			entities := masterentities.FieldActionEligibleVehicleItem{
				FieldActionEligibleVehicleItemSystemNumber: req.FieldActionEligibleVehicleItemSystemNumber,
				FieldActionEligibleVehicleSystemNumber:     value.FieldActionSystemNumber,
				LineTypeId:                                 req.LineTypeId,
				FieldActionEligibleVehicleItemLineNumber:   req.FieldActionEligibleVehicleItemLineNumber,
				ItemId:                                     req.ItemOperationCode,
				FieldActionFrt:                             req.FieldActionFrt,
			}

			err := tx.Save(&entities).Error

			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        err,
				}
			}
		}
	}

	return true, nil

}

func (r *FieldActionRepositoryImpl) ChangeStatusFieldAction(tx *gorm.DB, id int) (bool, *exceptions.BaseErrorResponse) {
	var entities masterentities.FieldAction

	result := tx.Model(&entities).
		Where("field_action_system_number = ?", id).
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

func (r *FieldActionRepositoryImpl) ChangeStatusFieldActionVehicle(tx *gorm.DB, id int) (bool, *exceptions.BaseErrorResponse) {
	var entities masterentities.FieldActionEligibleVehicle

	result := tx.Model(&entities).
		Where("field_action_eligible_vehicle_system_number = ?", id).
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

func (r *FieldActionRepositoryImpl) ChangeStatusFieldActionVehicleItem(tx *gorm.DB, id int) (bool, *exceptions.BaseErrorResponse) {
	var entities masterentities.FieldActionEligibleVehicleItem

	result := tx.Model(&entities).
		Where("field_action_eligible_vehicle_item_system_number = ?", id).
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
