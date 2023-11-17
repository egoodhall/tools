package cli

import (
	"github.com/egoodhall/tools/cmd/sshkeyd/pkg/sync"
)

type SyncCmd struct {
	Config sync.Config `embed:""`
}

func (cmd *SyncCmd) Run() error {
	return sync.Run(cmd.Config)
}
