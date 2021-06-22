package repositories

import (
	"telebot-trading/app/models"
	"telebot-trading/external/db"
)

func GetAdminByEmail(email string) *models.Admin {
	admin := new(models.Admin)
	res := db.GetDB().Where("email = ?", email).First(admin)
	if res.Error != nil {
		return nil
	}
	return admin
}
