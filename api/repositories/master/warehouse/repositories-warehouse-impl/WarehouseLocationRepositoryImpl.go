package masterwarehouserepositoryimpl

import (
	// masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	pagination "after-sales/api/payloads/pagination"
	masterwarehouserepository "after-sales/api/repositories/master/warehouse"
	utils "after-sales/api/utils"

	// utils "after-sales/api/utils"

	// masterwarehousegroupservice "after-sales/api/services/master/warehouse"
	masterwarehouseentities "after-sales/api/entities/master/warehouse"
	// "after-sales/api/payloads/pagination"

	"log"

	"gorm.io/gorm"
)

type WarehouseLocationImpl struct {
	DB *gorm.DB
}

func OpenWarehouseLocationImpl(db *gorm.DB) masterwarehouserepository.WarehouseLocationRepository {
	return &WarehouseLocationImpl{DB: db}
}

func (r *WarehouseLocationImpl) WithTrx(Trxhandle *gorm.DB) masterwarehouserepository.WarehouseLocationRepository {
	if Trxhandle == nil {
		log.Println("Transaction Database Not Found")
		return r
	}
	r.DB = Trxhandle
	return r
}

func (r *WarehouseLocationImpl) Save(request masterwarehousepayloads.GetWarehouseLocationResponse) (bool, error) {

	var warehouseMaster = masterwarehouseentities.WarehouseLocation{
		IsActive:                      utils.BoolPtr(request.IsActive),
		WarehouseLocationId:           request.WarehouseLocationId,
		CompanyId:                     request.CompanyId,
		WarehouseGroupId:              request.WarehouseGroupId,
		WarehouseLocationCode:         request.WarehouseLocationCode,
		WarehouseLocationName:         request.WarehouseLocationName,
		WarehouseLocationDetailName:   request.WarehouseLocationDetailName,
		WarehouseLocationPickSequence: request.WarehouseLocationPickSequence,
		WarehouseLocationCapacityInM3: request.WarehouseLocationCapacityInM3,
	}

	rows, err := r.DB.Model(&warehouseMaster).
		Save(&warehouseMaster).
		Rows()

	if err != nil {
		return false, err
	}

	defer rows.Close()

	return true, nil
}

func (r *WarehouseLocationImpl) GetById(warehouseLocationId int) (masterwarehousepayloads.GetWarehouseLocationResponse, error) {

	var entities masterwarehouseentities.WarehouseLocation
	var warehouseLocationResponse masterwarehousepayloads.GetWarehouseLocationResponse

	rows, err := r.DB.Model(&entities).
		Where(masterwarehousepayloads.GetWarehouseLocationResponse{
			WarehouseLocationId: warehouseLocationId,
		}).
		Scan(&warehouseLocationResponse).
		// Find(&warehouseMasterResponse).
		Rows()

	if err != nil {
		return warehouseLocationResponse, err
	}

	defer rows.Close()

	return warehouseLocationResponse, nil
}

func (r *WarehouseLocationImpl) GetAll(request masterwarehousepayloads.GetAllWarehouseLocationRequest, pages pagination.Pagination) (pagination.Pagination, error) {
	var entities []masterwarehouseentities.WarehouseLocation
	var warehouseLocationResponse []masterwarehousepayloads.GetAllWarehouseLocationResponse

	tempRows := r.DB.
		Model(&masterwarehouseentities.WarehouseLocation{}).
		Where("warehouse_location_code like ?", "%"+request.WarehouseLocationCode+"%").
		Where("warehouse_location_name like ?", "%"+request.WarehouseLocationName+"%").
		Where("warehouse_location_detail_name like ?", "%"+request.WarehouseLocationDetailName+"%")

	if request.IsActive != "" {
		tempRows = tempRows.Where("is_active = ?", request.IsActive)
	}

	if request.CompanyId != "" {
		tempRows = tempRows.Where("company_id = ?", request.CompanyId)
	}

	rows, err := tempRows.
		Scopes(pagination.Paginate(entities, &pages, tempRows)).
		Scan(&warehouseLocationResponse).
		Rows()

	if err != nil {
		return pages, err
	}

	defer rows.Close()

	pages.Rows = warehouseLocationResponse
	return pages, nil
}

func (r *WarehouseLocationImpl) ChangeStatus(warehouseLocationId int) (masterwarehousepayloads.GetWarehouseLocationResponse, error) {
	var entities masterwarehouseentities.WarehouseLocation
	var warehouseLocationPayloads masterwarehousepayloads.GetWarehouseLocationResponse

	rows, err := r.DB.Model(&entities).
		Where(masterwarehousepayloads.GetWarehouseLocationResponse{
			WarehouseLocationId: warehouseLocationId,
		}).
		Update("is_active", gorm.Expr("1 ^ is_active")).
		Rows()

	if err != nil {
		log.Panic((err.Error()))
	}

	defer rows.Close()

	rows, err = r.DB.Model(&entities).
		Where(masterwarehousepayloads.GetWarehouseLocationResponse{
			WarehouseLocationId: warehouseLocationId,
		}).
		// Find(&warehouseMasterPayloads).
		Scan(&warehouseLocationPayloads).
		Rows()

	if err != nil {
		return warehouseLocationPayloads, err
	}

	defer rows.Close()

	return warehouseLocationPayloads, nil
}
