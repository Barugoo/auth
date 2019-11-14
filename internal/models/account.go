package models

type Account struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
	Secret2FA    string `json:"secret_2fa"`
	IsActive     bool   `json:"is_active"`
}

func (a *Account) Has2FA() bool {
	return len(a.Secret2FA) > 0
}
