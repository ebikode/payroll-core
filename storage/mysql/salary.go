package storage

import (
	"github.com/biezhi/gorm-paginator/pagination"
	md "github.com/ebikode/payroll-core/model"
)

// DBSalaryStorage encapsulates DB Connection Model
type DBSalaryStorage struct {
	*MDatabase
}

// NewDBSalaryStorage Initialize Salary Storage
func NewDBSalaryStorage(db *MDatabase) *DBSalaryStorage {
	return &DBSalaryStorage{db}
}

// Get Fetch Single Salary fron DB
func (cdb *DBSalaryStorage) Get(id uint) *md.Salary {
	salary := md.Salary{}
	// Select resource from database
	err := cdb.db.
		Preload("Employee").
		Where("salaries.id=?", id).First(&salary).Error

	if salary.ID < 1 || err != nil {
		return nil
	}

	return &salary
}

// GetByEmployeeID Fetch Single Salary fron DB
func (cdb *DBSalaryStorage) GetByEmployeeID(id string) *md.Salary {
	salary := md.Salary{}
	// Select resource from database
	err := cdb.db.
		Preload("Employee").
		Where("employee_id=?", id).First(&salary).Error

	if salary.ID < 1 || err != nil {
		return nil
	}

	return &salary
}

// GetAll Fetch all salaries from DB
func (cdb *DBSalaryStorage) GetAll(page, limit int) []*md.Salary {
	var salaries []*md.Salary

	pagination.Paging(&pagination.Param{
		DB: cdb.db.
			Preload("Employee").
			Order("created_at desc").
			Find(&salaries),
		Page:    page,
		Limit:   limit,
		OrderBy: []string{"created_at desc"},
	}, &salaries)

	return salaries

}

// GetEmployeeSalaries Fetch all employee' salaries from DB
func (cdb *DBSalaryStorage) GetEmployeeSalaries(employeeID string) []*md.Salary {
	var salaries []*md.Salary

	cdb.db.
		Preload("Employee").
		Where("employee_id=?", employeeID).
		Find(&salaries)
	return salaries
}

// Store Add a new salary
func (cdb *DBSalaryStorage) Store(p md.Salary) (*md.Salary, error) {

	salary := p

	err := cdb.db.Create(&salary).Error

	if err != nil {
		return nil, err
	}
	return cdb.Get(salary.ID), nil
}

// Update a salary
func (cdb *DBSalaryStorage) Update(salary *md.Salary) (*md.Salary, error) {

	err := cdb.db.Save(&salary).Error

	if err != nil {
		return nil, err
	}

	return salary, nil
}

// Delete a salary
func (cdb *DBSalaryStorage) Delete(c md.Salary, isPermarnant bool) (bool, error) {

	var err error
	if isPermarnant {
		err = cdb.db.Unscoped().Delete(c).Error
	}
	if !isPermarnant {
		err = cdb.db.Delete(c).Error
	}

	if err != nil {
		return false, err
	}

	return true, nil
}
