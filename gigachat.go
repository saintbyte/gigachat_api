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
	"github.com/nu7hatch/gouuid"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
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

func (g *Gigachat) getExpiresAtFromFile() int64 {
	data, err := os.ReadFile(g.getExpiresFile())
	if err != nil {
		return 0
	}
	i, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return 0
	}
	return i
}
func (g *Gigachat) getTokenFromFile() string {
	data, err := os.ReadFile(g.getTokenFile())
	if err != nil {
		return ""
	}
	return string(data)
}
func (g Gigachat) getExpiresFile() string {
	filename, exists := os.LookupEnv(GigaChatExpiresFileEnv)
	if !exists {
		return ".gigachat_expires"
	}
	return filename
}

func (g Gigachat) getTokenFile() string {
	filename, exists := os.LookupEnv(GigaChatTokenFileEnv)
	if !exists {
		return ".gigachat_token"
	}
	return filename
}

func (g *Gigachat) setExpiresAtToFile(value int64) {
	fh, _ := os.OpenFile(g.getExpiresFile(), os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0666)
	fh.WriteString(strconv.FormatInt(value, 10))
	defer fh.Close()
}
func (g *Gigachat) setTokenToFile(value string) {
	fh, _ := os.OpenFile(g.getTokenFile(), os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0666)
	fh.WriteString(value)
	defer fh.Close()
}
func (g *Gigachat) getCurrentToken() string {
	expAt := g.getExpiresAtFromFile()
	token := g.getTokenFromFile()
	apochNow := time.Now().Unix()
	timeDelta := apochNow - (expAt / 1000)
	if timeDelta > 0 {
		newExpAt, token2 := g.Auth()
		g.setExpiresAtToFile(newExpAt)
		g.setTokenToFile(token2)
		token = token2
	}
	return token
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

func (g *Gigachat) getAuthData() string {
	value, exists := os.LookupEnv(GigaChatAuthData)
	if exists {
		return value
	}
	if g.AuthData != "" {
		return g.AuthData
	}
	return ""
}

// Auth Авторизация для получения токена для запросов.
func (g *Gigachat) Auth() (int64, string) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	u, err := uuid.NewV4()
	request, _ := http.NewRequest("POST", GigaChatOauthUrl, bytes.NewBufferString("scope=GIGACHAT_API_PERS"))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("RqUID", u.String())
	request.Header.Set("Authorization", "Basic "+g.getAuthData())
	client := &http.Client{}
	response, e := client.Do(request)

	if e != nil {
		log.Fatal(e)
	}
	if response.StatusCode != http.StatusOK {
		return 0, ""
	}
	log.Println(response.StatusCode)
	if response.StatusCode != http.StatusOK {
		return 0, ""
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(body))
	defer response.Body.Close()

	var result TokenResponse
	err2 := json.Unmarshal(body, &result)
	if err2 != nil {
		log.Fatal(err2)
	}
	os.Setenv(GigaChatToken, result.AccessToken)
	return result.ExpiresAt, result.AccessToken
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
		Model:             GigaChatModel,
		MaxTokens:         GigaChatMaxTokens,
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
		log.Fatal(err)
	}
	return result.Choices[0].Message.Content, nil
}

// Ask Просто спросить у модели
func (g *Gigachat) Ask(input string) (string, error) {
	return g.ChatCompletions([]MessageRequest{
		{
			Role:    "user",
			Content: input,
		},
	})
}
