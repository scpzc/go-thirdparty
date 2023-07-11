package oauth

import (
	"errors"
	"github.com/scpzc/go-thirdparty/result"
	"github.com/scpzc/go-thirdparty/utils"
)

// 微信公众号登录
type AuthWxMp struct {
	BaseRequest
}

func NewAuthWxMp(conf *AuthConfig) *AuthWxMp {
	authRequest := &AuthWxMp{}
	authRequest.Set(utils.RegisterSourceWechat, conf)

	authRequest.authorizeUrl = "https://api.weixin.qq.com/cgi-bin/qrcode/create"
	authRequest.TokenUrl = "https://api.weixin.qq.com/cgi-bin/token"
	authRequest.userInfoUrl = "https://api.weixin.qq.com/sns/userinfo"

	return authRequest
}

// 获取场景化二维码
func (a *AuthWxMp) GetSceneQrcode(accessToken,sceneCode string) (qrCode string, err error) {
	url := utils.NewUrlBuilder(a.authorizeUrl).
		AddParam("access_token", accessToken).
		Build()
	var data = map[string]any{
		"expire_seconds":600,
		"action_name":"QR_STR_SCENE",
		"action_info":map[string]any{
			"scene":map[string]any{
				"scene_str":sceneCode,
			},
		},
	}
	body, err := utils.Post(url,data)
	if err != nil {
		return
	}
	m := utils.JsonToMSS(body)
	if _, ok := m["errmsg"]; ok {
		return "", errors.New(m["errmsg"])
	}

	//生成场景化二维码
	qrCode = utils.NewUrlBuilder("https://mp.weixin.qq.com/cgi-bin/showqrcode").
		AddParam("ticket", m["ticket"]).
		Build()
	return qrCode, nil
}

// 获取token
func (a *AuthWxMp) GetAccessToken() (*result.TokenResult, error) {
	url := utils.NewUrlBuilder(a.TokenUrl).
		AddParam("grant_type", "client_credential").
		AddParam("appid", a.config.ClientId).
		AddParam("secret", a.config.ClientSecret).
		Build()

	body, err := utils.Post(url,nil)
	if err != nil {
		return nil, err
	}
	m := utils.JsonToMSS(body)
	if _, ok := m["errmsg"]; ok {
		return nil, errors.New(m["errmsg"])
	}
	token := &result.TokenResult{
		AccessToken:  m["access_token"],
		ExpireIn:     m["expires_in"],
	}
	if token.AccessToken == "" {
		return nil, errors.New("获取AccessToken数据为空！")
	}
	return token, nil
}



// 获取第三方用户信息
func (a *AuthWxMp) GetUserInfo(accessToken string, openId string) (*result.UserResult, error) {
	url := utils.NewUrlBuilder(a.userInfoUrl).
		AddParam("openid", openId).
		AddParam("access_token", accessToken).
		Build()

	body, err := utils.Get(url)
	if err != nil {
		return nil, err
	}
	m := utils.JsonToMSS(body)
	if _, ok := m["error"]; ok {
		return nil, errors.New(m["error_description"])
	}
	user := &result.UserResult{
		OpenId:    m["openid"],
		UnionId:   m["unionid"],
		UserName:  m["nickname"],
		NickName:  m["nickname"],
		AvatarUrl: m["headimgurl"],
		City:      m["city"],
		Province:  m["province"],
		Country:   m["country"],
		Language:  m["language"],
		Source:    a.registerSource,
		Gender:    utils.GetRealGender("sex").Desc,
	}
	if user.OpenId == "" {
		return nil, errors.New("获取用户信息为空！")
	}
	return user, nil
}
