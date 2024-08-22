package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"sync"

	"github.com/joho/godotenv"
	"theprimeagen.tv/claude/pkg/cmd"
)

type CommandResults struct {
    Stdout []string
    Stderr []string
    Code int
}

func NewCommandResults(name string, args []string, wait *sync.WaitGroup) CommandResults {
    stdout := []string{}
    stderr := []string{}
    var code int;
    ctx := context.Background()
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

func main() {
    godotenv.Load()

    outputFile := ""
	flag.StringVar(&outputFile, "output", "", "the output file for where claude should put the code")
	flag.Parse()

	args := flag.Args()
    wait := sync.WaitGroup{}

    for range 10 {
        wait.Add(1)
        res := NewCommandResults(args[0], args[1:], &wait)

        fmt.Printf("Claude test run exit %d\n", res.Code)
        if res.Code == 0 {
            break
        }
    }

}

