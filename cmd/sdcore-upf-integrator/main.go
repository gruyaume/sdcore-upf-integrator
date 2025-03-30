package main

import (
	"os"

	"github.com/gruyaume/goops/commands"
	"github.com/gruyaume/goops/environment"
	"github.com/gruyaume/sdcore-upf-integrator/internal/charm"
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
