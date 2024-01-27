package images

import (
	"coze-chat-proxy/bot/discord"
	"coze-chat-proxy/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Dalle3Resp struct {
	Data []Dalle3RespData `json:"data"`
}

type Dalle3RespData struct {
	RevisedPrompt string `json:"revised_prompt"`
	Url           string `json:"url"`
}

func dalle3(c *gin.Context, apiReq *Dalle3Req, bot *discord.DcBot, retryCount int) {
	sentMsg, err := bot.SendMessage(apiReq.Prompt)
	if err != nil {
		logger.Logger.Fatal(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": err.Error(),
		})
	}

	_, embedsChan, stopChan := bot.ReturnChainProcessed(sentMsg.ID)
	defer bot.CleanChans(sentMsg.ID)

	for {
		select {
		case embed := <-embedsChan:
			dalle3Data := Dalle3RespData{
				RevisedPrompt: apiReq.Prompt,
				Url:           embed.Image.URL,
			}
			dalle3Resp := &Dalle3Resp{
				Data: []Dalle3RespData{dalle3Data},
			}
			c.JSON(http.StatusOK, dalle3Resp)
		case <-stopChan:
			return
		}
	}
}