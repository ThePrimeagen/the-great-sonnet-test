package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"theprimeagen.tv/claude/pkg/ai"
	"theprimeagen.tv/claude/pkg/cmd"
)

type RunResults struct {
    Count int
    Prompt string
    RunCount int
    ClaudeTotalOutput []string
    Language string
}

func (r *RunResults) Save() {
    count := r.Count
    os.WriteFile("./data/count", []byte(fmt.Sprintf("%d", r.Count + 1)), 0644)
    os.WriteFile(fmt.Sprintf("./data/prompt%d", count), []byte(r.Prompt), 0644)

    os.WriteFile(fmt.Sprintf("./data/output%d", count), []byte(strings.Join(r.ClaudeTotalOutput, "\n---------------------------------------\n")), 0644)
}

func (r *RunResults) Push(result string) {
    r.ClaudeTotalOutput = append(r.ClaudeTotalOutput, result)
}

func NewRunResults(lang string) *RunResults {
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

    exports := getExports(lang)
    testRunner := getTestRunner(lang)

    prompt := strings.ReplaceAll(SYSTEM, "__LANGUAGE__", lang)
    prompt = strings.ReplaceAll(prompt, "__EXPORT__", exports)
    prompt = strings.ReplaceAll(prompt, "__TEST_RUNNER__", testRunner)

    return &RunResults{
        Count: count,
        Prompt: prompt,
        RunCount: 0,
        ClaudeTotalOutput: []string{},
        Language: lang,
    }
}

type CommandResults struct {
    Stdout []string
    Stderr []string
    Code int
}

func (c *CommandResults) String() string {
    out := strings.Join(c.Stdout, "\n")
    err := strings.Join(c.Stderr, "\n")

    return fmt.Sprintf("%s\n%s\n", out, err)
}

func NewCommandResults(ctx context.Context, name string, args []string, wait *sync.WaitGroup) *CommandResults {
    res := CommandResults{
        Code: -1,
        Stdout: []string{},
        Stderr: []string{},
    }

    cmdr := cmd.NewCmder(name, ctx).
        ApplyKVArgs(args).
        WithOutFn(func(b []byte) (int, error) {
            res.Stdout = append(res.Stdout, string(b))
            return len(b), nil
        }).
        WithErrFn(func(b []byte) (int, error) {
            res.Stderr = append(res.Stderr, string(b))
            return len(b), nil
        }).
        WithExitCodeFn(func(c int) {
            res.Code = c
            wait.Done()
        })

    go func() {
        err := cmdr.Run()
        slog.Error("cmdr.Run", "err", err)
    }()

    return &res
}

const EXPORT = "All functions used in tests are imported without default imports"

const SYSTEM = `You are a Staff Level Software Engineer with an incredible stock package
and your goal is to make the unit tests pass by providing the __LANGUAGE__ code based
on the errors you receive.

The errors provided in the user prompt will be the output from __TEST_RUNNER__

Your output will be used as a separate file, so make sure you create it as a module

Strip all markdown

__EXPORT__

No poetry please

format your output as valid __LANGUAGE__ ONLY.  Do not provide anything
but valid __LANGUAGE__ code, NO MARKDOWN, NO JSON, NO XML, JUST CODE

DO NOT HALLUCINATE

YOU SHALL WRITE ALL OF THE CODE
`

func run(ctx context.Context, claude *ai.ClaudeSonnet, name string, args []string, timeout int, output string) (*CommandResults, string) {
    wait := sync.WaitGroup{}
    wait.Add(1)

    fmt.Printf("Running Claude\n")
    res := NewCommandResults(ctx, name, args, &wait)
    wait.Wait()

    fmt.Printf("Test Runner(%d): %s\n", res.Code, res.String())
    if res.Code == 0 {
        return res, ""
    }

    out, err := claude.ReadWithTimeout(res.String(), time.Second * time.Duration(timeout))
    fmt.Printf("Claude response: %s\n", out)
    if err != nil {
        slog.Error("received error from claude", "err", err)
        os.Exit(1)
    }

    if strings.HasPrefix(out, "```") {
        out = out[1:len(out) - 1]
    }

    file, err := os.OpenFile(output, os.O_RDWR|os.O_CREATE, 0644)
    defer file.Close()

    err = os.WriteFile(output, []byte(out), 0644)
    if err != nil {
        slog.Error("error from writing claude output", "err", err, "output", out)
        os.Exit(1)
    }

    return res, out
}

func getExports(lang string) string {
    switch (lang) {
    case "JavaScript":
        return EXPORT
    case "Golang":
        return ""
    default:
        slog.Error("invalid language option.  Javascript or Golang", "lang", lang)
        os.Exit(1)
    }
    return ""
}

func getTestRunner(lang string) string {
    switch (lang) {
    case "JavaScript":
        return "vitest"
    case "Golang":
        return "go's built in test runner"
    default:
        slog.Error("invalid language option.  Javascript or Golang", "lang", lang)
        os.Exit(1)
    }
    return ""
}

func main() {
    godotenv.Load()

    outputFile := ""
	flag.StringVar(&outputFile, "output", "", "the output file for where claude should put the code")

    language := "JavaScript"
	flag.StringVar(&language, "lang", "JavaScript", "the language of the code")

    timeout := 5
	flag.IntVar(&timeout, "timeout", 5, "the amonut of seconds to wait for claude sonnet")
	flag.Parse()

    runResults := NewRunResults(language)
    defer runResults.Save()

	args := flag.Args()
    name := args[0]
    args = args[1:]
    ctx := context.Background()
    claude := ai.NewClaudeSonnet(runResults.Prompt, ctx)

    count := 0
    success := false
    for range 11 {
        count++
        res, claudeOutput := run(ctx, claude, name, args, timeout, outputFile)
        runResults.RunCount = count
        runResults.Push(claudeOutput)

        if res.Code == 0 {
            success = true
            break;
        }
    }

    status := "successful"
    if !success {
        status = "failed"
    }

    fmt.Printf("claude took %d tries and was %s\n", count, status)
}

