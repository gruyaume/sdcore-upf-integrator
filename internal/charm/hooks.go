package charm

import (
	"fmt"

	"github.com/gruyaume/goops"
	"github.com/gruyaume/goops/commands"
	"github.com/gruyaume/sdcore-upf-integrator/internal/integrations/fiveg_n4"
)

const (
	N4Integration = "fiveg_n4"
)

func isConfigValid(hookContext *goops.HookContext) error {
	n4Hostname, err := hookContext.Commands.ConfigGetString("n4-hostname")
	if err != nil {
		if err == commands.ErrConfigNotSet {
			return fmt.Errorf("n4-hostname config is not set")
		}
		return fmt.Errorf("n4-hostname config is not valid: %w", err)
	}
	if n4Hostname == "" {
		return fmt.Errorf("n4-hostname config is empty")
	}

	n4Port, err := hookContext.Commands.ConfigGetInt("n4-port")
	if err != nil {
		if err == commands.ErrConfigNotSet {
			return fmt.Errorf("n4-port config is not set")
		}
		return fmt.Errorf("n4-port config is not valid: %w", err)
	}
	if n4Port <= 0 {
		return fmt.Errorf("n4-port config is invalid: %d", n4Port)
	}

	return nil
}

func syncN4Integration(hookContext *goops.HookContext) error {
	n4Hostname, err := hookContext.Commands.ConfigGetString("n4-hostname")
	if err != nil {
		return fmt.Errorf("could not get n4-hostname config: %w", err)
	}
	n4Port, err := hookContext.Commands.ConfigGetInt("n4-port")
	if err != nil {
		return fmt.Errorf("could not get n4-port config: %w", err)
	}
	relationIDs, err := hookContext.Commands.RelationIDs(N4Integration)
	if err != nil {
		return fmt.Errorf("could not get relation IDs: %w", err)
	}
	for _, relationID := range relationIDs {
		err = fiveg_n4.PublishUPFN4Information(hookContext, relationID, n4Hostname, n4Port)
		if err != nil {
			return fmt.Errorf("could not publish UPF N4 information: %w", err)
		}
		hookContext.Commands.JujuLog(commands.Info, "Published UPF N4 information to relation ID:", relationID)
	}
	return nil
}

func isN4RelationCreated(hookContext *goops.HookContext) bool {
	relationIDs, err := hookContext.Commands.RelationIDs(N4Integration)
	if err != nil {
		return false
	}
	if len(relationIDs) == 0 {
		return false
	}
	return true
}

func HandleDefaultHook(hookContext *goops.HookContext) {
	isLeader, err := hookContext.Commands.IsLeader()
	if err != nil {
		hookContext.Commands.JujuLog(commands.Error, "Could not check if unit is leader:", err.Error())
		return
	}
	if !isLeader {
		hookContext.Commands.JujuLog(commands.Warning, "Unit is not leader")
		return
	}
	err = isConfigValid(hookContext)
	if err != nil {
		hookContext.Commands.JujuLog(commands.Warning, "Config is not valid:", err.Error())
		return
	}

	err = syncN4Integration(hookContext)
	if err != nil {
		hookContext.Commands.JujuLog(commands.Error, "Could not sync n4 integration:", err.Error())
		return
	}
}

func SetStatus(hookContext *goops.HookContext) {
	var status = commands.StatusActive
	var message = ""
	err := isConfigValid(hookContext)
	if err != nil {
		status = commands.StatusBlocked
		message = err.Error()
	} else if !isN4RelationCreated(hookContext) {
		status = commands.StatusBlocked
		message = "N4 relation not created"
	}
	err = hookContext.Commands.StatusSet(status, message)
	if err != nil {
		hookContext.Commands.JujuLog(commands.Error, "Could not set status:", err.Error())
		return
	}
	hookContext.Commands.JujuLog(commands.Info, "Status set to active")
}
