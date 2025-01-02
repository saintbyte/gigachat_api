package main

import (
	gigachat "github.com/saintbyte/gigachat_api"
	"log/slog"
	"os"
)

func main() {
	chat := gigachat.NewGigachat()

	uploadedFile, err := chat.UploadFile("LICENSE.txt")
	if err != nil {
		slog.Error("File upload error:", err)
		os.Exit(0)
	}

	messages := []gigachat.MessageRequest{}
	messages = append(messages, gigachat.MessageRequest{
		Role:    gigachat.GigaChatRoleSystem,
		Content: "Ты самый ИИ на планете и обязательно попробуешь найти функцию в входящем сообщении.",
	})
	messages = append(messages, gigachat.MessageRequest{
		Role:    gigachat.GigaChatRoleUser,
		Content: "Расскажи вкратце, что изложено в документе?", // Не используйте слова файл или аттач - не понимает
		Attachments: []string{
			uploadedFile.Id,
		},
	})

	result, err := chat.ChatCompletions(messages)
	if err != nil {
		slog.Error("Ask error:", err)
	}
	slog.Info("Result:", result)
	if result.Choices[0].FinishReason == gigachat.GigaChatFinishReasonFunctionCall {
		slog.Info("Function call")
		slog.Info("result.Choices[0]:", result.Choices[0])
	}
	answer := result.Choices[0].Message.Content
	slog.Info(answer)
}
