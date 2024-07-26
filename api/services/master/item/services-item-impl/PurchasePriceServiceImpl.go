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
	"encoding/json"
	"fmt"

	"github.com/labstack/gommon/log"
	"github.com/redis/go-redis/v9"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type PurchasePriceServiceImpl struct {
	PurchasePriceRepo masteritemrepository.PurchasePriceRepository
	DB                *gorm.DB
	RedisClient       *redis.Client // Redis client
}

func StartPurchasePriceService(PurchasePriceRepo masteritemrepository.PurchasePriceRepository, db *gorm.DB, redisClient *redis.Client) masteritemservice.PurchasePriceService {
	return &PurchasePriceServiceImpl{
		PurchasePriceRepo: PurchasePriceRepo,
		DB:                db,
		RedisClient:       redisClient,
	}
}

func (s *PurchasePriceServiceImpl) GetAllPurchasePrice(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, totalPages, totalRows, err := s.PurchasePriceRepo.GetAllPurchasePrice(tx, filterCondition, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, totalPages, totalRows, err
	}
	return results, totalPages, totalRows, nil
}

func (s *PurchasePriceServiceImpl) UpdatePurchasePrice(id int, req masteritempayloads.PurchasePriceRequest) (masteritementities.PurchasePrice, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.PurchasePriceRepo.UpdatePurchasePrice(tx, id, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return masteritementities.PurchasePrice{}, err
	}
	return results, nil
}

func (s *PurchasePriceServiceImpl) SavePurchasePrice(req masteritempayloads.PurchasePriceRequest) (masteritementities.PurchasePrice, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.PurchasePriceRepo.SavePurchasePrice(tx, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return masteritementities.PurchasePrice{}, err
	}
	return results, nil
}

func (s *PurchasePriceServiceImpl) GetPurchasePriceById(id int, pagination pagination.Pagination) (masteritempayloads.PurchasePriceResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.PurchasePriceRepo.GetPurchasePriceById(tx, id, pagination)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *PurchasePriceServiceImpl) AddPurchasePrice(req masteritempayloads.PurchasePriceDetailRequest) (masteritementities.PurchasePriceDetail, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.PurchasePriceRepo.AddPurchasePrice(tx, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return masteritementities.PurchasePriceDetail{}, err
	}
	return results, nil
}

func (s *PurchasePriceServiceImpl) GetAllPurchasePriceDetail(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, totalPages, totalRows, err := s.PurchasePriceRepo.GetAllPurchasePriceDetail(tx, filterCondition, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, totalPages, totalRows, err
	}
	return results, totalPages, totalRows, nil
}

func (s *PurchasePriceServiceImpl) GetPurchasePriceDetailById(id int) (masteritempayloads.PurchasePriceDetailResponses, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.PurchasePriceRepo.GetPurchasePriceDetailById(tx, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *PurchasePriceServiceImpl) UpdatePurchasePriceDetail(Id int, req masteritempayloads.PurchasePriceDetailRequest) (masteritementities.PurchasePriceDetail, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.PurchasePriceRepo.UpdatePurchasePriceDetail(tx, Id, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return masteritementities.PurchasePriceDetail{}, err
	}
	return results, nil
}

func (s *PurchasePriceServiceImpl) DeletePurchasePrice(id int, iddet []int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	deletemultiid, err := s.PurchasePriceRepo.DeletePurchasePrice(tx, id, iddet)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}
	return deletemultiid, nil
}

func (s *PurchasePriceServiceImpl) ChangeStatusPurchasePrice(Id int) (masteritementities.PurchasePrice, *exceptions.BaseErrorResponse) {

	tx := s.DB.Begin()

	// Ubah status
	entity, err := s.PurchasePriceRepo.ChangeStatusPurchasePrice(tx, Id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return masteritementities.PurchasePrice{}, err
	}
	return entity, nil
}

func (s *PurchasePriceServiceImpl) GenerateTemplateFile() (*excelize.File, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	// Generate template file
	f := excelize.NewFile()
	sheetName := "PurchasePrice"
	defer func() {
		if err := f.Close(); err != nil {
			// Handle file close error, if any
			log.Error(err)
		}
	}()

	// Create a new sheet.
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return f, &exceptions.BaseErrorResponse{Err: err, StatusCode: 500}
	}

	// Set value of a cell.
	f.SetCellValue(sheetName, "A1", "Item Code")
	f.SetCellValue(sheetName, "B1", "Item Name")
	f.SetCellValue(sheetName, "C1", "Purchase Price")
	f.SetColWidth(sheetName, "A", "C", 21.5)

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
	for col := 'A'; col <= 'C'; col++ {
		cell := string(col) + "1"
		f.SetCellStyle(sheetName, cell, cell, style)
	}

	// Fetch data for the template
	internalFilterCondition := []utils.FilterCondition{} // Adjust as needed
	paginate := pagination.Pagination{
		Limit: 3,
		Page:  1,
	}

	// Ensure GetAllPurchasePrice returns *exceptions.BaseErrorResponse
	results, _, _, errResp := s.PurchasePriceRepo.GetAllPurchasePrice(tx, internalFilterCondition, paginate)
	if errResp != nil {
		return f, errResp
	}

	data, err := ConvertPurchasePriceMapToStruct(results)
	if err != nil {
		return f, &exceptions.BaseErrorResponse{Err: err, StatusCode: 500}
	}

	for i, value := range data {
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", i+2), value.ItemCode)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", i+2), value.ItemName)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", i+2), value.PurchasePrice)
	}

	// Set active sheet of the workbook.
	f.SetActiveSheet(index)
	return f, nil
}

func ConvertPurchasePriceMapToStruct(maps []map[string]interface{}) ([]masteritempayloads.PurchasePriceByIdResponse, error) {
	var result []masteritempayloads.PurchasePriceByIdResponse

	// Marshal the maps into JSON
	jsonData, err := json.Marshal(maps)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON into the struct
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// func (s *PurchasePriceServiceImpl) UploadPreviewFile(rows [][]string) ([]masteritempayloads.PurchasePriceDetailRequest, *exceptions.BaseErrorResponse) {

// 	// Upload preview file
// 	var results []masteritempayloads.PurchasePriceDetailRequest
// 	for i, row := range rows {
// 		if i == 0 {
// 			continue
// 		}
// 		results = append(results, masteritempayloads.PurchasePriceDetailRequest{
// 			ItemCode:      row[0],
// 			ItemName:      row[1],
// 			PurchasePrice: row[2],
// 		})
// 	}
// 	return results, nil
// }

// func (s *PurchasePriceServiceImpl) ProcessDataUpload(req masteritempayloads.UploadRequest) (bool, *exceptions.BaseErrorResponse) {
// 	tx := s.DB.Begin()

// 	// Process data upload
// 	for _, value := range req.Data {
// 		_, err := s.PurchasePriceRepo.SavePurchasePrice(tx, value)
// 		if err != nil {
// 			return false, err
// 		}
// 	}
// 	return true, nil
// }
