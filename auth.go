package gigachat

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	uuid "github.com/nu7hatch/gouuid"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

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
