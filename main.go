package main

import (
	"fmt"
	"github.com/gocarina/gocsv"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"os"
	"smdtrns/utils"
)

type SMD struct {
	// Our example struct, you can use "-" to ignore a field
	Code        string `csv:"code"`              // code
	A           string `csv:"-"`                 // A
	Device      string `csv:"Type"`              // Type
	Function    string `csv:"Function"`          // Function
	Description string `csv:"Short description"` // Short description
	Case        string `csv:"Case"`              // Case
	Mnf         string `csv:"Mnf"`               // Mnf
}

var (
	//Log              = utils.Log
	SiteList         map[string]int
	chatID           int64
	telegramBotToken string
	config           *utils.Config
	configFile       string
	HelpMsg          = "Это простой мониторинг доступности сайтов. Он обходит сайты в списке и ждет что он ответит 200, если возвращается не 200 или ошибки подключения, то бот пришлет уведомления в групповой чат\n" +
		"Список доступных комманд:\n" +
		"/site_list - покажет список сайтов в мониторинге и их статусы (про статусы ниже)\n" +
		"/site_add [url] - добавит url в список мониторинга\n" +
		"/site_del [url] - удалит url из списка мониторинга\n" +
		"/help - отобразить это сообщение\n" +
		"\n" +
		"У сайтов может быть несколько статусов:\n" +
		"0 - никогда не проверялся (ждем проверки)\n" +
		"1 - ошибка подключения \n" +
		"200 - ОК-статус" +
		"все остальные http-коды считаются некорректными"
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
	utils.Log.Println("Running 1")

	runbot()

	smddataFile, err := os.OpenFile("test.txt", os.O_RDONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer smddataFile.Close()

	smdlist := []*SMD{}

	if err := gocsv.UnmarshalFile(smddataFile, &smdlist); err != nil { // Load clients from file
		panic(err)
	}
	for _, smd := range smdlist {
		fmt.Println("> ", smd.Device, " - ", smd.Description)
	}

	//if _, err := clientsFile.Seek(0, 0); err != nil { // Go to the start of the file
	//	panic(err)
	//}
	//
	//clients = append(clients, &Client{Id: "12", Name: "John", Age: "21"}) // Add clients
	//clients = append(clients, &Client{Id: "13", Name: "Fred"})
	//clients = append(clients, &Client{Id: "14", Name: "James", Age: "32"})
	//clients = append(clients, &Client{Id: "15", Name: "Danny"})
	//csvContent, err := gocsv.MarshalString(&clients) // Get all clients as CSV string
	////err = gocsv.MarshalFile(&clients, clientsFile) // Use this to save the CSV back to the file
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(csvContent) // Display all clients as CSV string
}
func runbot() {
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
		// Пользователь, который написал боту
		UserName := update.Message.From.UserName

		// ID чата/диалога.
		// Может быть идентификатором как чата с пользователем (тогда он равен UserID) так и публичного чата/канала
		ChatID := update.Message.Chat.ID

		// Текст сообщения
		Text := update.Message.Text

		utils.Log.Printf("[%s] %d %s", UserName, ChatID, Text)

		// Ответим пользователю его же сообщением
		reply := Text
		// Созадаем сообщение
		msg := tgbotapi.NewMessage(ChatID, reply)

		msg.ReplyToMessageID = update.Message.MessageID

		// и отправляем его
		bot.Send(msg)
	}

}
