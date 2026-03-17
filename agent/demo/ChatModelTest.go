package main

import (
	"context"
	"log"

	"github.com/cloudwego/eino/schema"
	"github.com/drama-generator/backend/agent/model"
)

func main() {
	chatModel := model.NewChatFireChatModel()
	ctx := context.Background()
	msg, err := chatModel.Generate(ctx, []*schema.Message{
		{
			Role:    schema.User,
			Content: "你好",
		},
	})
	if err != nil {
		log.Print(err)
	}
	log.Print(msg.Content)
}
