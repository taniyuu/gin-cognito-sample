package proxy

// AuthorizarProxy 認可操作を抽象化します
type AuthorizarProxy interface {
	ValidateJWT(token string) (sub, email string, err error)
}
