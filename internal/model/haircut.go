package model

import "time"

type Haircut struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"not null"`
	Description string
	Style       string             `gorm:"type:varchar(50)"`
	Length      string             `gorm:"type:varchar(20)"`
	Popularity  int                `gorm:"default:0"`
	FaceShapes  []HaircutFaceShape `gorm:"foreignKey:HaircutID"`
	HairTypes   []HaircutHairType  `gorm:"foreignKey:HaircutID"`
	Images      []HaircutImage     `gorm:"foreignKey:HaircutID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `gorm:"index"`
}

type HaircutFaceShape struct {
	HaircutID uint   `gorm:"primaryKey"`
	FaceShape string `gorm:"primaryKey;type:varchar(50)"`
}

type HaircutHairType struct {
	HaircutID uint   `gorm:"primaryKey"`
	HairType  string `gorm:"primaryKey;type:varchar(50)"`
}

type HaircutImage struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	HaircutID uint   `gorm:"not null"`
	URL       string `gorm:"not null"`
	Angle     string `gorm:"type:varchar(50)"`
	IsMain    bool   `gorm:"default:false"`
	CreatedAt time.Time
}
