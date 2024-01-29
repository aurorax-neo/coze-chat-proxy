package images

import (
	"coze-chat-proxy/bot/discord"
	"coze-chat-proxy/common"
	"coze-chat-proxy/logger"
	v1 "coze-chat-proxy/v1"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Dalle3Resp struct {
	Created int64            `json:"created"`
	Data    []Dalle3RespData `json:"data"`
}

type Dalle3RespData struct {
	Url     string `json:"url"`
	B64Json string `json:"b64_json"`
}

func dalle3(c *gin.Context, apiReq *Dalle3Req, bot *discord.ProxyBot) {
	sentMsg, err := bot.SendMessage(apiReq.Prompt)
	if err != nil {
		logger.Logger.Fatal(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": err.Error(),
		})
	}

	messageChans, stopChan := bot.ReturnChainProcessed(sentMsg.ID)
	defer bot.CleanChans(sentMsg.ID)

	// 定时器
	timer, err := v1.SetTimer(common.RequestOutTimeDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "abnormal timeout setting",
		})
		return
	}

	for {
		select {
		case messageChan := <-messageChans:
			_ = v1.TimerReset(timer, common.RequestOutTimeDuration)
			if len(messageChan.Embeds) == 0 {
				continue
			}
			dalle3Resp := &Dalle3Resp{
				Created: time.Now().Unix(),
				Data:    []Dalle3RespData{},
			}
			for _, embed := range messageChan.Embeds {
				if url := embed.Image.URL; url != "" {
					b64Json := ""
					//	下载图片
					bys := DownloadImg(url)
					if bys != nil {
						// 上传图片到图床
						if url_ := UploadImg2PicBed(bys); url_ != "" {
							url = url_
						}
						//	获取图片的base64
						b64Json = base64.StdEncoding.EncodeToString(bys)
					}
					dalle3Data := Dalle3RespData{
						Url:     url,
						B64Json: b64Json,
					}
					dalle3Resp.Created = time.Now().Unix()
					dalle3Resp.Data = append(dalle3Resp.Data, dalle3Data)
				}
			}
			c.JSON(http.StatusOK, dalle3Resp)
			return
		case <-timer.C:
			c.JSON(http.StatusInternalServerError, gin.H{
				"detail": "time out",
			})
			return
		case <-stopChan:
			c.JSON(http.StatusInternalServerError, gin.H{
				"detail": "no data generated",
			})
			return
		}
	}
}
