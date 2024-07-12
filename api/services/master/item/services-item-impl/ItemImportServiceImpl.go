package masteritemserviceimpl

import (
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	masteritemservice "after-sales/api/services/master/item"
	"after-sales/api/utils"
	"errors"
	"fmt"
	"strconv"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type ItemImportServiceImpl struct {
	itemImportRepo masteritemrepository.ItemImportRepository
	DB             *gorm.DB
}

// ProcessDataUpload implements masteritemservice.ItemImportService.
func (s *ItemImportServiceImpl) ProcessDataUpload(req masteritempayloads.ItemImportUploadRequest) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	saveSuccess, err := s.itemImportRepo.SaveItemImport(tx, req)
	defer helper.CommitOrRollback(tx, err)

	if err != nil {
		return saveSuccess, err
	}

	return true, nil
}

// UploadPreviewFile implements masteritemservice.ItemImportService.
func (s *ItemImportServiceImpl) UploadPreviewFile(rows [][]string) ([]masteritempayloads.ItemImportUploadResponse, *exceptions.BaseErrorResponse) {
	response := []masteritempayloads.ItemImportUploadResponse{}

	for index, value := range rows {
		data := masteritempayloads.ItemImportUploadResponse{}
		var failedQtyParse error
		var failedOrderaParse error
		if index > 0 {
			data.ItemCode = value[0]
			data.SupplierCode = value[1]
			data.ItemAliasCode = value[2]
			data.ItemAliasName = value[3]
			data.OrderQtyMultiplier, failedQtyParse = strconv.ParseFloat(value[4], 64)

			if failedQtyParse != nil {
				return response, &exceptions.BaseErrorResponse{Err: errors.New("make sure moq value is correct"), StatusCode: 400}
			}

			data.RoyaltyFlag = value[5]
			data.OrderConversion, failedOrderaParse = strconv.ParseFloat(value[6], 64)

			if failedOrderaParse != nil {
				return response, &exceptions.BaseErrorResponse{Err: errors.New("make sure order conversion value is correct"), StatusCode: 400}

			}

			response = append(response, data)
		} else {
			if value[0] != "Part_Number" || value[1] != "Supplier_Code" || value[2] != "Part_Number_Alias" || value[3] != "Part_Name_Alias" || value[4] != "MOQ" || value[5] != "Royalty" || value[6] != "Order_Conversion" && len(value) == 7 {
				return response, &exceptions.BaseErrorResponse{Err: errors.New("make sure header is correct"), StatusCode: 400}
			}
		}
	}

	return response, nil

}

// GenerateTemplateFile implements masteritemservice.ItemImportService.
func (s *ItemImportServiceImpl) GenerateTemplateFile() (*excelize.File, *exceptions.BaseErrorResponse) {
	//Initiate db for get data example
	tx := s.DB.Begin()

	f := excelize.NewFile()
	sheetName := "ItemImportMaster"
	defer func() {
		if err := f.Close(); err != nil {
			return
		}
	}()
	// Create a new sheet.
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return f, &exceptions.BaseErrorResponse{Err: err, StatusCode: 500}
	}
	// Set value of a cell.USPG_GMITEM2_INSERT
	f.SetCellValue(sheetName, "A1", "Part_Number")
	f.SetCellValue(sheetName, "B1", "Supplier_Code")
	f.SetCellValue(sheetName, "C1", "Part_Number_Alias")
	f.SetCellValue(sheetName, "D1", "Part_Name_Alias")
	f.SetCellValue(sheetName, "E1", "MOQ")
	f.SetCellValue(sheetName, "F1", "Royalty")
	f.SetCellValue(sheetName, "G1", "Order_Conversion")
	f.SetColWidth(sheetName, "A", "G", 21.5)

	// Create a style with bold font and border
	style, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "left"},
		Font: &excelize.Font{
			Bold: true,
		},
		Border: []excelize.Border{
			{
				Type:  "left",
				Color: "000000",
				Style: 1,
			},
			{
				Type:  "top",
				Color: "000000",
				Style: 1,
			},
			{
				Type:  "bottom",
				Color: "000000",
				Style: 1,
			},
			{
				Type:  "right",
				Color: "000000",
				Style: 1,
			},
		},
	})
	if err != nil {
		return f, &exceptions.BaseErrorResponse{Err: err, StatusCode: 500}
	}

	// Apply the style to the header cells
	for col := 'A'; col <= 'G'; col++ {
		cell := string(col) + "1"
		f.SetCellStyle(sheetName, cell, cell, style)
	}

	// Get data example

	id := []int{}

	internalFilterCondition := map[string]string{}
	externalFilterCondition := map[string]string{}

	paginate := pagination.Pagination{
		Limit: 3,
		Page:  0,
	}

	internalCriteria := utils.BuildFilterCondition(internalFilterCondition)
	externalCriteria := utils.BuildFilterCondition(externalFilterCondition)

	paginatedData, _, _, errorgetalldata := s.itemImportRepo.GetAllItemImport(tx, internalCriteria, externalCriteria, paginate)
	defer helper.CommitOrRollback(tx, errorgetalldata)

	data, _ := masteritempayloads.ConvertItemImportMapToStruct(paginatedData)

	for _, value := range data {
		id = append(id, value.ItemImportId)
	}

	for i := 0; i < len(id); i++ {

		result, _ := s.itemImportRepo.GetItemImportbyId(tx, id[i])

		f.SetCellValue(sheetName, fmt.Sprintf("A%d", i+2), result.ItemCode)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", i+2), result.SupplierCode)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", i+2), result.ItemAliasCode)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", i+2), result.ItemAliasName)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", i+2), result.OrderQtyMultiplier)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", i+2), result.RoyaltyFlag)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", i+2), result.OrderConversion)

	}

	// Set active sheet of the workbook.
	f.SetActiveSheet(index)
	return f, nil

}

// GetItemImportbyItemIdandSupplierId implements masteritemservice.ItemImportService.
func (s *ItemImportServiceImpl) GetItemImportbyItemIdandSupplierId(itemId int, supplierId int) (masteritempayloads.ItemImportByIdResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.itemImportRepo.GetItemImportbyItemIdandSupplierId(tx, itemId, supplierId)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

// GetItemImportbyId implements masteritemservice.ItemImportService.
func (s *ItemImportServiceImpl) GetItemImportbyId(Id int) (masteritempayloads.ItemImportByIdResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.itemImportRepo.GetItemImportbyId(tx, Id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

// GetAllItemImport implements masteritemservice.ItemImportService.
func (s *ItemImportServiceImpl) GetAllItemImport(internalFilter []utils.FilterCondition, externalFilter []utils.FilterCondition, pages pagination.Pagination) ([]map[string]any, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, totalPages, totalRows, err := s.itemImportRepo.GetAllItemImport(tx, internalFilter, externalFilter, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, totalPages, totalRows, err
	}
	return results, totalPages, totalRows, nil
}

// SaveItemImport implements masteritemservice.ItemImportService.
func (s *ItemImportServiceImpl) SaveItemImport(req masteritempayloads.ItemImportUploadRequest) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.itemImportRepo.SaveItemImport(tx, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

// UpdateItemImport implements masteritemservice.ItemImportService.
func (s *ItemImportServiceImpl) UpdateItemImport(req masteritementities.ItemImport) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.itemImportRepo.UpdateItemImport(tx, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

func StartItemImportService(ItemImportrepo masteritemrepository.ItemImportRepository, db *gorm.DB) masteritemservice.ItemImportService {
	return &ItemImportServiceImpl{
		itemImportRepo: ItemImportrepo,
		DB:             db,
	}
}
