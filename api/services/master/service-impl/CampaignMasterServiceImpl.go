package masterserviceimpl

import (
	mastercampaignmasterentities "after-sales/api/entities/master/campaign_master"
	exceptions "after-sales/api/exceptions"
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

func (s *CampaignMasterServiceImpl) PostCampaignMaster(req masterpayloads.CampaignMasterPost) (mastercampaignmasterentities.CampaignMaster, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.CampaignMasterRepo.PostCampaignMaster(tx, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *CampaignMasterServiceImpl) PostCampaignDetailMaster(req masterpayloads.CampaignMasterDetailPayloads) (int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.CampaignMasterRepo.PostCampaignDetailMaster(tx, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (s *CampaignMasterServiceImpl) PostCampaignMasterDetailFromHistory(id int, idhead int) (int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.CampaignMasterRepo.PostCampaignMasterDetailFromHistory(tx, id, idhead)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *CampaignMasterServiceImpl) ChangeStatusCampaignMaster(id int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.CampaignMasterRepo.ChangeStatusCampaignMaster(tx, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *CampaignMasterServiceImpl) ActivateCampaignMasterDetail(ids string, id int) (bool,int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result,id, err := s.CampaignMasterRepo.ActivateCampaignMasterDetail(tx, ids, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return result,0, err
	}
	return result,id, err
}

func (s *CampaignMasterServiceImpl) DeactivateCampaignMasterDetail(ids string, id int) (bool,int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result,id, err := s.CampaignMasterRepo.DeactivateCampaignMasterDetail(tx, ids, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return result,0, err
	}
	return result,id, nil
}

func (s *CampaignMasterServiceImpl) GetByIdCampaignMaster(id int) (map[string]interface{}, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.CampaignMasterRepo.GetByIdCampaignMaster(tx, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *CampaignMasterServiceImpl) GetByIdCampaignMasterDetail(id int, idhead int) (map[string]interface{}, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.CampaignMasterRepo.GetByIdCampaignMasterDetail(tx, id, idhead)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *CampaignMasterServiceImpl) GetAllCampaignMasterCodeAndName(pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.CampaignMasterRepo.GetAllCampaignMasterCodeAndName(tx, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *CampaignMasterServiceImpl) GetAllCampaignMaster(filtercondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{},int,int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result,totalpages,totalrows, err := s.CampaignMasterRepo.GetAllCampaignMaster(tx, filtercondition, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return result,0,0, err
	}
	return result,totalpages,totalrows, nil
}

func (s *CampaignMasterServiceImpl) GetAllCampaignMasterDetail(pages pagination.Pagination, id int) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, page, limit, err := s.CampaignMasterRepo.GetAllCampaignMasterDetail(tx, pages, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return result, 0, 0, err
	}
	return result, page, limit, nil
}

func (s *CampaignMasterServiceImpl) UpdateCampaignMasterDetail(id int,linetypeid int, req masterpayloads.CampaignMasterDetailPayloads) (int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.CampaignMasterRepo.UpdateCampaignMasterDetail(tx, id, linetypeid, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (s *CampaignMasterServiceImpl) GetAllPackageMasterToCopy(pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.CampaignMasterRepo.GetAllPackageMasterToCopy(tx, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *CampaignMasterServiceImpl) SelectFromPackageMaster(id int, idhead int) (int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.CampaignMasterRepo.SelectFromPackageMaster(tx, id, idhead)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return 0, err
	}
	return result, nil
}
