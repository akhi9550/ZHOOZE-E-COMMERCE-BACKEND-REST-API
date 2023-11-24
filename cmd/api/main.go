package main

import (
	"Zhooze/cmd/api/docs"
	_ "Zhooze/cmd/docs"
	"Zhooze/config"
	di "Zhooze/pkg/di"
	"log"
	"os"
)

// @title Go + Gin Zhooze E-Commerce API
// @version 1.0.0
// @description Zhooze is an E-commerce platform to purchase and sell shoes
// @contact.name API Support
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @host localhost:8000
// @BasePath /
// @query.collection.format multi

func main() {

	// // swagger 2.0 Meta Information
	docs.SwaggerInfo.Title = "Zhooze - E-commerce"
	docs.SwaggerInfo.Description = "Zhooze - E-commerce"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:3000"
	docs.SwaggerInfo.BasePath = ""
	docs.SwaggerInfo.Schemes = []string{"http"}

	config, configErr := config.LoadConfig()
	if configErr != nil {
		log.Fatal("cannot load config: ", configErr)
	}

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	server, productRepo, diErr := di.InitializeAPI(config)

	if diErr != nil {
		log.Fatal("cannot start server: ", diErr)
	} else {
		server.Start(productRepo, infoLog, errorLog)
	}

}

// cfig, err := config.LoadConfig()
// if err != nil {
// 	log.Fatalf("Error loading the config file")
// }
// fmt.Println(cfig)
// db, err := db.ConnectDatabase(cfig)
// if err != nil {
// 	log.Fatalf("Error connecting to the database:%v", err)
// }
// router := gin.Default()
// router.LoadHTMLGlob("template/*")
// userGroup := router.Group("/user")
// adminGroup := router.Group("/admin")
// routes.UserRoutes(userGroup, db)
// routes.AdminRoutes(adminGroup, db)
// router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
// router.Static("uploads", "./uploads")

// listenAdder := fmt.Sprintf("%s:%s", cfig.DBPort, cfig.DBHost)
// fmt.Printf("Starting server on %s..\n", cfig.BASE_URL)
// if err := router.Run(cfig.BASE_URL); err != nil {
// 	log.Fatalf("Error starting server on %s:%v", listenAdder, err)
// }
// }
