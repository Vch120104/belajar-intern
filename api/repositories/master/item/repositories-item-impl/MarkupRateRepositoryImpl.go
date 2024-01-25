package masteritemrepositoryimpl

import (
	masteritementities "after-sales/api/entities/master/item"
	masteritempayloads "after-sales/api/payloads/master/item"
	masteritemrepository "after-sales/api/repositories/master/item"
	"after-sales/api/utils"
	"log"
	"reflect"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MarkupRateRepositoryImpl struct {
	myDB *gorm.DB
}

func StartMarkupRateRepositoryImpl(db *gorm.DB) masteritemrepository.MarkupRateRepository {
	return &MarkupRateRepositoryImpl{myDB: db}
}

func (r *MarkupRateRepositoryImpl) WithTrx(trxHandle *gorm.DB) masteritemrepository.MarkupRateRepository {
	if trxHandle == nil {
		log.Println("Transaction Database Not Found!")
		return r
	}
	r.myDB = trxHandle
	return r
}

func (r *MarkupRateRepositoryImpl) GetAllMarkupRate(filterCondition []utils.FilterCondition) ([]map[string]interface{}, error) {
	var responses []masteritempayloads.MarkupRateListResponse
	var getOrderTypeResponse []masteritempayloads.OrderTypeResponse
	var c *gin.Context
	var internalServiceFilter, externalServiceFilter []utils.FilterCondition
	var orderTypeName string
	responseStruct := reflect.TypeOf(masteritempayloads.MarkupRateResponse{})

	for i := 0; i < len(filterCondition); i++ {
		flag := false
		for j := 0; j < responseStruct.NumField(); j++ {
			if filterCondition[i].ColumnField == responseStruct.Field(j).Tag.Get("parent_entity")+"."+responseStruct.Field(j).Tag.Get("json") {
				internalServiceFilter = append(internalServiceFilter, filterCondition[i])
				flag = true
				break
			}
			if !flag {
				externalServiceFilter = append(externalServiceFilter, filterCondition[i])
			}
		}
	}

	//apply external services filter
	for i := 0; i < len(externalServiceFilter); i++ {
		orderTypeName = externalServiceFilter[i].ColumnValue
	}

	// define table struct
	tableStruct := masteritempayloads.MarkupRateListResponse{}
	//define join table
	joinTable := utils.CreateJoinSelectStatement(r.myDB, tableStruct)
	//apply filter
	whereQuery := utils.ApplyFilter(joinTable, internalServiceFilter)
	//apply pagination and execute
	rows, err := whereQuery.Scan(&responses).Rows()

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if len(responses) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	orderTypeUrl := "http://10.1.32.26:8000/general-service/api/general/order-type-filter?order_type_name=" + orderTypeName

	errUrlMarkupRate := utils.Get(c, orderTypeUrl, &getOrderTypeResponse, nil)

	if errUrlMarkupRate != nil {
		return nil, errUrlMarkupRate
	}

	joinedData := utils.DataFrameInnerJoin(responses, getOrderTypeResponse, "OrderTypeId")

	return joinedData, nil
}

func (r *MarkupRateRepositoryImpl) GetMarkupRateById(Id int) (masteritempayloads.MarkupRateResponse, error) {
	entities := masteritementities.MarkupRate{}
	response := masteritempayloads.MarkupRateResponse{}

	rows, err := r.myDB.Model(&entities).
		Where(masteritementities.MarkupRate{
			MarkupRateId: Id,
		}).
		First(&response).
		Rows()

	if err != nil {
		return response, err
	}

	defer rows.Close()

	return response, nil
}

func (r *MarkupRateRepositoryImpl) SaveMarkupRate(request masteritempayloads.MarkupRateRequest) (bool, error) {
	entities := masteritementities.MarkupRate{
		IsActive:       true,
		MarkupRateId:   request.MarkupRateId,
		MarkupMasterId: request.MarkupMasterId,
		OrderTypeId:    request.OrderTypeId,
		MarkupRate:     request.MarkupRate,
	}

	err := r.myDB.Save(&entities).Error

	if err != nil {
		return false, err
	}

	return true, nil
}

