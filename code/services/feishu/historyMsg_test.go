package feishu

import (
	"github.com/k0kubun/pp/v3"
	"start-feishubot/initialization"
	"testing"
)

func TestClient_GetHistoryMsg(t *testing.T) {
	config := initialization.LoadConfig("../../config.yaml")
	Client := NewClient(config.FeishuAppId, config.FeishuAppSecret)
	reqMsg := &HistoryMsgRequest{
		ChatId:    "oc_76bcd5398733bd921487388a22e47daa",
		StartTime: "1686370000",
		EndTime:   "1686377855",
		Size:      50,
	}
	resp, err := Client.GetHistoryMsg(*reqMsg)
	if err != nil {
		t.Errorf("GetHistoryMsg() error = %v", err)
	}
	pp.Println(resp)
}

func TestClient_GetCustomHistoryMsg(t *testing.T) {
	config := initialization.LoadConfig("../../config.yaml")
	Client := NewClient(config.FeishuAppId, config.FeishuAppSecret)
	reqMsg := &CustomHistoryMsg{
		ChatId: "oc_76bcd5398733bd921487388a22e47daa",
		During: 8200,
		Size:   50,
	}
	resp, err := Client.GetCustomHistoryMsg(*reqMsg)
	if err != nil {
		t.Errorf("GetHistoryMsg() error = %v", err)
	}
	pp.Println(resp)
}
