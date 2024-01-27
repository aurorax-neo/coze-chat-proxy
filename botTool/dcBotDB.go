package botTool

import (
	"coze-chat-proxy/bot/discord"
)

type DcBotDB struct {
	index  int
	dcBots []*discord.DcBot
}

func NewDcBotDB() *DcBotDB {
	return &DcBotDB{
		index:  0,
		dcBots: make([]*discord.DcBot, 0),
	}
}

func (dcBotDB *DcBotDB) AddBot(bot *discord.DcBot) {
	if dcBotDB.dcBots == nil {
		dcBotDB.dcBots = make([]*discord.DcBot, 0)
	}
	dcBotDB.dcBots = append(dcBotDB.dcBots, bot)
}

func (dcBotDB *DcBotDB) GetBot() *discord.DcBot {
	// 如果账号已经被使用完毕，那么就从头开始使用
	if dcBotDB.index == len(dcBotDB.dcBots) {
		dcBotDB.index = 0
	} else if dcBotDB.index > len(dcBotDB.dcBots) { // 如果账号索引超过了账号数量，那么就返回nil
		return nil
	}
	bot := dcBotDB.dcBots[dcBotDB.index]
	dcBotDB.index++
	return bot
}
