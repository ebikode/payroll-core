package model

// Salary - a struct to rep plan database model
type Salary struct {
	BaseIntModel
	EmployeeID string    `json:"employee_id" gorm:"not null;type:varchar(20)"`
	Salary     float64   `json:"salary" gorm:"type:float(10,2)"`
	Pension    float64   `json:"pension" gorm:"type:float(5,2)"`
	Paye       float64   `json:"paye" gorm:"type:float(5,2)"`
	Nsitf      float64   `json:"nsitf" gorm:"type:float(5,2)"`
	Nhf        float64   `json:"nhf" gorm:"type:float(5,2)"`
	Itf        float64   `json:"itf" gorm:"type:float(5,2)"`
	Employee   *Employee `json:"employee"`
}
