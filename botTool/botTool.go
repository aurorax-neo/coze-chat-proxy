package botTool

import (
	"coze-chat-proxy/bot/discord"
	"coze-chat-proxy/common"
	"coze-chat-proxy/config"
	"coze-chat-proxy/logger"
	"encoding/json"
	"os"
)

var botToolInstance *BotTool

func init() {
	botToolInstance = GetShareChatToolInstance()
	botToolInstance.__GetBots1BotFile()
}

type BotTool struct {
	ProxyBotPool map[string]*ProxyBotDB
}

func GetShareChatToolInstance() *BotTool {
	if botToolInstance == nil {
		botToolInstance = &BotTool{}
	}
	return botToolInstance
}

func (botTool *BotTool) __GetBots1BotFile() {
	// 读取json文件内容
	bytes, err := os.ReadFile(config.CONFIG.BotConfig)
	if err != nil {
		logger.Logger.Error(err.Error())
	}

	// Define a slice to hold the proxyBots
	var dcBots []discord.ProxyBot

	// Unmarshal the JSON data into the proxyBots slice
	err = json.Unmarshal(bytes, &dcBots)
	if err != nil {
		logger.Logger.Error(err.Error())
	}
	botTool.ProxyBotPool = make(map[string]*ProxyBotDB)
	// Get the context
	ctx := common.GetContext()
	// Loop through the proxyBots
	for _, i := range dcBots {
		// Create a new ProxyBot
		bot := discord.NewProxyBot(i)
		// Start the bot
		bot.StartProxyBot(ctx)
		// Add the bot to the botDB
		db, exists := botTool.ProxyBotPool[i.Model]
		if !exists {
			// If not, create a new ProxyBotDB
			db = NewDcBotDB()
			// Add the ProxyBotDB to the botPool
			botTool.ProxyBotPool[i.Model] = db
		}
		// Add the bot to the botDB
		db.AddBot(bot)
	}
}

func (botTool *BotTool) GetBotByModel(model string) *discord.ProxyBot {
	db, exists := botTool.ProxyBotPool[model]
	if !exists {
		return nil
	}
	bot := db.GetBot()
	logger.Logger.Info("model: " + bot.Model + " coze_bot id: " + bot.CozeBotId + " guild_id: " + bot.GuildId + " channel_id: " + bot.ChannelID)
	return bot
}
