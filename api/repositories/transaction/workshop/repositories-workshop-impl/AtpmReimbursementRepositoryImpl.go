package transactionworkshoprepositoryimpl

import (
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshoprepository "after-sales/api/repositories/transaction/workshop"
	"after-sales/api/utils"
	generalserviceapiutils "after-sales/api/utils/general-service"
	"errors"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type AtpmReimbursementRepositoryImpl struct {
}

func OpenAtpmReimbursementRepositoryImpl() transactionworkshoprepository.AtpmReimbursementRepository {
	return &AtpmReimbursementRepositoryImpl{}
}

// uspg_atReimbursement_Select
// IF @Option = 2
func (r *AtpmReimbursementRepositoryImpl) GetAll(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {

	var entities []transactionworkshopentities.AtpmReimbursement

	tx = utils.ApplyFilter(tx.Model(&transactionworkshopentities.AtpmReimbursement{}), filterCondition)

	tx.Scopes(pagination.Paginate(&pages, tx)).Find(&entities)

	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Data not found",
				Err:        tx.Error,
			}
		}
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get data",
			Err:        tx.Error,
		}
	}

	if len(entities) == 0 {
		pages.Rows = []map[string]interface{}{}
		return pages, nil
	}

	var results []map[string]interface{}
	for _, entity := range entities {

		companyResponses, companyErr := generalserviceapiutils.GetCompanyDataById(entity.CompanyId)
		if companyErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: companyErr.StatusCode,
				Message:    "Failed to fetch company data from internal service",
				Err:        companyErr.Err,
			}
		}

		result := map[string]interface{}{
			"claim_system_number":         entity.ClaimSystemNumber,
			"company_id":                  entity.CompanyId,
			"company_name":                companyResponses.CompanyName,
			"reimbursement_status_id":     entity.ReimbursementStatusId,
			"invoice_system_number":       entity.InvoiceSystemNumber,
			"invoice_document_number":     entity.InvoiceDocumentNumber,
			"invoice_date":                entity.InvoiceDate,
			"tax_invoice_system_number":   entity.TaxInvoiceSystemNumber,
			"tax_invoice_document_number": entity.TaxInvoiceDocumentNumber,
			"tax_invoice_date":            entity.TaxInvoiceDate,
			"kwitansi_system_number":      entity.KwitansiSystemNumber,
			"kwitansi_document_number":    entity.KwitansiDocumentNumber,
		}

		results = append(results, result)
	}

	pages.Rows = results
	return pages, nil
}

// uspg_atReimbursement_Insert
// IF @Option = 0
func (r *AtpmReimbursementRepositoryImpl) New(tx *gorm.DB, req transactionworkshoppayloads.AtpmReimbursementRequest) (transactionworkshopentities.AtpmReimbursement, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.AtpmReimbursement
	var invoiceDate, taxInvoiceDate time.Time

	// Check if InvoiceDate is non-zero
	if !req.InvoiceDate.IsZero() {
		invoiceDate = req.InvoiceDate
	} else {
		// If InvoiceDate is empty (zero value), set it to the time.Time zero value
		invoiceDate = time.Time{}
	}

	// Check if TaxInvoiceDate is non-zero
	if !req.TaxInvoiceDate.IsZero() {
		taxInvoiceDate = req.TaxInvoiceDate
	} else {
		// If TaxInvoiceDate is empty (zero value), set it to the time.Time zero value
		taxInvoiceDate = time.Time{}
	}

	var existingRecord transactionworkshopentities.AtpmReimbursement
	if err := tx.Where("claim_system_number = ?", req.ClaimSystemNumber).First(&existingRecord).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return transactionworkshopentities.AtpmReimbursement{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to check for existing record",
				Err:        err,
			}
		}

		entity = transactionworkshopentities.AtpmReimbursement{
			ClaimSystemNumber:        req.ClaimSystemNumber,
			CompanyId:                req.CompanyId,
			ReimbursementStatusId:    req.ReimbursementStatusId,
			InvoiceSystemNumber:      req.InvoiceSystemNumber,
			InvoiceDocumentNumber:    req.InvoiceDocumentNumber,
			InvoiceDate:              invoiceDate,
			TaxInvoiceSystemNumber:   req.TaxInvoiceSystemNumber,
			TaxInvoiceDocumentNumber: req.TaxInvoiceDocumentNumber,
			TaxInvoiceDate:           taxInvoiceDate,
			KwitansiSystemNumber:     req.KwitansiSystemNumber,
			KwitansiDocumentNumber:   req.KwitansiDocumentNumber,
		}

		if err := tx.Create(&entity).Error; err != nil {
			return transactionworkshopentities.AtpmReimbursement{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to insert data",
				Err:        err,
			}
		}
	}

	return entity, nil
}

// uspg_atReimbursement_Update
// IF @Option = 1
func (r *AtpmReimbursementRepositoryImpl) Save(tx *gorm.DB, claimsysno int, req transactionworkshoppayloads.AtpmReimbursementUpdate) (transactionworkshopentities.AtpmReimbursement, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.AtpmReimbursement
	var invoiceDate, taxInvoiceDate time.Time

	if !req.InvoiceDate.IsZero() {
		invoiceDate = req.InvoiceDate
	} else {
		invoiceDate = time.Time{}
	}

	if !req.TaxInvoiceDate.IsZero() {
		taxInvoiceDate = req.TaxInvoiceDate
	} else {
		taxInvoiceDate = time.Time{}
	}

	if err := tx.Where("claim_system_number = ?", claimsysno).First(&entity).Error; err != nil {
		return transactionworkshopentities.AtpmReimbursement{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Reimbursement data not found",
			Err:        err,
		}
	}

	var claimVehicle transactionworkshopentities.AtpmClaimVehicle
	if err := tx.Where("claim_system_number = ?", claimsysno).First(&claimVehicle).Error; err != nil {
		return transactionworkshopentities.AtpmReimbursement{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Claim vehicle data not found",
			Err:        err,
		}
	}

	if claimVehicle.ClaimStatusId != 10 {
		return transactionworkshopentities.AtpmReimbursement{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Claim document is already closed",
		}
	}

	entity.ClaimSystemNumber = req.ClaimSystemNumber
	entity.CompanyId = req.CompanyId
	entity.ReimbursementStatusId = req.ReimbursementStatusId
	entity.InvoiceSystemNumber = req.InvoiceSystemNumber
	entity.InvoiceDocumentNumber = req.InvoiceDocumentNumber
	entity.InvoiceDate = invoiceDate
	entity.TaxInvoiceSystemNumber = req.TaxInvoiceSystemNumber
	entity.TaxInvoiceDocumentNumber = req.TaxInvoiceDocumentNumber
	entity.TaxInvoiceDate = taxInvoiceDate
	entity.KwitansiSystemNumber = req.KwitansiSystemNumber
	entity.KwitansiDocumentNumber = req.KwitansiDocumentNumber

	if err := tx.Save(&entity).Error; err != nil {
		return transactionworkshopentities.AtpmReimbursement{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to update reimbursement data",
			Err:        err,
		}
	}

	return entity, nil
}

// uspg_atReimbursement_Update
// IF @Option = 1
func (r *AtpmReimbursementRepositoryImpl) Submit(tx *gorm.DB, claimsysno int) (bool, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.AtpmReimbursement
	if err := tx.Where("claim_system_number = ?", claimsysno).First(&entity).Error; err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Reimbursement data not found",
			Err:        err,
		}
	}

	var claimVehicle transactionworkshopentities.AtpmClaimVehicle
	if err := tx.Where("claim_system_number = ?", claimsysno).First(&claimVehicle).Error; err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Claim vehicle data not found",
			Err:        err,
		}
	}

	// 'Draft' status
	if claimVehicle.ClaimStatusId != 10 {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Document status is not Draft",
		}
	}

	// Set the claim status to "Approval Payment"
	claimVehicle.ClaimStatusId = 20

	if err := tx.Model(&claimVehicle).Update("claim_status_id", claimVehicle.ClaimStatusId).Error; err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to update claim vehicle status",
			Err:        err,
		}
	}

	// Set the reimbursement status to "Approved"
	entity.ReimbursementStatusId = 2

	if err := tx.Model(&entity).Update("reimbursement_status_id", entity.ReimbursementStatusId).Error; err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to update reimbursement status",
			Err:        err,
		}
	}

	return true, nil
}
