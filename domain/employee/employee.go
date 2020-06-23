package employee

import (
	md "github.com/ebikode/payroll-core/model"
)

// ValidationFields struct to return for validation
type ValidationFields struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Password  string `json:"password,omitempty"`
	Username  string `json:"username,omitempty"`
	Email     string `json:"email,omitempty"`
	Avatar    string `json:"avatar,omitempty"`
	Thumb     string `json:"thumb,omitempty"`
}

//EmployeeRepository Repository provides access to the employee storage.
type EmployeeRepository interface {
	GetDashbordData(string) *md.EmployeeDashbordData
	// Get returns the employee with given ID.
	Get(string) *md.Employee
	// returns the public employee with given ID.
	GetPubEmployee(string) *md.PubEmployee
	GetActivePubEmployees() []*md.PubEmployee
	// GetEmployeeByPhone returns the employee with given phone number.
	GetEmployeeByEmail(string) *md.Employee
	// GetEmployees returns all employees paginated.
	GetEmployees(int, int) []*md.PubEmployee
	// Authenticate a employee
	Authenticate(string) (*md.Employee, error)
	// Saves a given employee to the repository.
	Store(md.Employee) (md.Employee, error)
	// Update a given employee in the repository.
	Update(*md.Employee) (*md.Employee, error)
	// Delete a given employee in the repository.
	Delete(md.Employee, bool) (bool, error)
}
