package cli

type CommonFlags struct {
	GithubUsers        []string `name:"github"`
	GitlabUsers        []string `name:"gitlab"`
	Urls               []string `name:"url"`
	AuthorizedKeysFile string   `name:"keys-file" short:"f" required:"" default:"$HOME/.ssh/authorized_keys"`
}
