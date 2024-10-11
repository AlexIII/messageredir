package main

import (
	"io"
	"log"
	"messageredir/cmd/messageredir/api/controllers"
	"messageredir/cmd/messageredir/api/middleware"
	"messageredir/cmd/messageredir/api/services"
	"messageredir/cmd/messageredir/app"
	"messageredir/cmd/messageredir/db/repo"
	"net/http"
	"os"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gorilla/mux"
	"github.com/natefinch/lumberjack"
)

const ConfigFileName = "messageredir.yaml"

type Command struct {
	app    *app.Config
	db     repo.DbRepo
	msg    services.MessageService
	config *app.Config
}

func (ctx Command) start(message *tgbotapi.Message) {
	user := ctx.db.GetOrCreateUser(message.Chat.ID, message.From.UserName, ctx.config.UserTokenLength)
	ctx.msg.SendStr(user.ChatId, "You are good to go!\nYour token: "+user.Token)
}

func (ctx Command) end(message *tgbotapi.Message) {
	ctx.db.DeleteUser(message.Chat.ID)
	ctx.msg.SendStr(message.Chat.ID, "You were erased from the system. Goodbye!")
}

func processCommand(cmd Command, message *tgbotapi.Message) bool {
	commands := map[string]func(message *tgbotapi.Message){
		"start": cmd.start,
		"end":   cmd.end,
	}
	cmdName := message.Command()
	if run, found := commands[cmdName]; found {
		log.Println("Processing command", cmdName, message.Chat.ID)
		run(message)
		return true
	}
	return false
}

func setupLogging() {
	logFile := &lumberjack.Logger{
		Filename:   "app.log",
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
		Compress:   true,
	}
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)
}

func main() {
	setupLogging()
	log.Println("App starting...")
	config := app.LoadConfig(ConfigFileName)

	// Init DB
	var dbRepo repo.DbRepo = repo.NewDbRepoGorm(config.DbFileName)

	// Init Telegram bot
	bot, err := tgbotapi.NewBotAPI(config.TgBotToken)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on telegram account %s", bot.Self.UserName)

	tgMsgService := services.NewMessageServiceTelegram()

	// Init HTTP server
	go func() {
		messageController := controllers.NewMessageController(config, tgMsgService)
		r := mux.NewRouter()
		r.HandleFunc("/{user_token}/smstourlforwarder", messageController.SmsToUrlForwarder).Methods("POST")

		http.Handle("/", middleware.UserAuth(&config, dbRepo, middleware.Logger(r)))

		portStr := ":" + strconv.Itoa(config.RestApiPort)
		tlsOn := config.TlsCertFile != "" && config.TlsKeyFile != ""
		serve := func() error {
			if tlsOn {
				return http.ListenAndServeTLS(portStr, config.TlsCertFile, config.TlsKeyFile, nil)
			} else {
				return http.ListenAndServe(portStr, nil)
			}
		}

		log.Println("Starting server on", config.RestApiPort, "TLS:", tlsOn)
		if err := serve(); err != nil {
			log.Fatal(err)
		}
	}()

	// Start listening to messages
	go func() {
		cmdCtx := Command{app: &config, db: dbRepo, msg: tgMsgService, config: &config}
		updConf := tgbotapi.NewUpdate(0)
		updConf.Timeout = 60
		updates := bot.GetUpdatesChan(updConf)
		for {
			select {
			case msg := <-tgMsgService.GetOutMessages():
				bot.Send(tgbotapi.NewMessage(msg.ChatId, msg.Message))

			case update := <-updates:
				if update.Message != nil { // Got a message
					if config.LogUserMessages {
						log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
					}
					if !processCommand(cmdCtx, update.Message) {
						log.Println("Unknown command")
						tgMsgService.SendStr(update.Message.Chat.ID, "Hi! This is redir bot.\nSend /start command to get your token or /end to leave the service.")
					}
				}
			}
		}
	}()

	for {
		time.Sleep(time.Hour)
	}
}
