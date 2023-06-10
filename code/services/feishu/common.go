package feishu

import lark "github.com/larksuite/oapi-sdk-go/v3"

type BotConfig struct {
	AppId     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
}

type Client struct {
	Client *lark.Client
}

func NewClient(appId string, appSecret string) *Client {
	return &Client{
		Client: lark.NewClient(appId, appSecret),
	}
}
