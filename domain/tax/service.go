package tax

import (
	md "github.com/ebikode/payroll-core/model"
	tr "github.com/ebikode/payroll-core/translation"
)

// TaxService  provides tax operations
type TaxService interface {
	GetTax(uint) *md.Tax
	GetTaxes(int, int) []*md.Tax
	GetEmployeeTaxes(string) []*md.Tax
	CreateTax(md.Tax) (*md.Tax, tr.TParam, error)
	UpdateTax(*md.Tax) (*md.Tax, tr.TParam, error)
}

type service struct {
	txRepo TaxRepository
}

// NewService creates a tax service with the necessary dependencies
func NewService(
	txRepo TaxRepository,
) TaxService {
	return &service{txRepo}
}

// Get a tax
func (s *service) GetTax(id uint) *md.Tax {
	return s.txRepo.Get(id)
}

// GetTaxes Get all taxs from DB
//
// @userType == admin | employee
func (s *service) GetTaxes(page, limit int) []*md.Tax {
	return s.txRepo.GetAll(page, limit)
}

// GetEmployeeTaxes Get all taxs from DB
//
// @userType == admin | employee
func (s *service) GetEmployeeTaxes(employeeID string) []*md.Tax {
	return s.txRepo.GetEmployeeTaxes(employeeID)
}

// CreateTax Creates New tax
func (s *service) CreateTax(c md.Tax) (*md.Tax, tr.TParam, error) {

	tax, err := s.txRepo.Store(c)

	if err != nil {
		tParam := tr.TParam{
			Key:          "error.resource_creation_error",
			TemplateData: nil,
			PluralCount:  nil,
		}

		return tax, tParam, err
	}

	return tax, tr.TParam{}, nil

}

// UpdateTax update existing tax
func (s *service) UpdateTax(c *md.Tax) (*md.Tax, tr.TParam, error) {
	tax, err := s.txRepo.Update(c)

	if err != nil {
		tParam := tr.TParam{
			Key:          "error.resource_update_error",
			TemplateData: nil,
			PluralCount:  nil,
		}

		return tax, tParam, err
	}

	return tax, tr.TParam{}, nil

}
