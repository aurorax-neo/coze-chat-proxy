package chat

import (
	"coze-chat-proxy/bot/discord"
	"coze-chat-proxy/common"
	"coze-chat-proxy/logger"
	v1 "coze-chat-proxy/v1"
	"coze-chat-proxy/v1/chat/apireq"
	"coze-chat-proxy/v1/chat/apiresp"
	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"time"
)

func gpt(c *gin.Context, apiReq *apireq.Req, bot *discord.ProxyBot) {

	// 遍历 req.Messages 拼接 newMessages
	newMessages := ""
	for _, message := range apiReq.Messages {
		newMessages += message.Content + "\n"
	}
	apiReq.NewMessages = newMessages

	sentMsg, err := bot.SendMessage(newMessages)
	if err != nil {
		logger.Logger.Fatal(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": err.Error(),
		})
	}

	messageChan, stopChan := bot.ReturnChainProcessed(sentMsg.ID)
	defer bot.CleanChans(sentMsg.ID)

	// 定时器
	timer, err := v1.SetTimer(common.RequestOutTimeDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "abnormal timeout setting",
		})
		return
	}

	// 流式返回
	if apiReq.Stream {
		__CompletionsStream(c, apiReq, messageChan, stopChan, timer)
	} else { // 非流式回应
		__CompletionsNoStream(c, apiReq, messageChan, stopChan, timer)
	}
}

func __CompletionsStream(c *gin.Context, apiReq *apireq.Req, messageChan chan *discordgo.MessageUpdate, stopChan chan string, timer *time.Timer) {
	// 响应id
	id := v1.GenerateID(29)
	strLen := 0
	c.Stream(func(w io.Writer) bool {
		select {
		case message := <-messageChan:
			_ = v1.TimerReset(timer, common.StreamRequestOutTime)
			// 如果回复为空则返回
			reply := message.Content
			// 如果回复为空则返回
			if reply == "" || strLen == len(reply) || len(reply) == 0 {
				return true
			}
			// 保留和 messageTemp 不同部分
			content := reply[strLen:]
			// 更新 strLen
			strLen = len(reply)
			for _, v := range content {
				apiRespObj := &apiresp.StreamObj{}
				// id
				apiRespObj.ID = id
				// created
				apiRespObj.Created = time.Now().Unix()
				// object
				apiRespObj.Object = "chat.completion.chunk"
				// choices
				delta := apiresp.StreamDeltaObj{
					Content: string(v),
				}
				choices := apiresp.StreamChoiceObj{
					Delta: delta,
				}
				apiRespObj.Choices = append(apiRespObj.Choices, choices)
				// model
				apiRespObj.Model = apiReq.Model
				// 生成响应
				bytes, err := v1.Obj2Bytes(apiRespObj)
				if err != nil {
					logger.Logger.Debug(err.Error())
					return true
				}
				c.SSEvent("", " "+string(bytes))
			}
			return true // 继续保持流式连接
		case <-timer.C:
			c.SSEvent("", " [DONE]")
			return false
		case <-stopChan:
			apiRespObj := &apiresp.StreamObj{}
			// id
			apiRespObj.ID = id
			// created
			apiRespObj.Created = time.Now().Unix()
			// object
			apiRespObj.Object = "chat.completion.chunk"
			// choices
			delta := apiresp.StreamDeltaObj{
				Content: "",
			}
			choices := apiresp.StreamChoiceObj{
				Delta:        delta,
				FinishReason: "stop",
			}
			apiRespObj.Choices = append(apiRespObj.Choices, choices)
			// model
			apiRespObj.Model = apiReq.Model
			// 生成响应
			bytes, err := v1.Obj2Bytes(apiRespObj)
			if err != nil {
				logger.Logger.Debug(err.Error())
				return true
			}
			c.SSEvent("", " "+string(bytes))
			c.SSEvent("", " [DONE]")
			return false // 关闭流式连接
		}
	})

}

func __CompletionsNoStream(c *gin.Context, apiReq *apireq.Req, replyChan chan *discordgo.MessageUpdate, stopChan chan string, timer *time.Timer) {
	content := ""
	for {
		select {
		case message := <-replyChan:
			_ = v1.TimerReset(timer, common.RequestOutTime)
			// 如果回复为空则返回
			reply := message.Content
			// 如果回复为空则返回
			if reply == "" || len(reply) == 0 {
				continue
			}
			content = reply
		case <-timer.C:
			apiRespObj := &apiresp.JsonObj{}
			// 返回响应
			c.JSON(http.StatusOK, apiRespObj)
			return
		case <-stopChan:
			completionTokens := CountTokens(content)
			promptTokens := CountTokens(apiReq.NewMessages)
			totalTokens := completionTokens + promptTokens

			apiRespObj := &apiresp.JsonObj{}
			// id
			apiRespObj.ID = v1.GenerateID(29)
			// created
			apiRespObj.Created = time.Now().Unix()
			// object
			apiRespObj.Object = "chat.completion"
			// model
			apiRespObj.Model = apiReq.Model
			// usage
			usage := apiresp.JsonUsageObj{
				PromptTokens:     promptTokens,
				CompletionTokens: completionTokens,
				TotalTokens:      totalTokens,
			}
			apiRespObj.Usage = usage
			// choices
			message := apiresp.JsonMessageObj{
				Role:    "assistant",
				Content: content,
			}
			choice := apiresp.JsonChoiceObj{
				Message:      message,
				FinishReason: "stop",
				Index:        0,
			}
			apiRespObj.Choices = append(apiRespObj.Choices, choice)
			// 返回响应
			c.JSON(http.StatusOK, apiRespObj)
			return // 退出 goroutine
		}
	}
}
