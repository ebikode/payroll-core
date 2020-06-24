package storage

import (
	"github.com/biezhi/gorm-paginator/pagination"
	md "github.com/ebikode/payroll-core/model"
)

// DBTaxStorage encapsulates DB Connection Model
type DBTaxStorage struct {
	*MDatabase
}

// NewDBTaxStorage Initialize Tax Storage
func NewDBTaxStorage(db *MDatabase) *DBTaxStorage {
	return &DBTaxStorage{db}
}

// Get Fetch Single Tax fron DB
func (cdb *DBTaxStorage) Get(id uint) *md.Tax {
	tax := md.Tax{}
	// Select resource from database
	err := cdb.db.
		Preload("Payroll").
		Preload("Payroll.Employee").
		Preload("Payroll.Employee.Salary").
		Joins("JOIN payrolls as pr ON pr.id = taxes.payroll_id").
		Where("taxes.id=?", id).First(&tax).Error

	if tax.ID < 1 || err != nil {
		return nil
	}

	return &tax
}

// GetAll Fetch all taxes from DB
func (cdb *DBTaxStorage) GetAll(page, limit int) []*md.Tax {
	var taxes []*md.Tax

	pagination.Paging(&pagination.Param{
		DB: cdb.db.
			Preload("Payroll").
			Preload("Payroll.Employee").
			Preload("Payroll.Employee.Salary").
			Joins("JOIN payrolls as pr ON pr.id = taxes.payroll_id").
			Order("created_at desc").
			Find(&taxes),
		Page:    page,
		Limit:   limit,
		OrderBy: []string{"created_at desc"},
	}, &taxes)

	return taxes

}

// GetEmployeeTaxes Fetch all employee' taxes from DB
func (cdb *DBTaxStorage) GetEmployeeTaxes(employeeID string) []*md.Tax {
	var taxes []*md.Tax

	cdb.db.
		Preload("Payroll").
		Preload("Payroll.Employee").
		Preload("Payroll.Employee.Salary").
		Joins("JOIN payrolls as pr ON pr.id = taxes.payroll_id").
		Where("employee_id=?", employeeID).
		Find(&taxes)
	return taxes
}

// Store Add a new tax
func (cdb *DBTaxStorage) Store(p md.Tax) (*md.Tax, error) {

	tax := p

	err := cdb.db.Create(&tax).Error

	if err != nil {
		return nil, err
	}
	return cdb.Get(tax.ID), nil
}

// Update a tax
func (cdb *DBTaxStorage) Update(tax *md.Tax) (*md.Tax, error) {

	err := cdb.db.Save(&tax).Error

	if err != nil {
		return nil, err
	}

	return tax, nil
}

// Delete a tax
func (cdb *DBTaxStorage) Delete(c md.Tax, isPermarnant bool) (bool, error) {

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
