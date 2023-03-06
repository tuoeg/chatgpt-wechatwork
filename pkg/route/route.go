package route

import (
	"server/pkg/controller"

	"github.com/gin-gonic/gin"
)

func InstallRoutes(r *gin.Engine) {
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())

	rootGroup := r.Group("/api/v1")
	// AuthRequired middleware that provide basic auth
	//rootGroup.Use(middleware.BasicAuthMiddleware())

	{
		// toDoController := controller.NewToDoController()
		// rootGroup.GET("/todo/get", toDoController.GetToDo)
		con := controller.NewController()
		rootGroup.GET("/message", con.Auth)
		rootGroup.POST("/message", con.Message)
	}
}
