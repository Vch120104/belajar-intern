package repositoriespointprospectingimpl

import (
	transactionunitentities "after-sales/api/entities/transaction/Unit"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionunitpayloads "after-sales/api/payloads/transaction/unit"
	repositories "after-sales/api/repositories/transaction/unit"
	"after-sales/api/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PointProspectingTransactionRepositoryImpl struct {
}

func NewPointProspectingTransactionRepositoryImpl() repositories.PointProspectingRepository {
	return &PointProspectingTransactionRepositoryImpl{}
}

func (r *PointProspectingTransactionRepositoryImpl) GetAllCompanyData(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var datas []transactionunitpayloads.PointProspectingResponse

	query := tx.Model(&transactionunitentities.GmComp0{}).
		Distinct("gmComp0.company_code as CompanyCode, gmComp0.company_name as CompanyName, " + "CASE gmComp0.record_status WHEN 'A' THEN 'Active' WHEN 'D' THEN 'Deactive' END as Status").
		Joins("JOIN gmRef as a on gmComp0.company_code = a.company_code").
		Joins("INNER JOIN gmEmp1 on gmEmp1.company_code = gmComp0.company_code")

	whereQ := utils.ApplyFilter(query, filterCondition)
	paginatedQuery := whereQ.Scopes(pagination.Paginate(&pages, whereQ)).Order("gmComp0.company_code")

	err := paginatedQuery.Scan(&datas).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error when scanning data",
			Err:        err,
		}
	}

	pages.Rows = datas

	logrus.Debug("result", datas)
	return pages, nil
}

func (r *PointProspectingTransactionRepositoryImpl) GetAllSalesRepresentative(tx *gorm.DB, filteredCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var datas []transactionunitpayloads.GetAllSalesRepresentativeResponse
	query := tx.Model(&transactionunitentities.GmEmp{}).
		Distinct("gmEmp.employee_no as EmployeeNo, gmEmp.employee_name as EmployeeName, gmComp0.company_name as CompanyName").
		Joins("INNER JOIN gmComp0 on gmComp0.company_code = gmEmp.company_code").
		Joins("INNER JOIN comGenVariable on comGenVariable.value = gmEmp.job_position AND comGenVariable.variable in(?, ?, ?, ?, ?, ?)",
			"JOB_POS_SALES_REP", "JOB_POS_SALES_COUNTER", "JOB_POS_SALES_COUNTER_DUTRO", "JOB_POS_SALES_COUNTER_HINO", "JOB_POS_SALES_COUNTER_MIX",
			"JOB_POS_SALES_COUNTER_RANGER")

	whereQ := utils.ApplyFilter(query, filteredCondition)
	paginatedQuery := whereQ.Scopes(pagination.Paginate(&pages, whereQ)).Order("gmEmp.employee_no")

	err := paginatedQuery.Scan(&datas).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error when scanning data",
			Err:        err,
		}
	}

	pages.Rows = datas
	logrus.Debug("datas", pages.Rows)
	fmt.Println()
	return pages, nil
}

func (r *PointProspectingTransactionRepositoryImpl) GetSalesByCompanyCode(tx *gorm.DB, companyCode float64, filteredCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var gmEmp []transactionunitpayloads.GetAllSalesRepresentativeResponse

	query := tx.Model(&transactionunitentities.GmEmp{}).
		Distinct("gmEmp.employee_no as EmployeeNo, gmEmp.employee_name as EmployeeName, gmComp0.company_name as CompanyName").
		Joins("INNER JOIN gmComp0 on gmComp0.company_code = gmEmp.company_code").
		Joins("INNER JOIN comGenVariable on comGenVariable.value = gmEmp.job_position AND comGenVariable.variable in(?, ?, ?, ?, ?, ?)",
			"JOB_POS_SALES_REP", "JOB_POS_SALES_COUNTER", "JOB_POS_SALES_COUNTER_DUTRO", "JOB_POS_SALES_COUNTER_HINO", "JOB_POS_SALES_COUNTER_MIX",
			"JOB_POS_SALES_COUNTER_RANGER").
		Where("gmEmp.company_code = ?", companyCode)

	whereQ := utils.ApplyFilter(query, filteredCondition)
	paginatedQuery := whereQ.Scopes(pagination.Paginate(&pages, whereQ)).Order("gmEmp.employee_no")

	err := paginatedQuery.Scan(&gmEmp).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error when scanning data",
			Err:        err,
		}
	}
	pages.Rows = gmEmp
	return pages, nil
}

func (r *PointProspectingTransactionRepositoryImpl) Process(tx *gorm.DB, request transactionunitpayloads.ProcessRequest) (bool, *exceptions.BaseErrorResponse) {
	err := tx.Model(&transactionunitentities.UtPointProspecting{}).
		Where("company_code = ? AND period_month = ? AND period_year = ?", request.CompanyCode, request.PeriodMonth, request.PeriodYear).
		Delete(&transactionunitentities.UtPointProspecting{}).Error

	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error when deleting data",
			Err:        err,
		}
	}
	month1stDate, err := time.Parse("2006-01-02", fmt.Sprintf("%s-%s-01", request.PeriodYear, request.PeriodMonth))
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error when parsing date",
			Err:        err,
		}
	}
	monthEndDate := month1stDate.AddDate(0, 1, -1)

	res := []transactionunitentities.UtPointProspecting{}

	query := tx.Model(&transactionunitentities.RtInvoice0{}).
	Distinct(`
        rtInvoice0.company_code, ?, ?, b.sales_rep_code, f.prospect_system_no, d.spm_system_no, rtInvoice0.vehicle_brand,
        CASE WHEN ISNULL(f.prospect_date, '') <> '' THEN f.prospect_date ELSE 0 END,
        CASE WHEN ISNULL(f.prospect_stage, '') <> '' THEN f.prospect_stage ELSE 0 END,
        CASE WHEN ISNULL(f.prospect_title_prefix, '') <> '' THEN f.prospect_title_prefix ELSE 0 END,
        CASE WHEN ISNULL(f.prospect_name, '') <> '' THEN f.prospect_name ELSE 0 END,
        CASE WHEN ISNULL(f.prospect_title_suffix, '') <> '' THEN f.prospect_title_suffix ELSE 0 END,
        CASE WHEN ISNULL(f.prospect_src_code, '') <> '' THEN f.prospect_src_code ELSE 0 END,
        CASE WHEN ISNULL(f.prospect_type, '') <> '' THEN f.prospect_type ELSE 0 END,
        CASE WHEN ISNULL(f.prospect_ref, '') <> '' THEN f.prospect_ref ELSE 0 END,
        CASE WHEN ISNULL(f.prospect_notes, '') <> '' THEN f.prospect_notes ELSE 0 END,
        CASE WHEN ISNULL(f.buying_budget, 0) <> 0 THEN f.buying_budget ELSE 0 END,
        CASE WHEN ISNULL(f.buying_plan, '') <> '' THEN f.buying_plan ELSE 0 END,
        CASE WHEN ISNULL(f.fund_type, '') <> '' THEN f.fund_type ELSE 0 END,
        CASE WHEN ISNULL(f.model_code, '') <> '' THEN f.model_code ELSE 0 END,
        CASE WHEN ISNULL(f.variant_code, '') <> '' THEN f.variant_code ELSE 0 END,
        CASE WHEN ISNULL(f.prospect_address_1, '') <> '' OR ISNULL(f.prospect_address_2, '') <> '' OR ISNULL(f.prospect_address_3, '') <> '' THEN 1 ELSE 0 END,
        CASE WHEN ISNULL(f.prospect_village_code, '') <> '' THEN 1 ELSE 0 END,
        CASE WHEN ISNULL(f.prospect_mobile_phone, '') <> '' THEN 1 ELSE 0 END,
        CASE WHEN ISNULL(f.prospect_email_address, '') <> '' THEN 1 ELSE 0 END,
        CASE WHEN ISNULL(f.prospect_website, '') <> '' THEN 1 ELSE 0 END,
        CASE WHEN ISNULL(f.prospect_phone_no, '') <> '' THEN 1 ELSE 0 END,
        CASE WHEN ISNULL(f.prospect_fax_no, '') <> '' THEN 1 ELSE 0 END,
        CASE WHEN ISNULL(f.prospect_gender, '') <> '' THEN 1 ELSE 0 END,
        CASE WHEN ISNULL(f.biz_type, '') <> '' THEN 1 ELSE 0 END,
        CASE WHEN ISNULL(f.biz_group, '') <> '' THEN 1 ELSE 0 END,
        CASE WHEN ISNULL(f.contact_person, '') <> '' THEN 1 ELSE 0 END,
        CASE WHEN ISNULL(f.contact_gender, '') <> '' THEN 1 ELSE 0 END,
        CASE WHEN ISNULL(f.contact_job_title, '') <> '' THEN 1 ELSE 0 END,
        CASE WHEN ISNULL(f.contact_mobile_phone, '') <> '' THEN 1 ELSE 0 END,
        CASE WHEN ISNULL(f.contact_email_address, '') <> '' THEN 1 ELSE 0 END,
        CASE WHEN ISNULL(f.test_drv_date_schedule, '') <> '' THEN 1 ELSE 0 END,
        CASE WHEN ISNULL(f.test_drv_date_actual, '') <> '' THEN 1 ELSE 0 END,
        CASE WHEN ISNULL(f.competitor_model, '') <> '' THEN 1 ELSE 0 END,
        CASE WHEN ISNULL(c.option_code, '') <> '' THEN 5 ELSE 0 END,
        CASE WHEN ISNULL(c.stage_cc_date, '') < ISNULL(c.stage_ch_date, '') THEN 5 ELSE 0 END,
        CASE WHEN ISNULL(c.stage_ch_date, '') < ISNULL(c.stage_p_date, '') THEN 5 ELSE 0 END,
        CASE WHEN ISNULL(c.stage_p_date, '') < ISNULL(c.stage_hp_date, '') THEN 5 ELSE 0 END,
        CASE WHEN ISNULL(c.stage_hp_date, '') < ISNULL(c.stage_do_date, '') THEN 5 ELSE 0 END,
        CASE WHEN ISNULL(g.follow_up_note, '') <> '' THEN 2 ELSE 0 END,
        CASE WHEN ISNULL(g.result_note, '') <> '' THEN 2 ELSE 0 END,
        CASE WHEN ISNULL(d.spm_stage_date, '') <> '' THEN 1 ELSE 0 END,
        CASE WHEN ISNULL(d.spm_remark, '') <> '' THEN 1 ELSE 0 END,
        CASE WHEN ISNULL(d.order_by_tax_reg_no, '') <> '' THEN 1 ELSE 0 END,
        CASE WHEN ISNULL(d.order_by_tax_reg_date, '') <> '' THEN 1 ELSE 0 END,
        CASE WHEN ISNULL(d.order_by_tax_name, '') <> '' THEN 1 ELSE 0 END,
        CASE WHEN ISNULL(d.order_by_tax_address_1, '') <> '' THEN 1 ELSE 0 END,
        CASE WHEN ISNULL(d.order_by_tax_address_2, '') <> '' THEN 1 ELSE 0 END,
        CASE WHEN ISNULL(d.order_by_tax_address_3, '') <> '' THEN 1 ELSE 0 END,
        CASE WHEN ISNULL(d.order_by_tax_village_code, '') <> '' THEN 1 ELSE 0 END,
        CASE WHEN ISNULL(d.order_by_tax_subdistrict_code, '') <> '' THEN 1 ELSE 0 END,
        CASE WHEN ISNULL(d.order_by_tax_municipality_code, '') <> '' THEN 1 ELSE 0 END,
        CASE WHEN ISNULL(d.order_by_tax_province_code, '') <> '' THEN 1 ELSE 0 END,
        CASE WHEN ISNULL(d.order_by_tax_city_code, '') <> '' THEN 1 ELSE 0 END,
        CASE WHEN ISNULL(d.order_by_tax_zip_code, '') <> '' THEN 1 ELSE 0 END,
        CASE WHEN ISNULL(d.pkp_status, '') <> '' THEN 1 ELSE 0 END,
        CASE WHEN ISNULL(d.pkp_date, '') <> '' THEN 1 ELSE 0 END,
        CASE WHEN ISNULL(d.pkp_no, '') <> '' THEN 1 ELSE 0 END,
        CASE WHEN ISNULL(d.pkp_type, '') <> '' THEN 1 ELSE 0 END,
        CASE WHEN ISNULL(d.corp_biz_type, '') <> '' THEN 1 ELSE 0 END,
        CASE WHEN ISNULL(d.corp_biz_group, '') <> '' THEN 1 ELSE 0 END,
        CASE WHEN ISNULL(d.corp_web_site, '') <> '' THEN 1 ELSE 0 END,
        CASE WHEN ISNULL(d.corp_po_no, '') <> '' THEN 1 ELSE 0 END,
        CASE WHEN ISNULL(d.corp_po_date, '') <> '' THEN 1 ELSE 0 END,
        CASE WHEN ISNULL(d.corp_contact_name, '') <> '' THEN 1 ELSE 0 END,
        CASE WHEN ISNULL(d.corp_contact_gender, '') <> '' THEN 1 ELSE 0 END,
        CASE WHEN ISNULL(d.corp_contact_job_title, '') <> '' THEN 1 ELSE 0 END,
        CASE WHEN ISNULL(d.corp_mobile_phone, '') <> '' THEN 1 ELSE 0 END,
        CASE WHEN ISNULL(d.corp_email_address, '') <> '' THEN 1 ELSE 0 END,
        CASE WHEN ISNULL(d.corr_prefix, '') <> '' THEN 1 ELSE 0 END,
        CASE WHEN ISNULL(d.corr_name, '') <> '' THEN 1 ELSE 0 END,
        CASE WHEN ISNULL(d.corr_suffix, '') <> '' THEN 1 ELSE 0 END,
        CASE WHEN ISNULL(d.corr_gender, '') <> '' THEN 1 ELSE 0 END,
        CASE WHEN ISNULL(d.corr_job_title, '') <> '' THEN 1 ELSE 0 END,
        CASE WHEN ISNULL(d.corr_address_1, '') <> '' OR ISNULL(d.corr_address_2, '') <> '' OR ISNULL(d.corr_address_3, '') <> '' THEN 1 ELSE 0 END,
        CASE WHEN ISNULL(d.corr_village_code, '') <> '' THEN 1 ELSE 0 END,
        CASE WHEN ISNULL(d.corr_phone_no, '') <> '' THEN 1 ELSE 0 END,
        CASE WHEN ISNULL(d.corr_fax_no, '') <> '' THEN 1 ELSE 0 END,
        CASE WHEN ISNULL(d.corr_mobile_phone, '') <> '' THEN 1 ELSE 0 END,
        CASE WHEN ISNULL(d.corr_email_address, '') <> '' THEN 1 ELSE 0 END,
        CASE WHEN ISNULL(d.remark, '') <> '' THEN d.remark ELSE 0 END,
        CASE WHEN ISNULL(c.user_birthday, '') <> '' THEN c.user_birthday ELSE 0 END,
        CASE WHEN ISNULL(c.user_religion, '') <> '' THEN 1 ELSE 0 END,
        0, ?, ?, NULL, NULL
    `, request.CompanyCode, request.PeriodYear, request.PeriodMonth, request.CreationUserId, time.Now()).
		Joins("inner join rtinvoice1 b on b.inv_sys_no = rtInvoice0.inv_sys_no and rtInvoice0.inv_type = ? and isnull(b.is_return, 0) = 0 and rtInvoice0.trx_type = ? and b.ref_type = ? and rtInvoice0.inv_status = ?", "AI01", "SU01", "SPM", "20").
		Joins("inner join utSPM1 c on c.spm_system_no = b.ref_sys_no and c.spm_line_no = b.ref_line_no and c.company_code = rtInvoice0.company_code").
		Joins("inner join utSPM0 d on d.spm_system_no = c.spm_system_no and d.company_code = rtInvoice0.company_code and d.spm_status = ?", "20").
		Joins("inner join utProspect1 e on e.prospect_system_no = c.prospect_system_no and e.prospect_line = c.prospect_line").
		Joins("inner join utProspect0 f on f.prospect_system_no = e.prospect_system_no and f.company_code = rtInvoice0.company_code").
		Joins("left join utProspect2 g on g.prospect_system_no = e.prospect_system_no and g.prospect_line = e.prospect_line").
		Where("rtInvoice0.company_code = ? and rtInvoice0.inv_date between ? and ?", request.CompanyCode, month1stDate, monthEndDate).
		Order("rtInvoice0.company_code").
		Scan(&res)

	if query.Error != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error when scanning data",
			Err:        query.Error,
		}
	}
	return true, nil

}
