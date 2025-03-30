package main

import (
	"os"

	"github.com/gruyaume/certificates-operator/internal/charm"
	"github.com/gruyaume/goops/commands"
	"github.com/gruyaume/goops/environment"
)

func main() {
	hookCommand := &commands.HookCommand{}
	execEnv := &environment.ExecutionEnvironment{}
	logger := commands.NewLogger(hookCommand)
	hookName := environment.JujuHookName(execEnv)
	if hookName != "" {
		charm.HandleDefaultHook(hookCommand, logger)
		charm.SetStatus(hookCommand, logger)
		os.Exit(0)
	}
}
