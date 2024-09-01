package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"
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

    cmdr := NewCmder(name, ctx).
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

