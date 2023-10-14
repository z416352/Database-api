package ginservices

import (
	"Database/internal/controllers"

	"github.com/gin-gonic/gin"
)

func Gin_handler(port string) {
	router := gin.Default()
	router.Use(restrictToLocalNetwork())

	router.POST("/prices", controllers.InsertData())
	router.GET("/prices/:symbol/:timeframe", controllers.GetData())
	router.GET("/mongodb/:dbname", controllers.GetDBExist())

	router.Run(":" + port)

}
