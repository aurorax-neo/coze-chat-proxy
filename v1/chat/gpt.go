package chat

import (
	"coze-chat-proxy/bot/discord"
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

func gpt(c *gin.Context, apiReq *apireq.Req, bot *discord.ProxyBot, retryCount int) {

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

	// 流式返回
	if apiReq.Stream {
		__CompletionsStream(c, apiReq, messageChan, stopChan)
	} else { // 非流式回应
		__CompletionsNoStream(c, apiReq, messageChan, stopChan)
	}
}

func __CompletionsStream(c *gin.Context, apiReq *apireq.Req, messageChan chan *discordgo.MessageUpdate, stopChan chan string) {
	// 响应id
	id := v1.GenerateID(29)
	strLen := 0
	c.Stream(func(w io.Writer) bool {
		select {
		case message := <-messageChan:
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
			apiRespObj := &apiresp.StreamObj{}
			// id
			apiRespObj.ID = id
			// created
			apiRespObj.Created = time.Now().Unix()
			// object
			apiRespObj.Object = "chat.completion.chunk"
			// choices
			delta := apiresp.StreamDeltaObj{
				Content: content,
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
			c.SSEvent("", string(bytes))
			return true // 继续保持流式连接
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
			c.SSEvent("", string(bytes))
			c.SSEvent("", "[DONE]")
			return false // 关闭流式连接
		}
	})

}

func __CompletionsNoStream(c *gin.Context, apiReq *apireq.Req, replyChan chan *discordgo.MessageUpdate, stopChan chan string) {
	content := ""
	for {
		select {
		case message := <-replyChan:
			// 如果回复为空则返回
			reply := message.Content
			// 如果回复为空则返回
			if reply == "" || len(reply) == 0 {
				continue
			}
			content = reply
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
			apiRespObj.Object = "chat.completion.chunk"
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

var (
	gptTurboReqStr = `
	{
	  "action": "next",
	  "messages": [
		{
		  "id": "aaa2b2cc-e7e9-47c5-8171-0ff8a6d9d6d3",
		  "author": {
			"role": "user"
		  },
		  "content": {
			"content_type": "text",
			"parts": [
			  ""
			]
		  },
		  "metadata": {}
		}
	  ],
	  "parent_message_id": "aaa1403d-c61e-4818-90e0-93a99465aec6",
	  "model": "gpt-4",
	  "timezone_offset_min": -480,
	  "suggestions": [
		""
	  ],
	  "history_and_training_disabled": true,
	  "conversation_mode": {
		"kind": "primary_assistant"
	  },
	  "force_paragen": false,
	  "force_rate_limit": false,
	  "arkose_token": ""
	}`
	gpt4ReqStr = `
	{
		"action": "next",
		"messages": [
			{
				"id": "aaa2b182-d834-4f30-91f3-f791fa953204",
				"author": {
					"role": "user"
				},
				"content": {
					"content_type": "text",
					"parts": [
						"画一只猫1231231231"
					]
				},
				"metadata": {}
			}
		],
		"parent_message_id": "aaa11581-bceb-46c5-bc76-cb84be69725e",
		"model": "gpt-4-gizmo",
		"timezone_offset_min": -480,
		"suggestions": [],
		"history_and_training_disabled": true,
		"conversation_mode": {
			"gizmo": {
				"gizmo": {
					"id": "g-YyyyMT9XH",
					"organization_id": "org-OROoM5KiDq6bcfid37dQx4z4",
					"short_url": "g-YyyyMT9XH-chatgpt-classic",
					"author": {
						"user_id": "user-u7SVk5APwT622QC7DPe41GHJ",
						"display_name": "ChatGPT",
						"selected_display": "name",
						"is_verified": true
					},
					"voice": {
						"id": "ember"
					},
					"display": {
						"name": "ChatGPT Classic",
						"description": "The latest version of GPT-4 with no additional capabilities",
						"welcome_message": "Hello",
						"profile_picture_url": "https://files.oaiusercontent.com/file-i9IUxiJyRubSIOooY5XyfcmP?se=2123-10-13T01%3A11%3A31Z&sp=r&sv=2021-08-06&sr=b&rscc=max-age%3D31536000%2C%20immutable&rscd=attachment%3B%20filename%3Dgpt-4.jpg&sig=ZZP%2B7IWlgVpHrIdhD1C9wZqIvEPkTLfMIjx4PFezhfE%3D",
						"categories": []
					},
					"share_recipient": "link",
					"updated_at": "2023-11-06T01:11:32.191060+00:00",
					"last_interacted_at": "2023-11-18T07:50:19.340421+00:00",
					"tags": [
						"public",
						"first_party"
					]
				},
				"tools": [],
				"files": [],
				"product_features": {
					"attachments": {
						"type": "retrieval",
						"accepted_mime_types": [
							"text/x-c",
							"text/html",
							"application/x-latext",
							"text/plain",
							"text/x-ruby",
							"text/x-typescript",
							"text/x-c++",
							"text/x-java",
							"text/x-sh",
							"application/vnd.openxmlformats-officedocument.presentationml.presentation",
							"text/x-script.python",
							"text/javascript",
							"text/x-tex",
							"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
							"application/msword",
							"application/pdf",
							"text/x-php",
							"text/markdown",
							"application/json",
							"text/x-csharp"
						],
						"image_mime_types": [
							"image/jpeg",
							"image/png",
							"image/gif",
							"image/webp"
						],
						"can_accept_all_mime_types": true
					}
				}
			},
			"kind": "gizmo_interaction",
			"gizmo_id": "g-YyyyMT9XH"
		},
		"force_paragen": false,
		"force_rate_limit": false,
		"arkose_token": ""
	}`
	apiRespStr = `{
		"id": "chatcmpl-LLKfuOEHqVW2AtHks7wAekyrnPAoj",
		"object": "chat.completion",
		"created": 1689864805,
		"model": "gpt-3.5-turbo",
		"usage": {
			"prompt_tokens": 0,
			"completion_tokens": 0,
			"total_tokens": 0
		},
		"choices": [
			{
				"message": {
					"role": "assistant",
					"content": "Hello! How can I assist you today?"
				},
				"finish_reason": "stop",
				"index": 0
			}
		]
	}`
	apiRespStrStream = `{
		"id": "chatcmpl-afUFyvbTa7259yNeDqaHRBQxH2PLH",
		"object": "chat.completion.chunk",
		"created": 1689867370,
		"model": "gpt-3.5-turbo",
		"choices": [
			{
				"delta": {
					"content": "Hello"
				},
				"index": 0,
				"finish_reason": null
			}
		]
	}`
	ApiRespStrStreamEnd = `{"id":"apirespid","object":"chat.completion.chunk","created":apicreated,"model": "apirespmodel","choices":[{"delta": {},"index": 0,"finish_reason": "stop"}]}`
)
