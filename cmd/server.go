package cmd

import (
	"fmt"
	"os"
	"server/pkg/route"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "server",
	Short: "server is a wecomchat server that interactive whit chatgpt",
	Long: `when user send message to wecomchat,the server transmit message to chatgpt,
			when chatgpt return result,callback wecomchat api`,
	Run: func(cmd *cobra.Command, args []string) {
		startServer()
	},
}

func startServer() {
	gin.SetMode(gin.Mode())
	r := gin.Default()
	route.InstallRoutes(r)
	serverBindAddr := fmt.Sprintf("%s:%d", "0.0.0.0", 8081)
	r.Run(serverBindAddr) // listen and serve
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
