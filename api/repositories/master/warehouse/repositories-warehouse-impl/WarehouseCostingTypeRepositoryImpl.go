package masterwarehouserepositoryimpl

import (
	masterwarehouseentities "after-sales/api/entities/master/warehouse"
	"after-sales/api/exceptions"
	masterwarehouserepository "after-sales/api/repositories/master/warehouse"
	"errors"
	"gorm.io/gorm"
	"net/http"
)

type WarehouseCostingTypeRepositoryImpl struct {
}

func (w *WarehouseCostingTypeRepositoryImpl) GetByCodeWarehouseCostingType(db *gorm.DB, CostingTypeCode string) (masterwarehouseentities.WarehouseCostingType, *exceptions.BaseErrorResponse) {
	CostingTypeEntities := masterwarehouseentities.WarehouseCostingType{}
	err := db.Model(&masterwarehouseentities.WarehouseCostingType{}).
		Where(masterwarehouseentities.WarehouseCostingType{WarehouseCostingTypeCode: CostingTypeCode}).
		First(&CostingTypeEntities).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return CostingTypeEntities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        errors.New("costing type with code " + CostingTypeCode + " not found"),
			}
		}
		return CostingTypeEntities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	return CostingTypeEntities, nil
}

func NewWarehouseCostingTypeRepositoryImpl() masterwarehouserepository.WarehouseCostingTypeRepository {
	return &WarehouseCostingTypeRepositoryImpl{}
}
