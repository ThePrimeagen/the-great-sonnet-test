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

type CommandResults struct {
    Stdout []string
    Stderr []string
    Code int
}

func NewCommandResults(ctx context.Context, name string, args []string, wait *sync.WaitGroup) CommandResults {
    stdout := []string{}
    stderr := []string{}
    var code int;
    cmdr := cmd.NewCmder(name, ctx).
        ApplyKVArgs(args).
        WithOutFn(func(b []byte) (int, error) {
            stdout = append(stdout, string(b))
            return len(b), nil
        }).
        WithErrFn(func(b []byte) (int, error) {
            stderr = append(stderr, string(b))
            return len(b), nil
        }).
        WithExitCodeFn(func(c int) {
            code = c
            wait.Done()
        })

    go func() {
        err := cmdr.Run()
        slog.Error("we got an error boys", "err", err)
        panic("crashing program since command runner failed to run")
    }()

    return CommandResults{
        Code: code,
        Stdout: stdout,
        Stderr: stderr,
    }
}

const SYSTEM = `Your goal is to make the unit tests pass by providing the JavaScript code based
on the errors you receive.

The errors provided in the user prompt will be the output from vitest

format your output as valid javascript ONLY.  Do not provide anything
but valid JavaScript code

Take a deep breath and DO NOT HALLUCINATE`

func run(ctx context.Context, claude *ai.ClaudeSonnet, name string, args []string, timeout int, output string) {
    wait := sync.WaitGroup{}
    wait.Add(1)
    res := NewCommandResults(ctx, name, args, &wait)

    fmt.Printf("Claude test run exit %d\n", res.Code)
    if res.Code == 0 {
        return
    }

    out, err := claude.ReadWithTimeout(strings.Join(res.Stderr, "\n"), time.Second * time.Duration(timeout))
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

    for range 11 {
        run(ctx, claude, name, args, timeout, outputFile)
    }

}

