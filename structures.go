package gigachat

// Ответ с токенос
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresAt   int64  `json:"expires_at"`
}

// Модель
type ModelItem struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	OwnedBy string `json:"owned_by"`
}

// Ответ на запрос списка моделей
type ModelsResponse struct {
	Object string      `json:"object"`
	Data   []ModelItem `json:"data"`
}

type MessageRequest struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Статистика использования
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens,omitempty"`
	TotalTokens      int `json:"total_tokens,omitempty"`
}

type MessageResponse struct {
	Role           string     `json:"role"`
	Content        string     `json:"content"`
	DataForContext []struct{} `json:"data_for_context"`
}

type ChoicesResponse struct {
	Message      MessageRequest `json:"message"`
	Index        int            `json:"index"`
	FinishReason string         `json:"finish_reason"`
}

type ChatCompletionRequest struct {
	Model             string           `json:"model"`
	Messages          []MessageRequest `json:"messages"`
	Stream            bool             `json:"stream"`
	RepetitionPenalty int              `json:"repetition_penalty"`
	Temperature       float32          `json:"temperature"`
	TopP              float32          `json:"top_p"`
	MaxTokens         int              `json:"max_tokens"`
	UpdateInterval    int              `json:"update_interval"`
}

type ChatCompletionResponse struct {
	Choices []ChoicesResponse `json:"choices"`
	Created int               `json:"created"`
	Model   string            `json:"model"`
	Usage   Usage             `json:"usage"`
	Object  string            `json:"object"`
}

type EmbeddingsRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
}

type EmbeddingsResponse struct {
	Object string `json:"object"`
	Data   []struct {
		Object    string    `json:"object"`
		Embedding []float32 `json:"embedding"`
		Index     int       `json:"index"`
		Usage     Usage     `json:"usage"`
	} `json:"data"`
	Model string `json:"model"`
}
