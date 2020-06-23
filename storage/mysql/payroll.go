package storage

import (
	"github.com/biezhi/gorm-paginator/pagination"
	md "github.com/ebikode/payroll-core/model"
	ut "github.com/ebikode/payroll-core/utils"
)

type DBPayrollStorage struct {
	*MDatabase
}

// Initialize Payroll Storage
func NewDBPayrollStorage(db *MDatabase) *DBPayrollStorage {
	return &DBPayrollStorage{db}
}

// Get payroll using employee_id and payroll_id
func (pdb *DBPayrollStorage) Get(employeeID, payrollID string) *md.Payroll {
	payroll := md.Payroll{}
	// Select resource from database
	err := pdb.db.
		Preload("Employee").
		Preload("Employee.Salary").
		Where("employee_id=? AND id=?", employeeID, payrollID).First(&payroll).Error

	if len(payroll.ID) < 1 || err != nil {
		return nil
	}

	return &payroll
}

// GetLastPayroll ...
func (pdb *DBPayrollStorage) GetLastPayroll() *md.Payroll {
	payroll := md.Payroll{}
	// Select resource from database
	err := pdb.db.
		Preload("Employee").
		Preload("Employee.Salary").
		Order("created_at").
		Limit(1).First(&payroll).Error

	if len(payroll.ID) < 1 || err != nil {
		return nil
	}

	return &payroll
}

// GetAll Get all payrolls
func (pdb *DBPayrollStorage) GetAll(page, limit int) []*md.Payroll {
	var payrolls []*md.Payroll
	// Select resource from database
	q := pdb.db.
		Preload("Employee").
		Preload("Employee.Salary")

	pagination.Paging(&pagination.Param{
		DB:      q.Order("created_at desc").Find(&payrolls),
		Page:    page,
		Limit:   limit,
		OrderBy: []string{"created_at desc"},
	}, &payrolls)

	return payrolls
}

// GetByEmployee Get all payrolls of a employee  form DB
func (pdb *DBPayrollStorage) GetByEmployee(employeeID string, page, limit int) []*md.Payroll {
	var payrolls []*md.Payroll
	// Select resource from database
	q := pdb.db.
		Preload("Employee").
		Preload("Employee.Salary")

	pagination.Paging(&pagination.Param{
		DB:      q.Where("employee_id=?", employeeID).Order("created_at desc").Find(&payrolls),
		Page:    page,
		Limit:   limit,
		OrderBy: []string{"created_at desc"},
	}, &payrolls)
	return payrolls
}

// GetByMonthYear ...
func (pdb *DBPayrollStorage) GetByMonthYear(month, year uint) []*md.Payroll {
	var payrolls []*md.Payroll
	// Select resource from database
	pdb.db.
		Preload("Employee").
		Preload("Employee.Salary").
		Where("month=? AND year=?", month, year).Order("created_at desc").Find(&payrolls)

	return payrolls
}

// Store Add a new payroll
func (pdb *DBPayrollStorage) Store(p md.Payroll) (*md.Payroll, error) {

	py := p

	err := pdb.db.Create(&py).Error

	if err != nil {
		return nil, err
	}
	return pdb.Get(py.EmployeeID, py.ID), nil
}

// Update a payroll
func (pdb *DBPayrollStorage) Update(payroll *md.Payroll) (*md.Payroll, error) {

	err := pdb.db.Save(&payroll).Error

	if err != nil {
		return nil, err
	}

	return payroll, nil
}

// UpdateStatus ...
func (pdb *DBPayrollStorage) UpdateStatus(status string, month, year int) {
	pdb.db.Model(md.Payroll{}).
		Where("status=? AND month=? AND year=?", ut.Pending, month, year).
		Updates(md.Payroll{Status: status})
	return
}

// UpdatePaymentStatus ...
func (pdb *DBPayrollStorage) UpdatePaymentStatus(status string, month, year int) {
	pdb.db.Model(md.Payroll{}).
		Where("payment_status=? AND month=? AND year=?", ut.Pending, month, year).
		Updates(md.Payroll{PaymentStatus: status})
	return
}

// Delete a payroll
func (pdb *DBPayrollStorage) Delete(p *md.Payroll, isPermarnant bool) (bool, error) {

	// var err error
	// if isPermarnant {
	// 	err = pdb.db.Unscoped().Delete(p).Error
	// }
	// if !isPermarnant {
	// 	err = pdb.db.Delete(p).Error
	// }

	// if err != nil {
	// 	return false, err
	// }

	return true, nil
}
