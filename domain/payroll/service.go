package payroll

import (
	md "github.com/ebikode/payroll-core/model"
	tr "github.com/ebikode/payroll-core/translation"
	ut "github.com/ebikode/payroll-core/utils"
)

// PayrollService provides payroll operations
type PayrollService interface {
	GetPayrollReports() []*md.PayrollReport
	GetPayrollAllMonthAndYear() []*md.PayrollMonthYear
	GetPayroll(string, string) *md.Payroll
	GetLastPayroll() *md.Payroll
	GetPayrolls(int, int) []*md.Payroll
	GetPayrollsByMonthYear(uint, uint) []*md.Payroll
	GetSinglePayrollByMonthYear(uint, uint) *md.Payroll
	GetEmployeePayrolls(string, int, int) []*md.Payroll
	CreatePayroll(md.Payroll) (*md.Payroll, tr.TParam, error)
	UpdatePayroll(*md.Payroll) (*md.Payroll, tr.TParam, error)
	UpdatePayrollStatus(string, int, int)
	UpdatePayrollPaymentStatus(string, int, int)
}

type service struct {
	pRepo PayrollRepository
}

// NewService creates p payroll service with the necessary dependencies
func NewService(
	pRepo PayrollRepository,
) PayrollService {
	return &service{pRepo}
}

/*
* Get single payroll of a employee
* @param employeeID => the ID of the employee whose payroll is needed
* @param payrollID => the ID of the payroll requested.
 */
func (s *service) GetPayroll(employeeID, payrollID string) *md.Payroll {
	return s.pRepo.Get(employeeID, payrollID)
}

/*
* GetLastPayroll
 */
func (s *service) GetLastPayroll() *md.Payroll {
	return s.pRepo.GetLastPayroll()
}

/*
* Get all payrolls
* @param page => the page number to return
* @param limit => limit per page to return
 */
func (s *service) GetPayrolls(page, limit int) []*md.Payroll {
	return s.pRepo.GetAll(page, limit)
}

func (s *service) GetPayrollsByMonthYear(month, year uint) []*md.Payroll {
	return s.pRepo.GetByMonthYear(month, year)
}

func (s *service) GetPayrollAllMonthAndYear() []*md.PayrollMonthYear {
	return s.pRepo.GetAllMonthAndYear()
}

func (s *service) GetSinglePayrollByMonthYear(month, year uint) *md.Payroll {
	return s.pRepo.GetSingleByMonthYear(month, year)
}

func (s *service) GetPayrollReports() []*md.PayrollReport {
	return s.pRepo.GetReports()
}

/*
* Get all payrolls of a employee
* @param employeeID => the ID of the employee whose payroll is needed
* @param page => the page number to return
* @param limit => limit per page to return
 */
func (s *service) GetEmployeePayrolls(employeeID string, page, limit int) []*md.Payroll {
	return s.pRepo.GetByEmployee(employeeID, page, limit)
}

// Create New payroll
func (s *service) CreatePayroll(p md.Payroll) (*md.Payroll, tr.TParam, error) {
	// Generate ID
	pID := ut.RandomBase64String(8, "MDpy")
	p.ID = pID

	payroll, err := s.pRepo.Store(p)

	if err != nil {
		tParam := tr.TParam{
			Key:          "error.resource_creation_error",
			TemplateData: nil,
			PluralCount:  nil,
		}

		return payroll, tParam, err
	}

	return payroll, tr.TParam{}, nil
}

// update existing payroll
func (s *service) UpdatePayroll(p *md.Payroll) (*md.Payroll, tr.TParam, error) {
	payroll, err := s.pRepo.Update(p)

	if err != nil {
		tParam := tr.TParam{
			Key:          "error.resource_update_error",
			TemplateData: nil,
			PluralCount:  nil,
		}

		return payroll, tParam, err
	}
	return payroll, tr.TParam{}, nil
}

// UpdatePayrollStatus existing payrolls status by month and year
func (s *service) UpdatePayrollStatus(status string, month, year int) {
	s.pRepo.UpdateStatus(status, month, year)
}

// UpdatePayrollPaymentStatus existing payrolls payment status by month and year
func (s *service) UpdatePayrollPaymentStatus(status string, month, year int) {
	s.pRepo.UpdatePaymentStatus(status, month, year)
}
