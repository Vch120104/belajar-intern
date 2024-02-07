package masterserviceimpl

import (
	masterpayloads "after-sales/api/payloads/master"
	masterrepository "after-sales/api/repositories/master"
	masterservice "after-sales/api/services/master"

	"gorm.io/gorm"
)

type IncentiveGroupImpl struct {
	incentiveGroupRepo masterrepository.IncentiveGroupRepository
}

func StartIncentiveGroup(incentiveGroupRepo masterrepository.IncentiveGroupRepository) masterservice.IncentiveGroupService {
	return &IncentiveGroupImpl{
		incentiveGroupRepo: incentiveGroupRepo,
	}
}
func (s *IncentiveGroupImpl) WithTrx(trxHandle *gorm.DB) masterservice.IncentiveGroupService {
	s.incentiveGroupRepo = s.incentiveGroupRepo.WithTrx(trxHandle)
	return s
}

func (s *IncentiveGroupImpl) GetAllIncentiveGroupIsActive() ([]masterpayloads.IncentiveGroupResponse, error) {
	results, err := s.incentiveGroupRepo.GetAllIncentiveGroupIsActive()
	if err != nil {
		return results, err
	}
	return results, nil
}
