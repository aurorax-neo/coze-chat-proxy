package botTool

import (
	"coze-chat-proxy/bot/discord"
	"coze-chat-proxy/common"
	"coze-chat-proxy/config"
	"coze-chat-proxy/logger"
	"encoding/json"
	"log"
	"os"
)

var botToolInstance *BotTool

func init() {
	botToolInstance = GetShareChatToolInstance()
	botToolInstance.__GetBots1BotFile()
}

type BotTool struct {
	BotDB map[string]*DcBotDB
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
		log.Fatal(err)
	}

	// Define a slice to hold the dcBots
	var dcBots []discord.DcBot

	// Unmarshal the JSON data into the dcBots slice
	err = json.Unmarshal(bytes, &dcBots)
	if err != nil {
		logger.Logger.Fatal(err.Error())
	}
	botTool.BotDB = make(map[string]*DcBotDB)
	// Get the context
	ctx := common.GetContext()
	// Loop through the dcBots
	for _, i := range dcBots {
		// Create a new DcBot
		bot := discord.NewDcBot(i)
		// Start the bot
		go bot.StartBot(ctx)
		// Add the bot to the botDB
		db, exists := botTool.BotDB[i.Model]
		if !exists {
			// If not, create a new DcBotDB
			db = NewDcBotDB()
			botTool.BotDB[i.Model] = db
		}
		// Add the bot to the botDB
		db.AddBot(bot)
	}
}

func (botTool *BotTool) GetBotByModel(model string) *discord.DcBot {
	db, exists := botTool.BotDB[model]
	if !exists {
		return nil
	}
	return db.GetBot()
}
