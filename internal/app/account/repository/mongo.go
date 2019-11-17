package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/barugoo/oscillo-auth/internal/app/models"
	"github.com/barugoo/oscillo-auth/internal/app/service"
)

type accountRepository struct {
	service service.AuthService
	db      *mongo.Collection
}

func NewAccountRepository(service service.AuthService, collection *mongo.Collection) AccountRepository {
	return &accountRepository{
		service: service,
		db:      collection,
	}
}

func (h *accountRepository) GetAccountByEmail(ctx context.Context, email string) (*models.Account, error) {
	span := h.service.StartSpan(ctx, "GetAccountByEmail")
	defer span.Finish()

	return h.getAccountByEmail(email)
}

func (h *accountRepository) getAccountByEmail(email string) (*models.Account, error) {
	var account *models.Account
	err := h.db.FindOne(context.TODO(), bson.D{{"email", email}}).Decode(&account)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (h *accountRepository) CreateAccount(ctx context.Context, account *models.Account) (*models.Account, error) {
	span := h.service.StartSpan(ctx, "CreateAccount")
	defer span.Finish()

	return h.createAccount(account)
}

func (h *accountRepository) createAccount(account *models.Account) (*models.Account, error) {
	result, err := h.db.InsertOne(context.TODO(), account)
	if err != nil {
		return nil, err
	}

	var acc *models.Account

	err = h.db.FindOne(context.TODO(), bson.D{{"id", result.InsertedID}}).Decode(&acc)
	if err != nil {
		return nil, err
	}
	return acc, nil
}

func (h *accountRepository) DeleteAccount(ctx context.Context, account *models.Account) (bool, error) {
	span := h.service.StartSpan(ctx, "DeleteAccount")
	defer span.Finish()

	return h.deleteAccount(account)
}

func (h *accountRepository) deleteAccount(account *models.Account) (bool, error) {
	_, err := h.db.DeleteOne(context.TODO(), bson.D{{"id", account.ID}})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (h *accountRepository) UpdateAccount(ctx context.Context, account *models.Account) (*models.Account, error) {
	span := h.service.StartSpan(ctx, "UpdateAccount")
	defer span.Finish()

	return h.updateAccount(account)
}

func (h *accountRepository) updateAccount(account *models.Account) (*models.Account, error) {
	_, err := h.db.UpdateOne(context.TODO(), bson.D{{"id", account.ID}}, account)
	if err != nil {
		return nil, err
	}
	return account, nil
}
