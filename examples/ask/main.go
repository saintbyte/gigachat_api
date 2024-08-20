package main

import (
	gigachat "github.com/saintbyte/gigachat_api"
	"log"
)

func main() {
	chat := gigachat.NewGigachat()
	aswer, err := chat.Ask("Сколько рыбы в море?")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(aswer)
}
