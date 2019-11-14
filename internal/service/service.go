package service

type AuthService interface {
	SetKV(key, value string) (bool, error)
	GetKV(key string) (string, error)
}

type authService struct {
	kv map[string]string
}

func (a *authService) SetKV(k, v string) (bool, error) {
	a.kv[k] = v
	return true, nil
}

func (a *authService) GetKV(k string) (string, error) {
	return a.kv[k], nil
}

func NewAuthService() AuthService {
	return &authService{
		kv: make(map[string]string),
	}
}
