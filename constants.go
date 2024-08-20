package gigachat

const (
	GigaChatTokenFileEnv           = "GIGACHAT_TOKEN_FILE"
	GigaChatExpiresFileEnv         = "GIGACHAT_EXPIRES_FILE"
	GigaChatToken                  = "GIGACHAT_TOKEN"
	GigaChatAuthData               = "GIGACHAT_AUTH_DATA"
	GigaChatOauthUrl               = "https://ngw.devices.sberbank.ru:9443/api/v2/oauth"
	GigaChatApiHost                = "gigachat.devices.sberbank.ru" //1 - gigachat.devices.sberbank.ru 2  gigachat-preview.devices.sberbank.ru
	GigaChatModelsPath             = "/api/v1/models"
	GigaChatChatCompletionPath     = "/api/v1/chat/completions"
	GigaChatEmbeddingsPath         = "/api/v1/embeddings"
	GigaChatModel                  = "GigaChat"
	GigaChatMaxTokens              = 16384
	MaxEmbeddingSize           int = 8192
)
