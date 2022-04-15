package proxy

// AuthorizarProxy 認可操作を抽象化します
type AuthorizarProxy interface {
	ValidateJWT(accessToken string) error
}
