package masterserviceimpl

import (
	exceptionsss_test "after-sales/api/expectionsss"
	"after-sales/api/helper"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	masterservice "after-sales/api/services/master"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type CampaignMasterServiceImpl struct {
	CampaignMasterRepo masterrepository.CampaignMasterRepository
	DB                 *gorm.DB
}

func startCampaignMasterService(CampaignMasterRepo masterrepository.CampaignMasterRepository, db *gorm.DB) masterservice.CampaignMasterService {
	return &CampaignMasterServiceImpl{
		CampaignMasterRepo: CampaignMasterRepo,
		DB:                 db,
	}
}

func (s *CampaignMasterServiceImpl) PostCampaignMaster(req masterpayloads.CampaignMasterPost)(bool,*exceptionsss_test.BaseErrorResponse){
	tx:= s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	_,err:= s.CampaignMasterRepo.PostCampaignMaster(tx,req)
	if err != nil{
		return false,err
	}
	return true,nil
}

func (s *CampaignMasterServiceImpl) PostCampaignDetailMaster(req masterpayloads.CampaignMasterDetailPayloads)(bool,*exceptionsss_test.BaseErrorResponse){
	tx:=s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	_,err:=s.CampaignMasterRepo.PostCampaignDetailMaster(tx,req)
	if err != nil{
		return false,err
	}
	return true,nil
}

func (s *CampaignMasterServiceImpl) ChangeStatusCampaignMaster(id int)(bool, *exceptionsss_test.BaseErrorResponse){
	tx :=s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	result,err := s.CampaignMasterRepo.ChangeStatusCampaignMaster(tx,id)
	if err != nil{
		return result,err
	}
	return result,nil
}

func (s *CampaignMasterServiceImpl) ActivateCampaignMasterDetail(ids string)(bool,*exceptionsss_test.BaseErrorResponse){
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	result,err := s.CampaignMasterRepo.ActivateCampaignMasterDetail(tx,ids)
	if err != nil{
		return result,err
	}
	return result,err
}

func (s *CampaignMasterServiceImpl) DeactivateCampaignMasterDetail(ids string)(bool,*exceptionsss_test.BaseErrorResponse){
	tx:=s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	result,err := s.CampaignMasterRepo.DeactivateCampaignMasterDetail(tx,ids)
	if err != nil{
		return result,err
	}
	return result,nil
}

func (s *CampaignMasterServiceImpl) GetByIdCampaignMaster(id int)(masterpayloads.CampaignMasterResponse,*exceptionsss_test.BaseErrorResponse){
	tx:= s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	result,err := s.CampaignMasterRepo.GetByIdCampaignMaster(tx,id)
	if err != nil{
		return result,err
	}
	return result,nil
}

func (s *CampaignMasterServiceImpl) GetByIdCampaignMasterDetail(id int)(masterpayloads.CampaignMasterDetailPayloads,*exceptionsss_test.BaseErrorResponse){
	tx:=s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	result,err:=s.CampaignMasterRepo.GetByIdCampaignMasterDetail(tx,id)
	if err != nil{
		return result,err
	}
	return result,nil
}

func (s *CampaignMasterServiceImpl) GetAllCampaignMasterCodeAndName()([]masterpayloads.GetHistory,*exceptionsss_test.BaseErrorResponse){
	tx:=s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	result,err:=s.CampaignMasterRepo.GetAllCampaignMasterCodeAndName(tx)
	if err != nil{
		return result,err
	}
	return result,nil
}

func (s *CampaignMasterServiceImpl) GetAllCampaignMaster(filtercondition []utils.FilterCondition,pages pagination.Pagination)([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse){
	tx:=s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	result,totalPages,totalRows,err:= s.CampaignMasterRepo.GetAllCampaignMaster(tx,filtercondition,pages)
	if err != nil{
		return result,0,0,err
	}
	return result,totalPages,totalRows,nil
}

func (s *CampaignMasterServiceImpl) GetAllCampaignMasterDetail(pages pagination.Pagination, id int)(pagination.Pagination, *exceptionsss_test.BaseErrorResponse){
	tx:=s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	result,err:=s.CampaignMasterRepo.GetAllCampaignMasterDetail(tx,pages,id)
	if err != nil{
		return result,err
	}
	return result,nil
}
