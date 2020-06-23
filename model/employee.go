package model

import "time"

//Employee - a struct to rep Employee account
type Employee struct {
	BaseModel
	Lang            string    `json:"lang" gorm:"not null;type:varchar(5);default:'en'"`
	FirstName       string    `json:"first_name" gorm:"not null;type:varchar(50);"`
	LastName        string    `json:"last_name" gorm:"not null;type:varchar(50);"`
	Position        string    `json:"position" gorm:"not null;type:varchar(100);"`
	Username        string    `json:"username" gorm:"not null;type:varchar(20);unique_index"`
	Address         string    `json:"address" gorm:"not null;type:varchar(150)"`
	About           string    `json:"about" gorm:"type:varchar(150)"`
	Email           string    `json:"email" gorm:"not null,type:varchar(100);unique_index"`
	EmailToken      string    `json:"email_token" gorm:"type:varchar(200)"`
	Password        string    `json:"password,omitempty" gorm:"not null;type:varchar(250)"`
	Pincode         string    `json:"pincode,omitempty" gorm:"not null;type:varchar(250)"`
	Phone           string    `json:"phone" gorm:"not null;type: varchar(20);unique_index"`
	AccountName     string    `json:"account_name" gorm:"not null;type:varchar(100)"`
	AccountNumber   string    `json:"account_number" gorm:"not null;type:int(10)"`
	BankName        string    `json:"bank_name" gorm:"not null;type: varchar(100)"`
	IsPincodeUsed   bool      `json:"is_pincode_used" gorm:"type:tinyint(1);default:0"`
	PincodeSentAt   time.Time `json:"pincode_sent_at"`
	IsPhoneVerified bool      `json:"is_phone_verified" gorm:"type:tinyint(1);default:0"`
	IsEmailVerified bool      `json:"is_email_verified" gorm:"type:tinyint(1);default:0"`
	Avatar          string    `json:"avatar" gorm:"type:varchar(500)"`
	Thumb           string    `json:"thumb" gorm:"type:varchar(500)"`
	Status          string    `json:"status" gorm:"type:enum('pending','active','suspended','resigned','fired','deleted');default:'pending'"`
	// Added for request body validation only
	Token  string  `json:"token,omitempty" gorm:"-"`
	Salary *Salary `json:"salary"  gorm:"foreignkey:EmployeeID"`
}

//PubEmployee - a struct to rep Employee account shown to others
// e.g Admin
type PubEmployee struct {
	BaseModel
	Lang            string  `json:"lang" gorm:"not null;type:varchar(5);default:'en'"`
	FirstName       string  `json:"first_name" gorm:"not null;type:varchar(50);unique_index"`
	LastName        string  `json:"last_name" gorm:"not null;type:varchar(50);unique_index"`
	Username        string  `json:"username" gorm:"not null;type:varchar(20);unique_index"`
	Address         string  `json:"address" gorm:"not null;type:varchar(150)"`
	About           string  `json:"about" gorm:"type:varchar(150)"`
	Email           string  `json:"email" gorm:"type:varchar(100)"`
	Phone           string  `json:"phone" gorm:"not null;type: varchar(20);unique_index"`
	AccountName     string  `json:"account_name" gorm:"not null;type:varchar(100)"`
	AccountNumber   string  `json:"account_number" gorm:"not null;type:int(10)"`
	BankName        string  `json:"bank_name" gorm:"not null;type: varchar(100)"`
	IsPhoneVerified bool    `json:"is_phone_verified" gorm:"type:tinyint(1);default:0"`
	IsEmailVerified bool    `json:"is_email_verified" gorm:"type:tinyint(1);default:0"`
	Avatar          string  `json:"avatar" gorm:"type:varchar(500)"`
	Thumb           string  `json:"thumb" gorm:"type:varchar(500)"`
	Status          string  `json:"status" gorm:"type:enum('pending','active','suspended','deleted');default:'pending'"`
	Salary          *Salary `json:"salary"  gorm:"foreignkey:EmployeeID"`
}

// EmployeeDashbordData - struct encapsulating admin dashboard data
type EmployeeDashbordData struct {
	GrossSalaryEarned float64 `json:"gross_salary_earned"`
	NetSalaryEarned   float64 `json:"net_salary_earned"`
	PensionPaid       float64 `json:"pension_paid"`
	PayePaid          float64 `json:"paye_paid"`
	NsitfPaid         float64 `json:"nsitf_paid"`
	NhfPaid           float64 `json:"nhf_paid"`
	ItfPaid           float64 `json:"itf_paid"`
}

// TableName Set PubEmployee's table name to be `Employees`
func (PubEmployee) TableName() string {
	return "employees"
}
