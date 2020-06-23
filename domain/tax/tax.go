package tax

import (
	md "github.com/ebikode/payroll-core/model"
)

// TaxRepository  provides access to the md.Tax storage.
type TaxRepository interface {
	// Get returns the tax with given ID.
	Get(uint) *md.Tax
	// Get returns all taxes.
	GetAll(int, int) []*md.Tax
	// GetEmployeeTaxes returns all taxes of a employee.
	GetEmployeeTaxes(string) []*md.Tax
	// Store a given tax to the repository.
	Store(md.Tax) (*md.Tax, error)
	// Update a given tax in the repository.
	Update(*md.Tax) (*md.Tax, error)
	// Delete a given tax in the repository.
	Delete(md.Tax, bool) (bool, error)
}
