package weixin

const (
	Url = "https://api.weixin.qq.com/cgi-bin/"
)

type Weixin struct {
	Request     Request
	AccessToken AccessToken
}

func New(token, appId, appSecret string) *Weixin {
	return &Weixin{
		Request:     Request{Token: token},
		AccessToken: AccessToken{AppId: appId, AppSecret: appSecret},
	}
}
