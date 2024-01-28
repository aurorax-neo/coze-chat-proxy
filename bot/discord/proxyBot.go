package discord

import (
	"context"
	"coze-chat-proxy/logger"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

type ProxyBot struct {
	Model                  string                                   `json:"model"`
	BotToken               string                                   `json:"bot_token"`
	CozeBotId              string                                   `json:"coze_bot_id"`
	GuildId                string                                   `json:"guild_id"`
	ChannelID              string                                   `json:"channel_id"`
	MessageUpdateChans     map[string]chan *discordgo.MessageUpdate `json:"-"`
	MessageUpdateStopChans map[string]chan string                   `json:"-"`
	Session                *discordgo.Session                       `json:"-"`
}

func NewProxyBot(bot ProxyBot) *ProxyBot {
	return &ProxyBot{
		Model:                  bot.Model,
		BotToken:               bot.BotToken,
		CozeBotId:              bot.CozeBotId,
		GuildId:                bot.GuildId,
		ChannelID:              bot.ChannelID,
		MessageUpdateStopChans: make(map[string]chan string),
		MessageUpdateChans:     make(map[string]chan *discordgo.MessageUpdate),
	}

}

func (proxyBot *ProxyBot) StartProxyBot(ctx context.Context) {
	// 判断是否已经存在
	if bot, ok := DcBotDBInstance.DB[proxyBot.BotToken]; ok {
		proxyBot.Session = bot.Session
		// 注册消息处理函数
		proxyBot.Session.AddHandler(proxyBot.messageUpdate)
	} else {
		dcBot := NeDcBot(proxyBot.BotToken)
		dcBot.StartBot(ctx, proxyBot)
		proxyBot.Session = dcBot.Session
	}
	logger.Logger.Info(fmt.Sprintf("ID: " + proxyBot.CozeBotId + " proxy bot is running"))
}

func (proxyBot *ProxyBot) messageUpdate(s *discordgo.Session, m *discordgo.MessageUpdate) {

	if m.Author != nil && s.State.User != nil && m.Author.ID == s.State.User.ID {
		return
	}

	// 检查消息是否是对我们的回复
	for _, mention := range m.Mentions {
		if mention.ID == s.State.User.ID {
			// Message
			if messageUpdateChan, exists := proxyBot.MessageUpdateChans[m.ReferencedMessage.ID]; exists {
				messageUpdateChan <- m
			}
			// Components
			if len(m.Message.Components) > 0 {
				stopChan := proxyBot.MessageUpdateStopChans[m.ReferencedMessage.ID]
				stopChan <- m.ReferencedMessage.ID
			}
			break
		}
	}

}

func (proxyBot *ProxyBot) SendMessage(message string) (*discordgo.Message, error) {
	if proxyBot.Session == nil {
		return nil, fmt.Errorf("proxyBot session not initialized")
	}

	// 添加@机器人逻辑
	sentMsg, err := proxyBot.Session.ChannelMessageSend(proxyBot.ChannelID, fmt.Sprintf("<@%s> %s", proxyBot.CozeBotId, message))
	if err != nil {
		return nil, fmt.Errorf("error sending message: %s", err)
	}
	return sentMsg, nil
}

func (proxyBot *ProxyBot) ReturnChainProcessed(msgId string) (chan *discordgo.MessageUpdate, chan string) {
	// 返回 Message
	messageUpdateChans := make(chan *discordgo.MessageUpdate)
	proxyBot.MessageUpdateChans[msgId] = messageUpdateChans
	// 返回停止消息
	messageUpdateStopChans := make(chan string)
	proxyBot.MessageUpdateStopChans[msgId] = messageUpdateStopChans
	return messageUpdateChans, messageUpdateStopChans
}

func (proxyBot *ProxyBot) CleanChans(msgId string) {
	delete(proxyBot.MessageUpdateChans, msgId)
	delete(proxyBot.MessageUpdateStopChans, msgId)
}
