package images

import (
	"coze-chat-proxy/botTool"
	"github.com/gin-gonic/gin"
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
		c.JSON(200, gin.H{
			"message": "无效的参数",
			"success": false,
		})
		return
	}

	// 获取 botTool
	botTool1 := botTool.GetShareChatToolInstance()
	// 从请求中获取 model
	bot := botTool1.GetBotByModel(apiReq.Model)
	if bot == nil {
		c.JSON(200, gin.H{
			"message": "model is not supported",
			"success": false,
		})
		return
	}
	dalle3(c, apiReq, bot, 1)
}
