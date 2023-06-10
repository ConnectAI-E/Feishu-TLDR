package feishu

import (
	"context"
	"fmt"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"start-feishubot/utils"
	"time"
)

type HistoryMsgRequest struct {
	ChatId    string `json:"chat_id"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Size      int    `json:"size" validate:"lte=50,gte=0"`
}

type CustomHistoryMsg struct {
	ChatId string `json:"chat_id"`
	During int    `json:"during"`
	Size   int    `json:"size"`
}

type MsgInfo struct {
	userName   string
	CreateTime string
	Content    string
}

func (c *Client) GetHistoryMsg(req HistoryMsgRequest) (
	[]MsgInfo, error) {
	reqest := larkim.NewListMessageReqBuilder().
		ContainerIdType("chat").
		ContainerId(req.ChatId).
		StartTime(req.StartTime).
		EndTime(req.EndTime).
		PageSize(req.Size).
		PageToken("").
		Build()
	resp, err := c.Client.Im.Message.List(context.Background(),
		reqest)
	if err != nil {
		return nil, err
	}
	Items := resp.Data.Items
	Items = filterTextMessage(Items)
	var msgList []MsgInfo
	for _, msg := range Items {
		msgList = append(msgList, getMsgField(msg))
	}
	return msgList, nil
}

func (c *Client) GetCustomHistoryMsg(req CustomHistoryMsg) (
	interface{}, error) {
	// 获取当前时间，到秒
	now := time.Now().Unix()
	//获取起始时间
	startTime := now - int64(req.During)
	// 转换为字符串

	msgList, err := c.GetHistoryMsg(HistoryMsgRequest{
		ChatId:    req.ChatId,
		StartTime: fmt.Sprintf("%v", startTime),
		EndTime:   fmt.Sprintf("%v", now),
		Size:      req.Size,
	})
	if err != nil {
		return nil, err
	}
	return msgList, nil

}

func filterTextMessage(msgList []*larkim.Message) []*larkim.Message {
	var newMsgList []*larkim.Message
	for _, msg := range msgList {
		if *msg.MsgType == "text" {
			newMsgList = append(newMsgList, msg)
		}
	}
	return newMsgList
}

func getMsgField(msg *larkim.Message) MsgInfo {
	var msgInfo MsgInfo
	msgInfo.userName = *msg.Sender.Id
	msgInfo.CreateTime = *msg.CreateTime
	msgInfo.Content = utils.ParseContent(*msg.Body.Content)
	return msgInfo
}
