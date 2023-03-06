package main

import (
	"fmt"
	"server/pkg/route"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.Mode())
	r := gin.Default()
	route.InstallRoutes(r)
	serverBindAddr := fmt.Sprintf("%s:%d", "0.0.0.0", 8081)
	r.Run(serverBindAddr) // listen and serve
}
