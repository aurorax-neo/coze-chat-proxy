package chat

import (
	"coze-chat-proxy/botTool"
	"coze-chat-proxy/v1/chat/apireq"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Completions(c *gin.Context) {
	// 从请求中获取参数
	apiReq := &apireq.Req{}

	err := json.NewDecoder(c.Request.Body).Decode(&apiReq)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"detail": "Invalid parameter",
		})
		return
	}
	// 获取 botTool
	botTool1 := botTool.GetShareChatToolInstance()
	// 从请求中获取 model
	bot := botTool1.GetBotByModel(apiReq.Model)
	if bot == nil {
		c.JSON(http.StatusOK, gin.H{
			"detail": "model is not supported",
		})
		return
	}
	gpt(c, apiReq, bot)
}
