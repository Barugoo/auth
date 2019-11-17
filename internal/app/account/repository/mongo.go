package repository

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/barugoo/oscillo-auth/internal/app/models"
	"github.com/barugoo/oscillo-auth/internal/app/service"
)

type accountRepository struct {
	service    service.AuthService
	collection *mongo.Collection
}

func NewAccountRepository(service service.AuthService, collection *mongo.Collection) AccountRepository {
	return &accountRepository{
		service:    service,
		collection: collection,
	}
}

func (h *accountRepository) GetAccountByEmail(ctx context.Context, email string) (*models.Account, error) {
	span := h.service.StartSpan(ctx, "GetAccountByEmail")
	defer span.Finish()

	account, err := h.getAccountByEmail(email)
	if err != nil {
		err = h.mapError(err)
	}
	return account, err
}

func (h *accountRepository) getAccountByEmail(email string) (*models.Account, error) {
	var account *models.Account
	err := h.collection.FindOne(context.TODO(), bson.D{{"email", email}}).Decode(&account)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (h *accountRepository) CreateAccount(ctx context.Context, account *models.Account) (*models.Account, error) {
	span := h.service.StartSpan(ctx, "CreateAccount")
	defer span.Finish()

	account, err := h.createAccount(account)
	if err != nil {
		err = h.mapError(err)
	}
	return account, err
}

func (h *accountRepository) createAccount(account *models.Account) (*models.Account, error) {
	result, err := h.collection.InsertOne(context.TODO(), account)
	if err != nil {
		return nil, err
	}

	var acc *models.Account

	err = h.collection.FindOne(context.TODO(), bson.D{{"id", result.InsertedID}}).Decode(&acc)
	if err != nil {
		return nil, err
	}
	return acc, nil
}

func (h *accountRepository) DeleteAccount(ctx context.Context, account *models.Account) (bool, error) {
	span := h.service.StartSpan(ctx, "DeleteAccount")
	defer span.Finish()

	ok, err := h.deleteAccount(account)
	if err != nil {
		err = h.mapError(err)
	}
	return ok, err
}

func (h *accountRepository) deleteAccount(account *models.Account) (bool, error) {
	_, err := h.collection.DeleteOne(context.TODO(), bson.D{{"id", account.ID}})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (h *accountRepository) UpdateAccount(ctx context.Context, account *models.Account) (*models.Account, error) {
	span := h.service.StartSpan(ctx, "UpdateAccount")
	defer span.Finish()

	account, err := h.updateAccount(account)
	if err != nil {
		err = h.mapError(err)
	}
	return account, err
}

func (h *accountRepository) updateAccount(account *models.Account) (*models.Account, error) {
	_, err := h.collection.UpdateOne(context.TODO(), bson.D{{"id", account.ID}}, account)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (h *accountRepository) mapError(err error) error {
	switch err {
	case mongo.ErrNoDocuments:
		return ErrAccountNotFound
	default:
		return fmt.Errorf("%v", err)
	}
}
