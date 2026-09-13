package main

import (
	"fmt"
	"setupwizard/cmd/command"
	"setupwizard/internal/adapters/left/console"
	"setupwizard/internal/adapters/right/json"
	"setupwizard/internal/core/service"
)

func main() {
	jRep := json.NewJsonRepository()
	appService := service.NewAppService(jRep)
	handler := console.NewHandler(appService)

	apps, err := handler.GetApps()
	if err != nil {
		fmt.Println(err)
		return
	}

	for i, app := range apps {
		fmt.Printf("%v:\n", i+1)
		fmt.Println(app.Name)
		fmt.Println(app.Link)
		fmt.Println("")
	}

	command.Execute()
}
