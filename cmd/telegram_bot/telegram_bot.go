package main

import (
	"os"

	"github.com/yazmeyaa/telegram_sticker_converter/internal/logger"
)

func main() {
	logger := logger.NewLogger(os.Stdout, "telegram_bot")
	logger.Info("Starting telegram bot")
}
