package schema

type ChatSessionResp struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type ChatResp struct {
	Choices          []ChatChoice `json:"choices"`
	Model            string       `json:"model"`
	PromptTokenUsage int          `json:"prompt_token_usage"`
	ChunkTokenUsage  int          `json:"chunk_token_usage"`
	Created          int64        `json:"created"`
	MessageID        int          `json:"message_id"`
	ParentID         int          `json:"parent_id"`
}

type ChatChoice struct {
	Index        int        `json:"index"`
	Delta        ChatDelta  `json:"delta"`
	FinishReason string     `json:"finish_reason,omitempty"`
	ToolCalls    []ToolCall `json:"tool_calls,omitempty"`
}

type ToolCall struct {
	Name      string `json:"name,omitempty"`
	Arguments string `json:"arguments,omitempty"`
}

type ChatDelta struct {
	Role    string `json:"role"`
	Content string `json:"content"`
	Type    string `json:"type"`
}
