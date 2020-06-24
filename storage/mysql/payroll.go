package storage

import (
	"fmt"

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

// GetReports ...
func (pdb *DBPayrollStorage) GetReports() []*md.PayrollReport {
	data := []*md.PayrollReport{}

	err := pdb.db.Table("payrolls").
		Select("sum(gross_salary) as gross_salary_paid, sum(net_salary) as net_salary_paid, sum(pension) as pension_paid, sum(paye) as paye_paid, sum(nsitf) as nsitf_paid, sum(nhf) as nhf_paid, sum(itf) as itf_paid, month as month, year as year").
		Joins("JOIN taxes as tax ON tax.payroll_id = payrolls.id").
		Where("payrolls.payment_status=?", ut.Success).
		Group("payrolls.year, payrolls.month DESC").
		Scan(&data).Error

	if err != nil {
		fmt.Println(err)
	}

	return data
}

// GetAllMonthAndYear ... this is used in filtering data on the frontend
func (pdb *DBPayrollStorage) GetAllMonthAndYear() []*md.PayrollMonthYear {
	data := []*md.PayrollMonthYear{}

	err := pdb.db.Table("payrolls").
		Select("month as month, year as year").
		Group("payrolls.year, payrolls.month DESC").
		Scan(&data).Error

	if err != nil {
		fmt.Println(err)
	}

	return data
}

// Get payroll using employee_id and payroll_id
func (pdb *DBPayrollStorage) Get(employeeID, payrollID string) *md.Payroll {
	payroll := md.Payroll{}
	// Select resource from database
	err := pdb.db.
		Preload("Tax").
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
		Preload("Tax").
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
		Preload("Tax").
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
		Preload("Tax").
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
		Preload("Tax").
		Preload("Employee").
		Preload("Employee.Salary").
		Where("month=? AND year=?", month, year).Order("created_at desc").Find(&payrolls)

	return payrolls
}

// GetSingleByMonthYear ...
func (pdb *DBPayrollStorage) GetSingleByMonthYear(month, year uint) *md.Payroll {
	payroll := md.Payroll{}
	// Select resource from database
	err := pdb.db.
		Preload("Tax").
		Preload("Employee").
		Preload("Employee.Salary").
		Where("month=? AND year=?", month, year).Order("created_at desc").
		First(&payroll).Error

	if len(payroll.ID) < 1 || err != nil {
		return nil
	}

	return &payroll
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
