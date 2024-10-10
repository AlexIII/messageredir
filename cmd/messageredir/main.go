package main

import (
	"io"
	"log"
	"messageredir/cmd/messageredir/api/controllers"
	"messageredir/cmd/messageredir/api/middleware"
	"messageredir/cmd/messageredir/api/services"
	"messageredir/cmd/messageredir/app"
	"messageredir/cmd/messageredir/db/models"
	"messageredir/cmd/messageredir/db/repo"
	"net/http"
	"os"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gorilla/mux"
	"github.com/natefinch/lumberjack"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const ConfigFileName = "messageredir.yaml"

type CmdContext struct {
	app.Context
	messageService services.MessageService
}

func cmdStart(ctx CmdContext, message *tgbotapi.Message) {
	user := repo.GetOrCreateUser(ctx.Db, message.Chat.ID, message.From.UserName, ctx.Config.UserTokenLength)
	ctx.messageService.SendStr(user.ChatId, "You are good to go!\nYour token: "+user.Token)
}

func cndEnd(ctx CmdContext, message *tgbotapi.Message) {
	repo.DeleteUser(ctx.Db, message.Chat.ID)
	ctx.messageService.SendStr(message.Chat.ID, "You were erased from the system. Goodbye!")
}

func processCommand(ctx CmdContext, message *tgbotapi.Message) bool {
	commands := map[string]func(ctx CmdContext, message *tgbotapi.Message){
		"start": cmdStart,
		"end":   cndEnd,
	}
	cmdName := message.Command()
	if cmd, found := commands[cmdName]; found {
		log.Println("Processing command", cmdName, message.Chat.ID)
		cmd(ctx, message)
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
	db, err := gorm.Open(sqlite.Open(config.DbFileName), &gorm.Config{})
	if err != nil {
		log.Panic("failed to connect to database")
	}
	db.AutoMigrate(&models.User{})

	// Init Telegram bot
	bot, err := tgbotapi.NewBotAPI(config.TgBotToken)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on telegram account %s", bot.Self.UserName)

	ctx := app.Context{Config: config, Db: db}
	tgMsgService := services.NewMessageServiceTelegram()

	// Init HTTP server
	go func() {
		messageController := controllers.NewMessageController(ctx, tgMsgService)
		r := mux.NewRouter()
		r.HandleFunc("/{user_token}/smstourlforwarder", messageController.SmsToUrlForwarder).Methods("POST")

		http.Handle("/", middleware.UserAuth(ctx, middleware.Logger(r)))

		portStr := ":" + strconv.Itoa(config.RestApiPort)
		tlsOn := ctx.Config.TlsCertFile != "" && ctx.Config.TlsKeyFile != ""
		serve := func() error {
			if tlsOn {
				return http.ListenAndServeTLS(portStr, ctx.Config.TlsCertFile, ctx.Config.TlsKeyFile, nil)
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
		cmdCtx := CmdContext{ctx, tgMsgService}
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
