package main

import (
	"context"

	"gitlab.com/subrotokumar/glitchr/transcoder/service"
)

func main() {
	worker := service.New()
	worker.Run(context.Background())
}
