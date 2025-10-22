package sync

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/egoodhall/tools/pkg/ssh"
)

type Config struct {
	RemoteKeyUrls      []string
	AuthorizedKeysFile string
	Prune              bool
}

func (cfg Config) ResolveAuthorizedKeys() ([]ssh.AuthorizedKey, error) {
	seenKeys := make(map[ssh.AuthorizedKey]struct{})
	allKeys := make([]ssh.AuthorizedKey, 0)

	if !cfg.Prune {
		localKeys, err := cfg.loadLocalKeys()
		if err != nil {
			return nil, fmt.Errorf("load local keys: %w", err)
		}
		for _, key := range localKeys {
			if _, seen := seenKeys[key.WithoutComment()]; !seen {
				seenKeys[key.WithoutComment()] = struct{}{}
				allKeys = append(allKeys, key)
			}
		}
		slog.Info("loaded local keys", "count", len(localKeys))
	}

	remoteKeys, err := cfg.loadRemoteKeys()
	if err != nil {
		return nil, fmt.Errorf("load remote keys: %w", err)
	}
	slog.Info("loaded remote keys", "count", len(remoteKeys))

	for _, key := range remoteKeys {
		if _, seen := seenKeys[key.WithoutComment()]; !seen {
			seenKeys[key.WithoutComment()] = struct{}{}
			allKeys = append(allKeys, key)
		}
	}

	return allKeys, nil
}

func (cfg *Config) loadLocalKeys() ([]ssh.AuthorizedKey, error) {
	seenKeys := make(map[ssh.AuthorizedKey]struct{})
	allKeys := make([]ssh.AuthorizedKey, 0)

	data, err := os.ReadFile(cfg.AuthorizedKeysFile)
	if os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("read authorized keys file: %w", err)
	}

	for _, line := range strings.Split(strings.Trim(string(data), "\n\t "), "\n") {
		if key, err := ssh.ParseAuthorizedKey(line); err != nil {
			return nil, err
		} else {
			seenKeys[key] = struct{}{}
			allKeys = append(allKeys, key)
		}
	}

	return allKeys, nil
}

func (cfg *Config) loadRemoteKeys() ([]ssh.AuthorizedKey, error) {
	seenKeys := make(map[ssh.AuthorizedKey]struct{})
	allKeys := make([]ssh.AuthorizedKey, 0)

	for _, url := range cfg.RemoteKeyUrls {
		keys, err := cfg.getRemoteKeys(url)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", url, err)
		}

		for _, key := range keys {
			if _, seen := seenKeys[key]; !seen {
				seenKeys[key.WithoutComment()] = struct{}{}
				key.Comment = url
				allKeys = append(allKeys, key)
			}
		}
	}

	return allKeys, nil
}

func (cfg *Config) getRemoteKeys(url string) ([]ssh.AuthorizedKey, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetch keys: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response: %d", response.StatusCode)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	keys := make([]ssh.AuthorizedKey, 0)
	for _, line := range strings.Split(strings.Trim(string(data), "\n\t "), "\n") {
		if key, err := ssh.ParseAuthorizedKey(line); err != nil {
			return nil, err
		} else {
			keys = append(keys, key)
		}
	}

	return keys, nil
}
