package cli

import (
	"os"
	"strings"
)

const (
	serviceName        = "com.github.egoodhall.tools.sshkeyd"
	serviceDescription = "Sync authorized SSH keys file from URLs"
)

type CommonFlags struct {
	GithubUsers        []string `name:"github"`
	GitlabUsers        []string `name:"gitlab"`
	Urls               []string `name:"url"`
	AuthorizedKeysFile string   `name:"keys-file" short:"f" required:"" default:"$HOME/.ssh/authorized_keys"`
}

func (flags CommonFlags) Args(args ...string) []string {
	if len(flags.GithubUsers) > 0 {
		args = append(args, "--github="+strings.Join(flags.GithubUsers, ","))
	}
	if len(flags.GitlabUsers) > 0 {
		args = append(args, "--gitlab="+strings.Join(flags.GitlabUsers, ","))
	}
	if len(flags.Urls) > 0 {
		args = append(args, "--url="+strings.Join(flags.Urls, ","))
	}
	return append(args,
		"--keys-file="+os.ExpandEnv(flags.AuthorizedKeysFile),
	)
}
