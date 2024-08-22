package ai

import (
	"context"
	"os"
	"time"

	"github.com/liushuangls/go-anthropic/v2"
)

var foo = 1337
var temperature float32 = 0.55

type ClaudeSonnet struct {
    client *anthropic.Client
    system string
    ctx context.Context
}

func NewClaudeSonnet(system string, ctx context.Context) *ClaudeSonnet {
    client := anthropic.NewClient(os.Getenv("ANTHROPIC_API_KEY"))
    return &ClaudeSonnet{
        ctx: ctx,
        client: client,
        system: system,
    }
}

func (o *ClaudeSonnet) chat(chat string, ctx context.Context) (string, error) {

    resp, err := o.client.CreateMessages(ctx, anthropic.MessagesRequest{
        Model: anthropic.ModelClaude3Sonnet20240229,
        Temperature: &temperature,
        MaxTokens: 1000,
        System: o.system,
        Messages: []anthropic.Message{
            {
                Role: anthropic.RoleUser,
                Content: []anthropic.MessageContent{
                    anthropic.NewTextMessageContent(chat),
                },
            },
        },
    })

    if err != nil {
        return "", err
    }

    return resp.Content[0].GetText(), nil
}

func (s *ClaudeSonnet) Name() string {
    return "anthropic"
}

func (s *ClaudeSonnet) ReadWithTimeout(p string, t time.Duration) (string, error) {
    ctx, cancel := context.WithCancel(s.ctx)
    go func() {
        <-time.NewTimer(t).C
        cancel()
    }()

    return s.chat(p, ctx)
}


