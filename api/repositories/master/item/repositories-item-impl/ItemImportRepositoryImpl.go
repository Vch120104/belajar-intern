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
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type ItemImportRepositoryImpl struct {
}

// SaveItemImport implements masteritemrepository.ItemImportRepository.
func (i *ItemImportRepositoryImpl) SaveItemImport(tx *gorm.DB, req masteritementities.ItemImport) (bool, *exceptions.BaseErrorResponse) {

	supplierResponse := masteritempayloads.SupplierResponse{}
	getSupplierbyIdUrl := config.EnvConfigs.GeneralServiceUrl + "supplier-master/" + strconv.Itoa(req.SupplierId)

	errGetSupplier := utils.Get(getSupplierbyIdUrl, &supplierResponse, nil)

	if supplierResponse == (masteritempayloads.SupplierResponse{}) {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New("supplier not found"),
		}
	}

	if errGetSupplier != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errGetSupplier,
		}
	}

	err := tx.Save(&req).Error

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusConflict,
				Err:        err,
			}
		} else {

			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
	}
	return true, nil
}

// GetItemImportbyItemIdandSupplierId implements masteritemrepository.ItemImportRepository.
func (i *ItemImportRepositoryImpl) GetItemImportbyItemIdandSupplierId(tx *gorm.DB, itemId int, supplierId int) (masteritempayloads.ItemImportByIdResponse, *exceptions.BaseErrorResponse) {
	model := masteritementities.ItemImport{}
	response := masteritempayloads.ItemImportByIdResponse{}
	supplierResponses := masteritempayloads.SupplierResponse{}

	query := tx.Model(&model).Select("mtr_item_import.*, Item.item_code AS item_code, Item.item_name AS item_name").Where(masteritementities.ItemImport{ItemId: itemId, SupplierId: supplierId}).
		InnerJoins("JOIN mtr_item Item ON mtr_item_import.item_id = Item.item_id", tx.Select(""))

	err := query.First(&response).Error

	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	supplierUrl := config.EnvConfigs.GeneralServiceUrl + "supplier-master/" + strconv.Itoa(response.SupplierId)

	if errSupplier := utils.Get(supplierUrl, &supplierResponses, nil); errSupplier != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New(""),
		}
	}

	response.SupplierName = supplierResponses.SupplierName
	response.SupplierCode = supplierResponses.SupplierCode

	return response, nil
}

// GetItemImportbyId implements masteritemrepository.ItemImportRepository.
func (i *ItemImportRepositoryImpl) GetItemImportbyId(tx *gorm.DB, Id int) (masteritempayloads.ItemImportByIdResponse, *exceptions.BaseErrorResponse) {
	model := masteritementities.ItemImport{}
	response := masteritempayloads.ItemImportByIdResponse{}
	supplierResponses := masteritempayloads.SupplierResponse{}

	query := tx.Model(&model).Select("mtr_item_import.*, Item.item_code AS item_code, Item.item_name AS item_name").Where(masteritementities.ItemImport{ItemImportId: Id}).
		InnerJoins("JOIN mtr_item Item ON mtr_item_import.item_id = Item.item_id", tx.Select(""))

	err := query.First(&response).Error

	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	supplierUrl := config.EnvConfigs.GeneralServiceUrl + "supplier-master/" + strconv.Itoa(response.SupplierId)

	if errSupplier := utils.Get(supplierUrl, &supplierResponses, nil); errSupplier != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New(""),
		}
	}

	response.SupplierName = supplierResponses.SupplierName
	response.SupplierCode = supplierResponses.SupplierCode

	return response, nil

}

// GetAllItemImport implements masteritemrepository.ItemImportRepository.

// |
// V
// ERROR!!!, failed to get supplier multi id from external, still on revision from general supplier-multi-ids, last updated (26 Aug 2024, by Kenth)
// ^
// |

func (i *ItemImportRepositoryImpl) GetAllItemImport(tx *gorm.DB, internalFilter []utils.FilterCondition, externalFilter []utils.FilterCondition, pages pagination.Pagination) ([]map[string]any, int, int, *exceptions.BaseErrorResponse) {
	model := masteritementities.ItemImport{}
	var responses []masteritempayloads.ItemImportResponse
	var supplierResponses []masteritempayloads.SupplierResponse
	var supplierMultipleId string
	var supplierCode string
	var supplierName string

	for _, values := range externalFilter {
		if values.ColumnField == "supplier_code" {
			supplierCode = values.ColumnValue
		} else {
			supplierName = values.ColumnValue
		}
	}

	if supplierCode != "" || supplierName != "" {
		supplierUrl := config.EnvConfigs.GeneralServiceUrl + "supplier-master?page=" + strconv.Itoa(0) + "&limit=" + strconv.Itoa(10000000) + "&supplier_code=" + supplierCode + "&supplier_name=" + supplierName

		if errSupplier := utils.Get(supplierUrl, &supplierResponses, nil); errSupplier != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        errors.New(""),
			}
		}

		if len(supplierResponses) == 0 {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        errors.New("supplier not found"),
			}
		}
		for _, value := range supplierResponses {
			supplierMultipleId += strconv.Itoa(value.SupplierId) + ","
		}
	}

	query := tx.Model(&model).Select("mtr_item_import.*, Item.item_code AS item_code, Item.item_name AS item_name").
		InnerJoins("JOIN mtr_item Item ON mtr_item_import.item_id = Item.item_id", tx.Select(""))

	if supplierCode != "" || supplierName != "" {
		query = query.Where("mtr_item_import.supplier_id IN (" + strings.TrimSuffix(supplierMultipleId, ",") + ")")
	}

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
			Err:        errors.New("query item import not found"),
		}
	}

	for _, value := range responses {
		supplierMultipleId += strconv.Itoa(value.SupplierId) + ","
	}

	supplierUrl := config.EnvConfigs.GeneralServiceUrl + "supplier-multi-id/" + supplierMultipleId

	if errSupplier := utils.Get(supplierUrl, &supplierResponses, nil); errSupplier != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errSupplier,
		}
	}

	if len(supplierResponses) == 0 {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New("supplier not found"),
		}
	}

	joinedDataSupplier, errdf := utils.DataFrameInnerJoin(responses, supplierResponses, "SupplierId")

	if errdf != nil {
		return joinedDataSupplier, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errdf,
		}
	}

	dataPaginate, totalPages, totalRows := pagination.NewDataFramePaginate(joinedDataSupplier, &pages)

	return dataPaginate, totalPages, totalRows, nil

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
	getSupplierbyIdUrl := config.EnvConfigs.GeneralServiceUrl + "supplier-master/" + strconv.Itoa(req.SupplierId)

	errGetSupplier := utils.Get(getSupplierbyIdUrl, &supplierResponse, nil)

	fmt.Print(supplierResponse)

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
