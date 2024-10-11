package main

import (
	"io"
	"log"
	"messageredir/cmd/messageredir/api/controllers"
	"messageredir/cmd/messageredir/api/middleware"
	"messageredir/cmd/messageredir/config"
	"messageredir/cmd/messageredir/db/repo"
	"messageredir/cmd/messageredir/services"
	"messageredir/cmd/messageredir/services/models"
	"messageredir/cmd/messageredir/strings"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/natefinch/lumberjack"
)

const ConfigFileName = "messageredir.yaml"

type App struct {
	config    *config.Config
	dbRepo    repo.DbRepo
	tgService services.TelegramService
}

func main() {
	config := config.Load(ConfigFileName)
	setupLogging(&config)
	log.Println("App starting...")

	dbRepo := repo.NewDbRepoGorm(config.DbFileName)
	tgService := services.StartTelegramService(config.TgBotToken)

	app := App{&config, dbRepo, tgService}
	go app.serveRest()
	app.serveBot()
}

func (app App) serveRest() {
	messageController := controllers.NewMessageController(app.config, app.tgService)
	r := mux.NewRouter()
	r.HandleFunc("/{user_token}/smstourlforwarder", messageController.SmsToUrlForwarder).Methods("POST")

	http.Handle("/", middleware.UserAuth(app.config, app.dbRepo, middleware.Logger(r)))

	portStr := ":" + strconv.Itoa(app.config.RestApiPort)
	tlsOn := app.config.TlsCertFile != "" && app.config.TlsKeyFile != ""
	serve := func() error {
		if tlsOn {
			return http.ListenAndServeTLS(portStr, app.config.TlsCertFile, app.config.TlsKeyFile, nil)
		} else {
			return http.ListenAndServe(portStr, nil)
		}
	}

	log.Println("Starting server on", app.config.RestApiPort, "TLS:", tlsOn)
	if err := serve(); err != nil {
		log.Fatal(err)
	}
}

func (app App) serveBot() {
	for {
		msg := <-app.tgService.GetReceiveChan()
		if app.config.LogUserMessages {
			log.Printf("[%s] %s", msg.Username, msg.Text)
		}
		if !app.Process(msg) {
			log.Println("Unknown command")
			app.tgService.Send(models.TelegramMessageOut{
				ChatId: msg.ChatId,
				Text:   strings.BotMsgHelp,
			})
		}
	}
}

func setupLogging(config *config.Config) {
	if config.LogFileName == "" {
		return
	}
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
