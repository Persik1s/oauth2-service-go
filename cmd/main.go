package main

import (
	"github.com/Persik1s/oauth2-service-go/internal"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		internal.ModuleApp,
	)

	app.Run()
}
