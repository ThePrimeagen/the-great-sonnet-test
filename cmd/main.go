package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"theprimeagen.tv/claude/pkg/ai"
	"theprimeagen.tv/claude/pkg/cmd"
)

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

const SYSTEM = `Your goal is to make the unit tests pass by providing the JavaScript code based
on the errors you receive.

The errors provided in the user prompt will be the output from vitest

Your output will be used as a separate file, so make sure you create it as a module

Strip all markdown

All functions used in tests are imported without default imports

No poetry please

format your output as valid javascript ONLY.  Do not provide anything
but valid JavaScript code

Take a deep breath and DO NOT HALLUCINATE`

func run(ctx context.Context, claude *ai.ClaudeSonnet, name string, args []string, timeout int, output string) bool {
    wait := sync.WaitGroup{}
    wait.Add(1)

    fmt.Printf("Running Claude\n")
    res := NewCommandResults(ctx, name, args, &wait)
    wait.Wait()

    fmt.Printf("Test Runner(%d): %s\n", res.Code, res.String())
    if res.Code == 0 {
        return true
    }

    out, err := claude.ReadWithTimeout(res.String(), time.Second * time.Duration(timeout))
    fmt.Printf("Claude response: %s\n", out)
    if err != nil {
        slog.Error("received error from claude", "err", err)
        os.Exit(1)
    }

    file, err := os.OpenFile(output, os.O_RDWR|os.O_CREATE, 0644)
    defer file.Close()

    err = os.WriteFile(output, []byte(out), 0777)
    if err != nil {
        slog.Error("error from writing claude output", "err", err, "output", out)
        os.Exit(1)
    }

    return false
}

func main() {
    godotenv.Load()

    outputFile := ""
	flag.StringVar(&outputFile, "output", "", "the output file for where claude should put the code")

    timeout := 5
	flag.IntVar(&timeout, "timeout", 5, "the amonut of seconds to wait for claude sonnet")
	flag.Parse()

	args := flag.Args()
    name := args[0]
    args = args[1:]
    ctx := context.Background()
    claude := ai.NewClaudeSonnet(SYSTEM, ctx)

    count := 0
    success := false
    for range 11 {
        count++
        if run(ctx, claude, name, args, timeout, outputFile) {
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

