package gigachat

const (
	GigaChatTokenFileEnv           = "GIGACHAT_TOKEN_FILE"   // Перемеменная среды с путем к файлу с токеном
	GigaChatExpiresFileEnv         = "GIGACHAT_EXPIRES_FILE" // Переменная среды с путем к файл где время устревания токена
	GigaChatToken                  = "GIGACHAT_TOKEN"        // Или токен берем из окружения
	GigaChatAuthData               = "GIGACHAT_AUTH_DATA"    // Данные дла авторизации чтоб получить токен
	GigaChatOauthUrl               = "https://ngw.devices.sberbank.ru:9443/api/v2/oauth"
	GigaChatApiHost                = "gigachat.devices.sberbank.ru" //1 - gigachat.devices.sberbank.ru 2  gigachat-preview.devices.sberbank.ru
	GigaChatModelsPath             = "/api/v1/models"
	GigaChatChatCompletionPath     = "/api/v1/chat/completions"
	GigaChatEmbeddingsPath         = "/api/v1/embeddings"
	GigaChatModel                  = "GigaChat"
	GigaChatMaxTokens              = 16384
	MaxEmbeddingSize           int = 8192
)
