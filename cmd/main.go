package main

import (
	"MKICS/cmd/cmd"
	"fmt"
	"os"
)

// @title MKICS
// @version 1.0
// @description  mkics
// @termsOfService http://swagger.io/terms/
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost
// @BasePath /api/v1

//go:generate swag init -o ./docs -g main.go -d ../../backend -g ../cmd/server/main.go
func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
