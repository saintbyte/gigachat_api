package main

import (
	gigachat "github.com/saintbyte/gigachat_api"
	"log/slog"
)

func main() {
	chat := gigachat.NewGigachat()
	messages := []gigachat.MessageRequest{}
	messages = append(messages, gigachat.MessageRequest{
		Role:    gigachat.GigaChatRoleSystem,
		Content: "Ты самый ИИ на планете и обязательно попробуешь найти функцию в входящем сообщении.",
	})
	messages = append(messages, gigachat.MessageRequest{
		Role:    gigachat.GigaChatRoleUser,
		Content: "Напиши код для python",
	})
	chat.Functions = []gigachat.Function{
		gigachat.Function{
			Name:        "think",
			Description: "Подумать о вечно",
			Parameters: gigachat.Parameters{
				Type: "object",
				Properties: map[string]gigachat.Property{
					"now_long": gigachat.Property{
						Type:        "string",
						Description: "Сколько надо думать?",
					},
				},
			},
			FewShotExamples: []gigachat.Example{
				gigachat.Example{
					Request: "Подумай в всем на свете минут 15",
					Params: map[string]string{
						"now_long": "15m",
					},
				},
				gigachat.Example{
					Request: "Подумай о каждом из нас полчаса",
					Params: map[string]string{
						"now_long": "30m",
					},
				},
			},
		},
	}
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
