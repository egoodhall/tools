package cli

import (
	"os"
	"strings"
	"time"
)

type CommonFlags struct {
	GithubUsers        []string      `name:"github"`
	GitlabUsers        []string      `name:"gitlab"`
	Urls               []string      `name:"url"`
	AuthorizedKeysFile string        `name:"keys-file" short:"f" required:"" default:"$HOME/.ssh/authorized_keys"`
	RefreshInterval    time.Duration `name:"interval" default:"0m" hidden:""`
}

func (flags CommonFlags) Args() []string {
	args := []string{
		"--keys-file=" + os.ExpandEnv(flags.AuthorizedKeysFile),
		"--interval=" + flags.RefreshInterval.String(),
	}
	if len(flags.GithubUsers) > 0 {
		args = append(args, "--github="+strings.Join(flags.GithubUsers, ","))
	}
	if len(flags.GitlabUsers) > 0 {
		args = append(args, "--gitlab="+strings.Join(flags.GitlabUsers, ","))
	}
	if len(flags.Urls) > 0 {
		args = append(args, "--url="+strings.Join(flags.Urls, ","))
	}
	return args
}
