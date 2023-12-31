package main

import (
	"os"

	"github.com/Finschia/finschia-sdk/server"
	svrcmd "github.com/Finschia/finschia-sdk/server/cmd"

	"github.com/Finschia/finschia-proxy/v2/app"
	"github.com/Finschia/finschia-proxy/v2/cmd/fnsad-proxy/cmd"
)

func main() {
	rootCmd, _ := cmd.NewRootCmd()

	if err := svrcmd.Execute(rootCmd, app.DefaultNodeHome); err != nil {
		switch e := err.(type) {
		case server.ErrorCode:
			os.Exit(e.Code)

		default:
			os.Exit(1)
		}
	}
}
