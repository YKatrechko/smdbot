package main

import (
	"fmt"
	"github.com/YKatrechko/smdbot/utils"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/YKatrechko/smdbot/dbase"
	"strings"
)

var (
	chatID  int64
	config  *utils.Config
	HelpMsg = "This bot search smd transistor by code\n" +
		"List of available commands:\n" +
		"*/search* `[code]` - search smd transistor by code\n" +
		"*/help* - show this message\n"
)

func initconf() {
	configFile, logFile := utils.Initopts()
	println("Configuration file:", configFile)

	config = utils.LoadConfig(configFile)
	utils.DefaultConfig(config)

	if logFile != "" {
		config.LogFile = logFile
	}

	if config.Token == "" {
		panic("Token can't be empty")
	}
	//utils.SaveConfig(config)
}

///
func main() {

	println("Initialization...")
	initconf()
	defer utils.Initlog(config.LogFile)()

	runbot()
}

func runbot() {
	db := dbase.InitDB(config.DBFile)
	defer db.Close()
	//CreateTable(db)

	utils.Log.Println("Running bot...")
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		utils.Log.Panic(err)
	}

	bot.Debug = true
	utils.Log.Printf("Authorized on account %s", bot.Self.UserName)

	// инициализируем канал, куда будут прилетать обновления от API
	var ucfg tgbotapi.UpdateConfig = tgbotapi.NewUpdate(0)
	ucfg.Timeout = 60
	updates, err := bot.GetUpdatesChan(ucfg)
	if err != nil {
		utils.Log.Println(err)
	}
	// В канал updates будут приходить все новые сообщения.
	for update := range updates {
		reply := ""
		if update.Message == nil {
			continue
		}
		// Пользователь, который написал боту
		UserName := update.Message.From.UserName

		// ID чата/диалога. Может быть идентификатором как чата с пользователем (тогда он равен UserID) так и публичного чата/канала
		ChatID := update.Message.Chat.ID

		utils.Log.Printf("Req [%s] %s", UserName, update.Message.Text)

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "search":
				code := update.Message.CommandArguments()
				if code == "" {
					reply = "*code* can't be empty to search\n"
					break
				}
				if strings.Contains(update.Message.CommandArguments(), " ") {
					code = strings.SplitN(code, " ", 2)[0]
				}
				utils.Log.Printf("Code [%s]", code)

				//update.Message.Text;
				readItems := dbase.ReadItemsByCode(db, code)
				utils.Log.Println(readItems)
				if len(readItems) == 0 {
					reply = fmt.Sprintf("device with code - *%s* isn't found", code)
				}
				for _, item := range readItems {
					str := fmt.Sprintf(
						"*Code* _%s_\n"+
							"*Device* _%s_\n"+
							"*Function* _%s_\n"+
							"*Description* _%s_\n",
						item.Code,
						item.Device,
						item.Function,
						item.Description)
					reply += str + "\n"
				}
				break
			case "help":
				reply = HelpMsg
			}
		} else {
			// Ответим пользователю его же сообщением
			reply = fmt.Sprintf("Your text: _%s_\n", update.Message.Text)
		}
		utils.Log.Printf("Resp: \n%s", reply)

		// Созадаем сообщение
		msg := tgbotapi.NewMessage(ChatID, reply)
		msg.ReplyToMessageID = update.Message.MessageID
		msg.ParseMode = "markdown"
		// и отправляем его
		bot.Send(msg)
	}

}
