package model

// Payroll a struct encapsulating payment
type Payroll struct {
	BaseModel
	EmployeeID    string    `json:"employee_id" gorm:"not null;type:varchar(20)"`
	ReferenceNo   string    `json:"reference_no" gorm:"not null;type:varchar(100)"`
	GrossSalary   float64   `json:"gross_salary" gorm:"type:float(10,2)"`
	NetSalary     float64   `json:"net_salary" gorm:"type:float(10,2)"`
	Currency      string    `json:"currency" gorm:"type:enum('NGN');default:'NGN'"`
	PaymentMethod string    `json:"payment_method" gorm:"type:enum('bank_transfer');default:'bank_transfer'"`
	Month         uint      `json:"month" gorm:"not null;type:int(2)"`
	Year          uint      `json:"year" gorm:"not null;type:int(4)"`
	PaymentStatus string    `json:"payment_status" gorm:"type:enum('success', 'failed', 'pending');default:'pending'"`
	Status        string    `json:"status" gorm:"type:enum('approved', 'rejected', 'pending');default:'pending'"`
	Employee      *Employee `json:"employee"`
}
