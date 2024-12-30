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
	Role        string   `json:"role"`
	Content     string   `json:"content"`
	Attachments []string `json:"attachments,omitempty"`
}

// Статистика использования
type Usage struct {
	// Данные об использовании модели.
	PromptTokens     int `json:"prompt_tokens"`               // int32 Количество токенов во входящем сообщении (роль user).
	CompletionTokens int `json:"completion_tokens,omitempty"` // int32  Количество токенов, сгенерированных моделью (роль assistant).
	TotalTokens      int `json:"total_tokens,omitempty"`      //int32 Общее количество токенов.
}

type MessageResponse struct {
	Role           string     `json:"role"`
	Content        string     `json:"content"`
	DataForContext []struct{} `json:"data_for_context"`
}

type ChoicesResponse struct {
	Message          MessageRequest       `json:"message"`
	Index            int                  `json:"index"`
	FinishReason     string               `json:"finish_reason"`
	FunctionsStateId string               `json:"functions_state_id,omitempty"`
	FunctionCall     FunctionCallResponse `json:"function_call,omitempty"`
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
	FunctionCall      string           `json:"function_call,omitempty"`
	Functions         []Function       `json:"functions,omitempty"`
}

type ChatCompletionResponse struct {
	Choices []ChoicesResponse `json:"choices"`
	Created int               `json:"created"`
	Model   string            `json:"model"`
	Usage   Usage             `json:"usage"`
	Object  string            `json:"object"`
}

type FunctionCallResponse struct {
	Name      string            `json:"name"`
	Arguments map[string]string `json:"arguments,omitempty"`
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

// Функции

type Function struct {
	Name             string       `json:"name"`
	Description      string       `json:"description"`
	Parameters       Parameters   `json:"parameters"`
	ReturnParameters []Parameters `json:"return_parameters,omitempty"`
	FewShotExamples  []Example    `json:"few_shot_examples,omitempty"`
}

type Parameters struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties"`
	Required   []string            `json:"required"`
}

type Property struct {
	Type        string   `json:"type"`
	Enum        []string `json:"enum,omitempty"`
	Description string   `json:"description"`
}

type Example struct {
	Request string            `json:"request"`
	Params  map[string]string `json:"params"`
}
type UploadedFile struct {
	Bytes        int    `json:"bytes"`
	CreatedAt    int    `json:"created_at"`
	Filename     string `json:"filename"`
	Id           string `json:"id"`
	Object       string `json:"object"`
	Purpose      string `json:"purpose"`
	AccessPolicy string `json:"access_policy"`
}
