package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sowjanyac93/bestWallet/modules/initializers"
	"github.com/sowjanyac93/bestWallet/modules/routes"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectDatabase()
	initializers.SyncDatabase()
}

func main() {
	router := gin.New()
	routes.Routes(router)
	router.Run()
}
