package storage

import (
	"fmt"

	"github.com/biezhi/gorm-paginator/pagination"
	md "github.com/ebikode/payroll-core/model"
	ut "github.com/ebikode/payroll-core/utils"
)

// DBEmployeeStorage ...
type DBEmployeeStorage struct {
	*MDatabase
}

// NewDBEmployeeStorage Initialize Employee Storage
func NewDBEmployeeStorage(db *MDatabase) *DBEmployeeStorage {
	return &DBEmployeeStorage{db}
}

// GetDashbordData ...
func (edb *DBEmployeeStorage) GetDashbordData(employeeID string) *md.EmployeeDashbordData {
	data := md.EmployeeDashbordData{}

	// Sum GrossSalaryPaid
	err := edb.db.Table("payrolls").
		Select("sum(gross_salary) as gross_salary_earned").
		Where("payment_status=? AND employee_id=?", ut.Success, employeeID).
		Scan(&data).Error

	// Sum NetSalaryPaid
	err = edb.db.Table("payrolls").
		Select("sum(net_salary) as net_salary_earned").
		Where("payment_status=? AND employee_id=?", ut.Success, employeeID).
		Scan(&data).Error

	// Sum PensionPaid
	err = edb.db.Table("payrolls").
		Select("sum(pension) as pension_paid").
		Joins("JOIN taxes as tax ON tax.payroll_id = payrolls.id").
		Where("payment_status=? AND employee_id=?", ut.Success, employeeID).
		Scan(&data).Error

	// Sum PayePaid
	err = edb.db.Table("payrolls").
		Select("sum(paye) as paye_paid").
		Joins("JOIN taxes as tax ON tax.payroll_id = payrolls.id").
		Where("payment_status=? AND employee_id=?", ut.Success, employeeID).
		Scan(&data).Error

	// Sum NsitfPaid
	err = edb.db.Table("payrolls").
		Select("sum(nsitf) as nsitf_paid").
		Joins("JOIN taxes as tax ON tax.payroll_id = payrolls.id").
		Where("payment_status=? AND employee_id=?", ut.Success, employeeID).
		Scan(&data).Error

	// Sum NhfPaid
	err = edb.db.Table("payrolls").
		Select("sum(nhf) as nhf_paid").
		Joins("JOIN taxes as tax ON tax.payroll_id = payrolls.id").
		Where("payment_status=? AND employee_id=?", ut.Success, employeeID).
		Scan(&data).Error

	// Sum ItfPaid
	err = edb.db.Table("payrolls").
		Select("sum(itf) as itf_paid").
		Joins("JOIN taxes as tax ON tax.payroll_id = payrolls.id").
		Where("payment_status=? AND employee_id=?", ut.Success, employeeID).
		Scan(&data).Error

	if err != nil {
		fmt.Println(err)
	}

	return &data
}

// Authenticate a employee
func (edb *DBEmployeeStorage) Authenticate(email string) (*md.Employee, error) {
	employee := &md.Employee{}
	q := edb.db.
		Preload("Salary")

	err := q.Where("email = ? AND is_email_verified = ?", email, true).First(&employee).Error

	if employee.ID == "" || err != nil {
		return nil, err
	}
	return employee, nil
}

// Get employee by ID
func (edb *DBEmployeeStorage) Get(id string) *md.Employee {
	employee := md.Employee{}
	// Select Employee
	err := edb.db.
		Preload("Salary").
		Where("id=?", id).
		First(&employee).Error

	if employee.ID == "" || err != nil {
		return nil
	}

	return &employee
}

//GetPubEmployee Gets public employee by ID
func (edb *DBEmployeeStorage) GetPubEmployee(id string) *md.PubEmployee {
	employee := md.PubEmployee{}
	// Select Employee
	err := edb.db.
		Preload("Salary").
		Where("id=?", id).
		First(&employee).Error

	if employee.ID == "" || err != nil {
		return nil
	}

	return &employee
}

//GetActivePubEmployees Gets public employee by ID
func (edb *DBEmployeeStorage) GetActivePubEmployees() []*md.PubEmployee {
	employee := []*md.PubEmployee{}
	// Select Employee
	edb.db.
		Preload("Salary").
		Where("status=?", ut.Active).
		Find(&employee)

	return employee
}

//GetEmployeeByEmail Gets employee by Email
func (edb *DBEmployeeStorage) GetEmployeeByEmail(email string) *md.Employee {
	employee := md.Employee{}

	// Select Employee
	err := edb.db.
		Preload("Salary").
		Where("email=?", email).
		First(&employee).Error

	if employee.ID == "" || err != nil {
		return nil
	}

	return &employee
}

// GetEmployees gets all employees with paging using @param page and limit
func (edb *DBEmployeeStorage) GetEmployees(page, limit int) []*md.PubEmployee {
	var employees []*md.PubEmployee

	q := edb.db.
		Preload("Salary")

	pagination.Paging(&pagination.Param{
		DB:      q.Find(&employees),
		Page:    page,
		Limit:   limit,
		OrderBy: []string{"created_at desc"},
	}, &employees)

	return employees
}

//Store a new employee
func (edb *DBEmployeeStorage) Store(u md.Employee) (md.Employee, error) {

	usr := &u

	err := edb.db.Create(&usr).Error

	if err != nil {
		return u, err
	}
	return u, nil
}

// Update a employee
func (edb *DBEmployeeStorage) Update(u *md.Employee) (*md.Employee, error) {

	err := edb.db.Save(&u).Error

	if err != nil {
		return u, err
	}

	return u, nil
}

// Delete a employee
func (edb *DBEmployeeStorage) Delete(u md.Employee, isPermarnant bool) (bool, error) {

	var err error
	if isPermarnant {
		err = edb.db.Unscoped().Delete(u).Error
	}
	if !isPermarnant {
		err = edb.db.Delete(u).Error
	}

	if err != nil {
		return false, err
	}

	return true, nil
}
