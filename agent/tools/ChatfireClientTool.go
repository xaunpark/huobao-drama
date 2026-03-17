package tools

import (
	"context"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

type ChatFireClientTool struct {
}

func NewChatFireClientTool() *ChatFireClientTool {
	return &ChatFireClientTool{}
}

func (t *ChatFireClientTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	locationParam := &schema.ParameterInfo{
		Type:     schema.DataType("string"),
		Required: true,
		Desc:     "location",
	}
	parameterInfo := map[string]*schema.ParameterInfo{
		"location": locationParam,
	}
	paramsOneOf := schema.NewParamsOneOfByParams(parameterInfo)
	toolInfo := &schema.ToolInfo{
		Name:        "weather_query_tool",
		Desc:        "tool for query location weather",
		ParamsOneOf: paramsOneOf,
	}
	return toolInfo, nil
}

func (t *ChatFireClientTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	// 执行逻辑：调用api或其他自定义处理逻辑
	// 将内容返回出去
	return "", nil
}
