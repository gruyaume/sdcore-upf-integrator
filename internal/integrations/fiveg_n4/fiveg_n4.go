package fiveg_n4

import (
	"fmt"

	"github.com/gruyaume/goops"
)

func PublishUPFN4Information(hookContext *goops.HookContext, relationID string, hostname string, port int) error {
	portStr := fmt.Sprintf("%d", port)
	relationData := map[string]string{
		"upf_hostname": hostname,
		"upf_port":     portStr,
	}
	err := hookContext.Commands.RelationSet(relationID, true, relationData)
	if err != nil {
		return fmt.Errorf("could not set relation data: %w", err)
	}
	return nil
}
