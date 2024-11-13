package test

import (
	"after-sales/api/config"
	transactionsparepartrepositoryimpl "after-sales/api/repositories/transaction/sparepart/repositories-sparepart-impl"
	transactionsparepartservice "after-sales/api/services/transaction/sparepart"
	transactionsparepartserviceimpl "after-sales/api/services/transaction/sparepart/services-sparepart-impl"
	"github.com/stretchr/testify/assert"
	"testing"
)

func setupBinningList() transactionsparepartservice.BinningListService {
	config.InitEnvConfigs(true, "")
	db := config.InitDB()
	rdb := config.InitRedis()
	GoodsReceiveRepositoryImpl := transactionsparepartrepositoryimpl.NewbinningListRepositoryImpl()
	GoodsReceiveServiceImpl := transactionsparepartserviceimpl.NewBinningListServiceImpl(GoodsReceiveRepositoryImpl, db, rdb)
	return GoodsReceiveServiceImpl
}

func TestDeleteBinningList(t *testing.T) {
	service := setupBinningList()
	BinningId := 10
	//BinningId := 3279
	res, err := service.DeleteBinningList(BinningId)
	if err != nil {
		assert.Nil(t, err, "failed to get error is not nil")
	}
	assert.True(t, res)
}
func TestDeleteBinningListDetailMultiId(t *testing.T) {
	service := setupBinningList()
	multiIdBinningListDetail := "20967,20968"

	res, err := service.DeleteBinningListDetailMultiId(multiIdBinningListDetail)
	if err != nil {
		assert.Nil(t, err, "failed to get error is not nil")
	}
	assert.True(t, res)
}
