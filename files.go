package gigachat

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Velocidex/go-magic/magic"
	"github.com/Velocidex/go-magic/magic_files"
	"io"
	"log"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
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

// Доступные для загрузки Mime типы
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

// Доступные для загрузки расширения файлов
func (g Gigachat) AllowedUploadFileExt() []string {
	return []string{
		"jpg",
		"jpeg",
		"png",
		"tiff",
		"bmp",
		"txt",
		"doc",
		"docx",
		"pdf",
	}
}

func (g *Gigachat) checkFileMimeType() string {
	handle := magic.NewMagicHandle(magic.MAGIC_NONE)
	defer handle.Close()

	// Load built in magic files
	magic_files.LoadDefaultMagic(handle)
	classification := handle.File("foobar.jpeg")
	return classification
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
func (g *Gigachat) getUploadData(fileReader io.Reader, fileName string) (*bytes.Buffer, *multipart.Writer, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		slog.Error("Ошибка создания части формы для файла:", err)
		return nil, nil, err
	}
	_, err = io.Copy(part, fileReader)
	if err != nil {
		slog.Error("Ошибка копирования файла в часть формы:", err)
		return nil, nil, err
	}
	_ = writer.WriteField("purpose", "general")
	err = writer.Close()
	if err != nil {
		slog.Error("Ошибка закрытия writer:", err)
		return nil, nil, err
	}
	return body, writer, err
}
func (g *Gigachat) UploadData(fileReader io.Reader, fileName string) (file UploadedFile, err error) {
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
	slog.Info("body:", body)
	var result UploadedFile
	err2 := json.Unmarshal(body, &result)
	if err2 != nil {
		log.Fatal(err2)
	}
	return result, nil
}

func (g *Gigachat) UploadFile(filePath string) (uploadedFile UploadedFile, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Ошибка открытия файла:", err)
		return
	}
	defer file.Close()
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
