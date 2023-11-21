package cli

import (
	"fmt"
	"os"

	"github.com/egoodhall/tools/cmd/sshkeyd/internal/sync"
	"github.com/egoodhall/tools/pkg/daemon"
	"github.com/kardianos/service"
)

type SyncCmd struct {
	Flags CommonFlags `embed:""`
	Write bool        `name:"write" short:"w" default:"true" negatable:""`
}

func (cmd *SyncCmd) Run() error {
	if cmd.Flags.RefreshInterval == 0 {
		return cmd.run(nil)
	}

	dmn, err := daemon.New("sshkeyd", daemon.Periodic(cmd.Flags.RefreshInterval, cmd.run))
	if err != nil {
		return err
	}

	return dmn.Run()
}

func (cmd *SyncCmd) run(logger service.Logger) error {
	keys, err := sync.Run(sync.Config{
		AuthorizedKeysFile: os.ExpandEnv(cmd.Flags.AuthorizedKeysFile),
		RemoteKeyUrls:      sync.GetUrls(cmd.Flags.GithubUsers, cmd.Flags.GitlabUsers, cmd.Flags.Urls),
		Prune:              true,
	})
	if err != nil {
		return err
	}

	if !cmd.Write {
		fmt.Print(keys)
		return nil
	}
	return os.WriteFile(os.ExpandEnv(cmd.Flags.AuthorizedKeysFile), []byte(keys), 0600)
}
