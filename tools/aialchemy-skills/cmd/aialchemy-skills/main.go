// Command aialchemy-skills installs skills published by the AIAlchemy LiteLLM
// marketplace into the native skill-discovery locations of Codex, Claude Code,
// Kiro, and Hermes.
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/app"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	os.Exit(app.Main(ctx, os.Args[1:], os.Getenv, os.Stdout, os.Stderr))
}
