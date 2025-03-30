package charm

import (
	"fmt"

	"github.com/gruyaume/goops/commands"
	"github.com/gruyaume/sdcore-upf-integrator/internal/integrations/fiveg_n4"
)

const (
	N4Integration = "fiveg_n4"
)

func isConfigValid(hookCommand *commands.HookCommand) error {
	n4Hostname, err := commands.ConfigGetString(hookCommand, "n4-hostname")
	if err != nil {
		if err == commands.ErrConfigNotSet {
			return fmt.Errorf("n4-hostname config is not set")
		}
		return fmt.Errorf("n4-hostname config is not valid: %w", err)
	}
	if n4Hostname == "" {
		return fmt.Errorf("n4-hostname config is empty")
	}

	n4Port, err := commands.ConfigGetInt(hookCommand, "n4-port")
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

func syncN4Integration(hookCommand *commands.HookCommand, logger *commands.Logger) error {
	n4Hostname, err := commands.ConfigGetString(hookCommand, "n4-hostname")
	if err != nil {
		return fmt.Errorf("could not get n4-hostname config: %w", err)
	}
	n4Port, err := commands.ConfigGetInt(hookCommand, "n4-port")
	if err != nil {
		return fmt.Errorf("could not get n4-port config: %w", err)
	}
	relationIDs, err := commands.RelationIDs(hookCommand, N4Integration)
	if err != nil {
		return fmt.Errorf("could not get relation IDs: %w", err)
	}
	for _, relationID := range relationIDs {
		err = fiveg_n4.PublishUPFN4Information(hookCommand, relationID, n4Hostname, n4Port)
		if err != nil {
			return fmt.Errorf("could not publish UPF N4 information: %w", err)
		}
		logger.Info("Published UPF N4 information to relation ID:", relationID)
	}
	return nil
}

func isN4RelationCreated(hookCommand *commands.HookCommand) bool {
	relationIDs, err := commands.RelationIDs(hookCommand, N4Integration)
	if err != nil {
		return false
	}
	if len(relationIDs) == 0 {
		return false
	}
	return true
}

func HandleDefaultHook(hookCommand *commands.HookCommand, logger *commands.Logger) {
	isLeader, err := commands.IsLeader(hookCommand)
	if err != nil {
		logger.Error("Could not check if unit is leader:", err.Error())
		return
	}
	if !isLeader {
		logger.Warning("Unit is not leader")
		return
	}
	err = isConfigValid(hookCommand)
	if err != nil {
		logger.Warning("Config is not valid:", err.Error())
		return
	}

	err = syncN4Integration(hookCommand, logger)
	if err != nil {
		logger.Error("Could not sync n4 integration:", err.Error())
		return
	}
}

func SetStatus(hookCommand *commands.HookCommand, logger *commands.Logger) {
	var status = commands.StatusActive
	var message = ""
	err := isConfigValid(hookCommand)
	if err != nil {
		status = commands.StatusBlocked
		message = err.Error()
	} else if !isN4RelationCreated(hookCommand) {
		status = commands.StatusBlocked
		message = "N4 relation not created"
	}
	err = commands.StatusSet(hookCommand, status, message)
	if err != nil {
		logger.Error("Could not set status:", err.Error())
		return
	}
	logger.Info("Status set to active")
}
