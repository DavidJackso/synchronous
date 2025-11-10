package main

import (
	"fmt"

	"github.com/rnegic/synchronous/internal/app"
)

func main() {

	app := app.New()

	err := app.Run()
	if err != nil {
		fmt.Println(err)
	}
}
