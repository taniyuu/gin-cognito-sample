package usecase

import (
	"context"
	"fmt"
	"log"

	"github.com/pkg/errors"

	"github.com/go-resty/resty/v2"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

// Auth ...
type Auth struct {
	jwk               *JWK
	jwkURL            string
	cognitoRegion     string
	cognitoUserPoolID string
}

// Config ...
type Config struct {
	CognitoRegion     string
	CognitoUserPoolID string
}

// JWK ...
type JWK struct {
	Keys []struct {
		Alg string `json:"alg"`
		E   string `json:"e"`
		Kid string `json:"kid"`
		Kty string `json:"kty"`
		N   string `json:"n"`
	} `json:"keys"`
}

// AuthnUsecase 認証操作を抽象化します
type AuthnUsecase interface {
	Authorization(ctx context.Context, accessToken string) error
}

type authnUsecase struct {
	jwk *JWK
}

// NewAuthenticationUsecase AuthenticationUsecaseを生成します
func NewAuthenticationUsecase(region, poolID string) AuthnUsecase {
	jwkURL := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json", region, poolID)
	jwk, err := cacheJWK(jwkURL)
	if err != nil {
		log.Fatal(err)
	}
	return &authnUsecase{jwk}
}

// アクセストークンを検証しアカウント情報を取得します
func (au *authnUsecase) Authorization(ctx context.Context, accessToken string) error {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("kid header not found")
		}

		log.Default().Println(kid, au.jwk.Keys)
		if !ok {
			return nil, fmt.Errorf("key with specified kid is not present in jwks")
		}
		var publickey interface{}
		err = keys.Raw(&publickey)
		if err != nil {
			return nil, fmt.Errorf("could not parse pubkey")
		}
		return publickey, nil
	})

	// token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
	// 	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
	// 		return nil, errors.WithStack(errors.New(fmt.Sprint("unexpected signing method: %v", token.Header["alg"])))
	// 	}
	// 	return []byte(config.Config.TokenSecret), nil
	// })

	if err != nil {
		return nil, errors.Wrap(err, errors.InvalidLoginSession, "failed to parse access token")
	}

	var accountID string
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if uid, ok := claims["accountID"].(string); ok {
			accountID = uid
		}
	}

	account := &model.LoginAccount{}
	if len(accountID) > 0 {
		if err = au.cr.Get(ctx, model.SessionKey(accountID, accessToken), account); err != nil {
			return nil, errors.Wrap(err, errors.InvalidLoginSession, "failed to authenticate account")
		} else {
			return account, nil
		}
	}

	return nil, errors.New(errors.InvalidLoginSession, "failed to get accountID from access token")
}

func cacheJWK(jwkURL string) (*JWK, error) {
	client := resty.New()
	result := new(JWK)
	resp, err := client.R().SetResult(result).Get("http://localhost:1314")
	log.Default().Print(resp)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return result, nil
}
