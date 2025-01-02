package gigachat

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"io"
	"log"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
)

/*
POST /files;
GET /files;
GET /files/{file};
POST /files/{file}/delete.
curl --location --request POST 'https://gigachat.devices.sberbank.ru/api/v1/files' \
--header 'Authorization: Bearer access_token' \
--form 'file=@"<путь_к_файлу>/example.jpeg"' \
--form 'purpose="general"'
*/

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {

	return quoteEscaper.Replace(s)
}

// AllowedUploadMimeTypes Доступные для загрузки Mime типы
func (g Gigachat) AllowedUploadMimeTypes() []string {
	return []string{
		"image/jpg",
		"image/jpeg",
		"image/png",
		"image/tiff",
		"image/bmp",
		"text/plain",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"application/msword",
		"application/pdf",
	}
}

// AllowedUploadFileExt Доступные для загрузки расширения файлов
func (g Gigachat) AllowedUploadFileExt() []string {
	return []string{
		".jpg",
		".jpeg",
		".png",
		".tiff",
		".bmp",
		".txt",
		".doc",
		".docx",
		".pdf",
	}
}

func (g *Gigachat) checkFile(filename string) (bool, error) {

	file, err := os.Open(filename)
	if err != nil {
		slog.Error("File open error: ", filename, err)
		return false, err
	}
	defer file.Close()
	ext := filepath.Ext(filename)
	validExt := false
	for _, b := range g.AllowedUploadFileExt() {
		if b == ext {
			validExt = true
			break
		}
	}
	if !validExt {
		return false, errors.New("invalid file extension")
	}
	return g.checkFileIoReader(file)
}

// Просто получить тип из ридера который из файла
func (g *Gigachat) getMimeType(fileReader io.Reader) (string, error) {
	mimetype.SetLimit(1024 * 1024) // 1MB
	mtype, err := mimetype.DetectReader(fileReader)
	if err != nil {
		slog.Error("DetectReader error: ", err)
		return "", err
	}
	simpleMimeType := mtype.String()

	if strings.Contains(simpleMimeType, ";") {
		simpleMimeType = strings.Split(simpleMimeType, ";")[0]
	}
	return simpleMimeType, nil
}
func (g *Gigachat) checkFileIoReader(fileReader io.Reader) (bool, error) {
	validMimeType := false
	simpleMimeType, err := g.getMimeType(fileReader)
	if err != nil {
		return false, err
	}
	for _, b := range g.AllowedUploadMimeTypes() {
		if b == simpleMimeType {
			validMimeType = true
			break
		}
	}
	return validMimeType, nil
}

func (g *Gigachat) postUploadRequest(url string, body io.Reader, writer *multipart.Writer) (*http.Request, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	request, err := http.NewRequest("POST", url, body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "Bearer "+g.getCurrentToken())
	if err != nil {
		return nil, err
	}
	return request, nil
}
func (g *Gigachat) getUploadData(file *os.File, fileName string) (*bytes.Buffer, *multipart.Writer, error) {
	var fileReader io.Reader
	fileReader = file
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, escapeQuotes(fileName)))

	mimeType, _ := g.getMimeType(fileReader)
	h.Set("Content-Type", mimeType)
	part, err := writer.CreatePart(h)
	file.Seek(0, 0)
	fileReader = file
	_, err = io.Copy(part, fileReader)
	if err != nil {
		slog.Error("Error in copy to form part:", err)
		return nil, nil, err
	}

	_ = writer.WriteField("purpose", "general")
	err = writer.Close()
	if err != nil {
		slog.Error("Writer close error:", err)
		return nil, nil, err
	}
	return body, writer, err
}
func (g *Gigachat) UploadData(fileReader *os.File, fileName string) (file UploadedFile, err error) {
	url := g.getRequestUrl(GigaChatChatFileUploadPath)
	data, writer, err := g.getUploadData(fileReader, fileName)
	request, err := g.postUploadRequest(url, data, writer)
	if err != nil {
		return UploadedFile{}, err
	}
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return UploadedFile{}, err
	}
	if response.StatusCode != http.StatusOK {
		return UploadedFile{}, errors.New("Response status: " + string(response.Status))
	}
	body, err := io.ReadAll(response.Body)
	defer response.Body.Close()
	var result UploadedFile
	err2 := json.Unmarshal(body, &result)
	if err2 != nil {
		log.Fatal(err2)
	}
	return result, nil
}

func (g *Gigachat) UploadFile(filePath string) (UploadedFile, error) {
	//Грузим файл в облако для сбера чтоб потом использовать его ID
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Ошибка открытия файла:", err)
		return UploadedFile{}, err
	}
	defer file.Close()
	checkResult, err := g.checkFile(filePath)
	if err != nil {
		slog.Error("File not valid:", err)
		return UploadedFile{}, err
	}
	if !checkResult {
		slog.Error("File not valid:", err)
		return UploadedFile{}, err
	}
	filename := filepath.Base(filePath)
	return g.UploadData(file, filename)
}

/*
func (g *Gigachat) ListFiles(fileReader io.Reader, fileName string) (file *File, err error) {
	url := g.getRequestUrl(GigaChatChatFileUploadPath)
}
func (g *Gigachat) GetFile(fileReader io.Reader, fileName string) (file *File, err error) {
	url := g.getRequestUrl(GigaChatChatFileUploadPath)
}

func (g *Gigachat) DeleteFiles(fileReader io.Reader, fileName string) (file *File, err error) {
	url := g.getRequestUrl(GigaChatChatFileUploadPath)
}
*/
