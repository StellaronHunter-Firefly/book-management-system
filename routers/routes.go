package routers

import (
	"book-management-system/controllers"
	"book-management-system/middlewares"
	"book-management-system/repositories"
	"book-management-system/services"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(middlewares.CORSMiddleware())

	// 初始化仓库和服务
	userRepo := repositories.NewUserRepository()
	bookRepo := repositories.NewCombinedBookRepository()

	authService := services.NewAuthService(userRepo)
	bookService := services.NewBookService(bookRepo)

	authController := controllers.NewAuthController(authService)
	bookController := controllers.NewBookController(bookService)

	// 公共路由
	api := router.Group("/api")
	{
		// 认证路由
		auth := api.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
		}

		// 公开的图书查询
		books := api.Group("/books")
		{
			books.GET("", bookController.GetAllBooks)
			books.GET("/search", bookController.SearchBooks)
			books.GET("/:id", bookController.GetBookByID)
			books.GET("/:id/availability", bookController.CheckAvailability)
		}
	}

	// 需要认证的路由
	authenticated := api.Group("")
	authenticated.Use(middlewares.AuthMiddleware())
	{
		// 用户相关
		user := authenticated.Group("/users")
		{
			user.GET("/profile", authController.GetProfile)
		}

		// 书籍借还（管理员和普通用户都可以）
		books := authenticated.Group("/books")
		{
			books.POST("/borrow", bookController.BorrowBook)
			books.POST("/return", bookController.ReturnBook)
			books.GET("/my-borrowed", bookController.GetMyBorrowedBooks)
			books.GET("/my-records", bookController.GetMyBorrowRecords)
		}

		// 管理员专用路由
		admin := authenticated.Group("/admin")
		admin.Use(middlewares.AdminOnly())
		{
			// 图书管理
			admin.POST("/books", bookController.CreateBook)
			admin.PUT("/books/:id", bookController.UpdateBook)
			admin.DELETE("/books/:id", bookController.DeleteBook)

			// 借阅记录管理
			admin.GET("/borrow-records", bookController.GetAllBorrowRecords)
		}
	}

	return router
}
