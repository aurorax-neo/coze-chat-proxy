package botTool

import (
	"coze-chat-proxy/bot/discord"
)

type ProxyBotDB struct {
	index     int
	proxyBots []*discord.ProxyBot
}

func NewDcBotDB() *ProxyBotDB {
	return &ProxyBotDB{
		index:     0,
		proxyBots: make([]*discord.ProxyBot, 0),
	}
}

func (dcBotDB *ProxyBotDB) AddBot(bot *discord.ProxyBot) {
	if dcBotDB.proxyBots == nil {
		dcBotDB.proxyBots = make([]*discord.ProxyBot, 0)
	}
	dcBotDB.proxyBots = append(dcBotDB.proxyBots, bot)
}

func (dcBotDB *ProxyBotDB) GetBot() *discord.ProxyBot {
	// 如果账号已经被使用完毕，那么就从头开始使用
	if dcBotDB.index == len(dcBotDB.proxyBots) {
		dcBotDB.index = 0
	} else if dcBotDB.index > len(dcBotDB.proxyBots) { // 如果账号索引超过了账号数量，那么就返回nil
		return nil
	}
	bot := dcBotDB.proxyBots[dcBotDB.index]
	dcBotDB.index++
	return bot
}
