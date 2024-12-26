package main

import (
	gigachat "github.com/saintbyte/gigachat_api"
	"log/slog"
)

func main() {
	chat := gigachat.NewGigachat()
	answer, err := chat.Ask("Сколько рыбы в море?")
	if err != nil {
		slog.Error("Ask error:", err)
	}
	slog.Info(answer)
}
