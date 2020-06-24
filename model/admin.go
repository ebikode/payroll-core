package model

//Admin - a struct to rep admin account
type Admin struct {
	BaseModel
	FirstName string `json:"first_name" gorm:"not null;type:varchar(20)"`
	LastName  string `json:"last_name" gorm:"not null;type:varchar(20)"`
	Email     string `json:"email" gorm:"type:varchar(100);unique_index"`
	Password  string `json:"password,omitempty" gorm:"not null;type:varchar(250)"`
	Phone     string `json:"phone" gorm:"not null;type: varchar(20);unique_index"`
	Avatar    string `json:"avatar" gorm:"type:varchar(500)"`
	Thumb     string `json:"thumb" gorm:"type:varchar(500)"`
	Role      string `json:"role" gorm:"type:enum('manager','super_admin');default:'manager'"`
	Status    string `json:"status" gorm:"type:enum('pending','active','suspended','deleted');default:'pending'"`
	Token     string `json:"token,omitempty" gorm:"-"`
}

// DashbordData - struct encapsulating admin dashboard data
type DashbordData struct {
	EmployeesCount        int64   `json:"employees_count"`
	ActiveEmployeesCount  int64   `json:"active_employees_count"`
	PendingEmployeesCount int64   `json:"pending_employees_count"`
	GrossSalaryPaid       float64 `json:"gross_salary_paid"`
	NetSalaryPaid         float64 `json:"net_salary_paid"`
	PensionPaid           float64 `json:"pension_paid"`
	PayePaid              float64 `json:"paye_paid"`
	NsitfPaid             float64 `json:"nsitf_paid"`
	NhfPaid               float64 `json:"nhf_paid"`
	ItfPaid               float64 `json:"itf_paid"`
}
