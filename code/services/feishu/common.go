package feishu

import (
	lark "github.com/larksuite/oapi-sdk-go/v3"
	"github.com/patrickmn/go-cache"
)

type BotConfig struct {
	AppId     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
}

type Client struct {
	Client *lark.Client
	Cache  *cache.Cache
}

var ClientInstance *Client

func NewClient(appId string, appSecret string) *Client {
	ClientInstance = &Client{
		Client: lark.NewClient(appId, appSecret),
		Cache:  cache.New(cache.NoExpiration, cache.NoExpiration),
	}
	return ClientInstance

}

func GetClient() *Client {
	return ClientInstance
}
