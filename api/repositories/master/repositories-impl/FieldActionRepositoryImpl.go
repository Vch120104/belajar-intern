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
	JoinTable := tx.Model(&entities).
		Select("mtr_field_action.*,faev.*").
		Joins("Join mtr_field_action_eligible_vehicle as faev ON faev.field_action_system_number=mtr_field_action.field_action_system_number")

	whereQuery := utils.ApplyFilter(JoinTable, filterCondition)
	err := whereQuery.Scopes(pagination.Paginate(&pages, JoinTable)).Order("mtr_field_action.field_action_system_number").Scan(&responses).Error

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
	rows, err := baseModelQuery.Scopes(pagination.Paginate(&pages, filterQuery)).Scan(&payloads).Rows()

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

func (r *FieldActionRepositoryImpl) GetAllFieldActionVehicleItemOperationDetailById(tx *gorm.DB, Id int, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	entities := []masterentities.FieldActionEligibleVehicleItemOperation{}
	payloads := []masterpayloads.FieldActionEligibleVehicleItemOperationResp{}
	// combinedpayloads := make([]map[string]interface{}, 0)
	// tableStruct := masterpayloads.FieldActionItemDetailResponse{}

	// baseModelQuery := utils.CreateJoinSelectStatement(tx, tableStruct).Where(masterentities.FieldActionEligibleVehicle{FieldActionEligibleVehicleSystemNumber: Id})
	baseModelQuery := tx.Model(&entities).Select("mtr_field_action_eligible_vehicle_item_operation.*,mtr_item.*,mtr_operation_code.*").
		Joins("JOIN mtr_mapping_item_operation ON mtr_mapping_item_operation.item_operation_id=mtr_field_action_eligible_vehicle_item_operation.item_operation_id").
		Joins("INNER JOIN mtr_item ON mtr_item.item_id=mtr_field_action_eligible_vehicle_item_operation.item_id").
		Joins("JOIN mtr_operation_code ON mtr_operation_code.operation_id=mtr_field_action_eligible_vehicle_item_operation.operation_id").
		Where(masterentities.FieldActionEligibleVehicleItemOperation{
			FieldActionEligibleVehicleSystemNumber: Id,
		})

	// filterQuery := utils.ApplyFilter(baseModelQuery, filterCondition)
	rows, err := baseModelQuery.Scopes(pagination.Paginate(&pages, baseModelQuery)).Scan(&payloads).Rows()

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

	// dataPaginate, totalPages, totalRows := pagination.NewDataFramePaginate(payloads, &pages)
	// return dataPaginate, totalPages, totalRows, nil
}

func (r *FieldActionRepositoryImpl) GetFieldActionVehicleItemDetailById(tx *gorm.DB, Id int) (masterpayloads.FieldActionEligibleVehicleItemOperationResp, *exceptions.BaseErrorResponse) {
	entities := masterentities.FieldActionEligibleVehicleItemOperation{}
	response := masterpayloads.FieldActionEligibleVehicleItemOperationResp{}

	err := tx.Model(&entities).Select("mtr_field_action_eligible_vehicle_item_operation.*,mtr_item.*,mtr_operation_code.*").
		Joins("JOIN mtr_mapping_item_operation ON mtr_mapping_item_operation.item_operation_id=mtr_field_action_eligible_vehicle_item_operation.item_operation_id").
		Joins("INNER JOIN mtr_item ON mtr_item.item_id=mtr_field_action_eligible_vehicle_item_operation.item_id").
		Joins("JOIN mtr_operation_code ON mtr_operation_code.operation_id=mtr_field_action_eligible_vehicle_item_operation.operation_id").
		Where(masterentities.FieldActionEligibleVehicleItemOperation{
			FieldActionEligibleVehicleItemOperationSystemNumber: Id,
		}).
		First(&response).Error

	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	return response, nil

}

func (r *FieldActionRepositoryImpl) PostFieldActionVehicleItemDetail(tx *gorm.DB, req masterpayloads.FieldActionEligibleVehicleItemOperationRequest, id int) (bool, *exceptions.BaseErrorResponse) {

	entities := masterentities.FieldActionEligibleVehicleItemOperation{
		FieldActionEligibleVehicleItemOperationSystemNumber: req.FieldActionEligibleVehicleItemOperationSystemNumber,
		FieldActionEligibleVehicleSystemNumber:              id,
		FieldActionFrt:                                      req.FieldActionFrt,
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

func (r *FieldActionRepositoryImpl) PostVehicleItemIntoAllVehicleDetail(tx *gorm.DB, headerId int, req masterpayloads.FieldActionEligibleVehicleItemOperationRequest) (bool, *exceptions.BaseErrorResponse) {

	// var entities masterentities.FieldActionEligibleVehicleItem
	var headerEntities masterentities.FieldActionEligibleVehicle

	var listId []masterentities.FieldActionEligibleVehicle
	// strid := strings.Split(headerId, ",")

	tx.Model(&headerEntities).Where("field_action_system_number = ?", headerId).Scan(&listId)

	for _, value := range listId {
		entities := masterentities.FieldActionEligibleVehicleItemOperation{
			FieldActionEligibleVehicleItemOperationSystemNumber: req.FieldActionEligibleVehicleItemOperationSystemNumber,
			FieldActionEligibleVehicleSystemNumber:              value.FieldActionSystemNumber,
			FieldActionFrt:                                      req.FieldActionFrt,
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
	var entities masterentities.FieldActionEligibleVehicleItemOperation

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
