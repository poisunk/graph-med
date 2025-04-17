package prompt

import (
	"context"
	"github.com/mark3labs/mcp-go/mcp"
	"graph-med/internal/service"
	"strings"
)

const (
	KGInquirePromptTemplate = `
你将作为基于知识图谱的医学问答系统，你具有查询知识图谱数据库数据的能力。请仔细阅读以下信息，并按照指示进行作答。
以下是你应该知道的信息：
<disease_attrs>
{{DISEASE_ATTRS}}
</disease_attrs>
在回答问题时，请遵循以下指南:
1. 确保回答内容与问题相关，不提供无关信息。
2. 若问题有多个部分，需逐一进行回答，如需要获取疾病信息，应先获取疾病有那些信息。
3. 你是一个基于知识图谱的问答系统，你的首要任务是解决用户的问题，你可以根据数据库中的相关知识和问题进行回答。
`
)

type KGInquirePrompt struct {
	kgService *service.KGService
}

func NewUserInquirePrompt(kgService *service.KGService) *KGInquirePrompt {
	return &KGInquirePrompt{
		kgService: kgService,
	}
}

func (p *KGInquirePrompt) Name() string {
	return "user_inquire_prompt"
}

func (p *KGInquirePrompt) Definition() mcp.Prompt {
	return mcp.NewPrompt(p.Name(),
		mcp.WithPromptDescription("知识图谱查询助手系统提示词"),
	)
}

func (p *KGInquirePrompt) ToolHandlerFunc(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	attrs := make([]string, len(p.kgService.GetMedicalAttrs()))
	for k, _ := range p.kgService.GetMedicalAttrs() {
		attrs = append(attrs, k)
	}

	content := strings.Replace(KGInquirePromptTemplate, "{{DISEASE_ATTRS}}", strings.Join(attrs, "\n"), 1)

	return mcp.NewGetPromptResult(
		"user_inquire_prompt",
		[]mcp.PromptMessage{
			mcp.NewPromptMessage(
				"system",
				mcp.NewTextContent(content),
			),
		},
	), nil
}
