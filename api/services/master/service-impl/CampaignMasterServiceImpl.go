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

func StartCampaignMasterService(CampaignMasterRepo masterrepository.CampaignMasterRepository, db *gorm.DB) masterservice.CampaignMasterService {
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

func (s *CampaignMasterServiceImpl) PostCampaignMasterDetailFromHistory(id int, idhead int)(bool,*exceptionsss_test.BaseErrorResponse){
	tx:=s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	result,err:= s.CampaignMasterRepo.PostCampaignMasterDetailFromHistory(tx,id,idhead)
	if err != nil{
		return result,err
	}
	return result,nil
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

func (s *CampaignMasterServiceImpl) ActivateCampaignMasterDetail(ids string, id int)(bool,*exceptionsss_test.BaseErrorResponse){
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	result,err := s.CampaignMasterRepo.ActivateCampaignMasterDetail(tx,ids, id)
	if err != nil{
		return result,err
	}
	return result,err
}

func (s *CampaignMasterServiceImpl) DeactivateCampaignMasterDetail(ids string, id int)(bool,*exceptionsss_test.BaseErrorResponse){
	tx:=s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	result,err := s.CampaignMasterRepo.DeactivateCampaignMasterDetail(tx,ids, id)
	if err != nil{
		return result,err
	}
	return result,nil
}

func (s *CampaignMasterServiceImpl) GetByIdCampaignMaster(id int)([]map[string]interface{},*exceptionsss_test.BaseErrorResponse){
	tx:= s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	result,err := s.CampaignMasterRepo.GetByIdCampaignMaster(tx,id)
	if err != nil{
		return result,err
	}
	return result,nil
}

func (s *CampaignMasterServiceImpl) GetByIdCampaignMasterDetail(id int,idhead int)(map[string]interface{},*exceptionsss_test.BaseErrorResponse){
	tx:=s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	result,err:=s.CampaignMasterRepo.GetByIdCampaignMasterDetail(tx,id,idhead)
	if err != nil{
		return result,err
	}
	return result,nil
}

func (s *CampaignMasterServiceImpl) GetAllCampaignMasterCodeAndName(pages pagination.Pagination)(pagination.Pagination,*exceptionsss_test.BaseErrorResponse){
	tx:=s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	result,err:=s.CampaignMasterRepo.GetAllCampaignMasterCodeAndName(tx,pages)
	if err != nil{
		return result,err
	}
	return result,nil
}

func (s *CampaignMasterServiceImpl) GetAllCampaignMaster(filtercondition []utils.FilterCondition,pages pagination.Pagination)(pagination.Pagination, *exceptionsss_test.BaseErrorResponse){
	tx:=s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	result,err:= s.CampaignMasterRepo.GetAllCampaignMaster(tx,filtercondition,pages)
	if err != nil{
		return result,err
	}
	return result,nil
}

func (s *CampaignMasterServiceImpl) GetAllCampaignMasterDetail(pages pagination.Pagination, id int)([]map[string]interface{},int,int, *exceptionsss_test.BaseErrorResponse){
	tx:=s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	result,page,limit,err:=s.CampaignMasterRepo.GetAllCampaignMasterDetail(tx,pages,id)
	if err != nil{
		return result,0,0,err
	}
	return result,page,limit,nil
}

func (s *CampaignMasterServiceImpl) UpdateCampaignMasterDetail(id int, req masterpayloads.CampaignMasterDetailPayloads)(bool,*exceptionsss_test.BaseErrorResponse){
	tx:=s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	result,err:=s.CampaignMasterRepo.UpdateCampaignMasterDetail(tx,id,req)
	if err!= nil{
		return false,err
	}
	return result,nil
}

func (s *CampaignMasterServiceImpl) GetAllPackageMasterToCopy(pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse){
	tx:=s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	result,err:= s.CampaignMasterRepo.GetAllPackageMasterToCopy(tx,pages)
	if err != nil{
		return result,err
	}
	return result,nil
}

func (s *CampaignMasterServiceImpl) SelectFromPackageMaster(id int, idhead int) (bool, *exceptionsss_test.BaseErrorResponse){
	tx:=s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	result,err:=s.CampaignMasterRepo.SelectFromPackageMaster(tx,id,idhead)
	if err != nil{
		return false,err 
	}
	return result,nil
}
 