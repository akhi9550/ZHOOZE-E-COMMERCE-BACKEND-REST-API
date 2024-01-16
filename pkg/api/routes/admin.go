package routes

import (
	"Zhooze/pkg/api/handlers"
	"Zhooze/pkg/api/middleware"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(r *gin.RouterGroup, adminHandler *handlers.AdminHandler, productHandler *handlers.ProductHandler, orderHandler *handlers.OrderHandler, couponHandler *handlers.CouponHandler, categoryHandler *handlers.CategoryHandler, offerHandler *handlers.OfferHandler) {

	r.POST("/adminlogin", adminHandler.LoginHandler)

	r.Use(middleware.AdminAuthMiddleware())
	{

		r.GET("/dashboard", adminHandler.DashBoard)
		r.GET("/sales-report", adminHandler.FilteredSalesReport)
		r.GET("/sales-report-date", adminHandler.SalesReportByDate)

		users := r.Group("/users")
		{
			users.GET("", adminHandler.GetUsers)
			users.PUT("/block", adminHandler.BlockUser)
			users.PUT("/unblock", adminHandler.UnBlockUser)
		}

		products := r.Group("/products")
		{
			products.GET("", productHandler.ShowAllProductsFromAdmin)
			products.POST("", productHandler.AddProducts)
			products.PUT("", productHandler.UpdateProduct)
			products.DELETE("", productHandler.DeleteProducts)
			products.GET("/search", productHandler.SearchProducts)
			products.POST("/upload-image", productHandler.UploadImage)
		}

		category := r.Group("/category")
		{
			category.GET("", categoryHandler.GetCategory)
			category.POST("", categoryHandler.AddCategory)
			category.PUT("", categoryHandler.UpdateCategory)
			category.DELETE("", categoryHandler.DeleteCategory)

		}

		order := r.Group("/order")
		{
			order.GET("", orderHandler.GetAllOrderDetailsForAdmin)
			order.GET("/approve", orderHandler.ApproveOrder)
			order.GET("/cancel", orderHandler.CancelOrderFromAdmin)
		}

		payment := r.Group("/payment-method")
		{
			payment.POST("", adminHandler.AddPaymentMethod)
			payment.GET("", adminHandler.ListPaymentMethods)
			payment.DELETE("", adminHandler.DeletePaymentMethod)
		}
		coupons := r.Group("/coupons")
		{
			coupons.POST("", couponHandler.AddCoupon)
			coupons.GET("", couponHandler.GetCoupon)
			coupons.PATCH("", couponHandler.ExpireCoupon)
		}
		offer := r.Group("/offer")
		{
			offer.POST("/product-offer", offerHandler.AddProdcutOffer)
			offer.GET("/product-offer", offerHandler.GetProductOffer)
			offer.DELETE("/product-offer", offerHandler.ExpireProductOffer)

			offer.POST("/category-offer", offerHandler.AddCategoryOffer)
			offer.GET("/category-offer", offerHandler.GetCategoryOffer)
			offer.DELETE("/category-offer", offerHandler.ExpireCategoryOffer)
		}

	}
}
