package aws

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"log"

	"github.com/taniyuu/gin-cognito-sample/domain/model"
	"github.com/taniyuu/gin-cognito-sample/domain/proxy"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/pkg/errors"
)

// AWS S3に対する操作を提供します
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

	mac := hmac.New(sha256.New, []byte(*cic.clientSecret))
	mac.Write([]byte(req.Email + *cic.clientID))
	secretHash := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	newUserData := &cognitoidentityprovider.SignUpInput{
		ClientId:       cic.clientID,
		Password:       aws.String(req.Password),
		SecretHash:     aws.String(secretHash),
		Username:       aws.String(req.Email),
		ClientMetadata: map[string]*string{"custom-attr": aws.String("日本語も送れる")},
	}

	suo, err := cic.idp.SignUpWithContext(ctx, newUserData)
	if err != nil {
		return "", errors.WithStack(err)
	}
	log.Default().Println(suo)
	log.Default().Println(base64.RawURLEncoding.EncodeToString([]byte(*suo.UserSub)))
	return *suo.UserSub, nil
}
