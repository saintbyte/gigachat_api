package main

import (
	gigachat "github.com/saintbyte/gigachat_api"
	"log"
)

func main() {
	chat := gigachat.NewGigachat()
	models, err := chat.GetModels()
	if err != nil {
		log.Fatal(err)
	}
	for _, model := range models {
		log.Println(model.ID)
		log.Println(model.OwnedBy)
		log.Println(model.Object)
		log.Println("-------------------------")
	}
}
