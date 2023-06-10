package oauth

import (
	"errors"
	"go-thirdparty/result"
	"go-thirdparty/utils"
	"strconv"
	"strings"
)

// QQ授权登录
type AuthQq struct {
	BaseRequest
}

func NewAuthQq(conf *AuthConfig) *AuthQq {
	authRequest := &AuthQq{}
	authRequest.Set(utils.RegisterSourceQQ, conf)

	authRequest.authorizeUrl = "https://graph.qq.com/oauth2.0/authorize"
	authRequest.TokenUrl = "https://graph.qq.com/oauth2.0/token"
	authRequest.openidUrl = "https://graph.qq.com/oauth2.0/me"
	authRequest.userInfoUrl = "https://graph.qq.com/user/get_user_info"

	return authRequest
}

// 获取登录地址
func (a *AuthQq) GetRedirectUrl(state string) string {
	url := utils.NewUrlBuilder(a.authorizeUrl).
		AddParam("response_type", "code").
		AddParam("client_id", a.config.ClientId).
		AddParam("redirect_uri", a.config.RedirectUrl).
		AddParam("state", a.GetState(state)).
		Build()
	return url
}

// 获取token
func (a *AuthQq) GetToken(code string) (*result.TokenResult, error) {
	url := utils.NewUrlBuilder(a.TokenUrl).
		AddParam("grant_type", "authorization_code").
		AddParam("code", code).
		AddParam("client_id", a.config.ClientId).
		AddParam("client_secret", a.config.ClientSecret).
		AddParam("redirect_uri", a.config.RedirectUrl).
		Build()

	body, err := utils.Get(url)
	if err != nil {
		return nil, err
	}
	if strings.Index(body, "callback") != -1 {
		body = body[strings.Index(body, "(")+1 : strings.LastIndex(body, ")")]
		m := utils.JsonToMSS(body)
		if _, ok := m["error"]; ok {
			return nil, errors.New(m["error_description"])
		}
	}

	m := utils.StrToMSS(body)
	token := &result.TokenResult{
		AccessToken:  m["access_token"],
		RefreshToken: m["refresh_token"],
		ExpireIn:     m["expires_in"],
	}
	return token, nil
}

// 获取用户openid
func (a *AuthQq) GetOpenid(accessToken string) (*result.Credentials, error) {
	url := utils.NewUrlBuilder(a.openidUrl).
		AddParam("access_token", accessToken).
		Build()

	body, err := utils.Get(url)
	if err != nil {
		return nil, err
	}
	body = body[strings.Index(body, "(")+1 : strings.LastIndex(body, ")")]
	m := utils.JsonToMSS(body)
	if _, ok := m["error"]; ok {
		return nil, errors.New(m["error_description"])
	}
	credentials := &result.Credentials{
		OpenId:  m["openid"],
		Unionid: m["unionid"],
	}
	return credentials, nil
}

// 获取用户信息
func (a *AuthQq) GetUserInfo(openId string, accessToken string) (*result.UserResult, error) {
	url := utils.NewUrlBuilder(a.userInfoUrl).
		AddParam("openid", openId).
		AddParam("access_token", accessToken).
		AddParam("oauth_consumer_key", a.config.ClientId).
		Build()

	body, err := utils.Get(url)
	if err != nil {
		return nil, err
	}
	m := utils.JsonToMSS(body)
	if _, ok := m["error"]; ok {
		return nil, errors.New(m["error_description"])
	}
	avatar := m["figureurl_qq_2"] //大小为100×100像素的QQ头像URL。需要注意，不是所有的用户都拥有QQ的100x100的头像，但40x40像素则是一定会有。
	if len(avatar) == 0 {
		avatar = m["figureurl_qq_1"] //大小为40×40像素的QQ头像URL。
	}
	user := &result.UserResult{
		NickName:  m["nickname"], //用户在QQ空间的昵称。
		AvatarUrl: avatar,
		Location:  m["province"] + m["city"],
		City:      m["city"],     //普通用户个人资料填写的城市
		Province:  m["province"], //普通用户个人资料填写的省份
		Source:    a.registerSource,
		Gender:    strconv.Itoa(utils.GetRealGender(m["gender"]).Code),
	}
	if m["ret"] != "0" { //ret	返回码  0: 正确返回
		return nil, errors.New("获取用户信息失败！")
	}
	return user, nil
}
