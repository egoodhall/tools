package cli

import (
	"os"

	"github.com/egoodhall/tools/pkg/daemon"
	"github.com/kardianos/service"
)

type ServiceCmd struct {
	Flags  CommonFlags `embed:""`
	Action string      `action:"" arg:"" enum:"install,start,stop,uninstall"`
}

func (cmd *ServiceCmd) Run() error {
	dmn, err := daemon.NewController(
		serviceName, serviceDescription,
		os.Args[0], cmd.Flags.Args("sync", "--interval=5m")...,
	)
	if err != nil {
		return err
	}
	return service.Control(dmn, cmd.Action)
}
