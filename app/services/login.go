package services

import (
	"errors"
	"telebot-trading/app/models"
	"telebot-trading/app/repositories"
)

func Login(email, password string) (*models.Admin, error) {
	admin := repositories.GetAdminByEmail(email)

	if admin != nil {
		if admin.VerifyPassword(password) {
			return admin, nil
		}
	}

	return nil, errors.New("email atau password salah")
}
