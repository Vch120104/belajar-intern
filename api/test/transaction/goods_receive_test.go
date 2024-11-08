package test

import (
	"after-sales/api/config"
	transactionsparepartrepositoryimpl "after-sales/api/repositories/transaction/sparepart/repositories-sparepart-impl"
	transactionsparepartservice "after-sales/api/services/transaction/sparepart"
	transactionsparepartserviceimpl "after-sales/api/services/transaction/sparepart/services-sparepart-impl"
	"github.com/stretchr/testify/assert"
	"testing"
)

func setupGoodsReceive() transactionsparepartservice.GoodsReceiveService {
	config.InitEnvConfigs(true, "")
	db := config.InitDB()
	rdb := config.InitRedis()
	GoodsReceiveRepositoryImpl := transactionsparepartrepositoryimpl.NewGoodsReceiveRepositoryImpl()
	GoodsReceiveServiceImpl := transactionsparepartserviceimpl.NewGoodsReceiveServiceImpl(GoodsReceiveRepositoryImpl, db, rdb)
	return GoodsReceiveServiceImpl
}
func TestSubmitGoodsReceive(t *testing.T) {
	goodsReceiveService := setupGoodsReceive()
	BinningToSubmit := 1046205

	res, err := goodsReceiveService.SubmitGoodsReceive(BinningToSubmit)
	assert.True(t, res, "result is false")
	assert.Nil(t, err, func() string {
		if err == nil {
			return "true"
		} else {
			if err.Err == nil {
				return err.Message
			} else {
				return err.Err.Error()
			}
		}
	})
}
