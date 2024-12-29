package gigachat

import (
	"io"
	"net/http"
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
func (g *Gigachat) postRequest(url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest("POST", url, body)
	return req, err
}

func (g *Gigachat) UploadFile(fileReader io.Reader, fileName string) (file *File, err error) {
	url := g.getRequestUrl(GigaChatChatFileUploadPath)
}
