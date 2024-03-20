package cmd

import (
	"apiserver/handler"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

func runApiServer(cmd *cobra.Command, args []string) {
	// cmd run
	route := gin.New()

	// http handler implement
	route.Use(gin.Logger())

	v1 := route.Group("/v1")
	{
		v1.GET("/helloworld", handler.HelloWorld())
	}
	if err := route.Run(":8080"); err != nil {
		panic("run api server error")
	}
}

var apiServerCmd = &cobra.Command{
	Use:   "apiserver",
	Short: "apiserver short description",
	Long:  "apiserver long description",
	Run:   runApiServer,
}

func init() {
	rootCmd.AddCommand(apiServerCmd)
}
