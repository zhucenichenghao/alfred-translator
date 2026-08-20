package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/zhucenichenghao/alfred-translator/app"
)

func main() {
	client := &http.Client{Timeout: 8 * time.Second}
	code := app.Run(context.Background(), os.Args, os.Getenv, os.Stdout, client)
	if code != 0 {
		os.Exit(code)
	}
}
