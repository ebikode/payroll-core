package model

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// BaseModel  Model base model definition, including fields `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`, which could be embedded in your models
//    type Employee struct {
//      BaseModel
//    }
type BaseModel struct {
	ID        string     `json:"id" gorm:"primary_key;type:varchar(20)"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" sql:"index"`
}

// BaseIntModel  definition, including fields `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`, which could be embedded in models
// This is for IDs that are integer
//    type Employee struct {
//      BaseIntModel
//    }
type BaseIntModel struct {
	ID        uint       `json:"id" gorm:"type:int(15) unsigned auto_increment;not null;primary_key"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" sql:"index"`
}

/*
EmployeeTokenData JWT claims struct for customer
*/
type EmployeeTokenData struct {
	EmployeeID string    `json:"customer_id"`
	DeviceID   string    `json:"device_id"`
	Username   string    `json:"username"`
	ExpireOn   time.Time `json:"expire_on"`
	jwt.StandardClaims
}

/*
AdminTokenData JWT claims struct for admin
*/
type AdminTokenData struct {
	AdminID  string    `json:"admin_id"`
	Role     string    `json:"role"`
	IP       string    `json:"ip"`
	ExpireOn time.Time `json:"expire_on"`
	jwt.StandardClaims
}
