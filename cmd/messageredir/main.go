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
	"messageredir/cmd/messageredir/userstrings"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/natefinch/lumberjack"
)

const configFileName = "messageredir.yaml"
const httpMaxBodySize = 200 * 1024 // 200 KB

type App struct {
	config   *config.Config
	db       repo.DbRepo
	telegram services.TelegramService
}

func main() {
	config := config.Load(configFileName)
	setupLogging(&config)
	log.Printf("App starting. Config: %+v", config)

	dbRepo := repo.NewDbRepoGorm(config.DbFileName)
	tgService := services.StartTelegramService(config.TgBotToken)

	app := App{&config, dbRepo, tgService}
	go app.serveRest()
	app.serveBot()
}

func (app App) serveRest() {
	messageController := controllers.NewMessageController(app.config, app.db, app.telegram)
	r := mux.NewRouter()
	r.HandleFunc("/{user_token}/smstourlforwarder", messageController.SmsToUrlForwarder).Methods("POST")

	http.Handle("/",
		middleware.Recover(
			middleware.UserAuth(app.config, app.db,
				middleware.Logger(
					middleware.BodyLimit(httpMaxBodySize,
						r)))))

	portStr := ":" + strconv.Itoa(app.config.RestApiPort)
	serve := func() error {
		if app.config.IsTlsEnabled() {
			return http.ListenAndServeTLS(portStr, app.config.TlsCertFile, app.config.TlsKeyFile, nil)
		} else {
			return http.ListenAndServe(portStr, nil)
		}
	}

	log.Println("Starting server on", app.config.RestApiPort, "TLS:", app.config.IsTlsEnabled())
	if err := serve(); err != nil {
		log.Fatal(err)
	}
}

func (app App) serveBot() {
	for {
		msg := <-app.telegram.GetReceiveChan()
		if app.config.LogUserMessages {
			log.Printf("[%s] %s", msg.Username, msg.Text)
		}
		if !app.Process(msg) {
			log.Println("Unknown command")
			app.telegram.Send(models.TelegramMessageOut{
				ChatId: msg.ChatId,
				Text:   userstrings.BotMsgHelp,
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
