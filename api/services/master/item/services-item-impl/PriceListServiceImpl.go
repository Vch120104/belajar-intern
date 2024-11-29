package masteritemserviceimpl

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	masteritemservice "after-sales/api/services/master/item"
	"after-sales/api/utils"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type PriceListServiceImpl struct {
	priceListRepo masteritemrepository.PriceListRepository
	itemService   masteritemservice.ItemService
	DB            *gorm.DB
	RedisClient   *redis.Client // Redis client
}

func StartPriceListService(priceListRepo masteritemrepository.PriceListRepository, itemService masteritemservice.ItemService, db *gorm.DB, redisClient *redis.Client) masteritemservice.PriceListService {
	return &PriceListServiceImpl{
		priceListRepo: priceListRepo,
		itemService:   itemService,
		DB:            db,
		RedisClient:   redisClient,
	}
}

// Duplicate implements masteritemservice.PriceListService.
func (s *PriceListServiceImpl) Duplicate(itemGroupId int, brandId int, currencyId int, date string) ([]masteritempayloads.PriceListItemResponses, *exceptions.BaseErrorResponse) {
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
	results, err := s.priceListRepo.Duplicate(tx, itemGroupId, brandId, currencyId, date)

	if err != nil {
		return results, err
	}
	return results, nil
}

// Download implements masteritemservice.PriceListService.
func (s *PriceListServiceImpl) Download(uploadRequest masteritempayloads.PriceListUploadDataRequest) (*excelize.File, *exceptions.BaseErrorResponse) {
	//Initiate db for get data example

	f := excelize.NewFile()
	sheetName := "PriceList"
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
	// Set HEADER
	f.SetCellValue(sheetName, "A1", "VEHICLE_BRAND")
	f.SetCellValue(sheetName, "B1", "CURRENCY_CODE")
	f.SetCellValue(sheetName, "C1", "EFFECTIVE_DATE")
	f.SetCellValue(sheetName, "D1", "ITEM_CODE")
	f.SetCellValue(sheetName, "E1", "PRICE_AMOUNT")
	f.SetColWidth(sheetName, "A", "E", 21.5)

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
	for col := 'A'; col <= 'E'; col++ {
		cell := string(col) + "1"
		f.SetCellStyle(sheetName, cell, cell, style)
	}

	//Get Data

	page := 0

	for {

		paginate := pagination.Pagination{
			Limit: 10,
			Page:  page,
		}
		result, _ := s.CheckPriceListItem(uploadRequest.ItemGroupId, uploadRequest.BrandId, uploadRequest.CurrencyId, uploadRequest.Date, paginate)

		data := result.Rows

		parsedDate, err := time.Parse("2006-01-02", uploadRequest.Date)
		if err != nil {
			return f, &exceptions.BaseErrorResponse{Err: err, StatusCode: 500}
		}
		formattedDate := parsedDate.Format("02/01/2006")

		if data, ok := data.([]masteritempayloads.PriceListItemResponses); ok {
			if len(data) == 0 {
				break
			}
			for i := 0; i < len(data); i++ {
				f.SetCellValue(sheetName, fmt.Sprintf("A%d", i+2), uploadRequest.BrandCode)
				f.SetCellValue(sheetName, fmt.Sprintf("B%d", i+2), uploadRequest.CurrencyCode)
				f.SetCellValue(sheetName, fmt.Sprintf("C%d", i+2), formattedDate)
				f.SetCellValue(sheetName, fmt.Sprintf("D%d", i+2), data[i].ItemCode)
				f.SetCellValue(sheetName, fmt.Sprintf("E%d", i+2), data[i].PriceListAmount)
			}
		}

		page++

	}

	// Set active sheet of the workbook.
	f.SetActiveSheet(index)
	return f, nil
}

// CheckPriceListItem implements masteritemservice.PriceListService.
func (s *PriceListServiceImpl) CheckPriceListItem(itemGroupId int, brandId int, currencyId int, date string, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
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
	results, err := s.priceListRepo.CheckPriceListItem(tx, itemGroupId, brandId, currencyId, date, pages)

	if err != nil {
		return results, err
	}
	return results, nil
}

// ProcessUploadFile implements masteritemservice.PriceListService.
func (s *PriceListServiceImpl) ProcessUploadFile(upload []masteritempayloads.PriceListProcessdDataRequest) (bool, *exceptions.BaseErrorResponse) {
	// Date string in "DD/MM/YYYY" format
	dateString := "14/05/2024"

	// Define the layout according to the expected format
	layout := "02/01/2006"

	// Parse the date string using the defined layout
	dateParse, err := time.Parse(layout, dateString)
	if err != nil {
		fmt.Println("Error parsing date:", err)

	}

	fmt.Println("Parsed date:", dateParse)
	for _, value := range upload {
		dateParse, _ := time.Parse("02/01/2006", value.Date)

		req := masteritempayloads.SavePriceListMultiple{
			BrandId:         value.BrandId,
			CurrencyId:      value.CurrencyId,
			ItemGroupId:     value.ItemGroupId,
			EffectiveDate:   dateParse,
			CompanyId:       value.CompanyId,
			PriceListCodeId: value.PriceListCodeId,

			Detail: []masteritempayloads.PriceListItemResponses{
				{ItemId: value.ItemId, ItemClassId: value.ItemClassId, PriceListAmount: value.PriceListAmount},
			},
		}

		_, err := s.SavePriceList(req)

		if err != nil {
			return false, err
		}

	}

	return true, nil

}

// UploadFile implements masteritemservice.PriceListService.
func (s *PriceListServiceImpl) UploadFile(rows [][]string, uploadRequest masteritempayloads.PriceListUploadDataRequest) ([]string, *exceptions.BaseErrorResponse) {

	result := []string{}
	req := []masteritempayloads.PriceListProcessdDataRequest{}

	for key := 1; key < len(rows); key++ {
		value := rows[key]
		if key != 0 {
			if !strings.EqualFold(value[0], uploadRequest.BrandCode) {
				result = append(result, fmt.Sprintf("Line %d : %s", key, "Vehicle Brand not match"))
			}
			if !strings.EqualFold(value[1], uploadRequest.CurrencyCode) {
				result = append(result, fmt.Sprintf("Line %d : %s", key, "Currency Code not match"))
			}

			//parse from excel
			parsedDate1, err := time.Parse("01-02-06", value[2])
			if err != nil {
				parsedDate1, err = time.Parse("02/01/2006", value[2])
				if err != nil {

					return result, &exceptions.BaseErrorResponse{StatusCode: 400, Err: err}
				}
			}
			//parse from json body
			parsedDate2, err := time.Parse("02/01/2006", uploadRequest.Date)
			if err != nil {

				return result, &exceptions.BaseErrorResponse{StatusCode: 400, Err: err}
			}

			if !(parsedDate1.Year() == parsedDate2.Year() && parsedDate1.Month() == parsedDate2.Month() && parsedDate1.Day() == parsedDate2.Day()) {
				result = append(result, fmt.Sprintf("Line %d : %s", key, "Effective Date not match"))
			}

			commonPriceList := false

			if uploadRequest.CompanyCode == "0" {
				commonPriceList = true
			}

			isItemCodeExist, itemId, itemClassId, _ := s.itemService.CheckItemCodeExist(value[3], uploadRequest.ItemGroupId, commonPriceList, uploadRequest.BrandId)
			if !isItemCodeExist {

				result = append(result, fmt.Sprintf("Line %d : %s", key, "Item Code not match"))
			} else {
				tx := s.DB.Begin()
				splitDate := strings.Split(uploadRequest.Date, "/")
				convertDate := splitDate[2] + "-" + splitDate[1] + "-" + splitDate[0]
				isPriceListExist, err := s.priceListRepo.CheckPriceListExist(tx, itemId, uploadRequest.BrandId, uploadRequest.CurrencyId, convertDate, uploadRequest.CompanyId)
				defer helper.CommitOrRollback(tx, err)
				if isPriceListExist {

					result = append(result, fmt.Sprintf("Line %d : %s", key, "Price List already exists"))
				} else {
					cleanValue := strings.ReplaceAll(value[4], ",", "")
					price, err := strconv.ParseFloat(cleanValue, 64)
					if err != nil {
						result = append(result, fmt.Sprintf("Line %d : %s", key, "Error read file"))
						return result, &exceptions.BaseErrorResponse{StatusCode: 500, Err: err}
					}

					model := masteritempayloads.PriceListProcessdDataRequest{
						BrandId:         uploadRequest.BrandId,
						ItemGroupId:     uploadRequest.ItemGroupId,
						CurrencyId:      uploadRequest.CurrencyId,
						Date:            uploadRequest.Date,
						PriceListCodeId: uploadRequest.PriceListCodeId,
						CompanyId:       uploadRequest.CompanyId,
						ItemId:          itemId,
						ItemClassId:     itemClassId,
						PriceListAmount: price,
					}
					req = append(req, model)

				}
			}
		}
	}

	//if no error massage inside result, then save all data
	fmt.Println(len(result))
	if len(result) == 0 {
		success, _ := s.ProcessUploadFile(req)
		if success {
			result = append(result, "All line validation completed..", "Finish save price list")
		}
	}

	return result, nil
}

// GenerateTemplateFile implements masteritemservice.PriceListService.
func (s *PriceListServiceImpl) GenerateDownloadTemplateFile() (*excelize.File, *exceptions.BaseErrorResponse) {
	//Initiate db for get data example

	f := excelize.NewFile()
	sheetName := "PriceList"
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

	// Set HEADER
	f.SetCellValue(sheetName, "A1", "VEHICLE_BRAND")
	f.SetCellValue(sheetName, "B1", "CURRENCY_CODE")
	f.SetCellValue(sheetName, "C1", "EFFECTIVE_DATE")

	//EXAMPLE DATE FORMAT
	styleH1, _ := f.NewStyle(&excelize.Style{Font: &excelize.Font{
		Color: "#FF0000",
	}})
	f.SetCellValue(sheetName, "H1", "*DATE FORMAT = 02/01/2006")
	f.SetCellStyle(sheetName, "H1", "H1", styleH1)

	f.SetCellValue(sheetName, "D1", "ITEM_CODE")
	f.SetCellValue(sheetName, "E1", "PRICE_AMOUNT")
	f.SetColWidth(sheetName, "A", "E", 21.5)

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
	for col := 'A'; col <= 'E'; col++ {
		cell := string(col) + "1"
		f.SetCellStyle(sheetName, cell, cell, style)
	}

	// Set active sheet of the workbook.
	f.SetActiveSheet(index)
	return f, nil
}

func (s *PriceListServiceImpl) GetPriceList(request masteritempayloads.PriceListGetAllRequest) ([]masteritempayloads.PriceListResponse, *exceptions.BaseErrorResponse) {
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
	results, err := s.priceListRepo.GetPriceList(tx, request)

	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *PriceListServiceImpl) GetPriceListById(Id int) (masteritempayloads.PriceListGetbyId, *exceptions.BaseErrorResponse) {
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
	results, err := s.priceListRepo.GetPriceListById(tx, Id)

	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *PriceListServiceImpl) SavePriceList(request masteritempayloads.SavePriceListMultiple) (int, *exceptions.BaseErrorResponse) {
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

	id, err := s.priceListRepo.SavePriceList(tx, request)

	if err != nil {
		return id, err
	}
	return id, nil
}

func (s *PriceListServiceImpl) ChangeStatusPriceList(Id int) (bool, *exceptions.BaseErrorResponse) {
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

	_, err = s.priceListRepo.GetPriceListById(tx, Id)

	if err != nil {
		return false, err
	}

	result, err := s.priceListRepo.ChangeStatusPriceList(tx, Id)

	if err != nil {
		return false, err
	}
	return result, nil
}

func (s *PriceListServiceImpl) GetAllPriceListNew(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
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
	result, total_page, total_rows, err := s.priceListRepo.GetAllPriceListNew(tx, filterCondition, pages)

	if err != nil {
		return nil, 0, 0, err
	}

	return result, total_page, total_rows, nil
}

func (s *PriceListServiceImpl) DeactivatePriceList(id string) (bool, *exceptions.BaseErrorResponse) {
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
	result, err := s.priceListRepo.DeactivatePriceList(tx, id)

	if err != nil {
		return false, err
	}

	return result, nil
}

func (s *PriceListServiceImpl) ActivatePriceList(id string) (bool, *exceptions.BaseErrorResponse) {
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
	result, err := s.priceListRepo.ActivatePriceList(tx, id)

	if err != nil {
		return false, err
	}

	return result, nil
}

func (s *PriceListServiceImpl) DeletePriceList(id string) (bool, *exceptions.BaseErrorResponse) {
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
	result, err := s.priceListRepo.DeletePriceList(tx, id)

	if err != nil {
		return false, err
	}

	return result, nil
}
