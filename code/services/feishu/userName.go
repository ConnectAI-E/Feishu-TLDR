package feishu

import (
	"context"
	larkcontact "github.com/larksuite/oapi-sdk-go/v3/service/contact/v3"
	"github.com/patrickmn/go-cache"
)

func (c *Client) GetSenderName(senderOpenId string) string {
	req := larkcontact.NewGetUserReqBuilder().
		UserId(senderOpenId).
		Build()
	userInfo, err := c.Client.Contact.User.Get(context.Background(), req)
	username := senderOpenId
	if err == nil && userInfo != nil && userInfo.Data != nil && userInfo.Data.User != nil {
		username = *userInfo.Data.User.Name
	} else {
		username = senderOpenId
	}
	return username
}

//cache

func (c *Client) GetSenderNameCache(senderOpenId string) string {
	username, found := c.Cache.Get(senderOpenId)
	if found {
		return username.(string)
	}
	username = c.GetSenderName(senderOpenId)
	c.Cache.Set(senderOpenId, username, cache.DefaultExpiration)
	return username.(string)
}
