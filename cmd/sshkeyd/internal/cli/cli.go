package cli

import "time"

type CommonFlags struct {
	GithubUsers        []string      `name:"github"`
	GitlabUsers        []string      `name:"gitlab"`
	Urls               []string      `name:"url"`
	AuthorizedKeysFile string        `name:"keys-file" short:"f" required:"" default:"$HOME/.ssh/authorized_keys"`
	RefreshInterval    time.Duration `name:"refresh-interval" short:"i" default:"0m"`
}
