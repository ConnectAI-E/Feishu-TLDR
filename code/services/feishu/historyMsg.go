package feishu

import (
	"bytes"
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
	string, error) {
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
		return "", err
	}
	Items := resp.Data.Items
	Items = filterTextMessage(Items)
	var buffer bytes.Buffer
	for _, msg := range Items {
		msgField := getMsgField(msg, c)
		msgStr := msgField.ToString()
		buffer.WriteString(msgStr)
	}
	return buffer.String(), nil
}

func (c *Client) GetCustomHistoryMsg(req CustomHistoryMsg) (
	string, error) {
	// 获取当前时间，到秒
	now := time.Now().Unix()
	//获取起始时间
	startTime := now - int64(req.During)
	// 转换为字符串

	msgStr, err := c.GetHistoryMsg(HistoryMsgRequest{
		ChatId:    req.ChatId,
		StartTime: fmt.Sprintf("%v", startTime),
		EndTime:   fmt.Sprintf("%v", now),
		Size:      req.Size,
	})
	if err != nil {
		return "", err
	}
	return msgStr, nil

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

func getMsgField(msg *larkim.Message, client *Client) MsgInfo {
	var msgInfo MsgInfo
	msgInfo.userName = client.GetSenderNameCache(*msg.Sender.Id)
	msgInfo.CreateTime = *msg.CreateTime
	msgInfo.Content = utils.ParseContent(*msg.Body.Content)
	return msgInfo
}

// ToString to string
// userName:CreateTime:Content||
func (m *MsgInfo) ToString() string {
	var buffer bytes.Buffer
	buffer.WriteString(m.userName)
	buffer.WriteString(":")
	buffer.WriteString(m.CreateTime)
	buffer.WriteString(":")
	buffer.WriteString(m.Content)
	buffer.WriteString("||")
	return buffer.String()
}
