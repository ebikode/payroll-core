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
	Tax           *Tax      `json:"tax" gorm:"foreignkey:PayrollID"`
}

// PayrollReport a struct to rep payroll monthly reports
type PayrollReport struct {
	Year            string  `json:"year"`
	Month           string  `json:"month"`
	GrossSalaryPaid float64 `json:"gross_salary_paid"`
	NetSalaryPaid   float64 `json:"net_salary_paid"`
	PensionPaid     float64 `json:"pension_paid"`
	PayePaid        float64 `json:"paye_paid"`
	NsitfPaid       float64 `json:"nsitf_paid"`
	NhfPaid         float64 `json:"nhf_paid"`
	ItfPaid         float64 `json:"itf_paid"`
}

// PayrollMonthYear a struct to rep payroll existing month and year combo
type PayrollMonthYear struct {
	Year  string `json:"year"`
	Month string `json:"month"`
}
