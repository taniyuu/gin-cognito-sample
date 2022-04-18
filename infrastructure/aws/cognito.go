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
func NewCognitoProxy(poolID, clientID, clientSecret string) proxy.UserProxy {
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
		ClientId:   cic.clientID,
		SecretHash: aws.String(cic.calcSecretHash(req.Email)),
		Username:   aws.String(req.Email),
		Password:   aws.String(req.Password),
		UserAttributes: []*cognitoidentityprovider.AttributeType{
			{Name: aws.String("name"), Value: aws.String(req.Name)},
		},
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
	aiao, err := cic.initiateAuthWithContext(ctx, req)
	if err != nil {
		return nil, err
	}
	// MFAなどの場合nilの可能性もあるので注意
	return &model.Token{
		IDToken:      *aiao.AuthenticationResult.IdToken,
		AccessToken:  *aiao.AuthenticationResult.AccessToken,
		RefreshToken: aiao.AuthenticationResult.RefreshToken}, nil
}

// Refresh トークンリフレッシュ
func (cic *cognitoIdpClient) Refresh(ctx context.Context, req *model.RefreshReq) (*model.Token, error) {
	iai := &cognitoidentityprovider.InitiateAuthInput{
		ClientId: cic.clientID,
		AuthFlow: aws.String(cognitoidentityprovider.AuthFlowTypeRefreshTokenAuth),
		AuthParameters: map[string]*string{
			"REFRESH_TOKEN": aws.String(req.RefreshToken),
			"SECRET_HASH":   aws.String(cic.calcSecretHash(req.Sub)),
		},
	}
	aiao, err := cic.idp.InitiateAuthWithContext(ctx, iai)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	log.Default().Println(aiao)
	// MFAなどの場合nilの可能性もあるので注意
	return &model.Token{
		IDToken:      *aiao.AuthenticationResult.IdToken,
		AccessToken:  *aiao.AuthenticationResult.AccessToken,
		RefreshToken: aiao.AuthenticationResult.RefreshToken}, nil
}

// ChangePassword パスワード変更
func (cic *cognitoIdpClient) ChangePassword(ctx context.Context, req *model.ChangePasswordReq) error {
	cpi := &cognitoidentityprovider.ChangePasswordInput{
		AccessToken:      aws.String(req.AccessToken),
		PreviousPassword: aws.String(req.PreviousPassword),
		ProposedPassword: aws.String(req.ProposedPassword),
	}
	_, err := cic.idp.ChangePasswordWithContext(ctx, cpi)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// ForgotPassword パスワード変更
func (cic *cognitoIdpClient) ForgotPassword(ctx context.Context, req *model.ForgotPasswordReq) error {
	fpi := &cognitoidentityprovider.ForgotPasswordInput{
		ClientId:   cic.clientID,
		SecretHash: aws.String(cic.calcSecretHash(req.Email)),
		Username:   aws.String(req.Email),
	}
	_, err := cic.idp.ForgotPasswordWithContext(ctx, fpi)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// ConfirmForgotPassword パスワード変更確認
func (cic *cognitoIdpClient) ConfirmForgotPassword(ctx context.Context, req *model.ConfirmForgotPasswordReq) error {
	cfpi := &cognitoidentityprovider.ConfirmForgotPasswordInput{
		ClientId:         cic.clientID,
		SecretHash:       aws.String(cic.calcSecretHash(req.Email)),
		Username:         aws.String(req.Email),
		ConfirmationCode: aws.String(req.Code),
		Password:         aws.String(req.Password),
	}
	_, err := cic.idp.ConfirmForgotPasswordWithContext(ctx, cfpi)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// GetProfile 属性取得
func (cic *cognitoIdpClient) GetProfile(ctx context.Context, req *model.GetProfileReq) (*model.User, error) {
	gui := &cognitoidentityprovider.GetUserInput{
		AccessToken: &req.AccessToken,
	}
	guo, err := cic.idp.GetUserWithContext(ctx, gui)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	log.Default().Println(guo)
	return cic.convertToUserModel(guo.UserAttributes), nil
}

// ChangeProfile 属性変更
func (cic *cognitoIdpClient) ChangeProfile(ctx context.Context, req *model.ChangeProfileReq) error {
	uuai := &cognitoidentityprovider.UpdateUserAttributesInput{
		AccessToken: aws.String(req.AccessToken),
		UserAttributes: []*cognitoidentityprovider.AttributeType{
			{Name: aws.String("name"), Value: aws.String(req.Name)},
		},
	}
	uuao, err := cic.idp.UpdateUserAttributesWithContext(ctx, uuai)
	if err != nil {
		return errors.WithStack(err)
	}
	log.Default().Println(uuao)
	return nil
}

// Signout ログアウト
func (cic *cognitoIdpClient) Signout(ctx context.Context, req *model.SignoutReq) error {
	rti := &cognitoidentityprovider.RevokeTokenInput{
		ClientId:     cic.clientID,
		ClientSecret: cic.clientSecret,
		Token:        aws.String(req.RefreshToken),
	}
	rto, err := cic.idp.RevokeTokenWithContext(ctx, rti)
	if err != nil {
		return errors.WithStack(err)
	}
	log.Default().Println(rto)
	return nil
}

// Invite 招待
func (cic *cognitoIdpClient) Invite(ctx context.Context, req *model.InviteReq) (string, error) {
	// 招待は２重送信を拒否する（アカウントの存在を確認してから送信する）
	acui := &cognitoidentityprovider.AdminCreateUserInput{
		UserPoolId: cic.poolID,
		Username:   aws.String(req.Email),
	}
	rto, err := cic.idp.AdminCreateUserWithContext(ctx, acui)
	if err != nil {
		return "", errors.WithStack(err)
	}
	log.Default().Println(rto)
	var sub string
	for _, attr := range rto.User.Attributes {
		if *attr.Name == "sub" {
			sub = *attr.Value
		}
	}
	return sub, nil
}

// RespondToInvitation 招待応答
func (cic *cognitoIdpClient) RespondToInvitation(ctx context.Context, req *model.RespondToInvitationReq) (*model.Token, error) {
	// 属性変更
	auuai := &cognitoidentityprovider.AdminUpdateUserAttributesInput{
		UserPoolId: cic.poolID,
		Username:   aws.String(req.Email),
		UserAttributes: []*cognitoidentityprovider.AttributeType{
			{Name: aws.String("name"), Value: aws.String(req.Name)},
			{Name: aws.String("email_verified"), Value: aws.String("true")}, // eメール確認済にする
		},
	}
	auuao, err := cic.idp.AdminUpdateUserAttributesWithContext(ctx, auuai)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	log.Default().Println(auuao)
	// ログイン
	aiao, err := cic.initiateAuthWithContext(ctx, &model.SigninReq{Email: req.Email, Password: req.ConfirmationCode})
	if err != nil {
		return nil, err
	}
	// パスワード変更
	rtaci := &cognitoidentityprovider.RespondToAuthChallengeInput{
		ClientId:      cic.clientID,
		ChallengeName: aiao.ChallengeName,
		Session:       aiao.Session,
		ChallengeResponses: map[string]*string{
			"USERNAME":     aws.String(req.Email),
			"NEW_PASSWORD": aws.String(req.Password),
			"SECRET_HASH":  aws.String(cic.calcSecretHash(req.Email)),
		},
	}
	rtaco, err := cic.idp.RespondToAuthChallengeWithContext(ctx, rtaci)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	log.Default().Println(rtaco)
	return &model.Token{
		IDToken:      *rtaco.AuthenticationResult.IdToken,
		AccessToken:  *rtaco.AuthenticationResult.AccessToken,
		RefreshToken: rtaco.AuthenticationResult.RefreshToken,
	}, nil
}

// GetUser subで検索
func (cic *cognitoIdpClient) GetUser(ctx context.Context, req *model.GetUserReq) (*model.User, error) {
	lui := &cognitoidentityprovider.ListUsersInput{
		UserPoolId: cic.poolID,
		Filter:     aws.String(fmt.Sprintf(`sub = "%s"`, req.Sub)),
	}
	luo, err := cic.idp.ListUsersWithContext(ctx, lui)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	log.Default().Println(luo)
	if len(luo.Users) == 0 {
		return nil, errors.WithStack(fmt.Errorf("user not found"))
	}
	return cic.convertToUserModel(luo.Users[0].Attributes), nil
}

func (cic *cognitoIdpClient) calcSecretHash(username string) string {
	mac := hmac.New(sha256.New, []byte(*cic.clientSecret))
	mac.Write([]byte(username + *cic.clientID))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func (cic *cognitoIdpClient) initiateAuthWithContext(ctx context.Context, req *model.SigninReq) (*cognitoidentityprovider.InitiateAuthOutput, error) {
	iai := &cognitoidentityprovider.InitiateAuthInput{
		ClientId: cic.clientID,
		AuthFlow: aws.String(cognitoidentityprovider.AuthFlowTypeUserPasswordAuth),
		AuthParameters: map[string]*string{
			"USERNAME":    aws.String(req.Email),
			"PASSWORD":    aws.String(req.Password),
			"SECRET_HASH": aws.String(cic.calcSecretHash(req.Email)),
		},
	}
	aiao, err := cic.idp.InitiateAuthWithContext(ctx, iai)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	log.Default().Println(aiao)
	return aiao, nil
}

func (cic *cognitoIdpClient) convertToUserModel(attrs []*cognitoidentityprovider.AttributeType) *model.User {
	u := new(model.User)
	for _, attr := range attrs {
		if *attr.Name == "email" {
			u.Email = *attr.Value
		}
		if *attr.Name == "name" {
			u.Name = *attr.Value
		}
	}
	return u
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

func (ca *cognitoAuthorizar) ValidateJWT(accessToken string) (string, error) {
	jt, err := jwt.Parse(
		[]byte(accessToken),
		jwt.WithKeySet(ca.jwk),
		jwt.WithValidate(true),
		jwt.WithIssuer(fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s", ca.region, ca.poolID)),
		jwt.WithAudience(ca.clientID),
		jwt.WithClaimValue("token_use", "id"),
	)
	if err != nil {
		return "", errors.WithStack(err)
	}
	log.Default().Printf("%+v", jt.PrivateClaims())
	return jt.Subject(), nil
}
