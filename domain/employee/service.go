package employee

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	md "github.com/ebikode/payroll-core/model"
	tr "github.com/ebikode/payroll-core/translation"
	ut "github.com/ebikode/payroll-core/utils"
	validation "github.com/go-ozzo/ozzo-validation"
	"golang.org/x/crypto/bcrypt"
)

// EmployeeService provides employee operations
type EmployeeService interface {
	GetEmployeeDashboardData(string) *md.EmployeeDashbordData
	GetEmployee(string) *md.Employee
	GetPubEmployee(string) *md.PubEmployee
	GetAllActivePubEmployee() []*md.PubEmployee
	GetEmployeeByEmail(string) *md.Employee
	GetAllEmployees(int, int) []*md.PubEmployee
	AuthenticateEmployee(string, string) (*md.Employee, tr.TParam, error)
	CreateEmployee(md.Employee) (md.Employee, tr.TParam, error)
	UpdateEmployee(*md.Employee) (*md.Employee, tr.TParam, error)
}

type service struct {
	employeeRepo EmployeeRepository
}

// NewService creates a employee service with the necessary dependencies
func NewService(
	employeeRepo EmployeeRepository,
) EmployeeService {
	return &service{employeeRepo}
}

// Authenticate a employee
func (s *service) AuthenticateEmployee(email string, password string) (*md.Employee, tr.TParam, error) {
	tParam := tr.TParam{
		Key:          "error.login_error",
		TemplateData: nil,
		PluralCount:  nil,
	}

	employee, err := s.employeeRepo.Authenticate(email)
	if err != nil {
		return nil, tParam, errors.New("Error")
	}

	isPasswordValid := ut.ValidatePassword(employee.Password, password)
	if !isPasswordValid {
		return nil, tParam, errors.New("Error")
	}

	emp := s.GetEmployee(employee.ID)

	return emp, tParam, nil

}

func (s *service) GetEmployeeDashboardData(id string) *md.EmployeeDashbordData {
	return s.employeeRepo.GetDashbordData(id)
}

// Get a employee
func (s *service) GetEmployee(id string) *md.Employee {
	return s.employeeRepo.Get(id)
}

// Get a public employee
func (s *service) GetPubEmployee(id string) *md.PubEmployee {
	return s.employeeRepo.GetPubEmployee(id)
}

// Get a public employee
func (s *service) GetAllActivePubEmployee() []*md.PubEmployee {
	return s.employeeRepo.GetActivePubEmployees()
}

// GetEmployeeByEmail Get a employee using their phone number
func (s *service) GetEmployeeByEmail(email string) *md.Employee {
	return s.employeeRepo.GetEmployeeByEmail(email)
}

// Get md.Employee
func (s *service) GetAllEmployees(page, limit int) []*md.PubEmployee {
	return s.employeeRepo.GetEmployees(page, limit)
}

// create new employee
func (s *service) CreateEmployee(u md.Employee) (md.Employee, tr.TParam, error) {

	token := ut.RandomBase64String(30, "")
	tokenCopy := token

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	hashedToken, _ := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)

	u.Password = string(hashedPassword)
	u.EmailToken = string(hashedToken)
	u.PincodeSentAt = time.Now().UTC()

	// Generate employee ID
	uID := ut.RandomBase64String(8, "pxpu")

	u.ID = uID

	employee, err := s.employeeRepo.Store(u)

	fmt.Printf("Error:: %s \n", err)

	if err != nil {
		tParam := tr.TParam{
			Key:          "error.resource_creation_error",
			TemplateData: nil,
			PluralCount:  nil,
		}

		return employee, tParam, err
	}

	employee.EmailToken = tokenCopy

	return employee, tr.TParam{}, nil
}

// update existing employee
func (s *service) UpdateEmployee(u *md.Employee) (*md.Employee, tr.TParam, error) {

	employee, err := s.employeeRepo.Update(u)

	if err != nil {
		tParam := tr.TParam{
			Key:          "error.resource_update_error",
			TemplateData: nil,
			PluralCount:  nil,
		}

		return employee, tParam, err
	}

	return employee, tr.TParam{}, nil

}

// Validate - Function for validating employee input during creation
func Validate(employee md.Employee, r *http.Request) error {
	return validation.ValidateStruct(&employee,
		// Phone cannot be empty, and th length must between 7 and 20
		validation.Field(&employee.Phone, ut.PhoneRule(r)...),
		validation.Field(&employee.Email, ut.EmailRule(r)...),
		validation.Field(&employee.FirstName, ut.NameRule(r)...),
		validation.Field(&employee.LastName, ut.NameRule(r)...),
		validation.Field(&employee.Username, ut.NameRule(r)...),
		validation.Field(&employee.Password, ut.PasswordRule(r)...),
	)
}

// ValidateUpdates - Function for validating employee input during update
func ValidateUpdates(employee md.Employee, r *http.Request) error {
	return validation.ValidateStruct(&employee,
		validation.Field(&employee.Email, ut.EmailRule(r)...),
		validation.Field(&employee.FirstName, ut.NameRule(r)...),
		validation.Field(&employee.LastName, ut.NameRule(r)...),
		validation.Field(&employee.Username, ut.NameRule(r)...),
		validation.Field(&employee.Avatar, ut.AvatarRule(r)...),
		validation.Field(&employee.Thumb, ut.AvatarRule(r)...),
	)
}
