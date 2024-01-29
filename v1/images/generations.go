package images

import (
	"coze-chat-proxy/botTool"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Dalle3Req struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	N      int    `json:"n"`
	Size   string `json:"size"`
}

func Generations(c *gin.Context) {
	// 从请求中获取参数
	apiReq := &Dalle3Req{}

	err := c.ShouldBindJSON(apiReq)
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
	dalle3(c, apiReq, bot, 1)
}
