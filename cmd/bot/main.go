package main

import (
	"context"
	"log"

	"github.com/V1merX/artifacts-bot/internal/app"
)

func main() {
	ctx := context.Background()

	app := app.New()

	if err := app.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
