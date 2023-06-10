package handlers

import (
	"github.com/k0kubun/pp/v3"
	"log"
	"start-feishubot/services/chatgpt"
	"start-feishubot/services/feishu"
	"start-feishubot/services/openai"
	"time"
)

type MessageAction struct { /*消息*/
	chatgpt *chatgpt.ChatGPT
}

func (m *MessageAction) MsgList(a *ActionInfo) string {
	//get recent message
	chatId := a.info.chatId
	client := feishu.GetClient()
	msgs, err := client.GetCustomHistoryMsg(feishu.
		CustomHistoryMsg{
		ChatId: *chatId,
		During: 60 * 60 * 60 * 20,
		Size:   50,
	})
	if err != nil {
		log.Println("get msg list error", err)
		return ""
	}
	return msgs
}

func (m *MessageAction) Execute(a *ActionInfo) bool {
	//get recent message
	msgs := m.MsgList(a)
	if msgs == "" {
		return false
	}
	pp.Println(msgs)
	cardId, err2 := sendOnProcess(a)
	if err2 != nil {
		return false
	}

	answer := ""
	chatResponseStream := make(chan string)
	done := make(chan struct{}) // 添加 done 信号，保证 goroutine 正确退出
	noContentTimeout := time.AfterFunc(10*time.Second, func() {
		pp.Println("no content timeout")
		close(done)
		err := updateFinalCard(*a.ctx, "请求超时", cardId)
		if err != nil {
			return
		}
		return
	})
	defer noContentTimeout.Stop()

	var msg []openai.Messages
	msg = append(msg, openai.Messages{
		Role: "system",
		Content: "下面我会给出一系列群聊对话，帮我总结出这段时间内群聊对话的主要内容," +
			"\n1-注意，直接给出分析结果即可，不要复述对话内容\n2." + "按照用户名对每个人表达的内容做出简短地总结性归纳",
	})

	msg = append(msg, openai.Messages{
		Role: "user", Content: msgs,
	})
	go func() {
		defer func() {
			if err := recover(); err != nil {
				err := updateFinalCard(*a.ctx, "聊天失败", cardId)
				if err != nil {
					return
				}
			}
		}()

		//log.Printf("UserId: %s , Request: %s", a.info.userId, msg)

		if err := m.chatgpt.StreamChat(*a.ctx, msg, chatResponseStream); err != nil {
			err := updateFinalCard(*a.ctx, "聊天失败", cardId)
			if err != nil {
				return
			}
			close(done) // 关闭 done 信号
		}

		close(done) // 关闭 done 信号
	}()
	ticker := time.NewTicker(700 * time.Millisecond)
	defer ticker.Stop() // 注意在函数结束时停止 ticker
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				err := updateTextCard(*a.ctx, answer, cardId)
				if err != nil {
					return
				}
			}
		}
	}()

	for {
		select {
		case res, ok := <-chatResponseStream:
			if !ok {
				return false
			}
			noContentTimeout.Stop()
			answer += res
			//pp.Println("answer", answer)
		case <-done: // 添加 done 信号的处理
			err := updateFinalCard(*a.ctx, answer, cardId)
			if err != nil {
				return false
			}
			ticker.Stop()
			msg := append(msg, openai.Messages{
				Role: "assistant", Content: answer,
			})
			a.handler.sessionCache.SetMsg(*a.info.sessionId, msg)
			close(chatResponseStream)
			//if new topic
			//if len(msg) == 2 {
			//	//fmt.Println("new topic", msg[1].Content)
			//	//updateNewTextCard(*a.ctx, a.info.sessionId, a.info.msgId,
			//	//	completions.Content)
			//}
			log.Printf("UserId: %s , Request: %s , Response: %s", a.info.userId, msg, answer)
			return false
		}
	}
}

func sendOnProcess(a *ActionInfo) (*string, error) {
	// send 正在处理中
	cardId, err := sendOnProcessCard(*a.ctx, a.info.sessionId, a.info.msgId)
	if err != nil {
		return nil, err
	}
	return cardId, nil

}
