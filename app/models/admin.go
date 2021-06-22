package models

import (
	"telebot-trading/app/helper"
	"time"
)

const ROLE_SUPER = "superadmin"
const ROLE_ADMIN = "admin"

type Admin struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:200"`
	Email     string `gorm:"unique"`
	Password  string `gorm:"size:200"`
	Role      string `gorm:"size:200"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (admin *Admin) VerifyPassword(pass string) bool {
	return helper.CheckHash(pass, admin.Password)
}
