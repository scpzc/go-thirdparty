# 第三方登录

 Go语言实现的第三方授权登录，整合QQ、微信、微信小程序、微博、抖音、支付宝等第三方平台的授权登录

由scpzc/go-thirdparty修改，目前只改了QQ授权登录，其它登录尚未验证

### QQ授权登录
```go
package main

import (
	"fmt"
	"github.com/scpzc/go-thirdparty/oauth"
)

func main()  {
    qqConf := &oauth.AuthConfig{
		ClientId:     "your app_id", 
		ClientSecret: "your app_secret", 
		RedirectUrl:  "your redirect_url", 
	}
	qqAuth := oauth.NewAuthQq(qqConf)
    //获取第三方登录地址
    qqAuth.GetRedirectUrl("sate") 
    //获取token信息
	tokenRes, err := qqAuth.GetToken(c.PostForm("code"))
    //获取openid
	openidRes, err := qqAuth.GetOpenid(tokenRes.AccessToken)
    //获取用户信息
    userInfo,err:=qqAuth.GetUserInfo(openidRes.OpenId,tokenRes.AccessToken)
}
```