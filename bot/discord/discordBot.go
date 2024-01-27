package discord

import (
	"context"
	"coze-chat-proxy/logger"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"syscall"
)

type DcBot struct {
	Model                  string                                   `json:"model"`
	BotToken               string                                   `json:"bot_token"`
	CozeBotId              string                                   `json:"coze_bot_id"`
	GuildId                string                                   `json:"guild_id"`
	ChannelID              string                                   `json:"channel_id"`
	MessageUpdateChans     map[string]chan *discordgo.MessageUpdate `json:"-"`
	MessageUpdateStopChans map[string]chan string                   `json:"-"`
	Session                *discordgo.Session                       `json:"-"`
}

func NewDcBot(bot DcBot) *DcBot {
	return &DcBot{
		Model:                  bot.Model,
		BotToken:               bot.BotToken,
		CozeBotId:              bot.CozeBotId,
		GuildId:                bot.GuildId,
		ChannelID:              bot.ChannelID,
		MessageUpdateStopChans: make(map[string]chan string),
		MessageUpdateChans:     make(map[string]chan *discordgo.MessageUpdate),
	}

}

func (dcBot *DcBot) StartBot(ctx context.Context) {
	var err error
	dcBot.Session, err = discordgo.New("Bot " + dcBot.BotToken)
	if err != nil {
		logger.Logger.Fatal(fmt.Sprint("error creating DcBot session,", err))
		return
	}

	// 注册消息处理函数
	dcBot.Session.AddHandler(dcBot.messageUpdate)

	// 打开websocket连接并开始监听
	err = dcBot.Session.Open()
	if err != nil {
		logger.Logger.Fatal(fmt.Sprint("error opening connection,", err))
		return
	}

	logger.Logger.Info(fmt.Sprint(dcBot.Model, " bot is now running. Press CTRL+C to exit."))

	go func() {
		<-ctx.Done()
		if err := dcBot.Session.Close(); err != nil {
			logger.Logger.Fatal(fmt.Sprint("error closing DcBot session,", err))
		}
	}()

	// 等待信号
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

func (dcBot *DcBot) messageUpdate(s *discordgo.Session, m *discordgo.MessageUpdate) {

	if m.Author != nil && s.State.User != nil && m.Author.ID == s.State.User.ID {
		return
	}

	// 检查消息是否是对我们的回复
	for _, mention := range m.Mentions {
		if mention.ID == s.State.User.ID {
			// Message
			if messageUpdateChan, exists := dcBot.MessageUpdateChans[m.ReferencedMessage.ID]; exists {
				messageUpdateChan <- m
			}
			// Components
			if len(m.Message.Components) > 0 {
				stopChan := dcBot.MessageUpdateStopChans[m.ReferencedMessage.ID]
				stopChan <- m.ReferencedMessage.ID
			}
			break
		}
	}

}

func (dcBot *DcBot) SendMessage(message string) (*discordgo.Message, error) {
	if dcBot.Session == nil {
		return nil, fmt.Errorf("dcBot session not initialized")
	}

	// 添加@机器人逻辑
	sentMsg, err := dcBot.Session.ChannelMessageSend(dcBot.ChannelID, fmt.Sprintf("<@%s> %s", dcBot.CozeBotId, message))
	if err != nil {
		return nil, fmt.Errorf("error sending message: %s", err)
	}
	return sentMsg, nil
}

func (dcBot *DcBot) ReturnChainProcessed(msgId string) (chan *discordgo.MessageUpdate, chan string) {
	// 返回 Message
	messageUpdateChans := make(chan *discordgo.MessageUpdate)
	dcBot.MessageUpdateChans[msgId] = messageUpdateChans
	// 返回停止消息
	messageUpdateStopChans := make(chan string)
	dcBot.MessageUpdateStopChans[msgId] = messageUpdateStopChans
	return messageUpdateChans, messageUpdateStopChans
}

func (dcBot *DcBot) CleanChans(msgId string) {
	delete(dcBot.MessageUpdateChans, msgId)
	delete(dcBot.MessageUpdateStopChans, msgId)
}

func (dcBot *DcBot) ChannelCreate(guildID, channelName string) (string, error) {
	// 创建新的频道
	st, err := dcBot.Session.GuildChannelCreate(guildID, channelName, discordgo.ChannelTypeGuildText)
	if err != nil {
		logger.Logger.Error(fmt.Sprint(context.Background(), fmt.Sprintf("创建频道时异常 %s", err.Error())))
		return "", err
	}
	fmt.Println("频道创建成功")
	return st.ID, nil
}

func (dcBot *DcBot) ChannelDel(channelId string) (string, error) {
	// 创建新的频道
	st, err := dcBot.Session.ChannelDelete(channelId)
	if err != nil {
		logger.Logger.Error(fmt.Sprint(context.Background(), fmt.Sprintf("删除频道时异常 %s", err.Error())))
		return "", err
	}
	fmt.Println("删除成功")
	return st.ID, nil
}
