package main

import (
	"context"
	"fmt"

	"github.com/abdheshnayak/syncmeet/syncmeet/framework"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(framework.Module, fx.NopLogger)

	if err := app.Start(context.Background()); err != nil {
		panic(err)
	}

	fmt.Println("App started")
	<-app.Done()
}
