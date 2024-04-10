package wechat

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"webbook/internal/domain"
)

var redirectUri = url.PathEscape("https://xxx.xxx1.com/oauth2/wechat/callback")

const (
	targetPattern = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"
	UrlPattern    = "https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redirect"
)

type Service interface {
	AuthURL(ctx context.Context, state string) (string, error)
	VerifyCode(ctx context.Context, code string) (domain.WechatInfo, error)
}

type WechatService struct {
	appId     string
	appSecret string
	client    *http.Client
}

func NewWechatService(appId, appSecret string) Service {
	return &WechatService{
		appId:     appId,
		appSecret: appSecret,
		client:    http.DefaultClient,
	}
}

func (w *WechatService) AuthURL(ctx context.Context, state string) (string, error) {

	return fmt.Sprintf(UrlPattern, w.appId, redirectUri, state), nil

}

func (w *WechatService) VerifyCode(ctx context.Context, code string) (domain.WechatInfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(targetPattern, w.appId, w.appSecret, code), nil)
	if err != nil {
		return domain.WechatInfo{}, err
	}

	// 适用于复杂请求
	resp, err := w.client.Do(req)
	if err != nil {
		return domain.WechatInfo{}, err
	}

	decoder := json.NewDecoder(resp.Body)
	var res Result
	err = decoder.Decode(&res)
	if err != nil {
		return domain.WechatInfo{}, err
	}

	if res.Errcode != 0 {
		return domain.WechatInfo{}, fmt.Errorf("微信返回错误响应，错误码：%d, 错误信息：%s", res.Errcode, res.Errmsg)
	}

	return domain.WechatInfo{
		OpenId:  res.Openid,
		UnionId: res.Unionid,
	}, nil
}

type Result struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`

	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Openid       string `json:"openid"`
	Scope        string `json:"scope"`
	Unionid      string `json:"unionid"`
}
