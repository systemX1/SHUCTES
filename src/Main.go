package main

import (
	"SHUCTES/src/config"
	"SHUCTES/src/database"
	"SHUCTES/src/log"
	"SHUCTES/src/router"
	"github.com/gin-gonic/gin"
)

func initApplication() (app *gin.Engine) {
	//ReleaseMode
	gin.SetMode(gin.ReleaseMode)
	app = gin.New()

	//init
	config.InitConfig()
	log.InitLog()
	database.InitDatabase()
	router.RegisterRouter(app)

	return app
}

func main() {
	app := initApplication()
	log.Logger.Infof("Starting...")
	if err := app.Run(config.Conf.Sev.Port); err != nil {
		log.Logger.Panicf("server run err: %s", err)
	}
	_ = database.DB.Close()
}
