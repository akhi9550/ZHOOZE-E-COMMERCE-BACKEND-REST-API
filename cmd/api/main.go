package main

import (
	"Zhooze/cmd/api/docs"
	config "Zhooze/pkg/config"
	di "Zhooze/pkg/di"
	"log"
	"os"

	_ "github.com/swaggo/files"
	_ "github.com/swaggo/gin-swagger"
)

func main() {

	// // swagger 2.0 Meta Information
	docs.SwaggerInfo.Title = "Zhooze - E-commerce"
	docs.SwaggerInfo.Description = "Zhooze is an E-commerce platform to purchasing and selling shoes"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:3000"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http","https"}

	config, configErr := config.LoadConfig()
	if configErr != nil {
		log.Fatal("cannot load config: ", configErr)
	}

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	server, diErr := di.InitializeAPI(config)

	if diErr != nil {

		log.Fatal("cannot start server: ", diErr)
	} else {
		server.Start(infoLog, errorLog)
	}

}
