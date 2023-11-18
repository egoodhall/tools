package cli

import (
	"fmt"
	"os"

	"github.com/egoodhall/tools/cmd/sshkeyd/internal/sync"
)

type SyncCmd struct {
	Flags CommonFlags `embed:""`
	Write bool        `name:"write" short:"w" default:"true"`
}

func (cmd *SyncCmd) Run() error {
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
