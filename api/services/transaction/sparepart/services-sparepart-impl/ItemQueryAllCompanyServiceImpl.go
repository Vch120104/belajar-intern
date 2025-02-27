package transactionsparepartserviceimpl

import (
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionsparepartrepository "after-sales/api/repositories/transaction/sparepart"
	transactionsparepartservice "after-sales/api/services/transaction/sparepart"
	"after-sales/api/utils"
	"fmt"
	"net/http"

	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
	excelize "github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type ItemQueryAllCompanyServiceImpl struct {
	ItemQueryAllCompanyRepository transactionsparepartrepository.ItemQueryAllCompanyRepository
	DB                            *gorm.DB
}

func NewItemQueryAllCompanyServiceImpl(
	itemQueryAllCompanyRepo transactionsparepartrepository.ItemQueryAllCompanyRepository,
	db *gorm.DB,
) transactionsparepartservice.ItemQueryAllCompanyService {
	return &ItemQueryAllCompanyServiceImpl{
		ItemQueryAllCompanyRepository: itemQueryAllCompanyRepo,
		DB:                            db,
	}
}

func (s *ItemQueryAllCompanyServiceImpl) GetAllItemQueryAllCompany(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        fmt.Errorf("panic recovered: %v", r),
			}
		} else if err != nil {
			tx.Rollback()
			logrus.Info("Transaction rollback due to error:", err)
		} else {
			if commitErr := tx.Commit().Error; commitErr != nil {
				logrus.WithError(commitErr).Error("Transaction commit failed")
				err = &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
				}
			}
		}
	}()

	response, responseErr := s.ItemQueryAllCompanyRepository.GetAllItemQueryAllCompany(tx, filterCondition, pages)
	if responseErr != nil {
		return response, responseErr
	}

	return response, nil
}

func (s *ItemQueryAllCompanyServiceImpl) GetItemQueryAllCompanyDownload(filterCondition []utils.FilterCondition) (*excelize.File, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        fmt.Errorf("panic recovered: %v", r),
			}
		} else if err != nil {
			tx.Rollback()
			logrus.Info("Transaction rollback due to error:", err)
		} else {
			if commitErr := tx.Commit().Error; commitErr != nil {
				logrus.WithError(commitErr).Error("Transaction commit failed")
				err = &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
				}
			}
		}
	}()

	f := excelize.NewFile()
	sheetName := "item_query_all_company"
	defer func() {
		f.DeleteSheet("Sheet1")
		if err := f.Close(); err != nil {
			log.Error(err)
		}
	}()

	// Create a new sheet.
	index, errExcel := f.NewSheet(sheetName)
	if errExcel != nil {
		return nil, &exceptions.BaseErrorResponse{
			Err:        errExcel,
			StatusCode: http.StatusInternalServerError,
		}
	}

	// Set value of a cell.
	f.SetCellValue(sheetName, "A1", "Company Code")
	f.SetCellValue(sheetName, "B1", "Company Name")
	f.SetCellValue(sheetName, "C1", "Item Code")
	f.SetCellValue(sheetName, "D1", "Item Name")
	f.SetCellValue(sheetName, "E1", "Moving Code")
	f.SetCellValue(sheetName, "F1", "Quantity On Hand")
	f.SetColWidth(sheetName, "A", "F", 21.5)

	// Create a style with bold font and border
	style, errExcel := f.NewStyle(&excelize.Style{
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
	if errExcel != nil {
		return nil, &exceptions.BaseErrorResponse{
			Err:        errExcel,
			StatusCode: http.StatusInternalServerError,
		}
	}

	// Apply the style to the header cells.
	for col := 'A'; col <= 'F'; col++ {
		cell := string(col) + "1"
		f.SetCellStyle(sheetName, cell, cell, style)
	}

	// Fetch data.
	response, err := s.ItemQueryAllCompanyRepository.GetItemQueryAllCompanyDownload(tx, filterCondition)
	if err != nil {
		return nil, err
	}

	// Assign data to cell
	for i, value := range response {
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", i+2), value.CompanyCode)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", i+2), value.CompanyName)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", i+2), value.ItemCode)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", i+2), value.ItemName)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", i+2), value.MovingCode)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", i+2), value.QuantityOnHand)
	}

	// Set active sheet of the workbook.
	f.SetActiveSheet(index)
	return f, nil
}
