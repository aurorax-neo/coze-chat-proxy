package chat

import (
	"coze-chat-proxy/logger"
	"github.com/pkoukk/tiktoken-go"
)

var (
	Tke *tiktoken.Tiktoken
)

func init() {
	// gpt-3.5-turbo encoding
	tke, err := tiktoken.GetEncoding("cl100k_base")
	if err != nil {
		logger.Logger.Error(err.Error())
	}
	Tke = tke

}

func CountTokens(text string) int {
	return len(Tke.Encode(text, nil, nil))
}
