package runner

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"theprimeagen.tv/claude/pkg/ai"
	"theprimeagen.tv/claude/pkg/cmd"
	"theprimeagen.tv/claude/pkg/prompt"
)

type Runner struct {
    Count int
    Timeout time.Duration
    OutputFile string

    RunCount int
    TotalOutput []string
    TotalCode []string
    TotalReasoning []string
    ClaudeReasoningOutput []string

    Name string
    Args []string

    Language string
    TestFile string

    Prompt string
    Code string
    Error string
    Reasoning string
    LastExitCode int
}

type RunnerParams struct {
    Name string
    Args []string
    Lang string
    TestFilePath string
    OutputFilePath string
}

func NewRunner(params RunnerParams, prompt string) *Runner {
    countBytes, err := os.ReadFile("./data/count")
    if err != nil {
        slog.Error("could not read data count", "err", err)
        os.Exit(1)
    }

    count, err := strconv.Atoi(strings.TrimSpace(string(countBytes)))
    if err != nil {
        slog.Error("error parsing countBytes", "err", err)
        os.Exit(1)
    }

    fileBytes, err := os.ReadFile(params.TestFilePath)
    if err != nil {
        slog.Error("could not read test file", "err", err)
        os.Exit(1)
    }

    file := string(fileBytes)

    return &Runner{
        Count: count,
        RunCount: 0,
        TotalOutput: []string{},
        TotalCode: []string{},
        TotalReasoning: []string{},
        Language: params.Lang,
        TestFile: file,
        Error: "",
        Reasoning: "",

        Prompt: prompt,
        Timeout: time.Second * 30,
        Name: params.Name,
        Args: params.Args,
        OutputFile: params.OutputFilePath,
        LastExitCode: -1,
    }
}

func (r *Runner) Save() {
    suffix := ".js"
    if r.Language == "Golang" {
        suffix = ".go"
    }

    count := r.Count
    os.WriteFile("./data/count", []byte(fmt.Sprintf("%d", r.Count + 1)), 0644)

    output := []string{}
    for i, o := range r.TotalOutput {
        output = append(output, "--------------", fmt.Sprintf("%d", i), "------------------")
        output = append(output, o)
        output = append(output, "\n")
        output = append(output, "\n")
        if i < len(r.TotalCode) {
            output = append(output, r.TotalCode[i])
            output = append(output, "\n")
            output = append(output, "\n")
        }

        if i < len(r.TotalReasoning) {
            output = append(output, r.TotalReasoning[i])
            output = append(output, "\n")
            output = append(output, "\n")
        }
    }

    os.WriteFile(fmt.Sprintf("./data/output%d", count), []byte(strings.Join(output, "")), 0644)
    os.WriteFile(fmt.Sprintf("./data/final.code%d%s", count, suffix), []byte(r.TotalOutput[len(r.TotalOutput)-1]), 0644)

    if len(r.ClaudeReasoningOutput) > 0 {
        os.WriteFile(fmt.Sprintf("./data/final.reason%d%s", count, suffix), []byte(r.ClaudeReasoningOutput[len(r.TotalOutput)-1]), 0644)
    }

    os.WriteFile(fmt.Sprintf("./data/test%d%s", count, suffix), []byte(r.TestFile), 0644)
}

func (r *Runner) ReasoningOutput(result string) {
    r.ClaudeReasoningOutput = append(r.ClaudeReasoningOutput, result)
}

func (r *Runner) Done() bool {
    return r.LastExitCode == 0
}

func (r *Runner) PrintResults() {
    if r.Done() {
        fmt.Printf("[38;2;74;246;38mS");
    } else {
        fmt.Printf("[38;2;237;67;55mF");
    }
}

func (r *Runner) RunTest(ctx context.Context) {
    wait := sync.WaitGroup{}
    wait.Add(1)

    res := cmd.NewCommandResults(ctx, r.Name, r.Args, &wait)
    wait.Wait()

    slog.Info("RunTest", "code", res.Code, "result", res.String())

    r.Error = strings.Join(res.Stderr, "\n")
    r.TotalOutput = append(r.TotalOutput, strings.Join(res.Stdout, "\n"))
    r.LastExitCode = res.Code
    r.RunCount++
}

func (r *Runner) ToPromptParams() prompt.PromptParams {
    params := prompt.CreatePromptParamsFromLanguage(r.Language)
    params.TestFile = r.TestFile
    params.Code = r.Code
    params.Error = r.Error
    params.Reasoning = r.Reasoning

    return params
}

func (r *Runner) RunCodeGen(ctx context.Context, ai ai.AI) {

    out, err := ai.ReadWithTimeout(
        prompt.CodeGenPrompt(r.ToPromptParams(), r.Prompt), r.Timeout)

    if err != nil {
        slog.Error("unable to receive code gen from claude", "err", err)
        os.Exit(1)
    }

    if strings.HasPrefix(out, "```") {
        out = out[1:len(out) - 1]
    }

    r.Code = out
    r.TotalCode = append(r.TotalCode, out)

    file, err := os.OpenFile(r.OutputFile, os.O_RDWR|os.O_CREATE, 0644)
    defer file.Close()

    err = os.WriteFile(r.OutputFile, []byte(out), 0644)
    if err != nil {
        slog.Error("error from writing claude output", "err", err, "output", out)
        os.Exit(1)
    }
}

func (r *Runner) RunReasoning(ctx context.Context, ai ai.AI) {

    reasoning, err := ai.ReadWithTimeout(
        prompt.ReasonPrompt(r.ToPromptParams()), r.Timeout)

    if err != nil {
        slog.Error("unable to receive reasoning from claude", "err", err)
        os.Exit(1)
    }

    slog.Info("Run Reasoning", "reasoning", reasoning)
    r.Reasoning = reasoning
    r.TotalReasoning = append(r.TotalReasoning, reasoning)

}
