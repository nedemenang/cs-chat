package main

import (
	"chat-api/infrastructure"
	"chat-api/infrastructure/common"
	"chat-api/infrastructure/database"
	"chat-api/infrastructure/log"
	"chat-api/infrastructure/router"
	"chat-api/infrastructure/validation"
	"os"
	"time"
)

func init() {

	common.LoadEnvVars()

}

func main() {
	var app = infrastructure.NewConfig().
		Name(os.Getenv("APP_NAME")).
		ContextTimeout(30 * time.Second).
		Logger(log.InstanceLogrusLogger).
		Validator(validation.InstanceGoPlayground).
		DbNoSQL(database.InstanceMongoDB)

	app.WebServerPort(os.Getenv("PORT")).
		WebServer(router.InstanceGin).
		Start()
}
