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

	"gorm.io/gorm"
)

type ItemImportRepositoryImpl struct {
}

// GetItemImportbyId implements masteritemrepository.ItemImportRepository.
func (i *ItemImportRepositoryImpl) GetItemImportbyId(tx *gorm.DB, Id int) (any, *exceptions.BaseErrorResponse) {
	model := masteritementities.ItemImport{}
	response := masteritempayloads.ItemImportByIdResponse{}
	supplierResponses := masteritempayloads.SupplierResponse{}

	query := tx.Model(&model).Select("mtr_item_import.*, Item.item_code AS item_code, Item.item_name AS item_name").Where(masteritementities.ItemImport{ItemImportId: Id}).
		InnerJoins("Item", tx.Select(""))

	err := query.First(&response).Error

	if err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	supplierUrl := config.EnvConfigs.GeneralServiceUrl + "supplier-master/" + strconv.Itoa(response.SupplierId)

	if errSupplier := utils.Get(supplierUrl, &supplierResponses, nil); errSupplier != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New(""),
		}
	}

	joinedDataSupplier := utils.DataFrameInnerJoin([]masteritempayloads.ItemImportByIdResponse{response}, []masteritempayloads.SupplierResponse{supplierResponses}, "SupplierId")

	if len(joinedDataSupplier) == 0 {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New(""),
		}
	}

	return joinedDataSupplier[0], nil

}

// GetAllItemImport implements masteritemrepository.ItemImportRepository.
func (i *ItemImportRepositoryImpl) GetAllItemImport(tx *gorm.DB, internalFilter []utils.FilterCondition, externalFilter []utils.FilterCondition, pages pagination.Pagination) ([]map[string]any, int, int, *exceptions.BaseErrorResponse) {
	model := masteritementities.ItemImport{}
	var responses []masteritempayloads.ItemImportResponse
	var supplierResponses []masteritempayloads.SupplierResponse
	var supplierCode string
	var supplierName string

	for _, values := range externalFilter {
		if values.ColumnField == "supplier_code" {
			supplierCode = values.ColumnValue
		} else {
			supplierName = values.ColumnValue
		}
	}

	query := tx.Model(&model).Select("mtr_item_import.*, Item.item_code AS item_code, Item.item_name AS item_name").
		InnerJoins("Item", tx.Select(""))

	whereQuery := utils.ApplyFilter(query, internalFilter)

	err := whereQuery.Scan(&responses).Error

	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if len(responses) == 0 {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New(""),
		}
	}

	supplierUrl := config.EnvConfigs.GeneralServiceUrl + "supplier-master?page=0&limit=10000&supplier_code=" + supplierCode + "&supplier_name=" + supplierName

	if errSupplier := utils.Get(supplierUrl, &supplierResponses, nil); errSupplier != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New(""),
		}
	}

	joinedDataSupplier := utils.DataFrameInnerJoin(responses, supplierResponses, "SupplierId")

	if len(joinedDataSupplier) == 0 {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New(""),
		}
	}

	dataPaginate, totalPages, totalRows := pagination.NewDataFramePaginate(joinedDataSupplier, &pages)

	return dataPaginate, totalPages, totalRows, nil

}

// SaveItemImport implements masteritemrepository.ItemImportRepository.
func (i *ItemImportRepositoryImpl) SaveItemImport(tx *gorm.DB, req masteritementities.ItemImport) (bool, *exceptions.BaseErrorResponse) {
	entities := masteritementities.ItemImport{
		SupplierId:         req.SupplierId,
		ItemId:             req.ItemId,
		OrderQtyMultiplier: req.OrderQtyMultiplier,
		ItemAliasCode:      req.ItemAliasCode,
		RoyaltyFlag:        req.RoyaltyFlag,
		ItemAliasName:      req.ItemAliasName,
		OrderConversion:    req.OrderConversion,
	}
	supplierResponse := masteritempayloads.SupplierResponse{}
	getSupplierbyIdUrl := config.EnvConfigs.GeneralServiceUrl + "api/general/supplier-master/" + strconv.Itoa(req.SupplierId)

	errGetSupplier := utils.Get(getSupplierbyIdUrl, &supplierResponse, nil)

	if errGetSupplier != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errGetSupplier,
		}
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

// UpdateItemImport implements masteritemrepository.ItemImportRepository.
func (i *ItemImportRepositoryImpl) UpdateItemImport(tx *gorm.DB, req masteritementities.ItemImport) (bool, *exceptions.BaseErrorResponse) {
	entities := masteritementities.ItemImport{
		ItemImportId:       req.ItemImportId,
		SupplierId:         req.SupplierId,
		ItemId:             req.ItemId,
		OrderQtyMultiplier: req.OrderQtyMultiplier,
		ItemAliasCode:      req.ItemAliasCode,
		RoyaltyFlag:        req.RoyaltyFlag,
		ItemAliasName:      req.ItemAliasName,
		OrderConversion:    req.OrderConversion,
	}

	supplierResponse := masteritempayloads.SupplierResponse{}
	getSupplierbyIdUrl := config.EnvConfigs.GeneralServiceUrl + "api/general/supplier-master/" + strconv.Itoa(req.SupplierId)

	errGetSupplier := utils.Get(getSupplierbyIdUrl, &supplierResponse, nil)

	if errGetSupplier != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errGetSupplier,
		}
	}

	err := tx.Updates(&entities).Where(masteritementities.ItemImport{ItemImportId: req.ItemImportId}).Error

	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return true, nil
}

func StartItemImportRepositoryImpl() masteritemrepository.ItemImportRepository {
	return &ItemImportRepositoryImpl{}
}
