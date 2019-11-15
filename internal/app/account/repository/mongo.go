package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/barugoo/oscillo-auth/internal/app/models"
)

type accountRepository struct {
	db *mongo.Collection
}

func NewAccountRepository(collection *mongo.Collection) AccountRepository {
	return &accountRepository{db: collection}
}

func (h *accountRepository) GetAccountByEmail(email string) (*models.Account, error) {
	var account *models.Account
	err := h.db.FindOne(context.TODO(), bson.D{{"email", email}}).Decode(&account)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (h *accountRepository) CreateAccount(account *models.Account) (*models.Account, error) {
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

func (h *accountRepository) DeleteAccount(account *models.Account) (bool, error) {
	_, err := h.db.DeleteOne(context.TODO(), bson.D{{"id", account.ID}})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (h *accountRepository) UpdateAccount(account *models.Account) (*models.Account, error) {
	_, err := h.db.UpdateOne(context.TODO(), bson.D{{"id", account.ID}}, account)
	if err != nil {
		return nil, err
	}
	return account, nil
}
