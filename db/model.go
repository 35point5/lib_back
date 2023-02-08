package db

import (
	"gorm.io/gorm"
	"time"
)

const BorrowLimit = 5
const CardLimit = 2
const Admin = 2
const Guest = 1

type Tabel_infor struct {
	Field   string
	Type    string
	Null    string
	Key     string
	Default string
	Extra   string
}

type User struct {
	gorm.Model
	Name     string `gorm:"unique;not null"`
	Role     int
	Limit    int
	Password string
	Cookie   string `gorm:"unique;not null"`
}

type Book struct {
	ISBN        string `gorm:"primarykey"`
	Title       string
	Author      string
	Publisher   string
	KeyWords    string
	Digest      string
	Category    string
	PublishTime string
	Remain      int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type Borrow struct {
	CardID    uint   `gorm:"primarykey"`
	ISBN      string `gorm:"primarykey"`
	Number    int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Card struct {
	gorm.Model
	Limit  int
	UserID uint
}
