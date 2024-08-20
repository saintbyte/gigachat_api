package gigachat

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nu7hatch/gouuid"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Gigachat struct {
}

func NewGigachat() *Gigachat {
	return &Gigachat{}
}

func (g *Gigachat) GetExpiresAtFromFile() int64 {
	data, err := os.ReadFile(g.GetExpiresFile())
	if err != nil {
		return 0
	}
	i, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return 0
	}
	return i
}
func (g *Gigachat) GetTokenFromFile() string {
	data, err := os.ReadFile(g.GetTokenFile())
	if err != nil {
		return ""
	}
	return string(data)
}
func (g Gigachat) GetExpiresFile() string {
	filename, exists := os.LookupEnv(GigaChatExpiresFileEnv)
	if !exists {
		return ".gigachat_expires"
	}
	return filename
}

func (g Gigachat) GetTokenFile() string {
	filename, exists := os.LookupEnv(GigaChatTokenFileEnv)
	if !exists {
		return ".gigachat_token"
	}
	return filename
}

func (g *Gigachat) SetExpiresAtToFile(value int64) {
	fh, _ := os.OpenFile(g.GetExpiresFile(), os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0666)
	fh.WriteString(strconv.FormatInt(value, 10))
	defer fh.Close()
}
func (g *Gigachat) SetTokenToFile(value string) {
	fh, _ := os.OpenFile(g.GetTokenFile(), os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0666)
	fh.WriteString(value)
	defer fh.Close()
}
func (g *Gigachat) GetCurrentToken() string {
	expAt := g.GetExpiresAtFromFile()
	token := g.GetTokenFromFile()
	apochNow := time.Now().Unix()
	timeDelta := apochNow - (expAt / 1000)
	if timeDelta > 0 {
		newExpAt, token2 := g.Auth()
		g.SetExpiresAtToFile(newExpAt)
		g.SetTokenToFile(token2)
		token = token2
	}
	return token
}
func (g Gigachat) GetRequestUrl(path string) string {
	return "https://" + GigaChatApiHost + path
}

func (g *Gigachat) GetRequest(url string) (*http.Request, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "Bearer "+g.GetCurrentToken())
	return request, nil
}
func (g *Gigachat) PostRequest(url string, body io.Reader) (*http.Request, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	request, err := http.NewRequest("POST", url, body)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "Bearer "+g.GetCurrentToken())
	if err != nil {
		return nil, err
	}
	return request, nil
}
func (g *Gigachat) Auth() (int64, string) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	u, err := uuid.NewV4()
	request, _ := http.NewRequest("POST", GigaChatOauthUrl, bytes.NewBufferString("scope=GIGACHAT_API_PERS"))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("RqUID", u.String())
	request.Header.Set("Authorization", "Basic "+os.Getenv(GigaChatAuthData))
	client := &http.Client{}
	response, e := client.Do(request)

	if e != nil {
		log.Fatal(e)
	}
	//if response.StatusCode != http.StatusOK {
	//	return "Так что-то пошло не так на удаленной стороне. Повтори вопрос.", nil
	//}
	fmt.Println(response.StatusCode)
	if response.StatusCode != http.StatusOK {
		return 0, ""
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(body))
	defer response.Body.Close()

	var result TokenResponse
	err2 := json.Unmarshal(body, &result)
	if err2 != nil {
		log.Fatal(err2)
	}
	os.Setenv(GigaChatToken, result.AccessToken)
	return result.ExpiresAt, result.AccessToken
}

func (g *Gigachat) GetModels() ([]ModelItem, error) {
	url := g.GetRequestUrl(GigaChatModelsPath)
	request, err := g.GetRequest(url)
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

func (g *Gigachat) Embeddings(input string) ([]float32, error) {
	url := g.GetRequestUrl(GigaChatEmbeddingsPath)
	var inputs []string
	inputs = append(inputs, input)
	jData, errJsonRequestEncode := json.Marshal(&EmbeddingsRequest{
		Model: "Embeddings",
		Input: inputs,
	})
	if errJsonRequestEncode != nil {
		return nil, errJsonRequestEncode
	}
	request, err := g.PostRequest(url, bytes.NewReader(jData))
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

func (g *Gigachat) ChatCompletions(input []string) (string, error) {
	url := g.GetRequestUrl(GigaChatChatCompletionPath)
	temperature := float32(1.0)
	var messages []MessageRequest
	for _, inputItem := range input {
		messages = append(messages, MessageRequest{
			Role:    "user",
			Content: inputItem,
		})
	}
	jData, errJsonRequestEncode := json.Marshal(&ChatCompletionRequest{
		Model:             GigaChatModel,
		MaxTokens:         GigaChatMaxTokens,
		Temperature:       temperature,
		Messages:          messages,
		Stream:            false,
		RepetitionPenalty: 1,
		TopP:              1,
		UpdateInterval:    0,
	})
	if errJsonRequestEncode != nil {
		return "", errJsonRequestEncode
	}
	request, err := g.PostRequest(url, bytes.NewReader(jData))
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
