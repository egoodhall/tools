package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/egoodhall/tools/cmd/sshkeyd/sync"
	"github.com/egoodhall/tools/pkg/ssh"
)

func main() { fatalIfError(Main) }

func Main() error {
	flags, err := parseFlags()
	if err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}

	cfg := sync.Config{
		RemoteKeyUrls:      sync.GetUrls(flag.Args()),
		AuthorizedKeysFile: flags.KeysFile,
		Prune:              flags.Prune,
	}
	keys, err := cfg.ResolveAuthorizedKeys()
	if err != nil {
		return fmt.Errorf("failed to resolve authorized keys: %w", err)
	}

	slog.Info("Loaded keys")
	for _, key := range keys {
		slog.Info("Key", "key", key)
	}

	fmt.Println(keys)
	if flags.Write {
		if err := writeKeys(flags.KeysFile, keys); err != nil {
			return fmt.Errorf("failed to write keys: %w", err)
		}
	}
	return nil
}

func writeKeys(file string, keys []ssh.AuthorizedKey) error {
	bld := new(strings.Builder)
	for _, key := range keys {
		bld.WriteString(key.String())
		bld.WriteRune('\n')
	}

	if err := os.WriteFile(file, []byte(bld.String()), 0644); err != nil {
		return fmt.Errorf("failed to write new keys file: %w", err)
	}
	return nil
}

func fatalIfError(main func() error) {
	if err := main(); err != nil {
		slog.Error("Fatal error", "error", err)
		os.Exit(1)
	}
}
