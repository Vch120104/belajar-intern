package masteroperationrepositoryimpl

import (
	masteroperationentities "after-sales/api/entities/master/operation"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	masteroperationrepository "after-sales/api/repositories/master/operation"

	"after-sales/api/utils"
	"log"

	"gorm.io/gorm"
)

type OperationSectionRepositoryImpl struct {
	myDB *gorm.DB
}

func StartOperationSectionRepositoryImpl(db *gorm.DB) masteroperationrepository.OperationSectionRepository {
	return &OperationSectionRepositoryImpl{myDB: db}
}

func (r *OperationSectionRepositoryImpl) WithTrx(trxHandle *gorm.DB) masteroperationrepository.OperationSectionRepository {
	if trxHandle == nil {
		log.Println("Transaction Database Not Found!")
		return r
	}
	r.myDB = trxHandle
	return r
}

func (r *OperationSectionRepositoryImpl) GetAllOperationSection() ([]masteroperationpayloads.OperationSectionResponse, error) {
	var OperationSections masteroperationentities.OperationSection
	var response []masteroperationpayloads.OperationSectionResponse

	rows, err := r.myDB.Model(&OperationSections).Scan(response).Rows()

	if err != nil {
		return response, err
	}

	defer rows.Close()

	return response, nil
}

func (r *OperationSectionRepositoryImpl) GetAllOperationSectionList(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, error) {
	entities := masteroperationentities.OperationSection{}
	var responses []masteroperationpayloads.OperationSectionListResponse
	// define table struct
	tableStruct := masteroperationpayloads.OperationSectionListResponse{}
	//define join table
	joinTable := utils.CreateJoinSelectStatement(r.myDB, tableStruct)
	//apply filter
	whereQuery := utils.ApplyFilter(joinTable, filterCondition)
	//apply pagination and execute
	rows, err := joinTable.Scopes(pagination.Paginate(&entities, &pages, whereQuery)).Scan(&responses).Rows()

	if err != nil {
		return pages, err
	}

	defer rows.Close()

	pages.Rows = responses

	return pages, nil

}

func (r *OperationSectionRepositoryImpl) GetOperationSectionName(GroupId int, SectionCode string) (masteroperationpayloads.OperationSectionNameResponse, error) {
	tableStruct := masteroperationpayloads.OperationSectionNameResponse{}

	joinTable := utils.CreateJoinSelectStatement(r.myDB, tableStruct)

	row, err := joinTable.Where("mtr_operation_group.operation_group_id = ?", GroupId).
		Where("mtr_operation_section.operation_section_code = ?", SectionCode).
		First(&tableStruct).Rows()

	if err != nil {
		return tableStruct, err
	}

	defer row.Close()

	return tableStruct, nil
}

func (r *OperationSectionRepositoryImpl) GetSectionCodeByGroupId(GroupId string) ([]masteroperationpayloads.OperationSectionCodeResponse, error) {
	tableStruct := masteroperationpayloads.OperationSectionCodeResponse{}
	var sliceTableStruct []masteroperationpayloads.OperationSectionCodeResponse

	joinTable := utils.CreateJoinSelectStatement(r.myDB, tableStruct)

	WhereQuery := joinTable.
		Where("mtr_operation_group.operation_group_id = ?", GroupId).
		Where("mtr_operation_section.is_active = 1")

	rows, err := WhereQuery.Scan(&sliceTableStruct).Rows()

	if err != nil {
		return sliceTableStruct, err
	}
	defer rows.Close()

	return sliceTableStruct, nil
}

func (r *OperationSectionRepositoryImpl) GetOperationSectionById(Id int) (masteroperationpayloads.OperationSectionResponse, error) {
	entities := masteroperationentities.OperationSection{}
	response := masteroperationpayloads.OperationSectionResponse{}

	rows, err := r.myDB.Model(&entities).
		Where(masteroperationentities.OperationSection{
			OperationSectionId: Id,
		}).
		First(&response).
		Rows()

	if err != nil {
		return response, err
	}

	defer rows.Close()

	return response, nil
}

func (r *OperationSectionRepositoryImpl) SaveOperationSection(request masteroperationpayloads.OperationSectionRequest) (bool, error) {
	entities := masteroperationentities.OperationSection{
		IsActive:                    request.IsActive,
		OperationSectionId:          request.OperationSectionId,
		OperationSectionCode:        request.OperationSectionCode,
		OperationGroupId:            request.OperationGroupId,
		OperationSectionDescription: request.OperationSectionDescription,
	}

	err := r.myDB.Save(&entities).Error

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *OperationSectionRepositoryImpl) ChangeStatusOperationSection(Id int) (bool, error) {
	var entities masteroperationentities.OperationSection

	result := r.myDB.Model(&entities).
		Where("operation_section_id = ?", Id).
		First(&entities)

	if result.Error != nil {
		return false, result.Error
	}

	// Toggle the IsActive value
	if entities.IsActive {
		entities.IsActive = false
	} else {
		entities.IsActive = true
	}

	result = r.myDB.Save(&entities)

	if result.Error != nil {
		return false, result.Error
	}

	return true, nil
}
