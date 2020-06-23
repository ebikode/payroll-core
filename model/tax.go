package model

// Tax - a struct to rep plan database model
type Tax struct {
	BaseIntModel
	PayrollID string   `json:"payroll_id" gorm:"not null;type:varchar(20)"`
	Pension   float64  `json:"pension" gorm:"type:float(10,2)"`
	Paye      float64  `json:"paye" gorm:"type:float(10,2)"`
	Nsitf     float64  `json:"nsitf" gorm:"type:float(10,2)"`
	Nhf       float64  `json:"nhf" gorm:"type:float(10,2)"`
	Itf       float64  `json:"Itf" gorm:"type:float(10,2)"`
	Payroll   *Payroll `json:"payroll"`
}
