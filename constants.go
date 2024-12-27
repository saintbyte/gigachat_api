package gigachat

// Авторизауия
const (
	GigaChatTokenFileEnv   = "GIGACHAT_TOKEN_FILE"   // Перемеменная среды с путем к файлу с токеном
	GigaChatExpiresFileEnv = "GIGACHAT_EXPIRES_FILE" // Переменная среды с путем к файл где время устревания токена
	GigaChatToken          = "GIGACHAT_TOKEN"        // Или токен берем из окружения
	GigaChatAuthData       = "GIGACHAT_AUTH_DATA"    // Данные дла авторизации чтоб получить токен
	GigaChatOauthUrl       = "https://ngw.devices.sberbank.ru:9443/api/v2/oauth"
)

// Ручки API
const (
	GigaChatApiHost            = "gigachat.devices.sberbank.ru" //1 - gigachat.devices.sberbank.ru 2  gigachat-preview.devices.sberbank.ru
	GigaChatModelsPath         = "/api/v1/models"
	GigaChatChatCompletionPath = "/api/v1/chat/completions"
	GigaChatEmbeddingsPath     = "/api/v1/embeddings"
)

// Настройки
const (
	GigaChatModel         = "GigaChat" // GigaChat, GigaChat-Pro, GigaChat-Max Если тестовый хост то GigaChat-Pro-preview
	GigaChatMaxTokens     = 16384
	MaxEmbeddingSize      = 8192
	GigaChatRoleUser      = "user"
	GigaChatRoleSystem    = "system"
	GigaChatRoleAssistant = "assistant"
	GigaChatRoleFunction  = "function" //  В сообщении с этой ролью передавайте в поле content валидный JSON-объект с результатами работы функции.
)

//Причины завершения.

const (
	GigaChatFinishReasonStop         = "stop"          // модель закончила формировать гипотезу и вернула полный ответ;
	GigaChatFinishReasonLength       = "length"        // достигнут лимит токенов в сообщении;
	GigaChatFinishReasonFunctionCall = "function_call" // указывает, что при запросе была вызвана встроенная функция или сгенерированы аргументы для пользовательской функции;
	GigaChatFinishReasonBlackList    = "blacklist"     //запрос попадает под тематические ограничения.
	GigaChatFinishReasonError        = "blacklist"     // ответ модели содержит невалидные аргументы пользовательской функции.
)

//Функции

const (
	GigaChatFunctionCallModeAuto = "auto"
	GigaChatFunctionCallModeNone = "none"
)
