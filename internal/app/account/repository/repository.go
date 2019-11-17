package repository

import (
	"context"

	models "github.com/barugoo/oscillo-auth/internal/app/account"
)

type AccountRepository interface {
	GetAccountByEmail(ctx context.Context, email string) (*models.Account, error)
	CreateAccount(ctx context.Context, account *models.Account) (*models.Account, error)
	DeleteAccount(ctx context.Context, account *models.Account) (bool, error)
	UpdateAccount(ctx context.Context, account *models.Account) (*models.Account, error)
}
