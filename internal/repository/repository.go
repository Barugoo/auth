package repository

import (
	"github.com/barugoo/oscillo-auth/internal/models"
)

type AccountRepository interface {
	GetAccountByEmail(email string) (*models.Account, error)
	CreateAccount(account *models.Account) (*models.Account, error)
	DeleteAccount(account *models.Account) (bool, error)
	UpdateAccount(account *models.Account) (*models.Account, error)
}
