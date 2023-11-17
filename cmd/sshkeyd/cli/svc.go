package cli

import (
	"time"

	"github.com/egoodhall/tools/cmd/sshkeyd/pkg/sync"
	"github.com/egoodhall/tools/pkg/daemon"
	"github.com/kardianos/service"
)

type ServiceCmd struct {
	GithubUsers        []string      `name:"github"`
	GitlabUsers        []string      `name:"gitlab"`
	Urls               []string      `name:"url"`
	AuthorizedKeysFile string        `name:"keys-file" short:"f" required:"" default:"$HOME/.ssh/authorized_keys"`
	RefreshInterval    time.Duration `name:"refresh-interval" short:"i" default:"5m"`
	Action             string        `action:"" arg:"" enum:"install,start,stop,uninstall"`
}

func (cmd *ServiceCmd) Run() error {
	dmn, err := daemon.New(daemon.Periodic(cmd.RefreshInterval, cmd.run))
	if err != nil {
		return err
	}
	return service.Control(dmn, cmd.Action)
}

func (cmd *ServiceCmd) run(logger service.Logger) error {
	return sync.Run(sync.Config{
		AuthorizedKeysFile: cmd.AuthorizedKeysFile,
		GithubUsers:        cmd.GithubUsers,
		GitlabUsers:        cmd.GitlabUsers,
		Urls:               cmd.Urls,
		Prune:              true,
		Save:               true,
	})
}
