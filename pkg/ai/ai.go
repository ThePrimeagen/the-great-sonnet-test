package ai

import (
	"context"
	"os"
	"time"

	"github.com/liushuangls/go-anthropic/v2"
	"github.com/sashabaranov/go-openai"
)

var foo = 1337
var startingTemp float32 = 0.3

type ClaudeSonnet struct {
    client *anthropic.Client
    system string
    ctx context.Context
    Temp float32
}

func NewClaudeSonnet(ctx context.Context) *ClaudeSonnet {
    client := anthropic.NewClient(os.Getenv("ANTHROPIC_API_KEY"))
    return &ClaudeSonnet{
        ctx: ctx,
        client: client,
        system: "",
        Temp: startingTemp,
    }
}

func (o *ClaudeSonnet) IncreaseTemp(amount float32) {
    if o.Temp + amount > 2 {
        o.Temp = 2
        return
    }

    o.Temp += amount
}

func (o *ClaudeSonnet) chat(chat string, ctx context.Context) (string, error) {

    resp, err := o.client.CreateMessages(ctx, anthropic.MessagesRequest{
        Model: anthropic.ModelClaude3Dot5Sonnet20240620,
        Temperature: &startingTemp,
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

func (c *ClaudeSonnet) SetTemp(temp float32) {
    c.Temp = temp
}

type OpenAIChat struct {
    client *openai.Client
    system string
    Temp float32
}

func (o *OpenAIChat) chat(chat string, ctx context.Context) (string, error) {
    resp, err := o.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
        Model: openai.GPT4o,
        Temperature: o.Temp,
        Seed: &foo,
        Messages: []openai.ChatCompletionMessage{
            {
                Role: openai.ChatMessageRoleUser,
                Content: chat,
            },
        },
    })
    if err != nil {
        return "", err
    }

    return resp.Choices[0].Message.Content, nil
}

func NewOpenAIChat(secret string, temp float32) *OpenAIChat {
    client := openai.NewClient(secret)
    return &OpenAIChat{
        client: client,
        Temp: temp,
    }
}

type StatefulOpenAIChat struct {
    ai *OpenAIChat
    ctx context.Context
}

func (s *StatefulOpenAIChat) Name() string {
    return "openai"
}

func NewStatefulOpenAIChat(ctx context.Context, temp float32) *StatefulOpenAIChat {
    return &StatefulOpenAIChat{
        ai: NewOpenAIChat(os.Getenv("OPENAI_API_KEY"), temp),
        ctx: ctx,
    }
}

func (s *StatefulOpenAIChat) prompt(p string, ctx context.Context) (string, error) {
    str, err := s.ai.chat(p, ctx)
    if err != nil {
        return "", err
    }
    return str, err
}

func (s *StatefulOpenAIChat) ReadWithTimeout(p string, t time.Duration) (string, error) {
    ctx, cancel := context.WithCancel(s.ctx)
    go func() {
        <-time.NewTimer(t).C
        cancel()
    }()

    return s.ai.chat(p, ctx)
}

func (s *StatefulOpenAIChat) SetTemp(temp float32) {
    s.ai.Temp = temp
}

type AI interface {
    ReadWithTimeout(p string, t time.Duration) (string, error)
    Name() string
    SetTemp(temp float32)
}
