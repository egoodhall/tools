package sync

import "fmt"

func GetUrls(githubUsers []string, gitlabUsers, extraUrls []string) []string {
	urls := make([]string, len(extraUrls))
	copy(urls, extraUrls)

	for _, user := range githubUsers {
		urls = append(urls, fmt.Sprintf("https://github.com/%s.keys", user))
	}

	for _, user := range gitlabUsers {
		urls = append(urls, fmt.Sprintf("https://gitlab.com/%s.keys", user))
	}

	return urls
}
