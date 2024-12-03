package masteritemrepositoryimpl

import (
	"after-sales/api/config"
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	"after-sales/api/utils"
	generalserviceapiutils "after-sales/api/utils/general-service"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type ItemImportRepositoryImpl struct {
}

func StartItemImportRepositoryImpl() masteritemrepository.ItemImportRepository {
	return &ItemImportRepositoryImpl{}
}

// SaveItemImport implements masteritemrepository.ItemImportRepository.
func (i *ItemImportRepositoryImpl) SaveItemImport(tx *gorm.DB, req masteritementities.ItemImport) (bool, *exceptions.BaseErrorResponse) {

	supplierResponse := masteritempayloads.SupplierResponse{}
	getSupplierbyIdUrl := config.EnvConfigs.GeneralServiceUrl + "supplier/" + strconv.Itoa(req.SupplierId)

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

	supplierUrl := config.EnvConfigs.GeneralServiceUrl + "supplier/" + strconv.Itoa(response.SupplierId)

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

	supplierUrl := config.EnvConfigs.GeneralServiceUrl + "supplier/" + strconv.Itoa(response.SupplierId)

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

func (i *ItemImportRepositoryImpl) GetAllItemImport(tx *gorm.DB, internalFilter []utils.FilterCondition, externalFilter []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	model := masteritementities.ItemImport{}
	var responses []masteritempayloads.ItemImportResponse
	var supplierCode string
	var supplierName string

	for _, values := range externalFilter {
		if values.ColumnField == "supplier_code" {
			supplierCode = values.ColumnValue
		} else if values.ColumnField == "supplier_name" {
			supplierName = values.ColumnValue
		}
	}

	if supplierCode != "" || supplierName != "" {
		params := generalserviceapiutils.SupplierMasterParams{
			SupplierCode: supplierCode,
			SupplierName: supplierName,
			Page:         0,
			Limit:        10000000,
		}

		suppliers, suppErr := generalserviceapiutils.GetAllSupplierMaster(params)
		if suppErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: suppErr.StatusCode,
				Err:        suppErr.Err,
			}
		}

		if len(suppliers) == 0 {
			internalFilter = append(internalFilter, utils.FilterCondition{
				ColumnField: "mtr_item_import.supplier_id",
				ColumnValue: "-1",
			})
		} else {
			var supplierIds []string
			for _, supplier := range suppliers {
				supplierIds = append(supplierIds, strconv.Itoa(supplier.SupplierId))
			}
			internalFilter = append(internalFilter, utils.FilterCondition{
				ColumnField: "mtr_item_import.supplier_id",
				ColumnValue: strings.Join(supplierIds, ","),
			})
		}
	}

	query := tx.Model(&model).
		Select("mtr_item_import.*, Item.item_code AS item_code, Item.item_name AS item_name").
		Joins("JOIN mtr_item Item ON mtr_item_import.item_id = Item.item_id")

	whereQuery := utils.ApplyFilter(query, internalFilter)

	err := whereQuery.Scopes(pagination.Paginate(&pages, whereQuery)).Scan(&responses).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if len(responses) == 0 {
		pages.Rows = []masteritempayloads.ItemImportResponse{}
		return pages, nil
	}

	supplierIds := []int{}
	for _, response := range responses {
		supplierIds = append(supplierIds, response.SupplierId)
	}

	supplierIds = utils.RemoveDuplicateIds(supplierIds)

	var supplierResponses []masteritempayloads.SupplierResponse
	if err := generalserviceapiutils.GetSupplierMasterByMultiId(supplierIds, &supplierResponses); err != nil {
		return pages, err
	}

	supplierMap := make(map[int]struct {
		SupplierName string
		SupplierCode string
	})
	for _, supplier := range supplierResponses {
		supplierMap[supplier.SupplierId] = struct {
			SupplierName string
			SupplierCode string
		}{
			SupplierName: supplier.SupplierName,
			SupplierCode: supplier.SupplierCode,
		}
	}

	type ItemImportWithSupplier struct {
		masteritempayloads.ItemImportResponse
		SupplierName string `json:"supplier_name"`
		SupplierCode string `json:"supplier_code"`
	}

	var combinedResponses []ItemImportWithSupplier

	for _, response := range responses {
		combined := ItemImportWithSupplier{
			ItemImportResponse: response,
		}

		if supplierInfo, ok := supplierMap[response.SupplierId]; ok {
			combined.SupplierName = supplierInfo.SupplierName
			combined.SupplierCode = supplierInfo.SupplierCode
		}

		combinedResponses = append(combinedResponses, combined)
	}

	pages.Rows = combinedResponses
	return pages, nil
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
	getSupplierbyIdUrl := config.EnvConfigs.GeneralServiceUrl + "supplier/" + strconv.Itoa(req.SupplierId)

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
