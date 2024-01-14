package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sowjanyac93/bestWallet/modules/controllers"
	"github.com/sowjanyac93/bestWallet/modules/middleware"
)

func Routes(router *gin.Engine) {
	// Public routes (no authentication required)
	router.POST("/login", controllers.Login)
	router.POST("/sign-up", controllers.Signup)
	router.POST("/validate", middleware.RequiredAuth, controllers.Validate)
	router.POST("/logout", controllers.Logout)

	authGroup := router.Group("/")
	authGroup.Use(middleware.RequiredAuth)
	{
		authGroup.POST("/am/add-account/:id", controllers.AddAccountForUserId)
		authGroup.POST("/am/add-self-account", controllers.AddSelfAccount)
		authGroup.DELETE("/am/delete-account/:id", controllers.DeleteAccount)
		authGroup.GET("/am/get-accounts", controllers.GetAccounts)
		authGroup.GET("/am/get-personal-account", controllers.GetPersonalAccount)

		authGroup.POST("/um/create-user", controllers.CreateUser)
		authGroup.DELETE("/um/delete-user/:id", controllers.DeleteUser)
		authGroup.GET("/um/get-users", controllers.GetUsers)

		authGroup.POST("/tm/deposit", controllers.DepositAmount)
		authGroup.POST("/tm/withdraw", controllers.WithdrawAmount)
		authGroup.POST("/tm/self-transfer", controllers.SelfTransfer)
		authGroup.GET("/tm/get-transactions", controllers.GetTransactions)
	}
}
