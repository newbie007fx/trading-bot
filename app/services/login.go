package services

import (
	"errors"
	"telebot-trading/app/models"
	"telebot-trading/app/repositories"
)

func Login(email, password string) (error, *models.Admin) {
	admin := repositories.GetAdminByEmail(email)

	if admin != nil {
		if admin.VerifyPassword(password) {
			return nil, admin
		}
	}

	return errors.New("Email atau password salah"), nil
}
