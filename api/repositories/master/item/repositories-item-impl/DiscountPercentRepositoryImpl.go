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

type DiscountPercentRepositoryImpl struct {
	myDB *gorm.DB
}

func StartDiscountPercentRepositoryImpl(db *gorm.DB) masteritemrepository.DiscountPercentRepository {
	return &DiscountPercentRepositoryImpl{myDB: db}
}

func (r *DiscountPercentRepositoryImpl) WithTrx(trxHandle *gorm.DB) masteritemrepository.DiscountPercentRepository {
	if trxHandle == nil {
		log.Println("Transaction Database Not Found!")
		return r
	}
	r.myDB = trxHandle
	return r
}

func (r *DiscountPercentRepositoryImpl) GetAllDiscountPercent(filterCondition []utils.FilterCondition) ([]map[string]interface{}, error) {
	var responses []masteritempayloads.DiscountPercentResponse
	var getOrderTypeResponse []masteritempayloads.OrderTypeResponse
	var c *gin.Context
	var internalServiceFilter, externalServiceFilter []utils.FilterCondition
	var orderTypeName string
	responseStruct := reflect.TypeOf(masteritempayloads.DiscountPercentResponse{})

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
	tableStruct := masteritempayloads.DiscountPercentListResponse{}
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

	errUrlDiscountPercent := utils.Get(c, orderTypeUrl, &getOrderTypeResponse, nil)

	if errUrlDiscountPercent != nil {
		return nil, errUrlDiscountPercent
	}

	joinedData := utils.DataFrameInnerJoin(responses, getOrderTypeResponse, "OrderTypeId")

	return joinedData, nil
}

func (r *DiscountPercentRepositoryImpl) GetDiscountPercentById(Id int) (masteritempayloads.DiscountPercentResponse, error) {
	entities := masteritementities.DiscountPercent{}
	response := masteritempayloads.DiscountPercentResponse{}

	rows, err := r.myDB.Model(&entities).
		Where(masteritementities.DiscountPercent{
			DiscountPercentId: Id,
		}).
		First(&response).
		Rows()

	if err != nil {
		return response, err
	}

	defer rows.Close()

	return response, nil
}

func (r *DiscountPercentRepositoryImpl) SaveDiscountPercent(request masteritempayloads.DiscountPercentResponse) (bool, error) {
	entities := masteritementities.DiscountPercent{
		IsActive:          request.IsActive,
		DiscountPercentId: request.DiscountPercentId,
		DiscountCodeId:    request.DiscountCodeId,
		OrderTypeId:       request.OrderTypeId,
		Discount:          request.Discount,
	}

	err := r.myDB.Save(&entities).Error

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *DiscountPercentRepositoryImpl) ChangeStatusDiscountPercent(Id int) (bool, error) {
	var entities masteritementities.DiscountPercent

	result := r.myDB.Model(&entities).
		Where("discount_percent_id = ?", Id).
		First(&entities)

	if result.Error != nil {
		return false, result.Error
	}

	if entities.IsActive {
		entities.IsActive = false
	} else {
		entities.IsActive = true
	}

	result = r.myDB.Save(&entities)

	if result.Error != nil {
		return false, result.Error
	}

	return true, nil
}