package salary

import (
	md "github.com/ebikode/payroll-core/model"
)

// Payload ...
type Payload struct {
	EmployeeID string  `json:"employee_id"`
	Salary     float64 `json:"salary"`
	Pension    float64 `json:"pension"`
	Paye       float64 `json:"paye"`
	Nsitf      float64 `json:"nsitf"`
	Nhf        float64 `json:"nhf"`
	Itf        float64 `json:"Itf"`
}

// ValidationFields struct to return for validation
type ValidationFields struct {
	EmployeeID string `json:"employee_id,omitempty"`
	Salary     string `json:"salary,omitempty"`
	Pension    string `json:"pension,omitempty"`
	Paye       string `json:"paye,omitempty"`
	Nsitf      string `json:"nsitf,omitempty"`
	Nhf        string `json:"nhf,omitempty"`
	Itf        string `json:"Itf,omitempty"`
}

// SalaryRepository  provides access to the md.Salary storage.
type SalaryRepository interface {
	// Get returns the plan with given ID.
	Get(uint) *md.Salary
	GetByEmployeeID(string) *md.Salary
	// Get returns all salaries.
	GetAll(int, int) []*md.Salary
	// Store a given plan to the repository.
	Store(md.Salary) (*md.Salary, error)
	// Update a given plan in the repository.
	Update(*md.Salary) (*md.Salary, error)
	// Delete a given plan in the repository.
	Delete(md.Salary, bool) (bool, error)
}
