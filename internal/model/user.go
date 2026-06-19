package model

import (
	"time"
)

type Role string

const (
	RoleUser   Role = "user"
	RoleBarber Role = "barber"
	RoleAdmin  Role = "admin"
)

type User struct {
	ID           uint   `gorm:"primaryKey;autoIncrement"`
	Email        string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
	Name         string `gorm:"not null"`
	Role         Role   `gorm:"type:varchar(20);default:'user'"`
	FaceShape    string `gorm:"type:varchar(50)"`
	HairType     string `gorm:"type:varchar(50)"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time `gorm:"index"`
}
