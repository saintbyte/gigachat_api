// Package gigachat Предоставляет доступ Gigachat
//
// Этот пакет сделан для того чтоб спрашивать у нейросети gigachat от сбера.
// Так и делать embedding
package gigachat

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"log"
	"log/slog"
	"net/http"
)

type Gigachat struct {
	ApiHost           string
	RepetitionPenalty int
	TopP              float32
	Model             string
	MaxTokens         int
	Temperature       float32
	AuthData          string
}

func NewGigachat() *Gigachat {
	return &Gigachat{
		ApiHost:           GigaChatApiHost,
		RepetitionPenalty: 1,
		TopP:              1.0,
		Model:             GigaChatModel,
		MaxTokens:         GigaChatMaxTokens,
		Temperature:       1,
		AuthData:          "",
	}
}

func (g *Gigachat) getRequestUrl(path string) string {
	return "https://" + g.ApiHost + path
}

func (g *Gigachat) getRequest(url string) (*http.Request, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "Bearer "+g.getCurrentToken())
	return request, nil
}

func (g *Gigachat) postRequest(url string, body io.Reader) (*http.Request, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	request, err := http.NewRequest("POST", url, body)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "Bearer "+g.getCurrentToken())
	if err != nil {
		return nil, err
	}
	return request, nil
}

// GetModels Получить список моделей.
func (g *Gigachat) GetModels() ([]ModelItem, error) {
	url := g.getRequestUrl(GigaChatModelsPath)
	request, err := g.getRequest(url)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, errors.New("Response status: " + string(response.Status))
	}
	body, err := io.ReadAll(response.Body)
	defer response.Body.Close()
	var result ModelsResponse
	err2 := json.Unmarshal(body, &result)
	if err2 != nil {
		log.Fatal(err2)
	}
	return result.Data, nil
}

// Embeddings получить вектора текста. Ограничение по количеству что-то типа 512.
func (g *Gigachat) Embeddings(input string) ([]float32, error) {
	url := g.getRequestUrl(GigaChatEmbeddingsPath)
	var inputs []string
	inputs = append(inputs, input)
	jData, errJsonRequestEncode := json.Marshal(&EmbeddingsRequest{
		Model: "Embeddings",
		Input: inputs,
	})
	if errJsonRequestEncode != nil {
		return nil, errJsonRequestEncode
	}
	request, err := g.postRequest(url, bytes.NewReader(jData))
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		defer response.Body.Close()
		return nil, errors.New("Response status: " + string(response.Status) + " " + string(body))
	}
	body, err := io.ReadAll(response.Body)
	defer response.Body.Close()
	var result EmbeddingsResponse
	err = json.Unmarshal(body, &result)
	return result.Data[0].Embedding, nil
}

// ChatCompletions Сдалать запрос к модели.
func (g *Gigachat) ChatCompletions(messages []MessageRequest) (string, error) {
	url := g.getRequestUrl(GigaChatChatCompletionPath)
	jData, errJsonRequestEncode := json.Marshal(&ChatCompletionRequest{
		Model:             g.Model,
		MaxTokens:         g.MaxTokens,
		Temperature:       g.Temperature,
		Messages:          messages,
		Stream:            false,
		RepetitionPenalty: g.RepetitionPenalty,
		TopP:              g.TopP,
		UpdateInterval:    0,
	})
	if errJsonRequestEncode != nil {
		return "", errJsonRequestEncode
	}
	request, err := g.postRequest(url, bytes.NewReader(jData))
	if err != nil {
		return "", err
	}
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	if response.StatusCode != http.StatusOK {

		return "", errors.New("Response status: " + string(response.Status))
	}
	body, err := io.ReadAll(response.Body)
	defer response.Body.Close()
	var result ChatCompletionResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		slog.Error("Json Unmarshal error:", err)
	}
	return result.Choices[0].Message.Content, nil
}

// Ask Просто спросить у модели
func (g *Gigachat) Ask(input string) (string, error) {
	return g.ChatCompletions([]MessageRequest{
		{
			Role:    GigaChatRoleUser,
			Content: input,
		},
	})
}
