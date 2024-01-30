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

type DcBotDB struct {
	DB map[string]*DcBot
}

var DcBotDBInstance *DcBotDB

func init() {
	DcBotDBInstance = &DcBotDB{
		DB: make(map[string]*DcBot),
	}
}

type DcBot struct {
	BotToken string
	Session  *discordgo.Session
}

func NeDcBot(token string) *DcBot {
	return &DcBot{
		BotToken: token,
	}
}

func (dcBot *DcBot) StartBot(ctx context.Context, proxyBot *ProxyBot) {
	// 判断是否已经存在
	if _, ok := DcBotDBInstance.DB[dcBot.BotToken]; ok {
		return
	}

	var err error
	dcBot.Session, err = discordgo.New("Bot " + dcBot.BotToken)
	if err != nil {
		logger.Logger.Error(fmt.Sprint("error creating DcBot bot,", err))
		return
	}
	// 添加到数据库
	DcBotDBInstance.DB[dcBot.BotToken] = dcBot

	// 打开websocket连接并开始监听
	err = dcBot.Session.Open()
	if err != nil {
		logger.Logger.Error(fmt.Sprint("error opening DcBot connection,", err))
		return
	}
	// 注册消息处理函数
	dcBot.Session.AddHandler(proxyBot.messageUpdate)

	// Get the bot user
	user := dcBot.Session.State.User
	if user != nil {
		logger.Logger.Info(fmt.Sprintf("ID: " + user.ID + " discord bot " + " is running"))
	}

	go func() {
		<-ctx.Done()
		if err := dcBot.Session.Close(); err != nil {
			logger.Logger.Error(fmt.Sprint("error closing DcBot bot,", err))
		}
	}()

	go func() {
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
		<-sc
	}()
}
