package transactionsparepartserviceimpl

import (
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionsparepartrepository "after-sales/api/repositories/transaction/sparepart"
	transactionsparepartservice "after-sales/api/services/transaction/sparepart"
)

type SupplySlipServiceImpl struct {
	supplySlipRepo transactionsparepartrepository.SupplySlipRepository
}

func StartSupplySlipService(supplySlipRepo transactionsparepartrepository.SupplySlipRepository) transactionsparepartservice.SupplySlipService {
	return &SupplySlipServiceImpl{
		supplySlipRepo: supplySlipRepo,
	}
}

func (s *SupplySlipServiceImpl) GetSupplySlipById(id int32) (transactionsparepartpayloads.SupplySlipResponse, error) {
	value, err := s.supplySlipRepo.GetSupplySlipById(id)
	if err != nil {
		return transactionsparepartpayloads.SupplySlipResponse{}, err
	}
	return value, nil
}

func (s *SupplySlipServiceImpl) GetSupplySlipDetailById(id int32) (transactionsparepartpayloads.SupplySlipDetailResponse, error) {
	value, err := s.supplySlipRepo.GetSupplySlipDetailById(id)
	if err != nil {
		return transactionsparepartpayloads.SupplySlipDetailResponse{}, err
	}
	return value, nil
}
