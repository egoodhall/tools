package cli

import (
	"os"
	"time"

	"github.com/egoodhall/tools/cmd/sshkeyd/internal/sync"
	"github.com/egoodhall/tools/pkg/daemon"
	"github.com/kardianos/service"
)

type ServiceCmd struct {
	Flags           CommonFlags   `embed:""`
	RefreshInterval time.Duration `name:"refresh-interval" short:"i" default:"5m"`
	Action          string        `action:"" arg:"" enum:"install,start,stop,uninstall"`
}

func (cmd *ServiceCmd) Run() error {
	dmn, err := daemon.New("sshkeyd", "Sync authorized SSH keys file from URLs", daemon.Periodic(cmd.RefreshInterval, cmd.run))
	if err != nil {
		return err
	}
	return service.Control(dmn, cmd.Action)
}

func (cmd *ServiceCmd) run(logger service.Logger) error {
	keys, err := sync.Run(sync.Config{
		AuthorizedKeysFile: os.ExpandEnv(cmd.Flags.AuthorizedKeysFile),
		RemoteKeyUrls:      sync.GetUrls(cmd.Flags.GithubUsers, cmd.Flags.GitlabUsers, cmd.Flags.Urls),
		Prune:              true,
	})
	if err != nil {
		return err
	}
	return os.WriteFile(os.ExpandEnv(cmd.Flags.AuthorizedKeysFile), []byte(keys), 0600)
}
