package aws

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"strings"

	"github.com/taniyuu/gin-cognito-sample/domain/model"
	"github.com/taniyuu/gin-cognito-sample/domain/proxy"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/pkg/errors"
)

// Amazon Cognitoに対する操作を提供します
type cognitoIdpClient struct {
	idp                            *cognitoidentityprovider.CognitoIdentityProvider
	poolID, clientID, clientSecret *string
}

// NewCognitoProxy AuthenticatorProxyを生成します
func NewCognitoProxy(poolID, clientID, clientSecret string) proxy.AuthenticatorProxy {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	return &cognitoIdpClient{
		cognitoidentityprovider.New(sess),
		&poolID, &clientID, &clientSecret,
	}
}

// Signup サインアップ
func (cic *cognitoIdpClient) Signup(ctx context.Context, req *model.CreateReq) (string, error) {
	// Cognitoはメールアドレス確認前でも2重にサインアップAPIを呼び出すとUsernameExistsExceptionになる
	// resendでは根本解決にならないので、確認前のユーザがいれば削除する
	getUser := &cognitoidentityprovider.AdminGetUserInput{
		UserPoolId: cic.poolID,
		Username:   aws.String(req.Email),
	}
	guo, _ := cic.idp.AdminGetUserWithContext(ctx, getUser)
	if guo != nil {
		// 存在した場合
		confirmed := false // メール確認属性
		for _, ua := range guo.UserAttributes {
			if *ua.Name == "email_verified" && *ua.Value == "true" {
				confirmed = true
			}
		}
		if !confirmed {
			// 削除 エラーは無視する
			adui := &cognitoidentityprovider.AdminDeleteUserInput{
				UserPoolId: cic.poolID,
				Username:   aws.String(req.Email),
			}
			cic.idp.AdminDeleteUserWithContext(ctx, adui)
		}
	}

	newUserData := &cognitoidentityprovider.SignUpInput{
		ClientId:       cic.clientID,
		SecretHash:     aws.String(cic.calcSecretHash(req.Email)),
		Username:       aws.String(req.Email),
		Password:       aws.String(req.Password),
		ClientMetadata: map[string]*string{"custom-attr": aws.String("日本語も送れる")},
	}

	suo, err := cic.idp.SignUpWithContext(ctx, newUserData)
	if err != nil {
		return "", errors.WithStack(err)
	}
	log.Default().Println(suo)
	return *suo.UserSub, nil
}

// ConfirmAndSigninReq 確認
func (cic *cognitoIdpClient) ConfirmAndSignin(ctx context.Context, req *model.ConfirmAndSigninReq) (*model.Token, error) {
	// 確認した後ログイン失敗の事象を回避するために一度ログインを試行する
	_, err := cic.Signin(ctx, &model.SigninReq{Email: req.Email, Password: req.Password})
	if err != nil && strings.HasPrefix(errors.Unwrap(err).Error(), "NotAuthorizedException") {
		return nil, errors.WithStack(err)
	}

	csi := &cognitoidentityprovider.ConfirmSignUpInput{
		ClientId:         cic.clientID,
		SecretHash:       aws.String(cic.calcSecretHash(req.Email)),
		Username:         aws.String(req.Email),
		ConfirmationCode: aws.String(req.ConfirmationCode),
	}
	_, err = cic.idp.ConfirmSignUpWithContext(ctx, csi)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	resp, err := cic.Signin(ctx, &model.SigninReq{Email: req.Email, Password: req.Password})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return resp, nil
}

// Signin ログイン
func (cic *cognitoIdpClient) Signin(ctx context.Context, req *model.SigninReq) (*model.Token, error) {
	aiai := &cognitoidentityprovider.AdminInitiateAuthInput{
		UserPoolId: cic.poolID,
		ClientId:   cic.clientID,
		AuthFlow:   aws.String(cognitoidentityprovider.AuthFlowTypeAdminUserPasswordAuth),
		AuthParameters: map[string]*string{
			"USERNAME":    aws.String(req.Email),
			"PASSWORD":    aws.String(req.Password),
			"SECRET_HASH": aws.String(cic.calcSecretHash(req.Email)),
		},
	}
	aiao, err := cic.idp.AdminInitiateAuthWithContext(ctx, aiai)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	log.Default().Println(aiao)
	// MFAなどの場合nilの可能性もあるので注意
	return &model.Token{
		IDToken:      *aiao.AuthenticationResult.IdToken,
		RefreshToken: aiao.AuthenticationResult.RefreshToken}, nil
}

// Refresh トークンリフレッシュ
func (cic *cognitoIdpClient) Refresh(ctx context.Context, req *model.RefreshReq) (*model.Token, error) {
	aiai := &cognitoidentityprovider.AdminInitiateAuthInput{
		UserPoolId: cic.poolID,
		ClientId:   cic.clientID,
		AuthFlow:   aws.String(cognitoidentityprovider.AuthFlowTypeRefreshTokenAuth),
		AuthParameters: map[string]*string{
			"REFRESH_TOKEN": aws.String(req.RefreshToken),
			"SECRET_HASH":   aws.String(cic.calcSecretHash(req.Sub)),
		},
	}
	aiao, err := cic.idp.AdminInitiateAuthWithContext(ctx, aiai)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	log.Default().Println(aiao)
	// MFAなどの場合nilの可能性もあるので注意
	return &model.Token{
		IDToken:      *aiao.AuthenticationResult.IdToken,
		RefreshToken: aiao.AuthenticationResult.RefreshToken}, nil
}

func (cic *cognitoIdpClient) calcSecretHash(username string) string {
	mac := hmac.New(sha256.New, []byte(*cic.clientSecret))
	mac.Write([]byte(username + *cic.clientID))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

// NewCognitoAuthorizar AuthorizarProxyを生成する

type cognitoAuthorizar struct {
	jwk                      jwk.Set
	region, poolID, clientID string
}

func NewCognitoAuthorizar(region, poolID, clientID string) proxy.AuthorizarProxy {
	jwkURL := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json", region, poolID)
	jset, err := jwk.Fetch(context.Background(), jwkURL)
	if err != nil {
		log.Fatal(err)
	}
	return &cognitoAuthorizar{
		jset, region, poolID, clientID,
	}
}

func (ca *cognitoAuthorizar) ValidateJWT(accessToken string) error {
	fmt.Println(ca.clientID)
	jt, err := jwt.Parse(
		[]byte(accessToken),
		jwt.WithKeySet(ca.jwk),
		jwt.WithValidate(true),
		jwt.WithIssuer(fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s", ca.region, ca.poolID)),
		jwt.WithAudience(ca.clientID),
		jwt.WithClaimValue("token_use", "id"),
	)
	log.Default().Printf("%+v", jt.PrivateClaims())
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
