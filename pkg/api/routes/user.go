package routes

import (
	"Zhooze/pkg/api/handlers"
	"Zhooze/pkg/api/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.RouterGroup, userHandler *handlers.UserHandler, otpHandler *handlers.OtpHandler, productHandler *handlers.ProductHandler, cartHandler *handlers.CartHandler, orderHandler *handlers.OrderHandler, paymentHandler *handlers.PaymentHandler, couponHandler *handlers.CouponHandler, offerHandler *handlers.OfferHandler, wishlistHandler *handlers.WishListHandler, walletHandler *handlers.WalletHandler) {

	r.POST("/signup", userHandler.UserSignup)
	r.POST("/userlogin", userHandler.Userlogin)

	r.POST("/send-otp", otpHandler.SendOtp)
	r.POST("/verify-otp", otpHandler.VerifyOtp)

	r.POST("/forgot-password", userHandler.ForgotPasswordSend)
	r.POST("/forgot-password-verify", userHandler.ForgotPasswordVerifyAndChange)

	r.GET("/razorpay", paymentHandler.MakePaymentRazorPay)
	r.GET("/update_status", paymentHandler.VerifyPayment)

	products := r.Group("/products")
	{
		products.GET("", productHandler.ShowAllProducts)
		products.POST("/filter", productHandler.FilterCategory)
		products.GET("/image", productHandler.ShowImages) //Individual Images

	}
	r.Use(middleware.UserAuthMiddleware())
	{
		address := r.Group("/address")
		{
			address.GET("", userHandler.GetAllAddress)
			address.POST("", userHandler.AddAddress)
			address.PUT("", userHandler.UpdateAddress)
			address.DELETE("", userHandler.DeleteAddressByID)
		}
		users := r.Group("/users")
		{
			users.GET("", userHandler.UserDetails)
			users.PUT("", userHandler.UpdateUserDetails)
			users.PUT("/changepassword", userHandler.ChangePassword)
		}

		wishlist := r.Group("/wishlist")
		{
			wishlist.POST("", wishlistHandler.AddToWishlist)
			wishlist.GET("", wishlistHandler.GetWishList)
			wishlist.DELETE("", wishlistHandler.RemoveFromWishlist)
		}

		cart := r.Group("/cart")
		{
			cart.POST("", cartHandler.AddToCart)
			cart.DELETE("", cartHandler.RemoveFromCart)
			cart.GET("", cartHandler.DisplayCart)
			cart.DELETE("/empty", cartHandler.EmptyCart)
			cart.PUT("/updatequantityadd", cartHandler.UpdateQuantityAdd)
			cart.PUT("/updatequantityless", cartHandler.UpdateQuantityless)

		}

		order := r.Group("/order")
		{
			order.POST("", orderHandler.OrderItemsFromCart)
			order.GET("", orderHandler.GetOrderDetails)
			order.GET("/checkout", orderHandler.CheckOut)
			order.GET("/place-order", orderHandler.PlaceOrderCOD)
			order.PUT("", orderHandler.CancelOrder)
		}
		r.POST("/coupon/apply", couponHandler.ApplyCoupon)
		r.GET("/referral/apply", userHandler.ApplyReferral)
		wallet := r.Group("/wallet")
		{
			wallet.GET("", walletHandler.GetWallet)
			wallet.GET("/history", walletHandler.WalletHistory)
		}
	}

}
