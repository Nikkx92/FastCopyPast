package main

import (
	a "FastCopyPast/internal/app"
	"FastCopyPast/internal/ui"
	"log"
	"os"

	"gioui.org/app"
)

var version = "1.0.0"

func main() {
	windowTitle := "FastCopyPast"

	go func() {
		a.SetOpacity(windowTitle)
	}()

	go func() {
		application := a.NewApp(windowTitle, version)

		u := ui.NewUI(application)

		if err := application.Run(u); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()

	app.Main()
}
