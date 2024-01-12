package masterwarehouserepositoryimpl

import (
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	masterwarehouserepository "after-sales/api/repositories/master/warehouse"
	utils "after-sales/api/utils"

	// masterwarehousegroupservice "after-sales/api/services/master/warehouse"
	masterwarehouseentites "after-sales/api/entities/master/warehouse"
	// "after-sales/api/payloads/pagination"

	"log"

	"gorm.io/gorm"
)

type WarehouseGroupImpl struct {
	DB *gorm.DB
}

func OpenWarehouseGroupImpl(db *gorm.DB) masterwarehouserepository.WarehouseGroupRepository {
	return &WarehouseGroupImpl{DB: db}
}

func (r *WarehouseGroupImpl) WithTrx(Trxhandle *gorm.DB) masterwarehouserepository.WarehouseGroupRepository {
	if Trxhandle == nil {
		log.Println("Transaction Database Not Found")
		return r
	}
	r.DB = Trxhandle
	return r
}

func (r *WarehouseGroupImpl) Save(request masterwarehousepayloads.GetWarehouseGroupResponse) (bool, error) {

	var warehouseGroup = masterwarehouseentites.WarehouseGroup{
		IsActive:           utils.BoolPtr(request.IsActive),
		WarehouseGroupId:   request.WarehouseGroupId,
		WarehouseGroupCode: request.WarehouseGroupCode,
		WarehouseGroupName: request.WarehouseGroupName,
		ProfitCenterId:     request.ProfitCenterId,
	}

	rows, err := r.DB.Model(&warehouseGroup).
		Save(&warehouseGroup).
		Rows()

	if err != nil {
		return false, err
	}

	defer rows.Close()

	return true, nil
}

func (r *WarehouseGroupImpl) GetById(warehouseGroupId int) (masterwarehousepayloads.GetWarehouseGroupResponse, error) {

	var entities masterwarehouseentites.WarehouseGroup
	var warehouseGroupResponse masterwarehousepayloads.GetWarehouseGroupResponse

	rows, err := r.DB.Model(&entities).
		Where(masterwarehousepayloads.GetWarehouseGroupResponse{
			WarehouseGroupId: warehouseGroupId,
		}).
		Find(&warehouseGroupResponse).
		Scan(&warehouseGroupResponse).
		Rows()

	if err != nil {
		return warehouseGroupResponse, err
	}

	defer rows.Close()

	return warehouseGroupResponse, nil
}

func (r *WarehouseGroupImpl) GetAll(request masterwarehousepayloads.GetAllWarehouseGroupRequest) ([]masterwarehousepayloads.GetWarehouseGroupResponse, error) {
	var entities []masterwarehouseentites.WarehouseGroup
	var warehouseGroupResponse []masterwarehousepayloads.GetWarehouseGroupResponse
	tempRows := r.DB.
		Model(&entities).
		Where("warehouse_group_code like ?", "%"+request.WarehouseGroupCode+"%").
		Where("warehouse_group_name like ?", "%"+request.WarehouseGroupName+"%")

	if request.IsActive != "" {
		tempRows = tempRows.Where("is_active = ?", request.IsActive)
	}

	rows, err := tempRows.
		Scan(&warehouseGroupResponse).
		Rows()

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return warehouseGroupResponse, nil
}

func (r *WarehouseGroupImpl) ChangeStatus(warehouseGroupId int) (masterwarehousepayloads.GetWarehouseGroupResponse, error) {
	var entities masterwarehouseentites.WarehouseGroup
	var warehouseGroupPayloads masterwarehousepayloads.GetWarehouseGroupResponse

	rows, err := r.DB.Model(&entities).
		Where(masterwarehousepayloads.GetWarehouseGroupResponse{
			WarehouseGroupId: warehouseGroupId,
		}).
		Update("is_active", gorm.Expr("1 ^ is_active")).
		Rows()

	if err != nil {
		log.Panic((err.Error()))
	}

	defer rows.Close()

	rows, err = r.DB.Model(&entities).
		Where(masterwarehousepayloads.GetWarehouseGroupResponse{
			WarehouseGroupId: warehouseGroupId,
		}).
		Find(&warehouseGroupPayloads).
		Scan(&warehouseGroupPayloads).
		Rows()

	if err != nil {
		return warehouseGroupPayloads, err
	}

	defer rows.Close()

	return warehouseGroupPayloads, nil
}
