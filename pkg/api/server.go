package http

import (
	"Zhooze/pkg/api/handlers"
	"log"

	"Zhooze/pkg/api/routes"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type ServerHTTP struct {
	engine *gin.Engine
}

func NewServerHTTP(userHandler *handlers.UserHandler, productHandler *handlers.ProductHandler, otpHandler *handlers.OtpHandler, adminHandler *handlers.AdminHandler, cartHandler *handlers.CartHandler, orderHandler *handlers.OrderHandler, couponHandler *handlers.CouponHandler, paymentHandler *handlers.PaymentHandler, categoryHandler *handlers.CategoryHandler, offerHandler *handlers.OfferHandler, wishlistHandler *handlers.WishListHandler, walletHandler *handlers.WalletHandler) *ServerHTTP {
	router := gin.New()

	router.LoadHTMLGlob("template/*.html")

	router.Use(gin.Logger())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routes.UserRoutes(router.Group("/user"), userHandler, otpHandler, productHandler, cartHandler, orderHandler, paymentHandler, couponHandler, offerHandler, wishlistHandler, walletHandler)
	routes.AdminRoutes(router.Group("/admin"), adminHandler, productHandler, orderHandler, couponHandler, categoryHandler, offerHandler)

	return &ServerHTTP{engine: router}
}

func (sh *ServerHTTP) Start(infoLog *log.Logger, errorLog *log.Logger) {
	infoLog.Printf("starting server on :3000")
	err := sh.engine.Run(":3000")
	errorLog.Fatal(err)
}
