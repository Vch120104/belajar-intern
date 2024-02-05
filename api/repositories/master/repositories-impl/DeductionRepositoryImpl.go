package masterrepositoryimpl

import (
	masterentities "after-sales/api/entities/master"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type DeductionRepositoryImpl struct {
}

func StartDeductionRepositoryImpl() masterrepository.DeductionRepository {
	return &DeductionRepositoryImpl{}
}

func (r *DeductionRepositoryImpl) GetAllDeduction(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, error) {

	entities := []masterentities.DeductionList{}

	baseModelQuery := tx.Model(&entities)

	wherequery := utils.ApplyFilter(baseModelQuery, filterCondition)

	rows, err := baseModelQuery.Scopes(pagination.Paginate(&entities, &pages, wherequery)).Scan(&entities).Rows()

	if len(entities) == 0 {
		return pages, gorm.ErrRecordNotFound
	}

	if err != nil {
		return pages, err
	}

	defer rows.Close()

	pages.Rows = entities

	return pages, nil
}

func (r *DeductionRepositoryImpl) GetDeductionById(tx *gorm.DB, Id int) (masterpayloads.DeductionListResponse, error) {

	entities := masterentities.DeductionList{}

	response := masterpayloads.DeductionListResponse{}

	rows, err := tx.Model(&entities).
		Where(masterentities.DeductionList{
			DeductionListId: int(Id),
		}).
		First(&response).
		Rows()

	if err != nil {
		return response, err
	}

	defer rows.Close()

	return response, nil
}

func (r *DeductionRepositoryImpl) GetAllDeductionDetail(tx *gorm.DB, pages pagination.Pagination, Id int) (pagination.Pagination, error) {

	entities := []masterentities.DeductionDetail{}

	response := []masterpayloads.DeductionDetailResponse{}

	baseModelQuery := tx.Model(&entities).
		Where(masterentities.DeductionDetail{
			DeductionListId: Id},
		)

	rows, err := baseModelQuery.Scopes(pagination.Paginate(&entities, &pages, baseModelQuery)).Scan(&response).Rows()

	if len(response) == 0 {
		return pages, gorm.ErrRecordNotFound
	}

	if err != nil {
		return pages, err
	}

	defer rows.Close()

	pages.Rows = response

	return pages, nil
}

func (r *DeductionRepositoryImpl) GetByIdDeductionDetail(tx *gorm.DB, Id int) (masterpayloads.DeductionDetailResponse, error) {

	entities := masterentities.DeductionDetail{}

	response := masterpayloads.DeductionDetailResponse{}

	rows, err := tx.Model(&entities).Where(
		masterentities.DeductionDetail{
			DeductionDetailId: int(Id),
		}).First(&response).Rows()

	if err != nil {
		return response, err
	}
	defer rows.Close()

	return response, nil
}

func (r *DeductionRepositoryImpl) SaveDeductionList(tx *gorm.DB, request masterpayloads.DeductionListResponse) (masterpayloads.DeductionListResponse, error) {

	entities := masterentities.DeductionList{
		DeductionName: request.DeductionName,
		EffectiveDate: request.EffectiveDate,
	}

	result := tx.Where(entities).Assign(entities).FirstOrCreate(&entities)

	if result.Error != nil {
		return masterpayloads.DeductionListResponse{}, result.Error
	}

	return request, nil
}

func (r *DeductionRepositoryImpl) SaveDeductionDetail(tx *gorm.DB, request masterpayloads.DeductionDetailResponse) (masterpayloads.DeductionDetailResponse, error) {
	condition := masterentities.DeductionDetail{
		DeductionListId:      request.DeductionListId,
		DeductionDetailLevel: request.DeductionDetailLevel,
	}

	entities := masterentities.DeductionDetail{
		DeductionDetailCode:  request.DeductionDetailCode,
		DeductionListId:      request.DeductionListId,
		DeductionDetailLevel: request.DeductionDetailLevel,
		DeductionPercent:     request.DeductionPercent,
	}

	result := tx.Where(entities).Assign(condition).FirstOrCreate(&entities)

	if result.Error != nil {
		return masterpayloads.DeductionDetailResponse{}, result.Error
	}

	return request, nil
}
