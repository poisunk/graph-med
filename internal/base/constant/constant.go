package constant

const (
	MODEL_DOUBAO_LITE_32K    = "doubao-1-5-lite-32k-250115"
	USER_QUERY_SYSTEM_PROMPT = `
以下是你应该知道的信息：
<disease_attrs>
{{DISEASE_ATTRS}}
</disease_attrs>
`
	USER_QUERY_USER_PROMPT = `
你将作为基于知识图谱的医学问答系统。请仔细阅读以下信息，并按照指示进行作答。
问题:
<question>
{{QUESTION}}
</question>
在回答问题时，请遵循以下指南:
1. 保持回答准确、清晰、简洁，避免使用过于复杂的医学术语，若必须使用，需进行简单解释。
2. 确保回答内容与问题相关，不提供无关信息。
3. 若问题有多个部分，需逐一进行回答。
`
	MAX_HISTORY_LENGTH = 50
)
