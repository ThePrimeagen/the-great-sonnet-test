package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"theprimeagen.tv/claude/pkg/ai"
	"theprimeagen.tv/claude/pkg/runner"
)

func main() {
    godotenv.Load()

    outputFile := ""
	flag.StringVar(&outputFile, "output", "", "the output file for where claude should put the code")

    logFile := ""
	flag.StringVar(&logFile, "log", "/tmp/tgst.log", "the log file of the programs output")

    promptStr := ""
	flag.StringVar(&promptStr, "prompt", "prompt/gencode1.prompt", "the prompt to use")

    language := "JavaScript"
	flag.StringVar(&language, "lang", "JavaScript", "the language of the code")

    temp64 := 0.3
	flag.Float64Var(&temp64, "temp", 0.3, "The temp of the model")

    model := "claude"
	flag.StringVar(&model, "model", "claude", "the model to use")

    temp := float32(temp64)

    testFileStr := ""
	flag.StringVar(&testFileStr, "tf", "", "the test file that is ran")
    flag.Parse()

	args := flag.Args()

    name := args[0]
    args = args[1:]
    prompt, err := os.ReadFile(promptStr)
    if err != nil {
        log.Fatalf("error reading prompt file: %s\n", err)
    }

    file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatal(err)
    }

    log.SetOutput(file)

    slog.Warn("prompt string", "prompt", string(prompt))
    runResults := runner.NewRunner(runner.RunnerParams{
        Lang: language,
        Name: name,
        Args: args,
        OutputFilePath: outputFile,
        TestFilePath: testFileStr,
    }, string(prompt))

    defer runResults.Save()

    ctx := context.Background()

    var aiModel ai.AI
    if model == "claude" {
        aiModel = ai.NewClaudeSonnet(ctx)
    } else {
        aiModel = ai.NewStatefulOpenAIChat(ctx, temp)
    }

    aiModel.SetTemp(temp)

    for range 20 {
        runResults.RunCodeGen(ctx, aiModel)
        runResults.RunTest(ctx)
        runResults.PrintResults()
        if runResults.Done() {
            break;
        }
        runResults.RunReasoning(ctx, aiModel)
    }
}

