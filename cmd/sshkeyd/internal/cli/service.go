package cli

import (
	"os"
	"time"

	"github.com/egoodhall/tools/pkg/daemon"
	"github.com/kardianos/service"
)

type ServiceCmd struct {
	Flags  CommonFlags `embed:""`
	Action string      `action:"" arg:"" enum:"install,start,stop,uninstall"`
}

func (cmd *ServiceCmd) Run() error {
	interval := cmd.Flags.RefreshInterval
	if interval == 0 {
		interval = 5 * time.Minute
	}

	dmn, err := daemon.NewController(
		"sshkeyd", "Sync authorized SSH keys file from URLs",
		os.Args[0], cmd.Flags.Args()...,
	)
	if err != nil {
		return err
	}
	return service.Control(dmn, cmd.Action)
}
