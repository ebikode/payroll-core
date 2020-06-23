package salary

import (
	"net/http"

	md "github.com/ebikode/payroll-core/model"
	tr "github.com/ebikode/payroll-core/translation"
	ut "github.com/ebikode/payroll-core/utils"
	validation "github.com/go-ozzo/ozzo-validation"
)

// SalaryService  provides salary operations
type SalaryService interface {
	GetSalary(uint) *md.Salary
	GetSalaryByEmployeeID(string) *md.Salary
	GetSalaries(int, int) []*md.Salary
	CreateSalary(md.Salary) (*md.Salary, tr.TParam, error)
	UpdateSalary(*md.Salary) (*md.Salary, tr.TParam, error)
}

type service struct {
	sRepo SalaryRepository
}

// NewService creates a salary service with the necessary dependencies
func NewService(
	sRepo SalaryRepository,
) SalaryService {
	return &service{sRepo}
}

// Get a salary
func (s *service) GetSalary(id uint) *md.Salary {
	return s.sRepo.Get(id)
}

// GetSalaries Get all salarys from DB
//
// @userType == admin | customer
func (s *service) GetSalaries(page, limit int) []*md.Salary {
	return s.sRepo.GetAll(page, limit)
}

func (s *service) GetSalaryByEmployeeID(employeeID string) *md.Salary {
	return s.sRepo.GetByEmployeeID(employeeID)
}

// CreateSalary Creates New salary
func (s *service) CreateSalary(c md.Salary) (*md.Salary, tr.TParam, error) {

	salary, err := s.sRepo.Store(c)

	if err != nil {
		tParam := tr.TParam{
			Key:          "error.resource_creation_error",
			TemplateData: nil,
			PluralCount:  nil,
		}

		return salary, tParam, err
	}

	return salary, tr.TParam{}, nil

}

// UpdateSalary update existing salary
func (s *service) UpdateSalary(c *md.Salary) (*md.Salary, tr.TParam, error) {
	salary, err := s.sRepo.Update(c)

	if err != nil {
		tParam := tr.TParam{
			Key:          "error.resource_update_error",
			TemplateData: nil,
			PluralCount:  nil,
		}

		return salary, tParam, err
	}

	return salary, tr.TParam{}, nil

}

// Validate Function for validating salary input
func Validate(salary Payload, r *http.Request) error {
	return validation.ValidateStruct(&salary,
		validation.Field(&salary.EmployeeID, ut.IDRule(r)...),
		validation.Field(&salary.Salary, ut.MoneyRule(r)...),
		validation.Field(&salary.Pension, ut.RequiredRule(r, "general.pension")...),
		validation.Field(&salary.Paye, ut.RequiredRule(r, "general.paye")...),
		validation.Field(&salary.Nsitf, ut.RequiredRule(r, "general.nsitf")...),
		validation.Field(&salary.Nhf, ut.RequiredRule(r, "general.nhf")...),
		validation.Field(&salary.Itf, ut.RequiredRule(r, "general.itf")...),
	)
}

// ValidateUpdates Function for validating salary update input
func ValidateUpdates(salary Payload, r *http.Request) error {
	return validation.ValidateStruct(&salary,
		validation.Field(&salary.Salary, ut.MoneyRule(r)...),
		validation.Field(&salary.Pension, ut.RequiredRule(r, "general.pension")...),
		validation.Field(&salary.Paye, ut.RequiredRule(r, "general.paye")...),
		validation.Field(&salary.Nsitf, ut.RequiredRule(r, "general.nsitf")...),
		validation.Field(&salary.Nhf, ut.RequiredRule(r, "general.nhf")...),
		validation.Field(&salary.Itf, ut.RequiredRule(r, "general.itf")...),
	)
}
