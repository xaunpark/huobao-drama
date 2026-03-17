package model

import (
	"context"
	"log"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino-ext/components/model/gemini"
	"github.com/cloudwego/eino-ext/components/model/openai"
)

func getApiConfig() {
	//
}

func NewChatFireChatModel() (*openai.ChatModel, error) {
	baseUrl := "https://api.chatfire.site/v1"
	apiKey := "sk-MMQEcMY6Nh2REIvtriLPdPZek4xHZPEhcfFGnQOhss9K2g6P"
	modelName := "gemini-3-flash-preview"

	// 初始化模型
	model, err := openai.NewChatModel(context.Background(), &openai.ChatModelConfig{
		APIKey:  apiKey,
		Model:   modelName,
		BaseURL: baseUrl,
	})
	if err != nil {
		log.Print(err)
	}
	return model, err
}

func NewGeminiChatModel() (*gemini.ChatModel, error) {
	return &gemini.ChatModel{}, nil
}

func NewArkChatModel() (*ark.ChatModel, error) {
	return &ark.ChatModel{}, nil
}
