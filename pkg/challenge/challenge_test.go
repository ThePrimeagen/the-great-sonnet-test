package challenge_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"theprimeagen.tv/claude/pkg/challenge"
)

func TestChallenge(t *testing.T) {
    for i := 0; i < 100; i++ {
        out := challenge.Challenge(i)
        require.Equal(t, out, 42)
    }
}

