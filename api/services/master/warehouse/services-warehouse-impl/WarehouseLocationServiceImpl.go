package masterwarehouseserviceimpl

import (
	// masterwarehousepayloads "after-sales/api/payloads/master/warehouse"

	masterwarehouseentities "after-sales/api/entities/master/warehouse"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	pagination "after-sales/api/payloads/pagination"
	masterwarehouserepository "after-sales/api/repositories/master/warehouse"
	masterwarehouseservice "after-sales/api/services/master/warehouse"
	"after-sales/api/utils"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
	// "log"
	// "after-sales/api/utils"
)

type WarehouseLocationServiceImpl struct {
	warehouseLocationRepo  masterwarehouserepository.WarehouseLocationRepository
	warehouseMasterService masterwarehouseservice.WarehouseMasterService
	DB                     *gorm.DB
	RedisClient            *redis.Client // Redis client
}

func OpenWarehouseLocationService(warehouseLocation masterwarehouserepository.WarehouseLocationRepository, warehouseMasterService masterwarehouseservice.WarehouseMasterService, db *gorm.DB, redisClient *redis.Client) masterwarehouseservice.WarehouseLocationService {
	return &WarehouseLocationServiceImpl{
		warehouseLocationRepo:  warehouseLocation,
		warehouseMasterService: warehouseMasterService,
		DB:                     db,
		RedisClient:            redisClient,
	}
}

// ProcessWarehouseLocationTemplate implements masterwarehouseservice.WarehouseLocationService.
func (s *WarehouseLocationServiceImpl) ProcessWarehouseLocationTemplate(req masterwarehousepayloads.ProcessWarehouseLocationTemplate, companyId int) (bool, *exceptions.BaseErrorResponse) {

	for _, value := range req.Data {
		fmt.Println(value.Validation)
		if value.Validation != "Ok" {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: 400,
				Err:        errors.New("failed to process data, check validation"),
			}
		}
	}

	for _, value := range req.Data {
		//GET GROUP ID AND WAREHOUSE ID

		groupId, warehouseId, _ := s.warehouseMasterService.GetWarehouseGroupAndMasterbyCodeandCompanyId(companyId, value.WarehouseCode)
		entities := masterwarehouseentities.WarehouseLocation{
			WarehouseId:                   warehouseId,
			WarehouseGroupId:              groupId,
			WarehouseLocationCode:         value.WarehouseLocationCode,
			WarehouseLocationName:         value.WarehouseLocationName,
			WarehouseLocationDetailName:   "", // FROM PROCEDURE DOCUMENTATION WHICH INSERT WITH DEFAULT VALUE
			WarehouseLocationPickSequence: 0,  // FROM PROCEDURE DOCUMENTATION WHICH INSERT WITH DEFAULT VALUE
			WarehouseLocationCapacityInM3: 0,  // FROM PROCEDURE DOCUMENTATION WHICH INSERT WITH DEFAULT VALUE
		}
		_, err := s.Save(entities)

		if err != nil {
			return false, err
		}

	}

	return true, nil
}

// UploadPreviewFile implements masterwarehouseservice.WarehouseLocationService.
func (s *WarehouseLocationServiceImpl) UploadPreviewFile(rows [][]string, companyId int) ([]masterwarehousepayloads.GetWarehouseLocationPreviewResponse, *exceptions.BaseErrorResponse) {
	response := []masterwarehousepayloads.GetWarehouseLocationPreviewResponse{}

	for index, value := range rows {
		data := masterwarehousepayloads.GetWarehouseLocationPreviewResponse{}

		if index > 0 {
			data.WarehouseCode = value[0]
			data.WarehouseLocationCode = value[1]
			data.WarehouseLocationName = value[2]

			//Check warehouseCode exist
			if warehouseCodeExist := s.warehouseMasterService.IsWarehouseMasterByCodeAndCompanyIdExist(companyId, value[0]); warehouseCodeExist {
				tx := s.DB.Begin()

				//Check warehouslocation exist
				isExist, err := s.warehouseLocationRepo.CheckIfLocationExist(tx, value[0], value[1], value[2])

				defer helper.CommitOrRollback(tx, err)

				if isExist {
					data.Validation = "Location is already exist"

				} else {

					data.Validation = "Ok"
				}
			} else {
				data.Validation = "Warehouse Code is invalid"
			}

			fmt.Print(data.Validation)

			response = append(response, data)
		} else {
			if value[0] != "WAREHOUSE_CODE" || value[1] != "LOCATION_CODE" || value[2] != "LOCATION_NAME" && len(value) == 3 {
				return response, &exceptions.BaseErrorResponse{Err: errors.New("make sure header is correct"), StatusCode: 400}
			}
		}
	}

	return response, nil
}

// GenerateTemplateFile implements masterwarehouseservice.WarehouseLocationService.
func (s *WarehouseLocationServiceImpl) GenerateTemplateFile() (*excelize.File, *exceptions.BaseErrorResponse) {
	//Initiate db for get data example
	tx := s.DB.Begin()

	f := excelize.NewFile()
	sheetName := "WarehouseLocation"
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
	f.SetCellValue(sheetName, "A1", "WAREHOUSE_CODE")
	f.SetCellValue(sheetName, "B1", "LOCATION_CODE")
	f.SetCellValue(sheetName, "C1", "LOCATION_NAME")
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

	//Get Data Example

	filter := map[string]string{}

	paginate := pagination.Pagination{
		Limit: 3,
		Page:  0,
	}

	filtercondition := utils.BuildFilterCondition(filter)

	result, errGetData := s.warehouseLocationRepo.GetAll(tx, filtercondition, paginate)
	defer helper.CommitOrRollback(tx, errGetData)

	data := result.Rows

	if sliceRows, ok := data.([]masterwarehousepayloads.GetAllWarehouseLocationResponse); ok {
		fmt.Print(data)
		for i := 0; i < len(sliceRows); i++ {
			f.SetCellValue(sheetName, fmt.Sprintf("A%d", i+2), sliceRows[i].WarehouseCode)
			f.SetCellValue(sheetName, fmt.Sprintf("B%d", i+2), sliceRows[i].WarehouseLocationCode)
			f.SetCellValue(sheetName, fmt.Sprintf("C%d", i+2), sliceRows[i].WarehouseLocationName)

		}
	}

	// Set active sheet of the workbook.
	f.SetActiveSheet(index)
	return f, nil

}

func (s *WarehouseLocationServiceImpl) Save(request masterwarehouseentities.WarehouseLocation) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	if request.WarehouseLocationId != 0 {
		_, err := s.warehouseLocationRepo.GetById(tx, request.WarehouseLocationId)

		if err != nil {
			return false, err
		}
	}

	save, err := s.warehouseLocationRepo.Save(tx, request)
	defer helper.CommitOrRollback(tx, err)

	if err != nil {
		return false, err
	}
	return save, err
}

func (s *WarehouseLocationServiceImpl) GetById(warehouseLocationId int) (masterwarehousepayloads.GetAllWarehouseLocationResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	get, err := s.warehouseLocationRepo.GetById(tx, warehouseLocationId)
	defer helper.CommitOrRollback(tx, err)

	if err != nil {
		return get, err
	}
	return get, nil
}

func (s *WarehouseLocationServiceImpl) GetAll(filter []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	get, err := s.warehouseLocationRepo.GetAll(tx, filter, pages)
	defer helper.CommitOrRollback(tx, err)

	if err != nil {
		return get, err
	}
	return get, nil
}

func (s *WarehouseLocationServiceImpl) ChangeStatus(warehouseLocationId int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	_, err := s.warehouseLocationRepo.GetById(tx, warehouseLocationId)

	if err != nil {
		return false, err
	}

	change_status, err := s.warehouseLocationRepo.ChangeStatus(tx, warehouseLocationId)
	defer helper.CommitOrRollback(tx, err)

	if err != nil {
		return change_status, err
	}
	return change_status, nil
}
