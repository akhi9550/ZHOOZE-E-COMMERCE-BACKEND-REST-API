package routes

import (
	"Zhooze/pkg/api/handlers"
	"Zhooze/pkg/api/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AdminRoutes(r *gin.RouterGroup, db *gorm.DB) *gin.RouterGroup {

	r.POST("/adminlogin", handlers.LoginHandler)

	r.Use(middleware.AdminAuthMiddleware())
	{

		r.GET("/dashboard", handlers.DashBoard)
		r.GET("/sales-report", handlers.FilteredSalesReport)
		r.GET("/sales-report-date", handlers.SalesReportByDate)

		users := r.Group("/users")
		{
			users.GET("", handlers.GetUsers)
			users.PUT("/block", handlers.BlockUser)
			users.PUT("/unblock", handlers.UnBlockUser)
		}

		products := r.Group("/products")
		{
			products.GET("", handlers.ShowAllProductsFromAdmin)
			products.POST("", handlers.AddProducts)
			products.PUT("", handlers.UpdateProduct) //update the product quantity
			products.DELETE("", handlers.DeleteProducts)
			products.POST("/upload-image", handlers.UploadImage)
		}

		category := r.Group("/category")
		{
			category.GET("", handlers.GetCategory)
			category.POST("", handlers.AddCategory)
			category.PUT("", handlers.UpdateCategory)
			category.DELETE("", handlers.DeleteCategory)

		}

		order := r.Group("/order")
		{
			order.GET("", handlers.GetAllOrderDetailsForAdmin)
			order.GET("/approve", handlers.ApproveOrder)
			order.GET("/cancel", handlers.CancelOrderFromAdmin)
		}

		payment := r.Group("/payment-method")
		{
			payment.POST("", handlers.AddPaymentMethod)
			payment.GET("", handlers.ListPaymentMethods)
			payment.DELETE("", handlers.DeletePaymentMethod)
		}
		coupons := r.Group("/coupons")
		{
			coupons.POST("", handlers.AddCoupon)
			coupons.GET("", handlers.GetCoupon)
			coupons.PATCH("", handlers.ExpireCoupon)
		}
		offer := r.Group("/offer")
		{
			offer.POST("/product-offer", handlers.AddProdcutOffer)
			offer.GET("/product-offer",handlers.GetProductOffer)
			offer.DELETE("/product-offer",handlers.ExpireProductOffer)

			offer.POST("/category-offer", handlers.AddCategoryOffer)
			offer.GET("/category-offer",handlers.GetCategoryOffer)
			offer.DELETE("/category-offer",handlers.ExpireCategoryOffer)
		}

	}
	return r
}
