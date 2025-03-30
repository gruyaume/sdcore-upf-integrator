package fiveg_n4

import (
	"fmt"

	"github.com/gruyaume/goops/commands"
)

func PublishUPFN4Information(hookCommand *commands.HookCommand, relationID string, hostname string, port int) error {
	portStr := fmt.Sprintf("%d", port)
	relationData := map[string]string{
		"upf_hostname": hostname,
		"upf_port":     portStr,
	}
	err := commands.RelationSet(hookCommand, relationID, true, relationData)
	if err != nil {
		return fmt.Errorf("could not set relation data: %w", err)
	}
	return nil
}
