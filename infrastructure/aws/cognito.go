package aws

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

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
	// panic("hoge")
	return &cognitoIdpClient{
		cognitoidentityprovider.New(sess),
		&poolID, &clientID, &clientSecret,
	}
}

// Signup サインアップ
func (cic *cognitoIdpClient) Signup(ctx context.Context, info *model.UserInfo) (string, error) {
	mac := hmac.New(sha256.New, []byte(*cic.clientSecret))
	mac.Write([]byte(info.Email + *cic.clientID))
	secretHash := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	newUserData := &cognitoidentityprovider.SignUpInput{
		ClientId:       cic.clientID,
		Password:       aws.String("Passw0rd!"),
		SecretHash:     aws.String(secretHash),
		Username:       aws.String(info.Email),
		ClientMetadata: map[string]*string{"test": aws.String("日本語")},
	}

	out, err := cic.idp.SignUp(newUserData)
	if err != nil {
		return "", errors.WithStack(err)
	}
	fmt.Println(out)
	return "hoge", nil
}
