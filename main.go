package main

import (
	"fmt"
	"github.com/YKatrechko/smdbot/utils"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/YKatrechko/smdbot/dbase"
	"strings"
)

var (
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
}

func main() {

	println("Initialization...")
	initconf()
	defer utils.Initlog(config)()

	runbot()
}

func runbot() {
	db := dbase.InitDB(config.DBFile)
	defer db.Close()

	utils.Log.Println("Running bot...")
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		utils.Log.Panic(err)
	}

	bot.Debug = config.Debug
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
		utils.Log.Printf("Req [%s](%s,%s) %s", update.Message.From.UserName,update.Message.From.FirstName,update.Message.From.LastName, update.Message.Text)
		ChatID := update.Message.Chat.ID

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
					break
				}
				str := fmt.Sprintf("``` |%6s |%14s |%22s |%7s\n",
					"code", "device", "function", "case")
				for _, item := range readItems {
					str += fmt.Sprintf(
						" |%6s |%14s |%22s |%7s\nDescr: %s\n",
						item.Code,
						item.Device,
						item.Function,
						item.Case,
						item.Description)
				}
				reply = str + "```"
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
