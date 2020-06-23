package payroll

import (
	md "github.com/ebikode/payroll-core/model"
)

// Payload Request data
type Payload struct {
	Month         int    `json:"month,omitempty"`
	Year          int    `json:"year,omitempty"`
	Status        string `json:"status,omitempty"`
	PaymentStatus string `json:"payment_status,omitempty"`
}

// ValidationFields  Error response data
type ValidationFields struct {
	Month         string `json:"month,omitempty"`
	Year          string `json:"year,omitempty"`
	Status        string `json:"status,omitempty"`
	PaymentStatus string `json:"account_id,omitempty"`
}

// PayrollRepository provides access to the Payroll storage.
type PayrollRepository interface {
	// Get returns the payroll with given ID.
	Get(string, string) *md.Payroll
	GetLastPayroll() *md.Payroll
	// returns all payrolls set with page and limit.
	GetAll(int, int) []*md.Payroll
	// returns the payrolls with given employeeID.
	GetByEmployee(string, int, int) []*md.Payroll
	GetByMonthYear(uint, uint) []*md.Payroll
	// Store a given employee payroll to the repository.
	Store(md.Payroll) (*md.Payroll, error)
	// Update a given payroll in the repository.
	Update(*md.Payroll) (*md.Payroll, error)
	UpdateStatus(string, int, int)
	UpdatePaymentStatus(string, int, int)
	// Delete a given payroll in the repository.
	Delete(*md.Payroll, bool) (bool, error)
}
